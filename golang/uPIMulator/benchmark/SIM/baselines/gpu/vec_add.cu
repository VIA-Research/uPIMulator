/* File:     vec_add.cu
 * Purpose:  Implement vector addition on a gpu using cuda
 *
 * Compile:  nvcc [-g] [-G] -o vec_add vec_add.cu
 * Run:      ./vec_add
 */

#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <math.h>

__global__ void Vec_add(unsigned int x[], unsigned int y[], unsigned int z[], int n) {
    int thread_id = blockIdx.x * blockDim.x + threadIdx.x;
    if (thread_id < n){
        z[thread_id] = x[thread_id] + y[thread_id];
    }
}


int main(int argc, char* argv[]) {
    int n, m;
    unsigned int *h_x, *h_y, *h_z;
    unsigned int *d_x, *d_y, *d_z;
    size_t size;

    /* Define vector length */
    n = 2621440;
    m = 320;
    size = m * n * sizeof(unsigned int);

    // Allocate memory for the vectors on host memory.
    h_x = (unsigned int*) malloc(size);
    h_y = (unsigned int*) malloc(size);
    h_z = (unsigned int*) malloc(size);

    for (int i = 0; i < n * m; i++) {
        h_x[i] = i+1;
        h_y[i] = n-i;
    }

    printf("Input size = %d\n", n * m);

    // Print original vectors.
    /*printf("h_x = ");
    for (int i = 0; i < m; i++){
        printf("%u ", h_x[i]);
    }
    printf("\n\n");
    printf("h_y = ");
    for (int i = 0; i < m; i++){
        printf("%u ", h_y[i]);
    }
    printf("\n\n");*/

    // Event creation
    cudaEvent_t start, stop;
    cudaEventCreate(&start);
    cudaEventCreate(&stop);
    float time1 = 0;

    /* Allocate vectors in device memory */
    cudaMalloc(&d_x, size);
    cudaMalloc(&d_y, size);
    cudaMalloc(&d_z, size);

    /* Copy vectors from host memory to device memory */
    cudaMemcpy(d_x, h_x, size, cudaMemcpyHostToDevice);
    cudaMemcpy(d_y, h_y, size, cudaMemcpyHostToDevice);
    
    // Start timer
    cudaEventRecord( start, 0 );

    /* Kernel Call */
    Vec_add<<<(n * m) / 256, 256>>>(d_x, d_y, d_z, n * m);

    // End timer
    cudaEventRecord( stop, 0 );
    cudaEventSynchronize( stop );
    cudaEventElapsedTime( &time1, start, stop );

    cudaMemcpy(h_z, d_z, size, cudaMemcpyDeviceToHost);
    /*printf("The sum is: \n");
    for (int i = 0; i < m; i++){
        printf("%u ", h_z[i]);
    }
    printf("\n");*/

    printf("Execution time = %f ms\n", time1);

    /* Free device memory */
    cudaFree(d_x);
    cudaFree(d_y);
    cudaFree(d_z);
    /* Free host memory */
    free(h_x);
    free(h_y);
    free(h_z);

    return 0;
}  /* main */
