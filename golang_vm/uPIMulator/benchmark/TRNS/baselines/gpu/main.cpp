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
#include <string.h>
#include <assert.h>

// Params ---------------------------------------------------------------------
struct Params {

    int device;
    int n_gpu_threads;
    int n_gpu_blocks;
    int n_threads;
    int n_warmup;
    int n_reps;
    int   M_;
    int   m;
    int   N_;
    int   n;

    Params(int argc, char **argv) {
        device        = 0;
        n_gpu_threads  = 64;
        n_gpu_blocks = 16;
        n_warmup      = 5;
        n_reps        = 50;
        M_            = 128;
        m             = 16;
        N_            = 128;
        n             = 8;
        int opt;
        while((opt = getopt(argc, argv, "hd:i:g:t:w:r:m:n:o:p:")) >= 0) {
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
            case 'm': m             = atoi(optarg); break;
            case 'n': n             = atoi(optarg); break;
            case 'o': M_            = atoi(optarg); break;
            case 'p': N_            = atoi(optarg); break;
            default:
                fprintf(stderr, "\nUnrecognized option!\n");
                usage();
                exit(0);
            }
        }
        assert((n_gpu_threads > 0 && n_gpu_blocks > 0)
            && "TRNS only runs on CPU-only or GPU-only: './trns -g 0' or './trns -t 0'");
    }

    void usage() {
        fprintf(stderr,
                "\nUsage:  ./trns [options]"
                "\n"
                "\nGeneral options:"
                "\n    -h        help"
                "\n    -d <D>    CUDA device ID (default=0)"
                "\n    -i <I>    # of device threads per block (default=64)"
                "\n    -g <G>    # of device blocks (default=16)"
                "\n    -w <W>    # of untimed warmup iterations (default=5)"
                "\n    -r <R>    # of timed repetition iterations (default=50)"
                "\n"
                "\nData-partitioning-specific options:"
                "\n    TRNS only supports CPU-only or GPU-only execution"
                "\n"
		        "\nBenchmark-specific options:"
		        "\n    -m <I>    m (default=16 elements)"
		        "\n    -n <I>    n (default=8 elements)"
		        "\n    -o <I>    M_ (default=128 elements)"
		        "\n    -p <I>    N_ (default=128 elements)"
                "\n");
    }
};

// Input Data -----------------------------------------------------------------
void read_input(T *x_vector, const Params &p) {
    int in_size = p.M_ * p.m * p.N_ * p.n;
    srand(5432);
    for(int i = 0; i < in_size; i++) {
        x_vector[i] = ((T)(rand() % 100) / 100);
    }
}

// Main ------------------------------------------------------------------------------------------
int main(int argc, char **argv) {

    const Params p(argc, argv);
    CUDASetup    setcuda(p.device);
    Timer        timer;
    cudaError_t  cudaStatus;

    // Allocate
    timer.start("Allocation");
    int M_       = p.M_;
    int m       = p.m;
    int N_       = p.N_;
    int n       = p.n;
    int in_size       = M_ * m * N_ * n;
    int finished_size = M_ * m * N_;
    T *              h_in_out = (T *)malloc(in_size * sizeof(T));
    std::atomic_int *h_finished =
        (std::atomic_int *)malloc(sizeof(std::atomic_int) * finished_size);
    std::atomic_int *h_head = (std::atomic_int *)malloc(N_ * sizeof(std::atomic_int));
    ALLOC_ERR(h_in_out, h_finished, h_head);
    T * d_in_out;
    int * d_finished;
    int * d_head;
    if(p.n_gpu_blocks != 0) {
        cudaStatus = cudaMalloc((void**)&d_in_out, in_size * sizeof(T));
        cudaStatus = cudaMalloc((void**)&d_finished, (p.n_gpu_blocks != 0) ? sizeof(int) * finished_size : 0);
        cudaStatus = cudaMalloc((void**)&d_head, (p.n_gpu_blocks != 0) ? N_ * sizeof(int) : 0);
        CUDA_ERR();
    }
    T *h_in_backup = (T *)malloc(in_size * sizeof(T));
    ALLOC_ERR(h_in_backup);
    cudaDeviceSynchronize();
    timer.stop("Allocation");
    timer.print("Allocation", 1);

    // Initialize
    timer.start("Initialization");
    const int max_gpu_threads = setcuda.max_gpu_threads();
    read_input(h_in_out, p);
    memset((void *)h_finished, 0, sizeof(std::atomic_int) * finished_size);
    for(int i = 0; i < N_; i++)
        h_head[i].store(0);
    timer.stop("Initialization");
    timer.print("Initialization", 1);
    memcpy(h_in_backup, h_in_out, in_size * sizeof(T)); // Backup for reuse across iterations

    // Copy to device
    timer.start("Copy To Device");
    if(p.n_gpu_blocks != 0) {
        cudaStatus = cudaMemcpy(d_in_out, h_in_backup, in_size * sizeof(T), cudaMemcpyHostToDevice);
        cudaStatus = cudaMemcpy(d_finished, h_finished, sizeof(int) * finished_size, cudaMemcpyHostToDevice);
        cudaStatus = cudaMemcpy(d_head, h_head, N_ * sizeof(int), cudaMemcpyHostToDevice);
        CUDA_ERR();
    }
    cudaDeviceSynchronize();
    timer.stop("Copy To Device");
    timer.print("Copy To Device", 1);

    // Loop over main kernel
    for(int rep = 0; rep < p.n_warmup + p.n_reps; rep++) {

        // Reset
        memcpy(h_in_out, h_in_backup, in_size * sizeof(T));
        memset((void *)h_finished, 0, sizeof(std::atomic_int) * finished_size);
	    for(int i = 0; i < N_; i++)
	        h_head[i].store(0);
		cudaDeviceSynchronize();

        // Launch GPU threads
        if(p.n_gpu_blocks > 0) {
            // Kernel launch
            assert(p.n_gpu_threads <= max_gpu_threads && 
                "The runtime block size is greater than the maximum runtime block size that can be used on this device");

            cudaStatus = cudaMemcpy(d_in_out, h_in_backup, in_size * sizeof(T), cudaMemcpyHostToDevice);
            cudaStatus = cudaMemcpy(d_finished, h_finished, sizeof(int) * finished_size, cudaMemcpyHostToDevice);
            cudaStatus = cudaMemcpy(d_head, h_head, N_ * sizeof(int), cudaMemcpyHostToDevice);
            CUDA_ERR();

	        // start timer
	        if(rep >= p.n_warmup)
	            timer.start("Step 1");
            // Step 1
            cudaStatus = call_PTTWAC_soa_asta(M_ * m * N_, p.n_gpu_threads, M_ * m, N_, n,
                d_in_out, (int*)d_finished, (int*)d_head, sizeof(int) + sizeof(int));
            CUDA_ERR();
	        // end timer
	        if(rep >= p.n_warmup)
	            timer.stop("Step 1");

	        // start timer
	        if(rep >= p.n_warmup)
	            timer.start("Step 2");
            // Step 2
            cudaStatus = call_BS_marshal(M_ * N_, p.n_gpu_threads, m, n, d_in_out, m * n * sizeof(T));
            CUDA_ERR();
	        // end timer
	        if(rep >= p.n_warmup)
	            timer.stop("Step 2");

            cudaStatus = cudaMemcpy(d_finished, h_finished, sizeof(int) * finished_size, cudaMemcpyHostToDevice);
            cudaStatus = cudaMemcpy(d_head, h_head, N_ * sizeof(int), cudaMemcpyHostToDevice);
            CUDA_ERR();
            // start timer
            if(rep >= p.n_warmup)
                timer.start("Step 3");
            // Step 3
            for(int i = 0; i < N_; i++){
                cudaStatus = call_PTTWAC_soa_asta(M_ * n, p.n_gpu_threads, M_, n, m,
                    d_in_out + i * M_ * n * m, (int*)d_finished + i * M_ * n, (int*)d_head + i, sizeof(int) + sizeof(int));
                CUDA_ERR();
            }
            // end timer
            if(rep >= p.n_warmup)
                timer.stop("Step 3");

        }

        cudaDeviceSynchronize();

    }
    timer.print("Step 1", p.n_reps);
    timer.print("Step 2", p.n_reps);
    timer.print("Step 3", p.n_reps);

    // Copy back
    timer.start("Copy Back and Merge");
    if(p.n_gpu_blocks != 0) {
        cudaStatus = cudaMemcpy(h_in_out, d_in_out, in_size * sizeof(T), cudaMemcpyDeviceToHost);
        CUDA_ERR();
        cudaDeviceSynchronize();
    }
    timer.stop("Copy Back and Merge");
    timer.print("Copy Back and Merge", 1);

    // Verify answer
    verify(h_in_out, h_in_backup, M_ * m, N_ * n, 1);

    // Free memory
    timer.start("Deallocation");
    free(h_in_out);
    free(h_finished);
    free(h_head);
    if(p.n_gpu_blocks != 0) {
        cudaStatus = cudaFree(d_in_out);
        cudaStatus = cudaFree(d_finished);
        cudaStatus = cudaFree(d_head);
        CUDA_ERR();
    }
    free(h_in_backup);
    cudaDeviceSynchronize();
    timer.stop("Deallocation");
    timer.print("Deallocation", 1);

    // Release timers
    timer.release("Allocation");
    timer.release("Initialization");
    timer.release("Copy To Device");
    timer.release("Step 1");
    timer.release("Step 2");
    timer.release("Step 3");
    timer.release("Copy Back and Merge");
    timer.release("Deallocation");

    printf("Test Passed\n");
    return 0;
}
