/**
* app.c
* HST-L Host Application Source File
*
*/
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <math.h>
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

#if ENERGY
#include <dpu_probe.h>
#endif

// Pointer declaration
static T* A;
static unsigned int* histo_host;
static unsigned int* histo;

// Create input arrays
static void read_input(T* A, const Params p) {

    char  dctFileName[100];
    FILE *File = NULL;

    // Open input file
    unsigned short temp;
    sprintf(dctFileName, p.file_name);
    if((File = fopen(dctFileName, "rb")) != NULL) {
        for(unsigned int y = 0; y < p.input_size; y++) {
            fread(&temp, sizeof(unsigned short), 1, File);
            A[y] = (unsigned int)ByteSwap16(temp);
            if(A[y] >= 4096)
                A[y] = 4095;
        }
        fclose(File);
    } else {
        printf("%s does not exist\n", dctFileName);
        exit(1);
    }
}

// Compute output in the host
static void histogram_host(unsigned int* histo, T* A, unsigned int bins, unsigned int nr_elements, int exp, unsigned int nr_of_dpus) {
    if(!exp){
        for (unsigned int i = 0; i < nr_of_dpus; i++) {
            for (unsigned int j = 0; j < nr_elements; j++) {
                T d = A[j];
                histo[i * bins + ((d * bins) >> DEPTH)] += 1;
            }
        }
    }
    else{
        for (unsigned int j = 0; j < nr_elements; j++) {
            T d = A[j];
            histo[(d * bins) >> DEPTH] += 1;
        }
    }
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

    // Allocate DPUs and load binary
    DPU_ASSERT(dpu_alloc(NR_DPUS, NULL, &dpu_set));
    DPU_ASSERT(dpu_load(dpu_set, DPU_BINARY, NULL));
    DPU_ASSERT(dpu_get_nr_dpus(dpu_set, &nr_of_dpus));
    printf("Allocated %d DPU(s)\n", nr_of_dpus);

    unsigned int i = 0;
    unsigned int input_size; // Size of input image
    unsigned int dpu_s = p.dpu_s;
    if(p.exp == 0)
        input_size = p.input_size * nr_of_dpus; // Size of input image
    else if(p.exp == 1)
        input_size = p.input_size; // Size of input image
	else
        input_size = p.input_size * dpu_s; // Size of input image

    const unsigned int input_size_8bytes = 
        ((input_size * sizeof(T)) % 8) != 0 ? roundup(input_size, 8) : input_size; // Input size per DPU (max.), 8-byte aligned
    const unsigned int input_size_dpu = divceil(input_size, nr_of_dpus); // Input size per DPU (max.)
    const unsigned int input_size_dpu_8bytes = 
        ((input_size_dpu * sizeof(T)) % 8) != 0 ? roundup(input_size_dpu, 8) : input_size_dpu; // Input size per DPU (max.), 8-byte aligned

    // Input/output allocation
    A = malloc(input_size_dpu_8bytes * nr_of_dpus * sizeof(T));
    T *bufferA = A;
    histo_host = malloc(p.bins * sizeof(unsigned int));
    histo = malloc(nr_of_dpus * p.bins * sizeof(unsigned int));

    // Create an input file with arbitrary data
    read_input(A, p);
    if(p.exp == 0){
        for(unsigned int j = 1; j < nr_of_dpus; j++){
            memcpy(&A[j * input_size_dpu_8bytes], &A[0], input_size_dpu_8bytes * sizeof(T));
        }
    }
    else if(p.exp == 2){
        for(unsigned int j = 1; j < dpu_s; j++)
            memcpy(&A[j * p.input_size], &A[0], p.input_size * sizeof(T));
    }

    // Timer declaration
    Timer timer;

    printf("NR_TASKLETS\t%d\tBL\t%d\tinput_size\t%u\n", NR_TASKLETS, BL, input_size);

    // Loop over main kernel
    for(int rep = 0; rep < p.n_warmup + p.n_reps; rep++) {
        memset(histo_host, 0, p.bins * sizeof(unsigned int));
        memset(histo, 0, nr_of_dpus * p.bins * sizeof(unsigned int));

        // Compute output on CPU (performance comparison and verification purposes)
        if(rep >= p.n_warmup)
            start(&timer, 0, rep - p.n_warmup);
        histogram_host(histo_host, A, p.bins, p.input_size, 1, nr_of_dpus);
        if(rep >= p.n_warmup)
            stop(&timer, 0);

        printf("Load input data\n");
        if(rep >= p.n_warmup)
            start(&timer, 1, rep - p.n_warmup);
        // Input arguments
        unsigned int kernel = 0;
        i = 0;
	    dpu_arguments_t input_arguments[NR_DPUS];
	    for(i=0; i<nr_of_dpus-1; i++) {
	        input_arguments[i].size=input_size_dpu_8bytes * sizeof(T); 
	        input_arguments[i].transfer_size=input_size_dpu_8bytes * sizeof(T); 
	        input_arguments[i].bins=p.bins;
	        input_arguments[i].kernel=kernel;
	    }
	    input_arguments[nr_of_dpus-1].size=(input_size_8bytes - input_size_dpu_8bytes * (NR_DPUS-1)) * sizeof(T); 
	    input_arguments[nr_of_dpus-1].transfer_size=input_size_dpu_8bytes * sizeof(T); 
	    input_arguments[nr_of_dpus-1].bins=p.bins;
	    input_arguments[nr_of_dpus-1].kernel=kernel;

        // Copy input arrays
        i = 0;
        DPU_FOREACH(dpu_set, dpu, i) {
            DPU_ASSERT(dpu_prepare_xfer(dpu, &input_arguments[i]));
        }
        DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, "DPU_INPUT_ARGUMENTS", 0, sizeof(input_arguments[0]), DPU_XFER_DEFAULT));
        DPU_FOREACH(dpu_set, dpu, i) {
            DPU_ASSERT(dpu_prepare_xfer(dpu, bufferA + input_size_dpu_8bytes * i));
        }
        DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, 0, input_size_dpu_8bytes * sizeof(T), DPU_XFER_DEFAULT));
        if(rep >= p.n_warmup)
            stop(&timer, 1);

        printf("Run program on DPU(s) \n");
        // Run DPU kernel
        if(rep >= p.n_warmup) {
            start(&timer, 2, rep - p.n_warmup);
            #if ENERGY
            DPU_ASSERT(dpu_probe_start(&probe));
            #endif
        }
        DPU_ASSERT(dpu_launch(dpu_set, DPU_SYNCHRONOUS));
        if(rep >= p.n_warmup) {
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

        printf("Retrieve results\n");
        i = 0;
        if(rep >= p.n_warmup)
            start(&timer, 3, rep - p.n_warmup);
        // PARALLEL RETRIEVE TRANSFER
        DPU_FOREACH(dpu_set, dpu, i) {
            DPU_ASSERT(dpu_prepare_xfer(dpu, histo + p.bins * i));
        }
        DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_FROM_DPU, DPU_MRAM_HEAP_POINTER_NAME, input_size_dpu_8bytes * sizeof(T), p.bins * sizeof(unsigned int), DPU_XFER_DEFAULT));
		
        // Final histogram merging
        for(i = 1; i < nr_of_dpus; i++){
            for(unsigned int j = 0; j < p.bins; j++){
                histo[j] += histo[j + i * p.bins];
            }			
        }		
        if(rep >= p.n_warmup)
            stop(&timer, 3);

    }

    // Print timing results
    printf("CPU ");
    print(&timer, 0, p.n_reps);
    printf("CPU-DPU ");
    print(&timer, 1, p.n_reps);
    printf("DPU Kernel ");
    print(&timer, 2, p.n_reps);
    printf("DPU-CPU ");
    print(&timer, 3, p.n_reps);

    #if ENERGY
    double energy;
    DPU_ASSERT(dpu_probe_get(&probe, DPU_ENERGY, DPU_AVERAGE, &energy));
    printf("DPU Energy (J): %f\t", energy);
    #endif	


    // Check output
    bool status = true;
    if(p.exp == 1) 
        for (unsigned int j = 0; j < p.bins; j++) {
            if(histo_host[j] != histo[j]){ 
                status = false;
#if PRINT
                printf("%u - %u: %u -- %u\n", j, j, histo_host[j], histo[j]);
#endif
            }
        }
    else if(p.exp == 2) 
        for (unsigned int j = 0; j < p.bins; j++) {
            if(dpu_s * histo_host[j] != histo[j]){ 
                status = false;
#if PRINT
                printf("%u - %u: %u -- %u\n", j, j, dpu_s * histo_host[j], histo[j]);
#endif
            }
        }
    else
        for (unsigned int j = 0; j < p.bins; j++) {
            if(nr_of_dpus * histo_host[j] != histo[j]){ 
                status = false;
#if PRINT
                printf("%u - %u: %u -- %u\n", j, j, nr_of_dpus * histo_host[j], histo[j]);
#endif
            }
        }
    if (status) {
        printf("[" ANSI_COLOR_GREEN "OK" ANSI_COLOR_RESET "] Outputs are equal\n");
    } else {
        printf("[" ANSI_COLOR_RED "ERROR" ANSI_COLOR_RESET "] Outputs differ!\n");
    }

    // Deallocation
    free(A);
    free(histo_host);
    free(histo);
    DPU_ASSERT(dpu_free(dpu_set));
	
    return status ? 0 : -1;
}
