/*
 * Copyright (c) 2016 University of Cordoba and University of Illinois
 * All rights reserved.
 *
 * Developed by:    IMPACT Research Group
 *                  University of Cordoba and University of Illinois
 *                  http://impact.crhc.illinois.edu/
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * with the Software without restriction, including without limitation the 
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
 * sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 *      > Redistributions of source code must retain the above copyright notice,
 *        this list of conditions and the following disclaimers.
 *      > Redistributions in binary form must reproduce the above copyright
 *        notice, this list of conditions and the following disclaimers in the
 *        documentation and/or other materials provided with the distribution.
 *      > Neither the names of IMPACT Research Group, University of Cordoba, 
 *        University of Illinois nor the names of its contributors may be used 
 *        to endorse or promote products derived from this Software without 
 *        specific prior written permission.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE 
 * CONTRIBUTORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS WITH
 * THE SOFTWARE.
 *
 */

#define _CUDA_COMPILER_

#include "support/common.h"
#include "support/partitioner.h"

// CUDA kernel ------------------------------------------------------------------------------------------
__global__ void Histogram_kernel(int size, int bins, int n_tasks, float alpha, unsigned int *data,
    unsigned int *histo
#ifdef CUDA_8_0
    , int *worklist
#endif
    ) {

    extern __shared__ unsigned int l_mem[];
    unsigned int* l_histo = l_mem;
#ifdef CUDA_8_0
    int* l_tmp = (int*)&l_histo[bins];
#endif
    
#ifdef CUDA_8_0
    Partitioner p = partitioner_create(n_tasks, alpha, worklist, l_tmp);
#else
    Partitioner p = partitioner_create(n_tasks, alpha);
#endif
    
    // Block and runtime index
    const int bx = blockIdx.x;
    const int tx = threadIdx.x;
    const int bD = blockDim.x;
    const int gD = gridDim.x;

    // Sub-histograms initialization
    for(int pos = tx; pos < bins; pos += bD) {
        l_histo[pos] = 0;
    }

    __syncthreads(); // Intra-block synchronization

    // Main loop
    for(int i = gpu_first(&p); gpu_more(&p); i = gpu_next(&p)) {
    
        // Global memory read
        unsigned int d = data[i * bD + tx];

        // Atomic vote in shared memory
        atomicAdd(&l_histo[((d * bins) >> 12)], 1);
    }

    __syncthreads(); // Intra-block synchronization

    // Merge per-block histograms and write to global memory
    for(int pos = tx; pos < bins; pos += bD) {
// Atomic addition in global memory
#ifdef CUDA_8_0
        atomicAdd_system(histo + pos, l_histo[pos]);
#else
        atomicAdd(histo + pos, l_histo[pos]);
#endif
    }
}

cudaError_t call_Histogram_kernel(int blocks, int threads, int size, int bins, int n_tasks, float alpha, 
    unsigned int *data, unsigned int *histo, int l_mem_size
#ifdef CUDA_8_0
    , int* worklist
#endif
    ){
    dim3 dimGrid(blocks);
    dim3 dimBlock(threads);
    Histogram_kernel<<<dimGrid, dimBlock, l_mem_size>>>(size, bins, n_tasks, alpha, 
        data, histo
#ifdef CUDA_8_0
        , worklist
#endif
        );
    cudaError_t err = cudaGetLastError();
    return err;
}
