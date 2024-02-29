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

#include "support/common.h"

extern __shared__ int l_mem[];

// GPU kernel ------------------------------------------------------------------------------------------
__global__ void PTTWAC_soa_asta(int A, int B, int b, T *input, int *finished, int *head) {

    int* done = l_mem;
    int* gid_ = &done[1];

    const int tid = threadIdx.x;
    int       m   = A * B - 1;

    if(tid == 0) // Dynamic fetch
        gid_[0] = atomicAdd(&head[0], 1);
    __syncthreads();

    while(gid_[0] < m) {
        int next_in_cycle = (gid_[0] * A) - m * (gid_[0] / B);
        if(next_in_cycle == gid_[0]) {
            if(tid == 0) // Dynamic fetch
                gid_[0] = atomicAdd(&head[0], 1);
            __syncthreads();
            continue;
        }
        T   data1, data2, data3, data4;
        int i = tid;
        if(i < b)
            data1 = input[gid_[0] * b + i];
        i += blockDim.x;
        if(i < b)
            data2 = input[gid_[0] * b + i];
        i += blockDim.x;
        if(i < b)
            data3 = input[gid_[0] * b + i];
        i += blockDim.x;
        if(i < b)
            data4 = input[gid_[0] * b + i];

        if(tid == 0) {
            //make sure the read is not cached
            done[0] = atomicAdd(&finished[gid_[0]], 0);
        }
        __syncthreads();

        for(; done[0] == 0; next_in_cycle = (next_in_cycle * A) - m * (next_in_cycle / B)) {
            T backup1, backup2, backup3, backup4;
            i = tid;
            if(i < b)
                backup1 = input[next_in_cycle * b + i];
            i += blockDim.x;
            if(i < b)
                backup2 = input[next_in_cycle * b + i];
            i += blockDim.x;
            if(i < b)
                backup3 = input[next_in_cycle * b + i];
            i += blockDim.x;
            if(i < b)
                backup4 = input[next_in_cycle * b + i];

            if(tid == 0) {
                done[0] = atomicExch(&finished[next_in_cycle], (int)1);
            }
            __syncthreads();

            if(!done[0]) {
                i = tid;
                if(i < b)
                    input[next_in_cycle * b + i] = data1;
                i += blockDim.x;
                if(i < b)
                    input[next_in_cycle * b + i] = data2;
                i += blockDim.x;
                if(i < b)
                    input[next_in_cycle * b + i] = data3;
                i += blockDim.x;
                if(i < b)
                    input[next_in_cycle * b + i] = data4;
            }
            i = tid;
            if(i < b)
                data1 = backup1;
            i += blockDim.x;
            if(i < b)
                data2 = backup2;
            i += blockDim.x;
            if(i < b)
                data3 = backup3;
            i += blockDim.x;
            if(i < b)
                data4 = backup4;
        }

        if(tid == 0) // Dynamic fetch
            gid_[0] = atomicAdd(&head[0], 1);
        __syncthreads();
    }
}

cudaError_t call_PTTWAC_soa_asta(int blocks, int threads, int A, int B, int b, T *input, 
    int *finished, int *head, int l_mem_size){
    dim3 dimGrid(blocks);
    dim3 dimBlock(threads);
    PTTWAC_soa_asta<<<dimGrid, dimBlock, l_mem_size>>>(A, B, b, input, 
        finished, head);
    cudaError_t err = cudaGetLastError();
    return err;
}

__global__ void BS_marshal(T *input, int tile_size, int width) {

  T* tile = (T*)l_mem;

  int tidx = threadIdx.x;
  int m = width*tile_size-1;
  int bid = blockIdx.x;

  input += tile_size*width*bid;
  for (int i = tidx; i < tile_size*width; i+=blockDim.x) {
    int next = (i * tile_size)-m*(i/width);
    tile[next] = input[i];
  }
  __syncthreads();
  for (int i = tidx; i < tile_size*width; i+=blockDim.x) {
    input[i] = tile[i];
  }
}

cudaError_t call_BS_marshal(int blocks, int threads, int m, int n, T *input, int l_mem_size){
    dim3 dimGrid(blocks);
    dim3 dimBlock(threads);
    BS_marshal<<<dimGrid, dimBlock, l_mem_size>>>(input, m, n);
    cudaError_t err = cudaGetLastError();
    return err;
}
