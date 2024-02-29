/**
* app.c
* TRNS Host Application Source File
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
#include <math.h>

#include "../support/common.h"
#include "../support/timer.h"
#include "../support/params.h"

// Define the DPU Binary path as DPU_BINARY here
#ifndef DPU_BINARY
#define DPU_BINARY "./bin/dpu_code"
#endif

#if ENERGY
#include <dpu_probe.h>
#endif

// Pointer declaration
static T* A_host;
static T* A_backup;
static T* A_result;

// Create input arrays
static void read_input(T* A, unsigned int nr_elements) {
    srand(0);
    printf("nr_elements\t%u\t", nr_elements);
    for (unsigned int i = 0; i < nr_elements; i++) {
        A[i] = (T) (rand());
    }
}

// Compute output in the host
static void trns_host(T* input, unsigned int A, unsigned int B, unsigned int b){
   T* output = (T*) malloc(sizeof(T) * A * B * b);
   unsigned int next;
   for (unsigned int j = 0; j < b; j++){
      for (unsigned int i = 0; i < A * B; i++){
         next = (i * A) - (A * B - 1) * (i / B);
         output[next * b + j] = input[i*b+j];
      }
   }
   for (unsigned int k = 0; k < A * B * b; k++){
      input[k] = output[k];
   }
   free(output);
}

// Main of the Host Application
int main(int argc, char **argv) {

    struct Params p = input_params(argc, argv);

    struct dpu_set_t dpu_set, dpu;
    uint32_t nr_of_dpus;
    
#if ENERGY
    struct dpu_probe_t probe;
    DPU_ASSERT(dpu_probe_init("energy_probe", &probe));
#endif

    unsigned int i = 0;
    unsigned int N_ = p.N_;
    const unsigned int n = p.n;
    const unsigned int M_ = p.M_;
    const unsigned int m = p.m;
    N_ = p.exp == 0 ? N_ * NR_DPUS : N_;

    // Input/output allocation
    A_host = malloc(M_ * m * N_ * n * sizeof(T));
    A_backup = malloc(M_ * m * N_ * n * sizeof(T));
    A_result = malloc(M_ * m * N_ * n * sizeof(T));
    T* done_host = malloc(M_ * n); // Host array to reset done array of step 3
    memset(done_host, 0, M_ * n);

    // Create an input file with arbitrary data
    read_input(A_host, M_ * m * N_ * n);
    memcpy(A_backup, A_host, M_ * m * N_ * n * sizeof(T));

    // Timer declaration
    Timer timer;

    printf("NR_TASKLETS\t%d\n", NR_TASKLETS);
    printf("M_\t%u, m\t%u, N_\t%u, n\t%u\n", M_, m, N_, n);

    // Loop over main kernel
    for(int rep = 0; rep < p.n_warmup + p.n_reps; rep++) {

        int timer_fix = 0;
        // Compute output on CPU (performance comparison and verification purposes)
        memcpy(A_host, A_backup, M_ * m * N_ * n * sizeof(T));
        if(rep >= p.n_warmup)
            start(&timer, 0, rep - p.n_warmup + timer_fix);
        trns_host(A_host, M_ * m, N_ * n, 1);
        if(rep >= p.n_warmup)
            stop(&timer, 0);

        unsigned int curr_dpu = 0;
        unsigned int active_dpus;
        unsigned int active_dpus_before = 0;
        unsigned int first_round = 1;

        while(curr_dpu < N_){
            // Allocate DPUs and load binary
            if((N_ - curr_dpu) > NR_DPUS){
                active_dpus = NR_DPUS;
            } else {
                active_dpus = (N_ - curr_dpu);
            }
            if((active_dpus_before != active_dpus) && (!(first_round))){
                DPU_ASSERT(dpu_free(dpu_set));
                DPU_ASSERT(dpu_alloc(active_dpus, NULL, &dpu_set));
                DPU_ASSERT(dpu_load(dpu_set, DPU_BINARY, NULL));
                DPU_ASSERT(dpu_get_nr_dpus(dpu_set, &nr_of_dpus));
                printf("Allocated %d DPU(s)\n", nr_of_dpus);
            } else if (first_round){
                DPU_ASSERT(dpu_alloc(active_dpus, NULL, &dpu_set));
                DPU_ASSERT(dpu_load(dpu_set, DPU_BINARY, NULL));
                DPU_ASSERT(dpu_get_nr_dpus(dpu_set, &nr_of_dpus));
                printf("Allocated %d DPU(s)\n", nr_of_dpus);
            }

            printf("Load input data (step 1)\n");
            if(rep >= p.n_warmup)
                start(&timer, 1, rep - p.n_warmup + timer_fix);
            // Load input matrix (step 1)
            for(unsigned int j = 0; j < M_ * m; j++){
                unsigned int i = 0;
                DPU_FOREACH(dpu_set, dpu) {
                    DPU_ASSERT(dpu_prepare_xfer(dpu, &A_backup[j * N_ * n + n * (i + curr_dpu)]));
                    i++;
                }
                DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, sizeof(T) * j * n, sizeof(T) * n, DPU_XFER_DEFAULT));
            }
            if(rep >= p.n_warmup)
                stop(&timer, 1);
            // Reset done array (for step 3)
            DPU_FOREACH(dpu_set, dpu) {
                DPU_ASSERT(dpu_prepare_xfer(dpu, done_host));
            }
            DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, M_ * m * n * sizeof(T), (M_ * n) / 8 == 0 ? 8 : M_ * n, DPU_XFER_DEFAULT));

            unsigned int kernel = 0;
            dpu_arguments_t input_arguments = {m, n, M_, kernel};
	        DPU_FOREACH(dpu_set, dpu, i) {
	            DPU_ASSERT(dpu_prepare_xfer(dpu, &input_arguments));
	        }
	        DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, "DPU_INPUT_ARGUMENTS", 0, sizeof(input_arguments), DPU_XFER_DEFAULT));
            printf("Run step 2 on DPU(s) \n");
            // Run DPU kernel
            if(rep >= p.n_warmup){
                start(&timer, 2, rep - p.n_warmup + timer_fix);
#if ENERGY
                DPU_ASSERT(dpu_probe_start(&probe));
#endif
            }
            DPU_ASSERT(dpu_launch(dpu_set, DPU_SYNCHRONOUS));
            if(rep >= p.n_warmup){
                stop(&timer, 2);
#if ENERGY
                DPU_ASSERT(dpu_probe_stop(&probe));
#endif
            }
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

            kernel = 1;
            dpu_arguments_t input_arguments2 = {m, n, M_, kernel};
	        DPU_FOREACH(dpu_set, dpu, i) {
	            DPU_ASSERT(dpu_prepare_xfer(dpu, &input_arguments2));
	        }
	        DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, "DPU_INPUT_ARGUMENTS", 0, sizeof(input_arguments2), DPU_XFER_DEFAULT));
            printf("Run step 3 on DPU(s) \n");
            // Run DPU kernel
            if(rep >= p.n_warmup){
                start(&timer, 3, rep - p.n_warmup + timer_fix);
#if ENERGY
                DPU_ASSERT(dpu_probe_start(&probe));
#endif
            }
            DPU_ASSERT(dpu_launch(dpu_set, DPU_SYNCHRONOUS));
            if(rep >= p.n_warmup){
                stop(&timer, 3);
#if ENERGY
                DPU_ASSERT(dpu_probe_stop(&probe));
#endif
            }
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
                start(&timer, 4, rep - p.n_warmup + timer_fix);
            DPU_FOREACH(dpu_set, dpu) {
                DPU_ASSERT(dpu_prepare_xfer(dpu, (T*)(&A_result[curr_dpu * m * n * M_])));
                curr_dpu++;
            }
            DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_FROM_DPU, DPU_MRAM_HEAP_POINTER_NAME, 0, sizeof(T) * m * n * M_, DPU_XFER_DEFAULT));
            if(rep >= p.n_warmup)
                stop(&timer, 4);

            if(first_round){
                first_round = 0;
            }
            timer_fix++;
        }
        DPU_ASSERT(dpu_free(dpu_set));

    }

    // Print timing results
    printf("CPU ");
    print(&timer, 0, p.n_reps);
    printf("CPU-DPU (Step 1) ");
    print(&timer, 1, p.n_reps);
    printf("Step 2 ");
    print(&timer, 2, p.n_reps);
    printf("Step 3 ");
    print(&timer, 3, p.n_reps);
    printf("DPU-CPU ");
    print(&timer, 4, p.n_reps);

    #if ENERGY
    double energy;
    DPU_ASSERT(dpu_probe_get(&probe, DPU_ENERGY, DPU_AVERAGE, &energy));
    printf("DPU Energy (J): %f\t", energy);
    #endif	

    // Check output
    bool status = true;
    for (i = 0; i < M_ * m * N_ * n; i++) {
        if(A_host[i] != A_result[i]){ 
            status = false;
#if PRINT
            printf("%d: %lu -- %lu\n", i, A_host[i], A_result[i]);
#endif
        }
    }
    if (status) {
        printf("[" ANSI_COLOR_GREEN "OK" ANSI_COLOR_RESET "] Outputs are equal\n");
    } else {
        printf("[" ANSI_COLOR_RED "ERROR" ANSI_COLOR_RESET "] Outputs differ!\n");
    }

    // Deallocation
    free(A_host);
    free(A_backup);
    free(A_result);
    free(done_host);
	
    return status ? 0 : -1;
}
