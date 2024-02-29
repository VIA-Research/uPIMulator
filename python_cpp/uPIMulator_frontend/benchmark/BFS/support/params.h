
#ifndef _PARAMS_H_
#define _PARAMS_H_

#include "common.h"
#include "utils.h"

static void usage() {
    PRINT(  "\nUsage:  ./program [options]"
            "\n"
            "\nBenchmark-specific options:"
            "\n    -f <F>    input matrix file name (default=data/roadNet-CA.txt)"
            "\n"
            "\nGeneral options:"
            "\n    -v <V>    verbosity"
            "\n    -h        help"
            "\n\n");
}

typedef struct Params {
  const char* fileName;
  unsigned int verbosity;
} Params;

static struct Params input_params(int argc, char **argv) {
    struct Params p;
    p.fileName      = "data/roadNet-CA.txt";
    p.verbosity     = 1;
    int opt;
    while((opt = getopt(argc, argv, "f:v:h")) >= 0) {
        switch(opt) {
            case 'f': p.fileName    = optarg;       break;
            case 'v': p.verbosity   = atoi(optarg); break;
            case 'h': usage(); exit(0);
            default:
                      PRINT_ERROR("Unrecognized option!");
                      usage();
                      exit(0);
        }
    }

    return p;
}

#endif

