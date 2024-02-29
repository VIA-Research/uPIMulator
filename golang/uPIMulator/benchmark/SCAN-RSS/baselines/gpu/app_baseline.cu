/*
* JGL@SAFARI
*/

/**
* GPU code with Thrust
*/
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <unistd.h>
#include <getopt.h>
#include <assert.h>

#include <iostream>
#include <fstream>
#include <cstdlib>
#include <ctime>
#include <cstdio>
#include <math.h>
#include <sys/time.h>

#include <vector>
#include <thrust/device_vector.h>
#include <thrust/host_vector.h>
#include <thrust/scan.h>
#include <thrust/copy.h>

#include "../../support/common.h"
#include "../../support/timer.h"

#define ANSI_COLOR_RED     "\x1b[31m"
#define ANSI_COLOR_GREEN   "\x1b[32m"
#define ANSI_COLOR_RESET   "\x1b[0m"

// Pointer declaration
static T* A;
static T* C;
static T* C2;

/**
* @brief creates input arrays
* @param nr_elements how many elements in input arrays
*/
static void read_input(T* A, unsigned int nr_elements) {
    //srand(0);
    printf("nr_elements\t%u\t", nr_elements);
    for (unsigned int i = 0; i < nr_elements; i++) {
        //A[i] = (T) (rand()) % 2;
        A[i] = i;
    }
}

/**
* @brief compute output in the host
*/
static void scan_host(T* C, T* A, unsigned int nr_elements) {
    C[0] = A[0];
    for (unsigned int i = 1; i < nr_elements; i++) {
        C[i] = C[i - 1] + A[i - 1];
    }
}

// Params ---------------------------------------------------------------------
typedef struct Params {
    unsigned int   input_size;
    int   n_warmup;
    int   n_reps;
    int   exp;
    int   n_threads;
}Params;

void usage() {
    fprintf(stderr,
        "\nUsage:  ./program [options]"
        "\n"
        "\nGeneral options:"
        "\n    -h        help"
        "\n    -w <W>    # of untimed warmup iterations (default=1)"
        "\n    -e <E>    # of timed repetition iterations (default=3)"
        "\n    -x <X>    Weak (0) or strong (1) scaling (default=0)"
        "\n    -t <T>    # of threads (default=8)"
        "\n"
        "\nBenchmark-specific options:"
        "\n    -i <I>    input size (default=640 * 3932160 elements)"
        "\n");
}

struct Params input_params(int argc, char **argv) {
    struct Params p;
    p.input_size    = 1258291200;
    p.n_warmup      = 1;
    p.n_reps        = 3;
    p.exp           = 0;
    p.n_threads     = 8;

    int opt;
    while((opt = getopt(argc, argv, "hi:w:e:x:t:")) >= 0) {
        switch(opt) {
        case 'h':
        usage();
        exit(0);
        break;
        case 'i': p.input_size    = atoi(optarg); break;
        case 'w': p.n_warmup      = atoi(optarg); break;
        case 'e': p.n_reps        = atoi(optarg); break;
        case 'x': p.exp           = atoi(optarg); break;
        case 't': p.n_threads     = atoi(optarg); break;
        default:
            fprintf(stderr, "\nUnrecognized option!\n");
            usage();
            exit(0);
        }
    }
    assert(p.n_threads > 0 && "Invalid # of threads!");

    return p;
}

/**
* @brief Main of the Host Application.
*/
int main(int argc, char **argv) {

    cudaDeviceProp device_properties;
    cudaGetDeviceProperties(&device_properties, 0);
    cudaSetDevice(0);

    struct Params p = input_params(argc, argv);

    unsigned int nr_of_dpus = 1;
    
    unsigned int i = 0;
    const unsigned int input_size = p.exp == 0 ? p.input_size * nr_of_dpus : p.input_size;

    // Input/output allocation
    A = (T*)malloc(input_size * sizeof(T));
    C = (T*)malloc(input_size * sizeof(T));
    C2 = (T*)malloc(input_size * sizeof(T));
    T *bufferA = A;
    T *bufferC = C2;

    // Create an input file with arbitrary data.
    read_input(A, input_size);

    // Timer declaration
    Timer timer;
    float time_gpu = 0;

    thrust::host_vector<T> h_output(input_size);

    // Loop over main kernel
    for(int rep = 0; rep < p.n_warmup + p.n_reps; rep++) {

        // Compute output on CPU (performance comparison and verification purposes)
        if(rep >= p.n_warmup)
            start(&timer, 0, rep - p.n_warmup);
        scan_host(C, A, input_size);
        if(rep >= p.n_warmup)
            stop(&timer, 0);


        // Event creation
        cudaEvent_t start, stop;
        cudaEventCreate(&start);
        cudaEventCreate(&stop);
        float time1 = 0;

        thrust::device_vector<T> d_input(input_size);
        cudaMemcpy(thrust::raw_pointer_cast(&d_input[0]), A, input_size * sizeof(T), cudaMemcpyHostToDevice);

        // Start timer
        cudaEventRecord( start, 0 );
        thrust::exclusive_scan(d_input.begin(),d_input.end(),d_input.begin());
        // End timer
        cudaEventRecord( stop, 0 );
        cudaEventSynchronize( stop );
        cudaEventElapsedTime( &time1, start, stop );
        time_gpu += time1;

        h_output = d_input;

        cudaEventDestroy(start);
        cudaEventDestroy(stop);
    }

    // Print timing results
    printf("CPU ");
    print(&timer, 0, p.n_reps);
    printf("Kernel (ms):");
    printf("%f\n", time_gpu / p.n_reps);

    // Check output
    bool status = true;
    for (i = 0; i < input_size; i++) {
        if(C[i] != h_output[i]){ 
            status = false;
            printf("%d: %lu -- %lu\n", i, C[i], h_output[i]);
        }
    }
    if (status) {
        printf("[" ANSI_COLOR_GREEN "OK" ANSI_COLOR_RESET "] Outputs are equal\n");
    } else {
        printf("[" ANSI_COLOR_RED "ERROR" ANSI_COLOR_RESET "] Outputs differ!\n");
    }

    // Deallocation
    free(A);
    free(C);
    free(C2);
	
    return 0;
}
