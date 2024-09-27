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

struct dpu_arguments_t {
    int n_size;
    int n_size_pad;
    int nr_rows;
    int max_rows;
};

struct dpu_info_t {
  int rows_per_dpu;
  int rows_per_dpu_pad;
  int prev_rows_dpu;
};

void init_data(int* A, int* B, int m_size, int n_size) {
	for (int i = 0; i < m_size * n_size; i++)
	{
		A[i] = i % 50;
	}

	for (int i = 0; i < n_size; i++)
	{
		B[i] = i % 50;
	}
}

void gemv_host(int* C, int* A, int* B, int m_size, int n_size) {
	for (int i = 0; i < m_size; i++)
	{
		C[i] = 0;
	}

	for (int m = 0; m < m_size; m++) {
		for (int n = 0; n < n_size; n++)
		{
			C[m] += A[m * n_size + n] * B[n];
		}
	}
}

int main() {
    int* A;
    int* B;
    int* C;
    int* C_dpu;

	struct dpu_set_t dpu_set;
	struct dpu_set_t dpu;
	int nr_of_dpus = NUM_DPUS;

	dpu_alloc(NR_DPUS, NULL, &dpu_set);
	dpu_load(dpu_set, DPU_BINARY, NULL);

    int i;
	int m_size = DATA_PREP_PARAMS;
	int n_size = 64;

	struct dpu_info_t* dpu_info = malloc(nr_of_dpus * sizeof(struct dpu_info_t));
	struct dpu_arguments_t *input_args = malloc(nr_of_dpus * sizeof(struct dpu_arguments_t));
	int max_rows_per_dpu = 0;
	int n_size_pad = n_size;
	if(n_size % 2 == 1)
	{
		n_size_pad++;
	}

	DPU_FOREACH(dpu_set, dpu, i) {
		int rows_per_dpu;
		int prev_rows_dpu = 0;
		int chunks = m_size / nr_of_dpus;
		rows_per_dpu = chunks;
		int rest_rows = m_size % nr_of_dpus;
		if (i < rest_rows) {
			rows_per_dpu++;
		}
		if (rest_rows > 0) {
			if (i >= rest_rows) {
				prev_rows_dpu = rest_rows * (chunks + 1) + (i - rest_rows) * chunks;
			}
			else {
				prev_rows_dpu = i * (chunks + 1);
			}
		} else {
			prev_rows_dpu = i * chunks;
		}

		int rows_per_dpu_pad = rows_per_dpu;
		if (rows_per_dpu_pad % 2 == 1) {
			rows_per_dpu_pad++;
        }
		if (rows_per_dpu_pad > max_rows_per_dpu) {
			max_rows_per_dpu = rows_per_dpu_pad;
        }

		dpu_info[i].rows_per_dpu = rows_per_dpu;
		dpu_info[i].rows_per_dpu_pad = rows_per_dpu_pad;
		dpu_info[i].prev_rows_dpu = prev_rows_dpu;

		input_args[i].n_size = n_size;
		input_args[i].n_size_pad = n_size_pad;
		input_args[i].nr_rows = rows_per_dpu;
	}

	A = malloc(max_rows_per_dpu * nr_of_dpus * n_size_pad * sizeof(int));
	B = malloc(n_size_pad * sizeof(int));
	C = malloc(max_rows_per_dpu * nr_of_dpus * sizeof(int));

	init_data(A, B, m_size, n_size);

	gemv_host(C, A, B, m_size, n_size);

    DPU_FOREACH(dpu_set, dpu, i) {
        input_args[i].max_rows = max_rows_per_dpu;

        dpu_prepare_xfer(dpu, &input_args[i]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, "DPU_INPUT_ARGUMENTS", 0, sizeof(struct dpu_arguments_t), DPU_XFER_DEFAULT);

    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &A[dpu_info[i].prev_rows_dpu * n_size]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, 0, max_rows_per_dpu * n_size_pad * sizeof(int), DPU_XFER_DEFAULT);

    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, B);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, max_rows_per_dpu * n_size_pad * sizeof(int) , n_size_pad * sizeof(int), DPU_XFER_DEFAULT);

    dpu_launch(dpu_set, DPU_SYNCHRONOUS);

    C_dpu = malloc(max_rows_per_dpu * nr_of_dpus * sizeof(int));

    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &C_dpu[i * max_rows_per_dpu]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_FROM_DPU, DPU_MRAM_HEAP_POINTER_NAME, max_rows_per_dpu * n_size_pad * sizeof(int) + n_size_pad * sizeof(int), max_rows_per_dpu * sizeof(int), DPU_XFER_DEFAULT);

	int status = 1;

	i = 0;
	for (int n = 0; n < nr_of_dpus; n++) {
		for (int j = 0; j < dpu_info[n].rows_per_dpu; j++) {
			if(C[i] != C_dpu[n * max_rows_per_dpu + j]) {
				status = 0;
			}
			i++;
		}
	}

	assert(status);

	free(A);
	free(B);
	free(C);
	free(C_dpu);
	dpu_free(dpu_set);

	return status ? 0 : -1;
}
