#ifndef _PARAMS_H_
#define _PARAMS_H_

#include "common.h"

typedef struct Params {
    unsigned int   max_rows;
    unsigned int   penalty;
    unsigned int   n_warmup;
    unsigned int   n_reps;
} Params;

static void usage() {
    fprintf(stderr,
            "\nUsage:  ./program [options]"
            "\n"
            "\nGeneral options:"
            "\n    -h        help"
            "\n    -w <W>    # of untimed warmup iterations (default=1)"
            "\n    -e <E>    # of timed repetition iterations (default=3)"
            "\n"
            "\nBenchmark-specific options:"
            "\n    -n <N>    size of sequence: length of the sequence"
            "\n    -p <P>    penalty: a positive integer"
            "\n");
}

struct Params input_params(int argc, char **argv) {
    struct Params p;
    p.n_warmup      = 1;
    p.n_reps        = 3;
    p.max_rows      = 256;
    p.penalty       = 1;

    int opt;
    while((opt = getopt(argc, argv, "hw:e:n:p:")) >= 0) {
        switch(opt) {
            case 'h':
                usage();
                exit(0);
                break;
            case 'w': p.n_warmup      = atoi(optarg); break;
            case 'e': p.n_reps        = atoi(optarg); break;
            case 'n': p.max_rows      = atoi(optarg); break;
            case 'p': p.penalty       = atoi(optarg); break;
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
