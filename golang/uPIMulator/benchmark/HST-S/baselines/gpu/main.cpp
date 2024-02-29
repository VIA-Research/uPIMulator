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

#include "support/cuda-setup.h"
#include "kernel.h"
#include "support/common.h"
#include "support/timer.h"
#include "support/verify.h"

#include <unistd.h>
#include <thread>
#include <assert.h>

// Params ---------------------------------------------------------------------
struct Params {

    int   device;
    int   n_gpu_threads;
    int   n_gpu_blocks;
    int   n_threads;
    int   n_warmup;
    int   n_reps;
    float alpha;
    int   in_size;
    int   n_bins;

    Params(int argc, char **argv) {
        device        = 0;
        n_gpu_threads  = 256;
        n_gpu_blocks = 16;
        n_threads     = 4;
        n_warmup      = 5;
        n_reps        = 50;
        alpha         = 0.2;
        in_size       = 1536 * 1024 * 640;
        n_bins        = 256;
        int opt;
        while((opt = getopt(argc, argv, "hd:i:g:t:w:r:a:n:b:")) >= 0) {
            switch(opt) {
            case 'h':
                usage();
                exit(0);
                break;
            case 'd': device        = atoi(optarg); break;
            case 'i': n_gpu_threads  = atoi(optarg); break;
            case 'g': n_gpu_blocks = atoi(optarg); break;
            case 't': n_threads     = atoi(optarg); break;
            case 'w': n_warmup      = atoi(optarg); break;
            case 'r': n_reps        = atoi(optarg); break;
            case 'a': alpha         = atof(optarg); break;
            case 'n': in_size       = atoi(optarg); break;
            case 'b': n_bins        = atoi(optarg); break;
            default:
                fprintf(stderr, "\nUnrecognized option!\n");
                usage();
                exit(0);
            }
        }
        if(alpha == 0.0) {
            assert(n_gpu_threads > 0 && "Invalid # of device threads!");
            assert(n_gpu_blocks > 0 && "Invalid # of device blocks!");
        } else if(alpha == 1.0) {
            assert(n_threads > 0 && "Invalid # of host threads!");
        } else if(alpha > 0.0 && alpha < 1.0) {
            assert(n_gpu_threads > 0 && "Invalid # of device threads!");
            assert(n_gpu_blocks > 0 && "Invalid # of device blocks!");
            assert(n_threads > 0 && "Invalid # of host threads!");
        } else {
#ifdef CUDA_8_0
            assert((n_gpu_threads > 0 && n_gpu_blocks > 0 || n_threads > 0) && "Invalid # of host + device workers!");
#else
            assert(0 && "Illegal value for -a");
#endif
        }
    }

    void usage() {
        fprintf(stderr,
                "\nUsage:  ./hsti [options]"
                "\n"
                "\nGeneral options:"
                "\n    -h        help"
                "\n    -d <D>    CUDA device ID (default=0)"
                "\n    -i <I>    # of device threads per block (default=256)"
                "\n    -g <G>    # of device blocks (default=16)"
                "\n    -t <T>    # of host threads (default=4)"
                "\n    -w <W>    # of untimed warmup iterations (default=5)"
                "\n    -r <R>    # of timed repetition iterations (default=50)"
                "\n"
                "\nData-partitioning-specific options:"
                "\n    -a <A>    fraction of input elements to process on host (default=0.2)"
#ifdef CUDA_8_0
                "\n              NOTE: Dynamic partitioning used when <A> is not between 0.0 and 1.0"
#else
                "\n              NOTE: <A> must be between 0.0 and 1.0"
#endif
                "\n"
                "\nBenchmark-specific options:"
                "\n    -n <N>    input size (default=1572864, i.e., 1536x1024)"
                "\n    -b <B>    # of bins in histogram (default=256)"
                "\n");
    }
};

// Input Data -----------------------------------------------------------------
void read_input(unsigned int *input, const Params &p) {

    char  dctFileName[100];
    FILE *File = NULL;

    // Open input file
    unsigned short temp;
    sprintf(dctFileName, "./input/image_VanHateren.iml");
    if((File = fopen(dctFileName, "rb")) != NULL) {
        for(int y = 0; y < p.in_size; y++) {
            int fr   = fread(&temp, sizeof(unsigned short), 1, File);
            input[y] = (unsigned int)ByteSwap16(temp);
            if(input[y] >= 4096)
                input[y] = 4095;
        }
        fclose(File);
    } else {
        printf("%s does not exist\n", dctFileName);
        exit(1);
    }
}

// Main ------------------------------------------------------------------------------------------
int main(int argc, char **argv) {

    Params p(argc, argv);
    CUDASetup    setcuda(p.device);
    Timer        timer;
    cudaError_t  cudaStatus;

    // Allocate buffers
    timer.start("Allocation");
    int n_tasks = divceil(p.in_size, p.n_gpu_threads);
#ifdef CUDA_8_0
    unsigned int *h_in;
    cudaStatus = cudaMallocManaged(&h_in, p.in_size * sizeof(unsigned int));
    std::atomic_uint *h_histo;
    cudaStatus = cudaMallocManaged(&h_histo, p.n_bins * sizeof(std::atomic_uint));
    unsigned int *    d_in     = h_in;
    std::atomic_uint *d_histo  = h_histo;
    std::atomic_int * worklist;
    cudaStatus = cudaMallocManaged(&worklist, sizeof(std::atomic_int));
#else
    unsigned int *    h_in          = (unsigned int *)malloc(p.in_size * sizeof(unsigned int));
    std::atomic_uint *h_histo       = (std::atomic_uint *)malloc(p.n_bins * sizeof(std::atomic_uint));
    unsigned int *    h_histo_merge = (unsigned int *)malloc(p.n_bins * sizeof(unsigned int));
    unsigned int *    d_in;
    cudaStatus = cudaMalloc((void**)&d_in, p.in_size * sizeof(unsigned int));
    unsigned int *    d_histo;
    cudaStatus = cudaMalloc((void**)&d_histo, p.n_bins * sizeof(unsigned int));
    ALLOC_ERR(h_in, h_histo, h_histo_merge);
#endif
    CUDA_ERR();
    cudaDeviceSynchronize();
    timer.stop("Allocation");
    timer.print("Allocation", 1);

    // Initialize
    timer.start("Initialization");
    const int max_gpu_threads = setcuda.max_gpu_threads();
    read_input(h_in, p);
#ifdef CUDA_8_0
    for(int i = 0; i < p.n_bins; i++) {
        h_histo[i].store(0);
    }
#else
    memset(h_histo, 0, p.n_bins * sizeof(unsigned int));
#endif
    cudaDeviceSynchronize();
    timer.stop("Initialization");
    timer.print("Initialization", 1);

#ifndef CUDA_8_0
    // Copy to device
    timer.start("Copy To Device");
    cudaStatus = cudaMemcpy(d_in, h_in, p.in_size * sizeof(unsigned int), cudaMemcpyHostToDevice);
    cudaStatus = cudaMemcpy(d_histo, h_histo, p.n_bins * sizeof(unsigned int), cudaMemcpyHostToDevice);
    cudaDeviceSynchronize();
    CUDA_ERR();
    timer.stop("Copy To Device");
    timer.print("Copy To Device", 1);
#endif

    // Loop over main kernel
    for(int rep = 0; rep < p.n_warmup + p.n_reps; rep++) {

        // Reset
#ifdef CUDA_8_0
        if(p.alpha < 0.0 || p.alpha > 1.0) { // Dynamic partitioning
            worklist[0].store(0);
        }
        for(int i = 0; i < p.n_bins; i++) {
            h_histo[i].store(0);
        }
#else
        memset(h_histo, 0, p.n_bins * sizeof(unsigned int));
        cudaStatus = cudaMemcpy(d_histo, h_histo, p.n_bins * sizeof(unsigned int), cudaMemcpyHostToDevice);
        cudaDeviceSynchronize();
        CUDA_ERR();
#endif

        if(rep >= p.n_warmup)
            timer.start("Kernel");

        p.n_gpu_blocks = p.in_size / p.n_gpu_threads;

        // Launch GPU threads
        // Kernel launch
        if(p.n_gpu_blocks > 0) {
            assert(p.n_gpu_threads <= max_gpu_threads && 
                "The runtime block size is greater than the maximum runtime block size that can be used on this device");
            cudaStatus = call_Histogram_kernel(p.n_gpu_blocks, p.n_gpu_threads, p.in_size, p.n_bins, n_tasks, 
                p.alpha, d_in, (unsigned int*)d_histo, p.n_bins * sizeof(unsigned int)
#ifdef CUDA_8_0
                + sizeof(int), (int*)worklist
#endif
                );
            CUDA_ERR();
        }

        // Launch CPU threads
        std::thread main_thread(run_cpu_threads, h_histo, h_in, p.in_size, p.n_bins, p.n_threads, p.n_gpu_threads,
            n_tasks, p.alpha
#ifdef CUDA_8_0
            , worklist
#endif
            );

        cudaDeviceSynchronize();
        main_thread.join();

        if(rep >= p.n_warmup)
            timer.stop("Kernel");
    }
    timer.print("Kernel", p.n_reps);

#ifndef CUDA_8_0
    // Copy back
    timer.start("Copy Back and Merge");
    cudaStatus = cudaMemcpy(h_histo_merge, d_histo, p.n_bins * sizeof(unsigned int), cudaMemcpyDeviceToHost);
    CUDA_ERR();
    cudaDeviceSynchronize();
    for(unsigned int i = 0; i < p.n_bins; ++i) {
        h_histo_merge[i] += (unsigned int)h_histo[i];
    }
    timer.stop("Copy Back and Merge");
    timer.print("Copy Back and Merge", 1);
#endif

    // Verify answer
#ifdef CUDA_8_0
    verify((unsigned int *)h_histo, h_in, p.in_size, p.n_bins);
#else
    verify((unsigned int *)h_histo_merge, h_in, p.in_size, p.n_bins);
#endif

    // Free memory
    timer.start("Deallocation");
#ifdef CUDA_8_0
    cudaStatus = cudaFree(h_in);
    cudaStatus = cudaFree(h_histo);
    cudaStatus = cudaFree(worklist);
#else
    free(h_in);
    free(h_histo);
    free(h_histo_merge);
    cudaStatus = cudaFree(d_in);
    cudaStatus = cudaFree(d_histo);
#endif
    CUDA_ERR();
    cudaDeviceSynchronize();
    timer.stop("Deallocation");
    timer.print("Deallocation", 1);

    // Release timers
    timer.release("Allocation");
    timer.release("Initialization");
    timer.release("Copy To Device");
    timer.release("Kernel");
    timer.release("Copy Back and Merge");
    timer.release("Deallocation");

    printf("Test Passed\n");
    return 0;
}
