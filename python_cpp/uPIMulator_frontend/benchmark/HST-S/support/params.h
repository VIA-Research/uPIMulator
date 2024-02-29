#ifndef _PARAMS_H_
#define _PARAMS_H_

#include "common.h"

typedef struct Params {
    unsigned int   input_size;
    unsigned int   bins;
    int   n_warmup;
    int   n_reps;
    const char *file_name;
    int  exp;
    int  dpu_s;
}Params;

static void usage() {
    fprintf(stderr,
        "\nUsage:  ./program [options]"
        "\n"
        "\nGeneral options:"
        "\n    -h        help"
        "\n    -w <W>    # of untimed warmup iterations (default=1)"
        "\n    -e <E>    # of timed repetition iterations (default=3)"
        "\n    -x <X>    Weak (0) or strong (1, 2) scaling (default=0)"
        "\n"
        "\nBenchmark-specific options:"
        "\n    -i <I>    input size (default=1536*1024 elements)"
        "\n    -b <B>    histogram size (default=256 bins)"
        "\n    -f <F>    input image file (default=../input/image_VanHateren.iml)"
        "\n");
}

struct Params input_params(int argc, char **argv) {
    struct Params p;
    p.input_size    = 1536 * 1024;
    p.bins          = 256;
    p.n_warmup      = 1;
    p.n_reps        = 3;
    p.exp           = 0;
    p.file_name     = "./input/image_VanHateren.iml";
    p.dpu_s         = 64;

    int opt;
    while((opt = getopt(argc, argv, "hi:b:w:e:f:x:z:")) >= 0) {
        switch(opt) {
        case 'h':
        usage();
        exit(0);
        break;
        case 'i': p.input_size    = atoi(optarg); break;
        case 'b': p.bins          = atoi(optarg); break;
        case 'w': p.n_warmup      = atoi(optarg); break;
        case 'e': p.n_reps        = atoi(optarg); break;
        case 'f': p.file_name     = optarg; break;
        case 'x': p.exp           = atoi(optarg); break;
        case 'z': p.dpu_s         = atoi(optarg); break;
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
