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

static uint64_t *A;
static uint64_t *B;
static uint64_t *C;
static uint64_t *C2;
static int pos;

bool pred(const uint64_t x){
  return (x % 2) == 0;
}


void  *create_test_file(unsigned int nr_elements) {
    //srand(0);

    A = (uint64_t*) malloc(nr_elements * sizeof(uint64_t));
    B = (uint64_t*) malloc(nr_elements * sizeof(uint64_t));
    C = (uint64_t*) malloc(nr_elements * sizeof(uint64_t));

    printf("nr_elements\t%u\t", nr_elements);
    for (int i = 0; i < nr_elements; i++) {
        //A[i] = (unsigned int) (rand());
        A[i] = i+1;
        B[i] = 0;
    }
}

/**
* @brief compute output in the host
*/
static int select_host(int size, int t) {
    pos = 0;
    C[pos] = A[pos];

    omp_set_num_threads(t);
    #pragma omp parallel for
    for(int my = 1; my < size; my++) {
        if(!pred(A[my])) {
            int p;
            #pragma omp atomic update
            pos++;
            p = pos;
            C[p] = A[my];
        }
    }
    return pos;
}

// Params ---------------------------------------------------------------------
typedef struct Params {
    char* dpu_type;
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
        "\n    -d <D>    DPU type (default=fsim)"
        "\n    -t <T>    # of threads (default=8)"
        "\n    -w <W>    # of untimed warmup iterations (default=2)"
        "\n    -e <E>    # of timed repetition iterations (default=5)"
        "\n"
        "\nBenchmark-specific options:"
        "\n    -i <I>    input size (default=8M elements)"
        "\n");
}

struct Params input_params(int argc, char **argv) {
    struct Params p;
    p.input_size    = 16 << 20;
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
        case 't': p.n_threads     = atoi(optarg); break;
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
    uint32_t accum = 0;
    int total_count;

    // Create an input file with arbitrary data.
    create_test_file(file_size);

    Timer timer;
    start(&timer, 0, 0);

    total_count = select_host(file_size, p.n_threads);

    stop(&timer, 0);

    printf("Total count = %d\t", total_count);

    printf("Kernel ");
    print(&timer, 0, 1);
    printf("\n");
    
    free(A);
    free(B);
    free(C);
    return 0;
}
