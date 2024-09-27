#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <dpu.h>
#include <dpu_log.h>
#include <unistd.h>
#include <getopt.h>
#include <assert.h>

#define NUM_DPUS 1
#define NUM_TASKLETS 1
#define DATA_PREP_PARAMS 64

struct dpu_arguments_t{
    int size;
	int kernel;
    long t_count;
};

struct dpu_results_t{
    long cycles;
    long t_count;
};

void read_input(long* A, int nr_elements) {
    for (int i = 0; i < nr_elements; i++) {
        A[i] = i % 100;
    }
}

long reduction_host(long* A, int nr_elements) {
    long count = 0;
    for (int i = 0; i < nr_elements; i++) {
        count += A[i];
    }
    return count;
}

int roundup(int n, int m) {
    return ((n / m) * m + m);
}

int divceil(int n, int m) {
    return ((n-1) / m + 1);
}

int main() {
    long* A;

    struct dpu_set_t dpu_set;
    struct dpu_set_t dpu;
    int nr_of_dpus = NUM_DPUS;

    dpu_alloc(NUM_DPUS, NULL, &dpu_set);
    dpu_load(dpu_set, DPU_BINARY, NULL);

    int input_size = DATA_PREP_PARAMS;
    int input_size_8bytes = ((input_size * sizeof(long)) % 8) != 0 ? roundup(input_size, 8) : input_size;
    int input_size_dpu = divceil(input_size, nr_of_dpus);
    int input_size_dpu_8bytes = ((input_size_dpu * sizeof(long)) % 8) != 0 ? roundup(input_size_dpu, 8) : input_size_dpu;

    A = malloc(input_size_dpu_8bytes * nr_of_dpus * sizeof(long));
    long *bufferA = A;
    long count = 0;
    long count_host = 0;

    read_input(A, input_size);

    count_host = reduction_host(A, input_size);

    count = 0;
    int kernel = 0;
    struct dpu_arguments_t* input_arguments = malloc(NUM_DPUS * sizeof(struct dpu_arguments_t));
    for(int i=0; i<nr_of_dpus-1; i++) {
        input_arguments[i].size=input_size_dpu_8bytes * sizeof(long);
        input_arguments[i].kernel=kernel;
    }
    input_arguments[nr_of_dpus-1].size=(input_size_8bytes - input_size_dpu_8bytes * (NUM_DPUS-1)) * sizeof(long);
    input_arguments[nr_of_dpus-1].kernel=kernel;

    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &input_arguments[i]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, "DPU_INPUT_ARGUMENTS", 0, sizeof(struct dpu_arguments_t), DPU_XFER_DEFAULT);
    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &bufferA[input_size_dpu_8bytes * i]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, 0, input_size_dpu_8bytes * sizeof(long), DPU_XFER_DEFAULT);

    dpu_launch(dpu_set, DPU_SYNCHRONOUS);

    struct dpu_results_t* results = malloc(nr_of_dpus * sizeof(struct dpu_results_t));
    long* results_count = malloc(nr_of_dpus * sizeof(long));

    struct dpu_results_t* results_retrieve = malloc(nr_of_dpus * NUM_TASKLETS * sizeof(struct dpu_results_t));

    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &results_retrieve[i * NUM_TASKLETS]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_FROM_DPU, "DPU_RESULTS", 0, NUM_TASKLETS * sizeof(struct dpu_results_t), DPU_XFER_DEFAULT);

    DPU_FOREACH(dpu_set, dpu, i) {
        for (int each_tasklet = 0; each_tasklet < NUM_TASKLETS; each_tasklet++) {
            if(each_tasklet == 0) {
                results[i].t_count = results_retrieve[i * NUM_TASKLETS + each_tasklet].t_count;
            }
        }

        count += results[i].t_count;
    }

    int status = 1;
    if(count != count_host) {
        status = 0;
    }

    assert(status);

    free(A);
    dpu_free(dpu_set);
	
    return status ? 0 : -1;
}
