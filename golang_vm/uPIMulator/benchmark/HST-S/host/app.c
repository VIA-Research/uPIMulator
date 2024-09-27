#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <math.h>
#include <dpu.h>
#include <dpu_log.h>
#include <unistd.h>
#include <getopt.h>
#include <assert.h>

#define NUM_DPUS 1
#define NUM_TASKLETS 1
#define DATA_PREP_PARAMS 64

#define DEPTH 12

struct dpu_arguments_t {
    int size;
    int transfer_size;
    int bins;
	int kernel;
};

void read_input(int* A, int input_size) {
    for (int i = 0; i < input_size; i++) {
        A[i] = i % 4096;
    }
}

void histogram_host(int* histo, int* A, int bins, int nr_elements, int exp, int nr_of_dpus) {
    if(!exp){
        for (int i = 0; i < nr_of_dpus; i++) {
            for (int j = 0; j < nr_elements; j++) {
                int d = A[j];
                histo[i * bins + ((d * bins) >> DEPTH)] += 1;
            }
        }
    }
    else{
        for (int j = 0; j < nr_elements; j++) {
            int d = A[j];
            histo[(d * bins) >> DEPTH] += 1;
        }
    }
}

int roundup(int n, int m) {
    return ((n / m) * m + m);
}

int divceil(int n, int m) {
    return ((n-1) / m + 1);
}

int main() {
    int* A;
    int* histo_host;
    int* histo;

    struct dpu_set_t dpu_set;
    struct dpu_set_t dpu;
    int nr_of_dpus = NUM_DPUS;
   
    dpu_alloc(NUM_DPUS, NULL, &dpu_set);
    dpu_load(dpu_set, DPU_BINARY, NULL);

    int input_size = DATA_PREP_PARAMS;

    int input_size_8bytes =  ((input_size * sizeof(int)) % 8) != 0 ? roundup(input_size, 8) : input_size;
    int input_size_dpu = divceil(input_size, nr_of_dpus);
    int input_size_dpu_8bytes = ((input_size_dpu * sizeof(int)) % 8) != 0 ? roundup(input_size_dpu, 8) : input_size_dpu;

    int bins = 256;
    A = malloc(input_size_dpu_8bytes * nr_of_dpus * sizeof(int));
    int *bufferA = A;
    histo_host = malloc(bins * sizeof(int));
    histo = malloc(nr_of_dpus * bins * sizeof(int));

    read_input(A, input_size);

    for (int i = 0; i < bins; i++) {
        histo_host[i] = 0;
    }

    for (int i = 0; i < nr_of_dpus * bins; i++) {
        histo[i] = 0;
    }

    histogram_host(histo_host, A, bins, input_size, 1, nr_of_dpus);

    int kernel = 0;
    struct dpu_arguments_t* input_arguments = malloc(NUM_DPUS * sizeof(struct dpu_arguments_t));
    for(int i=0; i<nr_of_dpus-1; i++) {
        input_arguments[i].size=input_size_dpu_8bytes * sizeof(int);
        input_arguments[i].transfer_size=input_size_dpu_8bytes * sizeof(int);
        input_arguments[i].bins=bins;
        input_arguments[i].kernel=kernel;
    }
    input_arguments[nr_of_dpus-1].size=(input_size_8bytes - input_size_dpu_8bytes * (NUM_DPUS-1)) * sizeof(int);
    input_arguments[nr_of_dpus-1].transfer_size=input_size_dpu_8bytes * sizeof(int);
    input_arguments[nr_of_dpus-1].bins=bins;
    input_arguments[nr_of_dpus-1].kernel=kernel;

    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &input_arguments[i]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, "DPU_INPUT_ARGUMENTS", 0, sizeof(struct dpu_arguments_t), DPU_XFER_DEFAULT);

    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &bufferA[input_size_dpu_8bytes * i]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, 0, input_size_dpu_8bytes * sizeof(int), DPU_XFER_DEFAULT);

    dpu_launch(dpu_set, DPU_SYNCHRONOUS);

    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &histo[bins * i]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_FROM_DPU, DPU_MRAM_HEAP_POINTER_NAME, input_size_dpu_8bytes * sizeof(int), bins * sizeof(int), DPU_XFER_DEFAULT);

    for(int i = 1; i < nr_of_dpus; i++){
        for(int j = 0; j < bins; j++){
            histo[j] += histo[j + i * bins];
        }
    }

    int status = 1;
    for (int j = 0; j < bins; j++) {
        if(histo_host[j] != histo[j]){
            status = 0;
            break;
        }
    }

    assert(status);

    free(A);
    free(histo_host);
    free(histo);
    dpu_free(dpu_set);
	
    return status ? 0 : -1;
}
