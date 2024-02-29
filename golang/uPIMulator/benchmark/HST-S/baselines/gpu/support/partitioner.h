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

#ifndef _PARTITIONER_H_
#define _PARTITIONER_H_

#ifndef _CUDA_COMPILER_
#include <iostream>
#endif

#if !defined(_CUDA_COMPILER_) && defined(CUDA_8_0)
#include <atomic>
#endif

// Partitioner definition -----------------------------------------------------

typedef struct Partitioner {

    int n_tasks;
    int cut;
    int current;
#ifndef _CUDA_COMPILER_
    int thread_id;
    int n_threads;
#endif


#ifdef CUDA_8_0
    // CUDA 8.0 support for dynamic partitioning
    int strategy;
#ifdef _CUDA_COMPILER_
    int *worklist;
    int *tmp;
#else
    std::atomic_int *worklist;
#endif
#endif

} Partitioner;

// Partitioning strategies
#define STATIC_PARTITIONING 0
#define DYNAMIC_PARTITIONING 1

// Create a partitioner -------------------------------------------------------

#ifdef _CUDA_COMPILER_
__device__
#endif
inline Partitioner partitioner_create(int n_tasks, float alpha
#ifndef _CUDA_COMPILER_
    , int thread_id, int n_threads
#endif
#ifdef CUDA_8_0
#ifdef _CUDA_COMPILER_
    , int *worklist
    , int *tmp
#else
    , std::atomic_int *worklist
#endif
#endif
    ) {
    Partitioner p;
    p.n_tasks = n_tasks;
#ifndef _CUDA_COMPILER_
    p.thread_id = thread_id;
    p.n_threads = n_threads;
#endif
    if(alpha >= 0.0 && alpha <= 1.0) {
        p.cut = p.n_tasks * alpha;
#ifdef CUDA_8_0
        p.strategy = STATIC_PARTITIONING;
#endif
    } else {
#ifdef CUDA_8_0
        p.strategy = DYNAMIC_PARTITIONING;
        p.worklist = worklist;
#ifdef _CUDA_COMPILER_
        p.tmp = tmp;
#endif
#endif
    }
    return p;
}

// Partitioner iterators: first() ---------------------------------------------

#ifndef _CUDA_COMPILER_

inline int cpu_first(Partitioner *p) {
#ifdef CUDA_8_0
    if(p->strategy == DYNAMIC_PARTITIONING) {
        p->current = p->worklist->fetch_add(1);
    } else
#endif
    {
        p->current = p->thread_id;
    }
    return p->current;
}

#else

__device__ inline int gpu_first(Partitioner *p) {
#ifdef CUDA_8_0
    if(p->strategy == DYNAMIC_PARTITIONING) {
        if(threadIdx.y == 0 && threadIdx.x == 0) {
            p->tmp[0] = atomicAdd_system(p->worklist, 1);
        }
        __syncthreads();
        p->current = p->tmp[0];
    } else
#endif
    {
        p->current = p->cut + blockIdx.x;
    }
    return p->current;
}

#endif

// Partitioner iterators: more() ----------------------------------------------

#ifndef _CUDA_COMPILER_

inline bool cpu_more(const Partitioner *p) {
#ifdef CUDA_8_0
    if(p->strategy == DYNAMIC_PARTITIONING) {
        return (p->current < p->n_tasks);
    } else
#endif
    {
        return (p->current < p->cut);
    }
}

#else

__device__ inline bool gpu_more(const Partitioner *p) {
    return (p->current < p->n_tasks);
}

#endif

// Partitioner iterators: next() ----------------------------------------------

#ifndef _CUDA_COMPILER_

inline int cpu_next(Partitioner *p) {
#ifdef CUDA_8_0
    if(p->strategy == DYNAMIC_PARTITIONING) {
        p->current = p->worklist->fetch_add(1);
    } else
#endif
    {
        p->current = p->current + p->n_threads;
    }
    return p->current;
}

#else

__device__ inline int gpu_next(Partitioner *p) {
#ifdef CUDA_8_0
    if(p->strategy == DYNAMIC_PARTITIONING) {
        if(threadIdx.y == 0 && threadIdx.x == 0) {
            p->tmp[0] = atomicAdd_system(p->worklist, 1);
        }
        __syncthreads();
        p->current = p->tmp[0];
    } else
#endif
    {
        p->current = p->current + gridDim.x;
    }
    return p->current;
}

#endif

#endif
