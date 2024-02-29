/**
* app.c
* CPU-DPU Communication Host Application Source File
*
*/
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <dpu.h>
#include <dpu_log.h>
#include <unistd.h>
#include <getopt.h>
#include <assert.h>

#include "../support/common.h"
#include "../support/timer.h"
#include "../support/params.h"

// Define the DPU Binary path as DPU_BINARY here
#ifndef DPU_BINARY
#define DPU_BINARY "./bin/dpu_code"
#endif

// Pointer declaration
static T* A;
static T* B;
static T* C;
static T* C2;

// Create input arrays
static void read_input(T* A, T* B, unsigned int nr_elements) {
    srand(0);
    printf("nr_elements\t%u\t", nr_elements);
    for (unsigned int i = 0; i < nr_elements; i++) {
        A[i] = (T) (rand());
        B[i] = A[i];
    }
}

// Main of the Host Application
int main(int argc, char **argv) {

    struct Params p = input_params(argc, argv);

    struct dpu_set_t dpu_set, dpu;
    uint32_t nr_of_dpus;
    
    // Allocate DPUs and load binary
    DPU_ASSERT(dpu_alloc(NR_DPUS, "nrThreadPerPool=8", &dpu_set));
    DPU_ASSERT(dpu_load(dpu_set, DPU_BINARY, NULL));
    DPU_ASSERT(dpu_get_nr_dpus(dpu_set, &nr_of_dpus));
    printf("Allocated %d DPU(s)\n", nr_of_dpus);

    unsigned int i = 0;
    unsigned int input_size = p.exp == 0 ? p.input_size * nr_of_dpus : p.input_size;

    // Input/output allocation
    A = malloc(input_size * sizeof(T));
    B = malloc(input_size * sizeof(T));
    C = malloc(input_size * sizeof(T));
    C2 = malloc(input_size * sizeof(T));
    T *bufferA = A;
    T *bufferC = C;

    // Create an input file with arbitrary data
    read_input(A, B, input_size);

    // Timer declaration
    Timer timer;

    printf("NR_TASKLETS\t%d\tBL\t%d\n", NR_TASKLETS, BL);

    // Loop over main kernel
    for(int rep = 0; rep < p.n_warmup + p.n_reps; rep++) {

        printf("Load input data\n");
        // Input arguments
        const unsigned int input_size_dpu = input_size / nr_of_dpus;
        // Copy input arrays
        if(rep >= p.n_warmup)
            start(&timer, 1, rep - p.n_warmup);
        i = 0;
#ifdef SERIAL
        DPU_FOREACH (dpu_set, dpu) {
            DPU_ASSERT(dpu_copy_to(dpu, DPU_MRAM_HEAP_POINTER_NAME, 0, bufferA + input_size_dpu * i, input_size_dpu * sizeof(T)));
            i++;
        }
#elif BROADCAST
        DPU_ASSERT(dpu_broadcast_to(dpu_set, DPU_MRAM_HEAP_POINTER_NAME, 0, bufferA, input_size_dpu * sizeof(T), DPU_XFER_DEFAULT));
#else
        DPU_FOREACH(dpu_set, dpu, i) {
            DPU_ASSERT(dpu_prepare_xfer(dpu, bufferA + input_size_dpu * i));
        }
        DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, 0, input_size_dpu * sizeof(T), DPU_XFER_DEFAULT));
#endif
        if(rep >= p.n_warmup)
            stop(&timer, 1);

        printf("Run program on DPU(s) \n");
        // Run DPU kernel
        if(rep >= p.n_warmup)
            start(&timer, 2, rep - p.n_warmup);
        //DPU_ASSERT(dpu_launch(dpu_set, DPU_SYNCHRONOUS));
        if(rep >= p.n_warmup)
            stop(&timer, 2);

#if PRINT
        {
            unsigned int each_dpu = 0;
            printf("Display DPU Logs\n");
            DPU_FOREACH (dpu_set, dpu) {
                printf("DPU#%d:\n", each_dpu);
                DPU_ASSERT(dpulog_read_for_dpu(dpu.dpu, stdout));
                each_dpu++;
            }
        }
#endif

        printf("Retrieve results\n");
        if(rep >= p.n_warmup)
            start(&timer, 3, rep - p.n_warmup);
        i = 0;
#ifdef SERIAL
        DPU_FOREACH (dpu_set, dpu) {
            DPU_ASSERT(dpu_copy_from(dpu, DPU_MRAM_HEAP_POINTER_NAME, 0, bufferC + input_size_dpu * i, input_size_dpu * sizeof(T)));
            i++;
        }
#else
        DPU_FOREACH(dpu_set, dpu, i) {
            DPU_ASSERT(dpu_prepare_xfer(dpu, bufferC + input_size_dpu * i));
        }
        DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_FROM_DPU, DPU_MRAM_HEAP_POINTER_NAME, 0, input_size_dpu * sizeof(T), DPU_XFER_DEFAULT));
#endif
        if(rep >= p.n_warmup)
            stop(&timer, 3);

    }

    // Print timing results
    printf("CPU-DPU ");
    print(&timer, 1, p.n_reps);
    double time_load = timer.time[1] / (1000 * p.n_reps);
    printf("CPU-DPU Bandwidth (GB/s): %f\n", (input_size * 8)/(time_load*1e6));
    printf("DPU Kernel ");
    print(&timer, 2, p.n_reps);
    printf("\n");
    printf("DPU-CPU ");
    print(&timer, 3, p.n_reps);
    double time_retrieve = timer.time[3] / (1000 * p.n_reps);
    printf("DPU-CPU Bandwidth (GB/s): %f\n", (input_size * 8)/(time_retrieve*1e6));

    // Check output
    bool status = true;
#ifdef BROADCAST
    for (i = 0; i < input_size/nr_of_dpus; i++) {
        if(B[i] != bufferC[i]){ 
            status = false;
#if PRINT
            printf("%d: %u -- %u\n", i, B[i], bufferA[i]);
#endif
        }
    }
#else
    for (i = 0; i < input_size; i++) {
        if(B[i] != bufferC[i]){ 
            status = false;
#if PRINT
            printf("%d: %u -- %u\n", i, B[i], bufferA[i]);
#endif
        }
    }
#endif
    if (status) {
        printf("[" ANSI_COLOR_GREEN "OK" ANSI_COLOR_RESET "] Outputs are equal\n");
    } else {
        printf("[" ANSI_COLOR_RED "ERROR" ANSI_COLOR_RESET "] Outputs differ!\n");
    }

    // Deallocation
    free(A);
    free(B);
    free(C);
    free(C2);
    DPU_ASSERT(dpu_free(dpu_set));
	
    return status ? 0 : -1;
}
