/**
* @file app.c
* @brief Template for a Host Application Source File.
*
*/
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <unistd.h>
#include <getopt.h>
#include <assert.h>
#include <stdint.h>

#include <omp.h>
#include "../../support/timer.h"

static int32_t *A;
static int32_t *B;
static int32_t *C;
static int32_t *C2; 

/**
* @brief creates a "test file" by filling a buffer of 64MB with pseudo-random values
* @param nr_elements how many 32-bit elements we want the file to be
* @return the buffer address
*/
void  *create_test_file(unsigned int nr_elements) {
    srand(0);
    printf("nr_elements\t%u\t", nr_elements);
    A = (uint32_t*) malloc(nr_elements * sizeof(uint32_t));
    B = (uint32_t*) malloc(nr_elements * sizeof(uint32_t));
    C = (uint32_t*) malloc(nr_elements * sizeof(uint32_t));
    
    for (int i = 0; i < nr_elements; i++) {
        A[i] = (int) (rand());
        B[i] = (int) (rand());
    }

}

/**
* @brief compute output in the host
*/
static void vector_addition_host(unsigned int nr_elements, int t) {
    omp_set_num_threads(t);
    #pragma omp parallel for
    for (int i = 0; i < nr_elements; i++) {
        C[i] = A[i] + B[i];
    }
}

// Params ---------------------------------------------------------------------
typedef struct Params {
    int   input_size;
    int   n_warmup;
    int   n_reps;
    int   n_threads;
}Params;

void usage() {
    fprintf(stderr,
        "\nUsage:  ./program [options]"
        "\n"
        "\nGeneral options:"
        "\n    -h        help"
        "\n    -t <T>    # of threads (default=8)"
        "\n    -w <W>    # of untimed warmup iterations (default=1)"
        "\n    -e <E>    # of timed repetition iterations (default=3)"
        "\n"
        "\nBenchmark-specific options:"
        "\n    -i <I>    input size (default=8M elements)"
        "\n");
}

struct Params input_params(int argc, char **argv) {
    struct Params p;
    p.input_size    = 16777216;
    p.n_warmup      = 1;
    p.n_reps        = 3;
    p.n_threads     = 5;

    int opt;
    while((opt = getopt(argc, argv, "hi:w:e:t:")) >= 0) {
        switch(opt) {
        case 'h':
        usage();
        exit(0);
        break;
        case 'i': p.input_size    = atoi(optarg); break;
        case 'w': p.n_warmup      = atoi(optarg); break;
        case 'e': p.n_reps        = atoi(optarg); break;
        case 't': p.n_threads        = atoi(optarg); break;
        default:
            fprintf(stderr, "\nUnrecognized option!\n");
            usage();
            exit(0);
        }
    }
    assert(p.n_threads > 0 && "Invalid # of ranks!");

    return p;
}

/**
* @brief Main of the Host Application.
*/
int main(int argc, char **argv) {

    struct Params p = input_params(argc, argv);

    const unsigned int file_size = p.input_size;

    // Create an input file with arbitrary data.
    create_test_file(file_size);

    Timer timer;
    start(&timer, 0, 0);

    vector_addition_host(file_size, p.n_threads);
	
    stop(&timer, 0);
    printf("Kernel ");
    print(&timer, 0, 1);
    printf("\n");

    free(A);
    free(B);
    free(C);

   return 0;
 }
