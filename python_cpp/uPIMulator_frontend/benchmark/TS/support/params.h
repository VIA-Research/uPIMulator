#ifndef _PARAMS_H_
#define _PARAMS_H_

#include "common.h"

// Params ---------------------------------------------------------------------
typedef struct Params {
  unsigned long  input_size_n;
  unsigned long  input_size_m;
  int  n_warmup;
  int  n_reps;
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
    "\n    -n <n>    n (TS length. Default=64K elements)"
    "\n    -m <m>    m (Query length. Default=256 elements)"
    "\n");
  }

  struct Params input_params(int argc, char **argv) {
    struct Params p;
    p.input_size_n  = 1 << 16;
    p.input_size_m  = 1 << 8;

    p.n_warmup      = 1;
    p.n_reps        = 3;

    int opt;
    while((opt = getopt(argc, argv, "hw:e:n:m:")) >= 0) {
      switch(opt) {
        case 'h':
        usage();
        exit(0);
        break;
        case 'w': p.n_warmup      = atoi(optarg); break;
        case 'e': p.n_reps        = atoi(optarg); break;
        case 'n': p.input_size_n  = atol(optarg); break;
        case 'm': p.input_size_m  = atol(optarg); break;
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
