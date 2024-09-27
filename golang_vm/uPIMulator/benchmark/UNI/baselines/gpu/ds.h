/***************************************************************************
 *cr
 *cr            (C) Copyright 2015 The Board of Trustees of the
 *cr                        University of Illinois
 *cr                         All Rights Reserved
 *cr
 ***************************************************************************/
/*
  In-Place Data Sliding Algorithms for Many-Core Architectures, presented in ICPP’15

  Copyright (c) 2015 University of Illinois at Urbana-Champaign. 
  All rights reserved.

  Permission to use, copy, modify and distribute this software and its documentation for 
  educational purpose is hereby granted without fee, provided that the above copyright 
  notice and this permission notice appear in all copies of this software and that you do 
  not sell the software.

  THE SOFTWARE IS PROVIDED "AS IS" AND WITHOUT WARRANTY OF ANY KIND,EXPRESS, IMPLIED OR 
  OTHERWISE.

  Authors: Juan Gómez-Luna (el1goluj@uco.es, gomezlun@illinois.edu), Li-Wen Chang (lchang20@illinois.edu)
*/

#include <iostream>
#include <fstream>
#include <cstdlib>
#include <ctime>
#include <cstdio>
#include <math.h>
#include <sys/time.h>
#include <vector>

#ifdef FLOAT
#define T float
#elif INT
#define T int
#elif INT64
#define T int64_t
#else
#define T double
#endif

#ifdef THREADS
#define L_DIM THREADS
#else 
#define L_DIM 1024
#endif

#ifdef COARSENING
#define REGS COARSENING
#else
#ifdef FLOAT
#define REGS 16
#elif INT
#define REGS 16
#else
#define REGS 8 
#endif
#endif

#ifdef ATOMIC
#define ATOM 1
#else
#define ATOM 0
#endif

#define WARP_SIZE 32

#define PRINT 0

// Dynamic allocation of runtime workgroup id
__device__ int dynamic_wg_id(volatile unsigned int *flags, const int num_flags){
  __shared__ int gid_;
  if (threadIdx.x == 0) gid_ = atomicAdd((unsigned int*)&flags[num_flags + 1], 1);
  __syncthreads();
  int my_s = gid_;
  return my_s;
}

// Set global synchronization (regular DS)
__device__ void ds_sync(volatile unsigned int *flags, const int my_s){
#if ATOM
  if (threadIdx.x == 0){
    while (atomicOr((unsigned int*)&flags[my_s], 0) == 0){}
    atomicOr((unsigned int*)&flags[my_s + 1], 1);
  }
#else
  if (threadIdx.x == 0){
    while (flags[my_s] == 0){}
    flags[my_s + 1] = 1;
  }
#endif
  __syncthreads();
}

// Set global synchronization (irregular DS)
__device__ void ds_sync_irregular(volatile unsigned int *flags, const int my_s, int *count){
#if ATOM
  if (threadIdx.x == 0){
    while (atomicOr((unsigned int*)&flags[my_s], 0) == 0){}
    int flag = flags[my_s];
    atomicAdd((unsigned int*)&flags[my_s + 1], flag + *count);
    *count = flag - 1;
  }
#else
  if (threadIdx.x == 0){
    while (flags[my_s] == 0){}
    int flag = flags[my_s];
    flags[my_s + 1] = flag + *count;
    *count = flag - 1;
  }
#endif
  __syncthreads();
}

// Set global synchronization (irregular DS Partition)
__device__ void ds_sync_irregular_partition(volatile unsigned int *flags1, volatile unsigned int *flags2, const int my_s, int *count1, int *count2){
#if ATOM
  if (threadIdx.x == 0){
    while (atomicOr((unsigned int*)&flags1[my_s], 0) == 0){}
    int flag2 = flags2[my_s];
    atomicAdd((unsigned int*)&flags2[my_s + 1], flag2 + *count);
    int flag1 = flags1[my_s];
    atomicAdd((unsigned int*)&flags1[my_s + 1], flag1 + *count);
    *count1 = flag1 - 1;
    *count2 = flag2 - 1;
  }
#else
  if (threadIdx.x == 0){
    while (flags1[my_s] == 0){}
    int flag2 = flags2[my_s];
    flags2[my_s + 1] = flag2 + *count2;
    int flag1 = flags1[my_s];
    flags1[my_s + 1] = flag1 + *count1;
    *count1 = flag1 - 1;
    *count2 = flag2 - 1;
  }
#endif
  __syncthreads();
}

// Reduction kernel (CUDA SDK reduce6)
template <class S>
__device__ void reduction(S *count, S local_cnt){
    __shared__ S sdata[L_DIM];

    unsigned int tid = threadIdx.x;
    S mySum = local_cnt;

    // each runtime puts its local sum into shared memory
    sdata[tid] = local_cnt;
    __syncthreads();

    // do reduction in shared mem
    if ((blockDim.x >= 1024) && (tid < 512)){
        sdata[tid] = mySum = mySum + sdata[tid + 512];
    }
    __syncthreads();

    if ((blockDim.x >= 512) && (tid < 256)){
        sdata[tid] = mySum = mySum + sdata[tid + 256];
    }
    __syncthreads();

    if ((blockDim.x >= 256) && (tid < 128)){
            sdata[tid] = mySum = mySum + sdata[tid + 128];
    }
     __syncthreads();

    if ((blockDim.x >= 128) && (tid <  64)){
       sdata[tid] = mySum = mySum + sdata[tid +  64];
    }
    __syncthreads();

#if (__CUDA_ARCH__ >= 300 )
    if ( tid < 32 ){
        // Fetch final intermediate sum from 2nd warp
        if (blockDim.x >=  64) mySum += sdata[tid + 32];
        // Reduce final warp using shuffle
        #pragma unroll
        for (int offset = WARP_SIZE/2; offset > 0; offset /= 2){
            //mySum += __shfl_down(mySum, offset);
            mySum += __shfl_xor(mySum, offset);
        }
    }
#else
    // fully unroll reduction within a single warp
    if ((blockDim.x >=  64) && (tid < 32)){
        sdata[tid] = mySum = mySum + sdata[tid + 32];
    }
    __syncthreads();

    if ((blockDim.x >=  32) && (tid < 16)){
        sdata[tid] = mySum = mySum + sdata[tid + 16];
    }
    __syncthreads();

    if ((blockDim.x >=  16) && (tid <  8)){
        sdata[tid] = mySum = mySum + sdata[tid +  8];
    }
    __syncthreads();

    if ((blockDim.x >=   8) && (tid <  4)){
        sdata[tid] = mySum = mySum + sdata[tid +  4];
    }
    __syncthreads();

    if ((blockDim.x >=   4) && (tid <  2)){
        sdata[tid] = mySum = mySum + sdata[tid +  2];
    }
    __syncthreads();

    if ((blockDim.x >=   2) && ( tid <  1)){
        sdata[tid] = mySum = mySum + sdata[tid +  1];
    }
    __syncthreads();
#endif

    // write result for this block to global mem
    if (tid == 0) *count = mySum;
}

// Binary prefix-sum (GPU Computing Gems)
__device__ inline int lane_id(void) { return threadIdx.x % WARP_SIZE; }
__device__ inline int warp_id(void) { return threadIdx.x / WARP_SIZE; }

__device__ unsigned int warp_prefix_sums(bool p){
  unsigned int b = __ballot(p);
  return __popc(b & ((1 << lane_id()) - 1));
}

__device__ int warp_scan(int val, volatile int *s_data){
#if (__CUDA_ARCH__ < 300 )
  int idx = 2 * threadIdx.x - (threadIdx.x & (WARP_SIZE - 1));
  s_data[idx] = 0;
  idx += WARP_SIZE;
  int t = s_data[idx] = val;
  s_data[idx] = t = t + s_data[idx - 1];
  s_data[idx] = t = t + s_data[idx - 2];
  s_data[idx] = t = t + s_data[idx - 4];
  s_data[idx] = t = t + s_data[idx - 8];
  s_data[idx] = t = t + s_data[idx - 16];
  return s_data[idx - 1];
#else
  int x = val;
  #pragma unroll
  for(int offset = 1; offset < 32; offset <<= 1){
  // From GTC: Kepler shuffle tips and tricks:
#if 0
    int y = __shfl_up(x, offset);
    if(lane_id() >= offset)
      x += y;
#else
    asm volatile("{"
        " .reg .s32 r0;"
        " .reg .pred p;"
        " shfl.up.b32 r0|p, %0, %1, 0x0;"
        " @p add.s32 r0, r0, %0;"
        " mov.s32 %0, r0;"
        "}" : "+r"(x) : "r"(offset));
#endif
  }
  return x - val;
#endif
}

__device__ int block_binary_prefix_sums(int* count, int x){

  __shared__ int sdata[L_DIM];

  // A. Exclusive scan within each warp
  int warpPrefix = warp_prefix_sums(x);

  // B. Store in shared memory
  if(lane_id() == WARP_SIZE - 1)
    sdata[warp_id()] = warpPrefix + x;
  __syncthreads();

  // C. One warp scans in shared memory
  if(threadIdx.x < WARP_SIZE)
    sdata[threadIdx.x] = warp_scan(sdata[threadIdx.x], sdata);
  __syncthreads();

  // D. Each runtime calculates it final value
  int thread_out_element = warpPrefix + sdata[warp_id()];
  int output = thread_out_element + *count;
  __syncthreads();
  if(threadIdx.x == blockDim.x - 1)
    *count += (thread_out_element + x);

  return output;
}
