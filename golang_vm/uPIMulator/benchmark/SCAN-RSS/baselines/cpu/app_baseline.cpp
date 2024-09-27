/*
* JGL@SAFARI
*/

/**
* CPU code with Thrust
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
#include <thrust/system/omp/execution_policy.h>
#include <thrust/system/omp/vector.h>

#include <omp.h>

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
        "\n    -i <I>    input size (default=8M elements)"
        "\n");
}

struct Params input_params(int argc, char **argv) {
    struct Params p;
    p.input_size    = 2 << 20;
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

    struct Params p = input_params(argc, argv);

    unsigned int nr_of_dpus = 1;
    
    unsigned int i = 0;
    const unsigned int input_size =  p.input_size;
    assert(input_size % (p.n_threads) == 0 && "Input size!");

    // Input/output allocation
    A = (T*)malloc(input_size * sizeof(T));
    C = (T*)malloc(input_size * sizeof(T));
    T *bufferA = A;

    // Create an input file with arbitrary data.
    read_input(A, input_size);

    // Timer declaration
    Timer timer;
    float time_gpu = 0;

    thrust::omp::vector<T> h_output(input_size);

    // Loop over main kernel
    for(int rep = 0; rep < p.n_warmup + p.n_reps; rep++) {

        // Compute output on CPU (performance comparison and verification purposes)
        if(rep >= p.n_warmup)
            start(&timer, 0, rep - p.n_warmup);
        scan_host(C, A, input_size);
        if(rep >= p.n_warmup)
            stop(&timer, 0);

        memcpy(thrust::raw_pointer_cast(&h_output[0]), A, input_size * sizeof(T));

        omp_set_num_threads(p.n_threads);

        if(rep >= p.n_warmup)
            start(&timer, 1, rep - p.n_warmup);
        thrust::exclusive_scan(thrust::omp::par, h_output.begin(),h_output.end(),h_output.begin());
        if(rep >= p.n_warmup)
            stop(&timer, 1);
    }

    // Print timing results
    printf("CPU ");
    print(&timer, 0, p.n_reps);
    printf("Kernel ");
    print(&timer, 1, p.n_reps);

    // Check output
    bool status = true;
    for (i = 0; i < input_size; i++) {
        if(C[i] != h_output[i]){ 
            status = false;
            //printf("%d: %lu -- %lu\n", i, C[i], h_output[i]);
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
	
    return 0;
}
