#include <stdio.h>
#include <stdlib.h>
#include <sys/time.h>
#include <cuda.h>
#include "../../support/common.h"

#define THREAD 128

__global__ void gemv(int m, int n, T *adim, T *b, T *d_ans);

void cgemv(int m, int n, T *adim, T *b, T *d_ans);

double gettime()
{
	struct timeval tv;
	gettimeofday(&tv, NULL);
	return tv.tv_sec + (double)tv.tv_usec*1.0e-6;
}

int main(int argc, char **argv)
{
	/* for CPU */
	int i, j;
	T **bdim; 
	T *c, *ans, *h_ans, *h_c;
	int n = 8192;
	int m = 20480;

	bdim = (T**) malloc(NUM_LAYERS * sizeof(T*));
	for(int l = 0; l < NUM_LAYERS; l++)
		bdim[l] = (T*)malloc(sizeof(T)*m*n);
	c = (T*)malloc(sizeof(T) *n);
	h_c = (T*)malloc(sizeof(T) *n);
	ans = (T*)malloc(sizeof(T) *m);
	h_ans = (T*)malloc(sizeof(T) *m);

	/* for GPU */
	T *d_bdim; 
	T *d_c, *d_ans;
	cudaMalloc((void **)&d_bdim, sizeof(T)*m*n);
	cudaMalloc((void **)&d_c, sizeof(T)*n);
	cudaMalloc((void **)&d_ans, sizeof(T)*m);

	for(i = 0; i < n; i++)
	{
		if(i % 50 < 48)
		{
			c[i] = 0;
			h_c[i] = 0;
		}
		else
		{
			c[i] = i % 2;
			h_c[i] = i % 2;
		}
	}
	for(int l = 0; l < NUM_LAYERS; l++)
		for(i = 0; i < n; i++)
		{
			for(j = 0; j < m; j++){
				if(j % 100 < 98)
				{

					bdim[l][i*m+j] = 0;
				}
				else
				{

					bdim[l][i*m+j] = (l + i) % 2;
				}
			}
		}

	for(j = 0; j < m; j++){
		ans[j] = 0;
		h_ans[j] = 0;
	}
	// Computation on the host for verification
	T* vector = c;
	T* output = ans;
	T* matrix;
	int mm = m;
	int nn = n;
	for(int l = 0; l < NUM_LAYERS; l++){
		matrix = bdim[l];
		cgemv(mm, nn, matrix, vector, output);
		vector = output;
                h_ans = output;
		mm = n; nn = m;
	}

	// Event creation
	cudaEvent_t start, stop;
	cudaEventCreate(&start);
	cudaEventCreate(&stop);
	float time1 = 0;
	float time2 = 0;
	cudaMemcpy(d_ans, h_ans, sizeof(T)*m, cudaMemcpyHostToDevice);
	cudaMemcpy(d_c, h_c, sizeof(T)*n, cudaMemcpyHostToDevice);

	vector = d_c;
	output = d_ans;
	mm = m;
	nn = n;
	for(int l = 0; l < NUM_LAYERS; l++){
		cudaMemcpy(d_bdim, bdim[l], sizeof(T)*m*n, cudaMemcpyHostToDevice);
		matrix = d_bdim;
		// Start timer
		cudaEventRecord( start, 0 );
		gemv<<<mm, THREAD>>>(mm, nn, matrix, vector, output);
		// End timer
		cudaEventRecord( stop, 0 );
		cudaEventSynchronize( stop );
		cudaEventElapsedTime( &time2, start, stop );
		time1 += time2;
		vector = output;
		d_ans = output;
		mm = n; nn = m;
	}

	cudaMemcpy(h_ans, d_ans, sizeof(T)*m, cudaMemcpyDeviceToHost);
	cudaMemcpy(h_c, d_c, sizeof(T)*n, cudaMemcpyDeviceToHost);

	for(i = 0; i < m; i++)
	{
		if(ans[i] != h_ans[i])
		printf("ERROR in Ans %d -> %d -- %d\n", i, ans[i], h_ans[i]);
        }

	for(i = 0; i < n; i++)
	{
		if(c[i] != h_c[i])
		printf("ERROR in C %d -> %d -- %d\n", i, c[i], h_c[i]);
	}
	printf("Execution time = %f ms\n", time1);


	for(int l = 0; l < NUM_LAYERS; l++)
		free(bdim[l]);


	free(bdim);
	free(c);
	free(ans);
	free(h_c);
	cudaFree(d_bdim);
	cudaFree(d_c);
	cudaFree(d_ans);
	cudaEventDestroy(start);
	cudaEventDestroy(stop);

	return 0;
} 

__global__ void gemv(int m, int n, T* adim, T* b, T* d_ans)
{
	int i;
	int div = n/THREAD;
	__shared__ T tmp[THREAD];

	tmp[threadIdx.x] = 0.0;

	for(i = 0; i < div; i++){
		tmp[threadIdx.x] += adim[blockIdx.x*n+i*THREAD+threadIdx.x] * b[i * THREAD + threadIdx.x];
	}
	if(threadIdx.x < m%THREAD)
		tmp[threadIdx.x] += adim[blockIdx.x*n+THREAD*div+threadIdx.x] * b[THREAD * div + threadIdx.x];

	__syncthreads();

	for(i = THREAD / 2; i > 31; i = i / 2)
	{
		if(threadIdx.x < i)
			tmp[threadIdx.x] += tmp[threadIdx.x + i];
		__syncthreads();
	}

	if(threadIdx.x < 16)
	{
		tmp[threadIdx.x] += tmp[threadIdx.x + 16];
		__syncthreads();
		tmp[threadIdx.x] += tmp[threadIdx.x + 8];
		__syncthreads();
		tmp[threadIdx.x] += tmp[threadIdx.x + 4];
		__syncthreads();
		tmp[threadIdx.x] += tmp[threadIdx.x + 2];
		__syncthreads();
		tmp[threadIdx.x] += tmp[threadIdx.x + 1];
		__syncthreads();
	}


	if(threadIdx.x == 0)
		d_ans[blockIdx.x] = max(0, tmp[0]);

}

void cgemv(int m, int n, T *adim, T *b, T *d_ans)
{
	int i, j;

	for(i = 0; i < m; i++){
		for(j = 0; j < n; j++)
			d_ans[i] += adim[i*n+j] * b[j];
		d_ans[i] = max(0, d_ans[i]);
	}

}
