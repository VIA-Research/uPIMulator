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
#include "../../support/timer.h"
#include "../../support/common.h"

T** A;
T* B;
T* C;

// Create input arrays
static void init_data(T** A, T* B, unsigned int m_size, unsigned int n_size){
    for (unsigned int l = 0; l < NUM_LAYERS; l++)
		for (unsigned int i = 0; i < m_size * n_size; i++){
			if(i % 100 < 98){
				A[l][i] = 0;
			}else{
				A[l][i] = (l+i) % 2;
			}
		}
	for (unsigned int i = 0; i < n_size; i++){
		if(i % 50 < 48){
			B[i] = 0;
		}
		else{
			B[i] = i % 2;
		}
	}
}

// Compute output in the host
static void mlp_host(T* C, T** A, T* B, unsigned int m_size, unsigned int n_size) {
	for (unsigned int nl = 0; nl < NUM_LAYERS; nl++){
		for (unsigned int m = 0; m < m_size; m++){
			C[m] = 0;
		}
		#pragma omp parallel for
		for (unsigned int m = 0; m < m_size; m++){
			for (unsigned int n = 0; n < n_size; n++){
				C[m] += A[nl][m * n_size + n] * B[n];
			}
			C[m] = max(0, C[m]);
		}
		for (unsigned int n = 0; n < n_size; n++){
			B[n] = C[n];
		}
	}
}

static uint64_t mlp_host_sum(uint64_t n_size, uint64_t m_size) {
  uint64_t sum = 0;
  for (uint64_t m = 0; m < n_size; m++){
    sum += B[m];
  }
  return sum;
}

// Params ---------------------------------------------------------------------
typedef struct Params {
  char* dpu_type;
  int   nr_of_ranks;
  int   input_size_n;
  int   input_size_m;
  int   n_warmup;
  int   n_reps;
}Params;

void usage() {
  fprintf(stderr,
    "\nUsage:  ./program [options]"
    "\n"
    "\nGeneral options:"
    "\n    -h        help"
    "\n    -d <D>    DPU type (default=fsim)"
    "\n    -r <R>    # of ranks (default=2)"
    "\n"
    "\nBenchmark-specific options:"
    "\n    -i <I>    input size (default=8M elements)"
    "\n");
  }

  struct Params input_params(int argc, char **argv) {
    struct Params p;
    p.dpu_type      = "fsim";
    p.nr_of_ranks   = 1;
    p.input_size_n  = 1 << 9;
    p.input_size_m  = 1 << 9;
    p.n_warmup      = 2;
    p.n_reps        = 3;

    int opt;
    while((opt = getopt(argc, argv, "hd:r:i:")) >= 0) {
      switch(opt) {
        case 'h':
        usage();
        exit(0);
        break;
        case 'd': p.dpu_type        = optarg; break;
        case 'r': p.nr_of_ranks     = atoi(optarg); break;
        case 'n': p.input_size_n    = atoi(optarg); break;
        case 'm': p.input_size_m    = atoi(optarg); break;
        default:
        fprintf(stderr, "\nUnrecognized option!\n");
        usage();
        exit(0);
      }
    }
    assert(p.nr_of_ranks > 0 && "Invalid # of ranks!");

    return p;
  }

  /**
  * @brief Main of the Host Application.
  */
  int main(int argc, char **argv) {

    struct Params p = input_params(argc, argv);
    uint64_t n_size = 8192;
    uint64_t m_size = 20480;

    Timer timer;
    A = malloc(NUM_LAYERS * sizeof(T*));
    for(int l = 0; l < NUM_LAYERS; l++)
        A[l] = malloc(n_size*m_size*sizeof(unsigned int));
    B = malloc(m_size*sizeof(unsigned int));
    C = malloc(m_size*sizeof(unsigned int));

    // Create an input file with arbitrary data.
    init_data(A, B, m_size, n_size);

    start(&timer, 0, 1);
    mlp_host(C, A, B, n_size, m_size);
    stop(&timer, 0);

    uint32_t sum = mlp_host_sum(n_size, m_size);
   
    printf("Kernel ");
    print(&timer, 0, 1);
    printf("\n");

    printf("SUM = %d \n", sum);

    for(int l = 0; l < NUM_LAYERS; l++)
        free(A[l]);
    free(A);
    free(B);
    free(C);

    return 0;
}
