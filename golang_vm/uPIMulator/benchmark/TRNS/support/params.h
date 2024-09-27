#ifndef _PARAMS_H_
#define _PARAMS_H_

#include "common.h"

typedef struct Params {
    unsigned int   M_;
    unsigned int   m;
    unsigned int   N_;
    unsigned int   n;
    int   n_warmup;
    int   n_reps;
    int  exp;
}Params;

static void usage() {
    fprintf(stderr,
        "\nUsage:  ./program [options]"
        "\n"
        "\nGeneral options:"
        "\n    -h        help"
        "\n    -w <W>    # of untimed warmup iterations (default=1)"
        "\n    -e <E>    # of timed repetition iterations (default=3)"
        "\n    -x <X>    Weak (0) or strong (1) scaling (default=0)"
        "\n"
        "\nBenchmark-specific options:"
        "\n    -m <I>    m (default=16 elements)"
        "\n    -n <I>    n (default=8 elements)"
        "\n    -o <I>    M_ (default=12288 elements)"
        "\n    -p <I>    N_ (default=1 elements)"
        "\n");
}

struct Params input_params(int argc, char **argv) {
    struct Params p;
    p.M_            = 12288;
    p.m             = 16;
    p.N_            = 1;
    p.n             = 8;
    p.n_warmup      = 1;
    p.n_reps        = 3;
    p.exp           = 0;

    int opt;
    while((opt = getopt(argc, argv, "hw:e:x:m:n:o:p:")) >= 0) {
        switch(opt) {
        case 'h':
        usage();
        exit(0);
        break;
        case 'w': p.n_warmup      = atoi(optarg); break;
        case 'e': p.n_reps        = atoi(optarg); break;
        case 'x': p.exp           = atoi(optarg); break;
        case 'm': p.m             = atoi(optarg); break;
        case 'n': p.n             = atoi(optarg); break;
        case 'o': p.M_            = atoi(optarg); break;
        case 'p': p.N_            = atoi(optarg); break;
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
