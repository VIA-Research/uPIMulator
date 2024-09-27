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

#include <cuda.h>
#include <cuda_runtime.h>
#include <fstream>

// Allocation error checking
#define ERR_1(v1)                                                                                                      \
    if(v1 == NULL) {                                                                                                   \
        fprintf(stderr, "Allocation error at %s, %d\n", __FILE__, __LINE__);                                           \
        exit(-1);                                                                                                      \
    }
#define ERR_2(v1,v2) ERR_1(v1) ERR_1(v2)
#define ERR_3(v1,v2,v3) ERR_2(v1,v2) ERR_1(v3)
#define ERR_4(v1,v2,v3,v4) ERR_3(v1,v2,v3) ERR_1(v4)
#define ERR_5(v1,v2,v3,v4,v5) ERR_4(v1,v2,v3,v4) ERR_1(v5)
#define ERR_6(v1,v2,v3,v4,v5,v6) ERR_5(v1,v2,v3,v4,v5) ERR_1(v6)
#define GET_ERR_MACRO(_1,_2,_3,_4,_5,_6,NAME,...) NAME
#define ALLOC_ERR(...) GET_ERR_MACRO(__VA_ARGS__,ERR_6,ERR_5,ERR_4,ERR_3,ERR_2,ERR_1)(__VA_ARGS__)

#define CUDA_ERR()                                                                                                     \
    if(cudaStatus != cudaSuccess) {                                                                                    \
        fprintf(stderr, "CUDA error: %s\n at %s, %d\n", cudaGetErrorString(cudaStatus), __FILE__, __LINE__);           \
        exit(-1);                                                                                                      \
    }

struct CUDASetup {

    cudaDeviceProp device_prop;

    CUDASetup(int device) {
        cudaError_t cudaStatus;
        cudaStatus = cudaSetDevice(device);
        CUDA_ERR();

        cudaStatus = cudaGetDeviceProperties(&device_prop, device);
        CUDA_ERR();
        fprintf(stderr, "%s\t", device_prop.name);

    }

    int max_gpu_threads() {
        return device_prop.maxThreadsPerBlock;
    }
};
