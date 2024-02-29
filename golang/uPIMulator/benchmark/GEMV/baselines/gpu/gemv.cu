#include <stdio.h>
#include <stdlib.h>
#include <sys/time.h>
#include <cuda.h>

#define THREAD 128

#define T int

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
int *bdim, *c, *ans, *h_ans;
//double start, stop;
//double cpu_time, gpu_time;
int n = 8192;
int m = 20480;

bdim = (T*)malloc(sizeof(T) *m*n);
c = (T*)malloc(sizeof(T) *n);
ans = (T*)malloc(sizeof(T) *m);
h_ans = (T*)malloc(sizeof(T) *m);

/* for GPU */
T *d_bdim, *d_c, *d_ans;
cudaMalloc((void **)&d_bdim, sizeof(T)*m*n);
cudaMalloc((void **)&d_c, sizeof(T)*n);
cudaMalloc((void **)&d_ans, sizeof(T)*m);

for(i = 0; i < n; i++)
{
c[i] = 1;
for(j = 0; j < m; j++)
bdim[i*m+j] = 1;
}

//start = gettime();
cgemv(m, n, bdim, c, ans);
//stop = gettime();
//cpu_time=stop - start;

// Event creation
cudaEvent_t start, stop;
cudaEventCreate(&start);
cudaEventCreate(&stop);
float time1 = 0;


cudaMemcpy(d_bdim, bdim, sizeof(T)*m*n, cudaMemcpyHostToDevice);
cudaMemcpy(d_c, c, sizeof(T)*n, cudaMemcpyHostToDevice);

// Start timer
cudaEventRecord( start, 0 );
//start = gettime();
gemv<<<m, THREAD>>>(m, n, d_bdim, d_c, d_ans);
//stop = gettime();
// End timer
cudaEventRecord( stop, 0 );
cudaEventSynchronize( stop );
cudaEventElapsedTime( &time1, start, stop );

//gpu_time=stop - start;

cudaMemcpy(h_ans, d_ans, sizeof(T)*m, cudaMemcpyDeviceToHost);

//printf("cpu_time : %.6f[sec]\n",cpu_time);
//printf("gpu_time : %.6f[sec]\n",gpu_time);
//printf("%f x\n", cpu_time / gpu_time);


for(i = 0; i < m; i++)
printf("%d -- %d\n", ans[i], h_ans[i]);

printf("Execution time = %f ms\n", time1);


free(bdim);
free(c);
free(ans);
free(h_ans);
cudaFree(d_bdim);
cudaFree(d_c);
cudaFree(d_ans);

return 0;
} 

__global__ void gemv(int m, int n, T* adim, T* b, T* d_ans)
{
int i;
int div = n/THREAD;
__shared__ T tmp[THREAD];

tmp[threadIdx.x] = 0.0;

for(i = 0; i < div; i++)
{
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
d_ans[blockIdx.x] = tmp[0];

}

void cgemv(int m, int n, T *adim, T *b, T *d_ans)
{
int i, j;

for(i = 0; i < m; i++)
for(j = 0; j < n; j++)
d_ans[i] += adim[i*n+j] * b[j];

}
