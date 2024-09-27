#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <dpu.h>
#include <dpu_log.h>
#include <unistd.h>
#include <getopt.h>
#include <assert.h>
#include <time.h>

#define NUM_DPUS 1
#define NUM_TASKLETS 1
#define DATA_PREP_PARAMS 64

struct dpu_arguments_t{
	long input_size;
	long slice_per_dpu;
	int kernel;
};

struct dpu_results_t {
    long found;
};

void create_test_file(long * input, long * querys, long  nr_elements, long nr_querys) {
	input[0] = 1;
	for (long i = 1; i < nr_elements; i++) {
		input[i] = input[i - 1] + 1;
	}
	for (long i = 0; i < nr_querys; i++) {
		querys[i] = i;
	}
}

long binarySearch(long * input, long * querys, long input_size, long num_querys)
{
	long result = -1;
	long r;
	for(long q = 0; q < num_querys; q++)
	{
		long l = 0;
		r = input_size;
		while (l <= r) {
			long m = l + (r - l) / 2;

			if (input[m] == querys[q]) {
			    result = m;
			}

			if (input[m] < querys[q]) {
			    l = m + 1;
			} else {
			    r = m - 1;
			}
		}
	}
	return result;
}

int main() {
	struct dpu_set_t dpu_set;
	struct dpu_set_t dpu;
	int nr_of_dpus = NUM_DPUS;
	long input_size = DATA_PREP_PARAMS;
	long num_querys = DATA_PREP_PARAMS / 8;
	long result_host = -1;
	long result_dpu  = -1;

	dpu_alloc(NUM_DPUS, NULL, &dpu_set);
	dpu_load(dpu_set, DPU_BINARY, NULL);

	if(num_querys % (nr_of_dpus * NUM_TASKLETS)) {
	    num_querys = num_querys + (nr_of_dpus * NUM_TASKLETS - num_querys % (nr_of_dpus * NUM_TASKLETS));
    }
    
	assert(num_querys % (nr_of_dpus * NUM_TASKLETS) == 0);

	long * input  = malloc((input_size) * sizeof(long));
	long * querys = malloc((num_querys) * sizeof(long));

	create_test_file(input, querys, input_size, num_querys);

	result_host = binarySearch(input, querys, input_size - 1, num_querys);

	long slice_per_dpu          = num_querys / nr_of_dpus;
	struct dpu_arguments_t input_arguments;
	input_arguments.input_size = input_size;
	input_arguments.slice_per_dpu = slice_per_dpu;
	input_arguments.kernel = 0;

    DPU_FOREACH(dpu_set, dpu, i)
    {
        dpu_prepare_xfer(dpu, &input_arguments);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, "DPU_INPUT_ARGUMENTS", 0, sizeof(struct dpu_arguments_t), DPU_XFER_DEFAULT);

    DPU_FOREACH(dpu_set, dpu, i)
    {
        dpu_prepare_xfer(dpu, input);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, 0, input_size * sizeof(long), DPU_XFER_DEFAULT);

    DPU_FOREACH(dpu_set, dpu, i)
    {
        dpu_prepare_xfer(dpu, &querys[slice_per_dpu * i]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, input_size * sizeof(long), slice_per_dpu * sizeof(long), DPU_XFER_DEFAULT);

    dpu_launch(dpu_set, DPU_SYNCHRONOUS);

    struct dpu_results_t* results_retrieve = malloc(nr_of_dpus * NUM_TASKLETS * sizeof(struct dpu_results_t));
    DPU_FOREACH(dpu_set, dpu, i)
    {
        dpu_prepare_xfer(dpu, &results_retrieve[i * NUM_TASKLETS]);
    }

    dpu_push_xfer(dpu_set, DPU_XFER_FROM_DPU, "DPU_RESULTS", 0, NUM_TASKLETS * sizeof(struct dpu_results_t), DPU_XFER_DEFAULT);

    DPU_FOREACH(dpu_set, dpu, i)
    {
        for(int each_tasklet = 0; each_tasklet < NUM_TASKLETS; each_tasklet++)
        {
            if(results_retrieve[i * NUM_TASKLETS + each_tasklet].found > result_dpu)
            {
                result_dpu = results_retrieve[i * NUM_TASKLETS + each_tasklet].found;
            }
        }
    }

	int status = (result_dpu == result_host);
	assert(status);

	free(input);
	dpu_free(dpu_set);

	return status ? 0 : 1;
}
