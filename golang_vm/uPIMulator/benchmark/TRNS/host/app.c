#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <dpu.h>
#include <dpu_log.h>
#include <unistd.h>
#include <getopt.h>
#include <assert.h>
#include <math.h>

#define NUM_DPUS 1
#define NUM_TASKLETS 1
#define DATA_PREP_PARAMS 64

struct dpu_arguments_t {
    int m;
    int n;
    int M_;
    int kernel;
};

void read_input(long* A, int nr_elements) {
    for (int i = 0; i < nr_elements; i++) {
        A[i] = i % 100;
    }
}

void trns_host(long* input, int A, int B, int b){
   long* output = malloc(sizeof(long) * A * B * b);
   int next;
   for (int j = 0; j < b; j++){
      for (int i = 0; i < A * B; i++){
         next = (i * A) - (A * B - 1) * (i / B);
         output[next * b + j] = input[i*b+j];
      }
   }
   for (int k = 0; k < A * B * b; k++){
      input[k] = output[k];
   }
   free(output);
}

int main() {
    long* A_host;
    long* A_backup;
    long* A_result;

    struct dpu_set_t dpu_set;
    struct dpu_set_t dpu;

    int nr_of_dpus = NUM_DPUS;

    int N_ = 64;
    int n = 8;
    int M_ = DATA_PREP_PARAMS;
    int m = 4;

    A_host = malloc(M_ * m * N_ * n * sizeof(long));
    A_backup = malloc(M_ * m * N_ * n * sizeof(long));
    A_result = malloc(M_ * m * N_ * n * sizeof(long));
    char* done_host = malloc(M_ * n * sizeof(char));

    for (int i = 0; i < M_ * n; i++) {
        done_host[i] = 0;
    }

    read_input(A_host, M_ * m * N_ * n);

    for (int i = 0; i < M_ * m * N_ * n; i++) {
        A_backup[i] = A_host[i];
    }

    trns_host(A_host, M_ * m, N_ * n, 1);

    int curr_dpu = 0;
    int active_dpus;
    int active_dpus_before = 0;
    int first_round = 1;

    while(curr_dpu < N_){
        if((N_ - curr_dpu) > NUM_DPUS){
            active_dpus = NUM_DPUS;
        } else {
            active_dpus = (N_ - curr_dpu);
        }

        if((active_dpus_before != active_dpus) && (!(first_round))){
            dpu_free(dpu_set);
            dpu_alloc(active_dpus, NULL, &dpu_set);
            dpu_load(dpu_set, DPU_BINARY, NULL);

            nr_of_dpus = active_dpus;
        } else if (first_round){
            dpu_alloc(active_dpus, NULL, &dpu_set);
            dpu_load(dpu_set, DPU_BINARY, NULL);

            nr_of_dpus = active_dpus;
        }

        for(int j = 0; j < M_ * m; j++){
            DPU_FOREACH(dpu_set, dpu, i) {
                dpu_prepare_xfer(dpu, &A_backup[j * N_ * n + n * (i + curr_dpu)]);
            }
            dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, sizeof(long) * j * n, sizeof(long) * n, DPU_XFER_DEFAULT);
        }

        DPU_FOREACH(dpu_set, dpu, i) {
            dpu_prepare_xfer(dpu, done_host);
        }
        dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, M_ * m * n * sizeof(long), (M_ * n) / 8 == 0 ? 8 : M_ * n, DPU_XFER_DEFAULT);

        int kernel = 0;
        struct dpu_arguments_t input_arguments;
        input_arguments.m = m;
        input_arguments.n = n;
        input_arguments.M_ = M_;
        input_arguments.kernel = kernel;

        DPU_FOREACH(dpu_set, dpu, i) {
            dpu_prepare_xfer(dpu, &input_arguments);
        }
        dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, "DPU_INPUT_ARGUMENTS", 0, sizeof(struct dpu_arguments_t), DPU_XFER_DEFAULT);

        dpu_launch(dpu_set, DPU_SYNCHRONOUS);

        kernel = 1;
        struct dpu_arguments_t input_arguments2;
        input_arguments2.m = m;
        input_arguments2.n = n;
        input_arguments2.M_ = M_;
        input_arguments2.kernel = kernel;

        DPU_FOREACH(dpu_set, dpu, i) {
            dpu_prepare_xfer(dpu, &input_arguments2);
        }
        dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, "DPU_INPUT_ARGUMENTS", 0, sizeof(struct dpu_arguments_t), DPU_XFER_DEFAULT);

        dpu_launch(dpu_set, DPU_SYNCHRONOUS);

        DPU_FOREACH(dpu_set, dpu, i) {
            dpu_prepare_xfer(dpu, &A_result[curr_dpu * m * n * M_]);
            curr_dpu++;
        }
        dpu_push_xfer(dpu_set, DPU_XFER_FROM_DPU, DPU_MRAM_HEAP_POINTER_NAME, 0, sizeof(long) * m * n * M_, DPU_XFER_DEFAULT);

        if(first_round){
            first_round = 0;
        }
    }

    dpu_free(dpu_set);

    int status = 1;
    for (int i = 0; i < M_ * m * N_ * n; i++) {
        if(A_host[i] != A_result[i]){ 
            status = 0;
            break;
        }
    }

    assert(status);

    free(A_host);
    free(A_backup);
    free(A_result);
    free(done_host);
	
    return status ? 0 : -1;
}
