#include <cuda.h>
#include <limits.h>
#include "binary_search.h"

#include <chrono>
#include <iostream>

#define BLOCKDIM 512
#define SEARCH_CHUNK 16
#define BLOCK_CHUNK (BLOCKDIM*SEARCH_CHUNK)


__global__ void search_kernel(const long int *arr,
    const long int len, const long int *querys, const long int num_querys, long int *res, bool *flag)
{
    int search;
    if(*flag == false) {
        int tid = threadIdx.x;
        __shared__ int s_arr[BLOCK_CHUNK];

        /* Since each value is being copied to shared memory, the rest of the
        following uncommented code is unncessary, since a direct comparison
        can be done at the time of copy below. */
        // for(int i = 0; i < BLOCKDIM; ++i) {
        //     int shared_loc = i*SEARCH_CHUNK + tid;
        //     int global_loc = shared_loc + BLOCK_CHUNK * blockIdx.x;
        //     if(arr[global_loc] == search) {
        //         *flag = true;
        //         *res = global_loc;
        //     }
        //     __syncthreads();
        // }

        /* Copy chunk of array that this entire block of threads will read
        from the slower global memory to the faster shared memory. */
        for(long int i = 0; i < SEARCH_CHUNK; ++i) {
            int shared_loc = tid*SEARCH_CHUNK + i;
            int global_loc = shared_loc + BLOCK_CHUNK * blockIdx.x;

            /* Make sure to stay within the bounds of the global array,
            else assign a dummy value. */
            if(global_loc < len) {
              s_arr[shared_loc] = arr[global_loc];
            }
            else {
              s_arr[shared_loc] = INT_MAX;
            }
        }
        __syncthreads();

        for(long int i = 0; i < num_querys; i++)
        {
            search = querys[i];
            /* For each runtime, set the initial search range. */
            int L = 0;
            int R = SEARCH_CHUNK - 1;
            int m = (L + R) / 2;

            /* Pointer to the part of the shared array for this runtime. */
            int *s_ptr = &s_arr[tid*SEARCH_CHUNK];

            /* Each runtime will search a chunk of the block array.
            Many blocks will not find a solution so the search must
            be allowed to fail on a per block basis. The loop will
            break (fail) when L >= R. */
            while(L <= R && *flag == false)
            {
                if(s_ptr[m] < search) {
                    L = m + 1;
                }
                else if(s_ptr[m] > search) {
                    R = m - 1;
                }
                else {
                    *flag = true;
                    *res = m += tid*SEARCH_CHUNK + BLOCK_CHUNK * blockIdx.x;
                }

                m = (L + R) / 2;
            }
        }
    }
}



int binary_search(const long int *arr, const long int len, const long int *querys, const long int num_querys)
{
    long int *d_arr, *d_querys, *d_res;
    bool *d_flag;

    size_t arr_size = len * sizeof(long int);
    size_t querys_size = num_querys * sizeof(long int);
    size_t res_size = sizeof(long int);
    size_t flag_size = sizeof(bool);

    cudaMalloc(&d_arr, arr_size);
    cudaMalloc(&d_querys, querys_size);
    cudaMalloc(&d_res, res_size);
    cudaMalloc(&d_flag, flag_size);

    cudaMemcpy(d_arr, arr, arr_size, cudaMemcpyHostToDevice);
    cudaMemcpy(d_querys, querys, querys_size, cudaMemcpyHostToDevice);
    cudaMemset(d_flag, 0, flag_size);

    /* Set res value to -1, so that if the function returns -1, that
    indicates an algorithm failure. */
    cudaMemset(d_res, -0x1, res_size);

    int blockSize = BLOCKDIM;
    int gridSize = (len-1)/BLOCK_CHUNK + 1;

    auto start = std::chrono::high_resolution_clock::now();
    search_kernel<<<gridSize,blockSize>>>(d_arr, len, d_querys, num_querys ,d_res, d_flag);
    cudaDeviceSynchronize();
    auto end = std::chrono::high_resolution_clock::now();
    std::cout << "Kernel Time: " <<
        std::chrono::duration_cast<std::chrono::milliseconds>(end-start).count() <<
        " ms" << std::endl;

    long int res;
    cudaMemcpy(&res, d_res, res_size, cudaMemcpyDeviceToHost);

    return res;
}
