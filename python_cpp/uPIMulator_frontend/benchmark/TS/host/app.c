/**
 * app.c
 * TS Host Application Source File
 *
 */
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

#if ENERGY
#include <dpu_probe.h>
#endif

#include "params.h"
#include "timer.h"

// Define the DPU Binary path as DPU_BINARY here
#define DPU_BINARY "./bin/ts_dpu"

#define MAX_DATA_VAL 127

static DTYPE tSeries[1 << 26];
static DTYPE query  [1 << 15];
static DTYPE AMean  [1 << 26];
static DTYPE ASigma [1 << 26];
static DTYPE minHost;
static DTYPE minHostIdx;

// Create input arrays
static DTYPE *create_test_file(unsigned int ts_elements, unsigned int query_elements) {
	srand(0);

	for (uint64_t i = 0; i < ts_elements; i++)
	{
		tSeries[i] = i % MAX_DATA_VAL;
	}

	for (uint64_t i = 0; i < query_elements; i++)
	{
		query[i] = i % MAX_DATA_VAL;
	}

	return tSeries;
}

// Compute output in the host
static void streamp(DTYPE* tSeries, DTYPE* AMean, DTYPE* ASigma, int ProfileLength,
		DTYPE* query, int queryLength, DTYPE queryMean, DTYPE queryStdDeviation)
{
	DTYPE distance;
	DTYPE dotprod;
	minHost    = INT32_MAX;
	minHostIdx = 0;

	for (int subseq = 0; subseq < ProfileLength; subseq++)
	{
		dotprod = 0;
		for(int j = 0; j < queryLength; j++)
		{
			dotprod += tSeries[j + subseq] * query[j];
		}

		distance = 2 * (queryLength - (dotprod - queryLength * AMean[subseq]
					* queryMean) / (ASigma[subseq] * queryStdDeviation));

		if(distance < minHost)
		{
			minHost = distance;
			minHostIdx = subseq;
		}
	}
}

static void compute_ts_statistics(unsigned int timeSeriesLength, unsigned int ProfileLength, unsigned int queryLength)
{
	double* ACumSum = malloc(sizeof(double) * timeSeriesLength);
	ACumSum[0] = tSeries[0];
	for (uint64_t i = 1; i < timeSeriesLength; i++)
		ACumSum[i] = tSeries[i] + ACumSum[i - 1];
	double* ASqCumSum = malloc(sizeof(double) * timeSeriesLength);
	ASqCumSum[0] = tSeries[0] * tSeries[0];
	for (uint64_t i = 1; i < timeSeriesLength; i++)
		ASqCumSum[i] = tSeries[i] * tSeries[i] + ASqCumSum[i - 1];
	double* ASum = malloc(sizeof(double) * ProfileLength);
	ASum[0] = ACumSum[queryLength - 1];
	for (uint64_t i = 0; i < timeSeriesLength - queryLength; i++)
		ASum[i + 1] = ACumSum[queryLength + i] - ACumSum[i];
	double* ASumSq = malloc(sizeof(double) * ProfileLength);
	ASumSq[0] = ASqCumSum[queryLength - 1];
	for (uint64_t i = 0; i < timeSeriesLength - queryLength; i++)
		ASumSq[i + 1] = ASqCumSum[queryLength + i] - ASqCumSum[i];
	double * AMean_tmp = malloc(sizeof(double) * ProfileLength);
	for (uint64_t i = 0; i < ProfileLength; i++)
		AMean_tmp[i] = ASum[i] / queryLength;
	double* ASigmaSq = malloc(sizeof(double) * ProfileLength);
	for (uint64_t i = 0; i < ProfileLength; i++)
		ASigmaSq[i] = ASumSq[i] / queryLength - AMean[i] * AMean[i];
	for (uint64_t i = 0; i < ProfileLength; i++)
	{
		ASigma[i] = sqrt(ASigmaSq[i]);
		AMean[i]  = (DTYPE) AMean_tmp[i];
	}

	free(ACumSum);
	free(ASqCumSum);
	free(ASum);
	free(ASumSq);
	free(ASigmaSq);
	free(AMean_tmp);
}

// Main of the Host Application
int main(int argc, char **argv) {

	// Timer declaration
	Timer timer;

	struct Params p = input_params(argc, argv);
	struct dpu_set_t dpu_set, dpu;
	uint32_t nr_of_dpus;

	// Allocate DPUs and load binary
	DPU_ASSERT(dpu_alloc(NR_DPUS, NULL, &dpu_set));
	DPU_ASSERT(dpu_load(dpu_set, DPU_BINARY, NULL));
	DPU_ASSERT(dpu_get_nr_dpus(dpu_set, &nr_of_dpus));

#if ENERGY
	struct dpu_probe_t probe;
	DPU_ASSERT(dpu_probe_init("energy_probe", &probe));
#endif

	unsigned long int ts_size =  p.input_size_n;
	const unsigned int query_length = p.input_size_m;

	// Size adjustment
	if(ts_size % (nr_of_dpus * NR_TASKLETS*query_length))
		ts_size = ts_size +  (nr_of_dpus * NR_TASKLETS * query_length - ts_size % (nr_of_dpus * NR_TASKLETS*query_length));

	// Create an input file with arbitrary data
	create_test_file(ts_size, query_length);
	compute_ts_statistics(ts_size, ts_size - query_length, query_length);

	DTYPE query_mean;
	double queryMean = 0;
	for(unsigned i = 0; i < query_length; i++) queryMean += query[i];
	queryMean /= (double) query_length;
	query_mean = (DTYPE) queryMean;

	DTYPE query_std;
	double queryStdDeviation;
	double queryVariance = 0;
	for(unsigned i = 0; i < query_length; i++)
	{
		queryVariance += (query[i] - queryMean) * (query[i] - queryMean);
	}
	queryVariance /= (double) query_length;
	queryStdDeviation = sqrt(queryVariance);
	query_std = (DTYPE) queryStdDeviation;

	DTYPE *bufferTS     = tSeries;
	DTYPE *bufferQ      = query;
	DTYPE *bufferAMean  = AMean;
	DTYPE *bufferASigma = ASigma;

	uint32_t slice_per_dpu = ts_size / nr_of_dpus;

	unsigned int kernel = 0;
	dpu_arguments_t input_arguments = {ts_size, query_length, query_mean, query_std, slice_per_dpu, 0, kernel};
	uint32_t mem_offset;

	dpu_result_t result;
	result.minValue = INT32_MAX;
	result.minIndex = 0;
	result.maxValue = 0;
	result.maxIndex = 0;

	for (int rep = 0; rep < p.n_warmup + p.n_reps; rep++) {

		if (rep >= p.n_warmup)
			start(&timer, 1, rep - p.n_warmup);
		uint32_t i = 0;

		DPU_FOREACH(dpu_set, dpu) {
			input_arguments.exclusion_zone = 0;

			DPU_ASSERT(dpu_copy_to(dpu, "DPU_INPUT_ARGUMENTS", 0, (const void *) &input_arguments, sizeof(input_arguments)));
			i++;
		}

		i = 0;
		mem_offset = 0;
		DPU_FOREACH(dpu_set, dpu, i)
		{
			DPU_ASSERT(dpu_prepare_xfer(dpu, bufferQ));
		}

		DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, 0, query_length * sizeof(DTYPE), DPU_XFER_DEFAULT));

		i = 0;

		mem_offset += query_length * sizeof(DTYPE);
		DPU_FOREACH(dpu_set, dpu, i) {
			DPU_ASSERT(dpu_prepare_xfer(dpu, bufferTS + slice_per_dpu * i));
		}

		DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, mem_offset,(slice_per_dpu + query_length)*sizeof(DTYPE), DPU_XFER_DEFAULT));

		mem_offset += ((slice_per_dpu + query_length) * sizeof(DTYPE));

		i = 0;
		DPU_FOREACH(dpu_set, dpu, i) {
			DPU_ASSERT(dpu_prepare_xfer(dpu, bufferAMean + slice_per_dpu * i));
		}

		DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, mem_offset, (slice_per_dpu + query_length)*sizeof(DTYPE), DPU_XFER_DEFAULT));

		i = 0;

		mem_offset += ((slice_per_dpu + query_length) * sizeof(DTYPE));

		DPU_FOREACH(dpu_set, dpu, i) {
			DPU_ASSERT(dpu_prepare_xfer(dpu, bufferASigma + slice_per_dpu * i));
		}

		DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, mem_offset, (slice_per_dpu + query_length)*sizeof(DTYPE), DPU_XFER_DEFAULT));

		if (rep >= p.n_warmup)
			stop(&timer, 1);

		// Run kernel on DPUs
		if (rep >= p.n_warmup)
		{
			start(&timer, 2, rep - p.n_warmup);
#if ENERGY
			DPU_ASSERT(dpu_probe_start(&probe));
#endif
		}

		DPU_ASSERT(dpu_launch(dpu_set, DPU_SYNCHRONOUS));

		if (rep >= p.n_warmup)
		{
			stop(&timer, 2);
#if ENERGY
			DPU_ASSERT(dpu_probe_stop(&probe));
#endif
		}

		dpu_result_t* results_retrieve[nr_of_dpus];

		if (rep >= p.n_warmup)
			start(&timer, 3, rep - p.n_warmup);

		DPU_FOREACH(dpu_set, dpu, i) {
			results_retrieve[i] = (dpu_result_t*)malloc(NR_TASKLETS * sizeof(dpu_result_t));
		}


		DPU_FOREACH(dpu_set, dpu, i) {
			DPU_ASSERT(dpu_prepare_xfer(dpu, results_retrieve[i]));
		}
		DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_FROM_DPU, "DPU_RESULTS", 0, NR_TASKLETS * sizeof(dpu_result_t), DPU_XFER_DEFAULT));

		i = 0;
		DPU_FOREACH(dpu_set, dpu, i) {
			for (unsigned int each_tasklet = 0; each_tasklet < NR_TASKLETS; each_tasklet++) {
				if(results_retrieve[i][each_tasklet].minValue < result.minValue && results_retrieve[i][each_tasklet].minValue > 0)
				{
					result.minValue = results_retrieve[i][each_tasklet].minValue;
					result.minIndex = (DTYPE)results_retrieve[i][each_tasklet].minIndex + (i * slice_per_dpu);
				}

			}
			free(results_retrieve[i]);
			i++;
		}

		if(rep >= p.n_warmup)
			stop(&timer, 3);


#if PRINT
		printf("LOGS\n");
		DPU_FOREACH(dpu_set, dpu) {
			DPU_ASSERT(dpu_log_read(dpu, stdout));
		}
#endif

		if (rep >= p.n_warmup)
			start(&timer, 4, rep - p.n_warmup);
		streamp(tSeries, AMean, ASigma, ts_size - query_length - 1, query, query_length, query_mean, query_std);
		if(rep >= p.n_warmup)
			stop(&timer, 4);
	}

#if ENERGY
	double acc_energy, avg_energy, acc_time, avg_time;
	DPU_ASSERT(dpu_probe_get(&probe, DPU_ENERGY, DPU_ACCUMULATE, &acc_energy));
	DPU_ASSERT(dpu_probe_get(&probe, DPU_ENERGY, DPU_AVERAGE, &avg_energy));
	DPU_ASSERT(dpu_probe_get(&probe, DPU_TIME, DPU_ACCUMULATE, &acc_time));
	DPU_ASSERT(dpu_probe_get(&probe, DPU_TIME, DPU_AVERAGE, &avg_time));
#endif

	// Print timing results
	printf("CPU Version Time (ms): ");
	print(&timer, 4, p.n_reps);
	printf("Inter-DPU Time (ms): ");
	print(&timer, 0, p.n_reps);
	printf("CPU-DPU Time (ms): ");
	print(&timer, 1, p.n_reps);
	printf("DPU Kernel Time (ms): ");
	print(&timer, 2, p.n_reps);
	printf("DPU-CPU Time (ms): ");
	print(&timer, 3, p.n_reps);

#if ENERGY
	printf("Energy (J): %f J\t", avg_energy);
#endif

	int status = (minHost == result.minValue);
	if (status) {
		printf("[" ANSI_COLOR_GREEN "OK" ANSI_COLOR_RESET "] results are equal\n");
	} else {
		printf("[" ANSI_COLOR_RED "ERROR" ANSI_COLOR_RESET "] results differ!\n");
	}

	DPU_ASSERT(dpu_free(dpu_set));

#if ENERGY
	DPU_ASSERT(dpu_probe_deinit(&probe));
#endif

	return 0;
}
