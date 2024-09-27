#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <dpu.h>
#include <dpu_log.h>
#include <unistd.h>
#include <getopt.h>
#include <assert.h>

struct dpu_arguments_t {
    int size;
    int transfer_size;
    int kernel;
};

#define NUM_DPUS 1
#define NUM_TASKLETS 1
#define DATA_PREP_PARAMS 1024

void read_input(int* A, int* B, int nr_elements) {
    for (int i = 0; i < nr_elements; i++) {
        A[i] = i;
        B[i] = i;
    }
}

void vector_addition_host(int* C, int* A, int* B, int nr_elements) {
    for (int i = 0; i < nr_elements; i++) {
        C[i] = A[i] + B[i];
    }
}

int roundup(int n, int m) {
    return ((n / m) * m + m);
}

int divceil(int n, int m) {
    return ((n-1) / m + 1);
}

int main() {
    int BL = 10;
    int BLOCK_SIZE_LOG2 = BL;
    int BLOCK_SIZE = (1 << BLOCK_SIZE_LOG2);

    int* A;
    int* B;
    int* C;
    int* C2;

    struct dpu_set_t dpu_set;
    struct dpu_set_t dpu;

    int nr_of_dpus = NUM_DPUS;
    int input_size = DATA_PREP_PARAMS;

    dpu_alloc(nr_of_dpus, NULL, &dpu_set);
    dpu_load(dpu_set, DPU_BINARY, NULL);

    int i = 0;

    int input_size_8bytes = ((input_size * sizeof(int)) % 8) != 0 ? roundup(input_size, 8) : input_size;
    int input_size_dpu = divceil(input_size, nr_of_dpus);
    int input_size_dpu_8bytes = ((input_size_dpu * sizeof(int)) % 8) != 0 ? roundup(input_size_dpu, 8) : input_size_dpu;

    A = malloc(input_size_dpu_8bytes * nr_of_dpus * sizeof(int));
    B = malloc(input_size_dpu_8bytes * nr_of_dpus * sizeof(int));
    C = malloc(input_size_dpu_8bytes * nr_of_dpus * sizeof(int));
    C2 = malloc(input_size_dpu_8bytes * nr_of_dpus * sizeof(int));
    int *bufferA = A;
    int *bufferB = B;
    int *bufferC = C2;

    read_input(A, B, input_size);

    vector_addition_host(C, A, B, input_size);

    int kernel = 0;
    struct dpu_arguments_t* input_arguments = malloc(sizeof(struct dpu_arguments_t) * nr_of_dpus);
    for(i=0; i<nr_of_dpus-1; i++) {
        input_arguments[i].size = input_size_dpu_8bytes * sizeof(int);
        input_arguments[i].transfer_size = input_size_dpu_8bytes * sizeof(int);
        input_arguments[i].kernel = kernel;
    }
    input_arguments[nr_of_dpus-1].size = (input_size_8bytes - input_size_dpu_8bytes * (nr_of_dpus-1)) * sizeof(int);
    input_arguments[nr_of_dpus-1].transfer_size = input_size_dpu_8bytes * sizeof(int);
    input_arguments[nr_of_dpus-1].kernel = kernel;

    i = 0;
    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &input_arguments[i]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, "DPU_INPUT_ARGUMENTS", 0, sizeof(struct dpu_arguments_t), DPU_XFER_DEFAULT);

    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &bufferA[input_size_dpu_8bytes * i]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, 0, input_size_dpu_8bytes * sizeof(int), DPU_XFER_DEFAULT);

    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &bufferB[input_size_dpu_8bytes * i]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, input_size_dpu_8bytes * sizeof(int), input_size_dpu_8bytes * sizeof(int), DPU_XFER_DEFAULT);

    dpu_launch(dpu_set, DPU_SYNCHRONOUS);

    i = 0;
    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &bufferC[input_size_dpu_8bytes * i]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_FROM_DPU, DPU_MRAM_HEAP_POINTER_NAME, input_size_dpu_8bytes * sizeof(int), input_size_dpu_8bytes * sizeof(int), DPU_XFER_DEFAULT);

    int status = 1;
    for (i = 0; i < input_size; i++) {
        if(C[i] != bufferC[i]){ 
            status = 0;
            break;
        }
    }

    assert(status);

    free(input_arguments);
    free(A);
    free(B);
    free(C);
    free(C2);
    dpu_free(dpu_set);
	
    return status ? 0 : -1;
}
