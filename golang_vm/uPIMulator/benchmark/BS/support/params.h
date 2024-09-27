#ifndef _PARAMS_H_
#define _PARAMS_H_

#include "common.h"

typedef struct Params {
  long  num_querys;
  unsigned   n_warmup;
  unsigned   n_reps;
}Params;

void usage() {
  fprintf(stderr,
    "\nUsage:  ./program [options]"
    "\n"
    "\nGeneral options:"
    "\n    -h        help"
    "\n    -w <W>    # of untimed warmup iterations (default=1)"
    "\n    -e <E>    # of timed repetition iterations (default=3)"
    "\n"
    "\nBenchmark-specific options:"
    "\n    -i <I>    problem size (default=2 queries)"
    "\n");
  }

  struct Params input_params(int argc, char **argv) {
    struct Params p;
    p.num_querys    = PROBLEM_SIZE;
    p.n_warmup      = 1;
    p.n_reps        = 3;

    int opt;
    while((opt = getopt(argc, argv, "h:i:w:e:")) >= 0) {
      switch(opt) {
        case 'h':
        usage();
        exit(0);
        break;
        case 'i': p.num_querys    = atol(optarg); break;
        case 'w': p.n_warmup      = atoi(optarg); break;
        case 'e': p.n_reps        = atoi(optarg); break; 
	default:
        	fprintf(stderr, "\nUnrecognized option!\n");
        	usage();
        	exit(0);
      }
    }
    assert(NR_DPUS > 0 && "Invalid # of dpus!");

    return p;
  }
  #endif
