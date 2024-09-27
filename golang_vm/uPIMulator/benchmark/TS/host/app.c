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
#include <time.h>

#define NUM_DPUS 1
#define NUM_TASKLETS 1
#define DATA_PREP_PARAMS 64

#define MAX_DATA_VAL 127

struct dpu_arguments_t {
	int ts_length;
    int query_length;
    int query_mean;
    int query_std;
    int slice_per_dpu;
    int exclusion_zone;
    int kernel;
};

struct dpu_result_t {
    int minValue;
    int minIndex;
    int maxValue;
    int maxIndex;
};

int *create_test_file(int* tSeries, int* query, int ts_elements, int query_elements) {
	for (int i = 0; i < ts_elements; i++)
	{
		tSeries[i] = i % MAX_DATA_VAL;
	}

	for (int i = 0; i < query_elements; i++)
	{
		query[i] = i % MAX_DATA_VAL;
	}

	return tSeries;
}

void streamp(int* tSeries, int* AMean, int* ASigma, int *minHost, int *minHostIdx, int ProfileLength,
		int* query, int queryLength, int queryMean, int queryStdDeviation)
{
	int distance;
	int dotprod;
	*minHost    = 2147483647;
	*minHostIdx = 0;

	for (int subseq = 0; subseq < ProfileLength; subseq++)
	{
		dotprod = 0;
		for(int j = 0; j < queryLength; j++)
		{
			dotprod += tSeries[j + subseq] * query[j];
		}

		distance = 2 * (queryLength - (dotprod - queryLength * AMean[subseq] * queryMean) / (ASigma[subseq] * queryStdDeviation));

		if(distance < *minHost)
		{
			*minHost = distance;
			*minHostIdx = subseq;
		}
	}
}

void compute_ts_statistics(int* tSeries, int* AMean, int* ASigma, int timeSeriesLength, int ProfileLength, int queryLength)
{
	int* ACumSum = malloc(sizeof(int) * timeSeriesLength);
	ACumSum[0] = tSeries[0];
	for (int i = 1; i < timeSeriesLength; i++) {
		ACumSum[i] = tSeries[i] + ACumSum[i - 1];
	}

	int* ASqCumSum = malloc(sizeof(int) * timeSeriesLength);
	ASqCumSum[0] = tSeries[0] * tSeries[0];
	for (int i = 1; i < timeSeriesLength; i++) {
		ASqCumSum[i] = tSeries[i] * tSeries[i] + ASqCumSum[i - 1];
    }

	int* ASum = malloc(sizeof(int) * ProfileLength);
	ASum[0] = ACumSum[queryLength - 1];
	for (int i = 0; i < timeSeriesLength - queryLength; i++) {
		ASum[i + 1] = ACumSum[queryLength + i] - ACumSum[i];
    }

	int* ASumSq = malloc(sizeof(int) * ProfileLength);
	ASumSq[0] = ASqCumSum[queryLength - 1];
	for (int i = 0; i < timeSeriesLength - queryLength; i++) {
		ASumSq[i + 1] = ASqCumSum[queryLength + i] - ASqCumSum[i];
    }

	int * AMean_tmp = malloc(sizeof(int) * ProfileLength);
	for (int i = 0; i < ProfileLength; i++) {
		AMean_tmp[i] = ASum[i] / queryLength;
    }

	int* ASigmaSq = malloc(sizeof(int) * ProfileLength);
	for (int i = 0; i < ProfileLength; i++) {
		ASigmaSq[i] = ASumSq[i] / queryLength - AMean[i] * AMean[i];
    }

	for (int i = 0; i < ProfileLength; i++)
	{
		ASigma[i] = sqrt(ASigmaSq[i]);
		AMean[i]  = AMean_tmp[i];
	}

	free(ACumSum);
	free(ASqCumSum);
	free(ASum);
	free(ASumSq);
	free(ASigmaSq);
	free(AMean_tmp);
}

int main() {
    int* tSeries = malloc((1 << 15) * sizeof(int));
    int* query = malloc((1 << 15) * sizeof(int));
    int* AMean = malloc((1 << 15) * sizeof(int));
    int* ASigma = malloc((1 << 15) * sizeof(int));
    int* minHost = malloc(sizeof(int));
    int* minHostIdx = malloc(sizeof(int));

	struct dpu_set_t dpu_set;
	struct dpu_set_t dpu;
	int nr_of_dpus = NUM_DPUS;

	dpu_alloc(NUM_DPUS, NULL, &dpu_set);
	dpu_load(dpu_set, DPU_BINARY, NULL);

	int ts_size = DATA_PREP_PARAMS;
	int query_length = 64;

	if(ts_size % (nr_of_dpus * NUM_TASKLETS*query_length)) {
		ts_size = ts_size +  (nr_of_dpus * NUM_TASKLETS * query_length - ts_size % (nr_of_dpus * NUM_TASKLETS*query_length));
    }

	create_test_file(tSeries, query, ts_size, query_length);
	compute_ts_statistics(tSeries, AMean, ASigma, ts_size, ts_size - query_length, query_length);

	int query_mean;
	int queryMean = 0;
	for(int i = 0; i < query_length; i++) {
	    queryMean += query[i];
	}
	queryMean /= query_length;
	query_mean = queryMean;

	int query_std;
	int queryStdDeviation;
	int queryVariance = 0;
	for(int i = 0; i < query_length; i++)
	{
		queryVariance += (query[i] - queryMean) * (query[i] - queryMean);
	}
	queryVariance /= query_length;
	queryStdDeviation = sqrt(queryVariance);
	query_std = queryStdDeviation;

	int *bufferTS     = tSeries;
	int *bufferQ      = query;
	int *bufferAMean  = AMean;
	int *bufferASigma = ASigma;

	int slice_per_dpu = ts_size / nr_of_dpus;

	int kernel = 0;
	struct dpu_arguments_t input_arguments;
    input_arguments.ts_length = ts_size;
    input_arguments.query_length = query_length;
    input_arguments.query_mean = query_mean;
    input_arguments.query_std = query_std;
    input_arguments.slice_per_dpu = slice_per_dpu;
    input_arguments.exclusion_zone = 0;
    input_arguments.kernel = kernel;

	struct dpu_result_t result;
	result.minValue = 2147483647;
	result.minIndex = 0;
	result.maxValue = 0;
	result.maxIndex = 0;

    DPU_FOREACH(dpu_set, dpu, i) {
        input_arguments.exclusion_zone = 0;

        dpu_copy_to(dpu, "DPU_INPUT_ARGUMENTS", 0, &input_arguments, sizeof(struct dpu_arguments_t));
    }

    int mem_offset = 0;
    DPU_FOREACH(dpu_set, dpu, i)
    {
        dpu_prepare_xfer(dpu, bufferQ);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, 0, query_length * sizeof(int), DPU_XFER_DEFAULT);

    mem_offset += query_length * sizeof(int);

    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &bufferTS[slice_per_dpu * i]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, mem_offset,(slice_per_dpu + query_length)*sizeof(int), DPU_XFER_DEFAULT);

    mem_offset += ((slice_per_dpu + query_length) * sizeof(int));

    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &bufferAMean[slice_per_dpu * i]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, mem_offset, (slice_per_dpu + query_length)*sizeof(int), DPU_XFER_DEFAULT);

    mem_offset += ((slice_per_dpu + query_length) * sizeof(int));

    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &bufferASigma[slice_per_dpu * i]);
    }

    dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, mem_offset, (slice_per_dpu + query_length)*sizeof(int), DPU_XFER_DEFAULT);

    dpu_launch(dpu_set, DPU_SYNCHRONOUS);

    struct dpu_result_t* results_retrieve = malloc(nr_of_dpus * NUM_TASKLETS * sizeof(struct dpu_result_t));

    DPU_FOREACH(dpu_set, dpu, i) {
        dpu_prepare_xfer(dpu, &results_retrieve[i * NUM_TASKLETS]);
    }
    dpu_push_xfer(dpu_set, DPU_XFER_FROM_DPU, "DPU_RESULTS", 0, NUM_TASKLETS * sizeof(struct dpu_result_t), DPU_XFER_DEFAULT);

    DPU_FOREACH(dpu_set, dpu, i) {
        for (int each_tasklet = 0; each_tasklet < NUM_TASKLETS; each_tasklet++) {
            if(results_retrieve[i * NUM_TASKLETS + each_tasklet].minValue < result.minValue && results_retrieve[i * NUM_TASKLETS + each_tasklet].minValue > 0)
            {
                result.minValue = results_retrieve[i * NUM_TASKLETS + each_tasklet].minValue;
                result.minIndex = results_retrieve[i * NUM_TASKLETS + each_tasklet].minIndex + (i * slice_per_dpu);
            }

        }
    }

    streamp(tSeries, AMean, ASigma, minHost, minHostIdx, ts_size - query_length - 1, query, query_length, query_mean, query_std);

	int status = (*minHost == result.minValue);
	assert(status);

	dpu_free(dpu_set);

	return 0;
}
