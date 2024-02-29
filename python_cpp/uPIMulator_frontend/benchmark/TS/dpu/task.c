/*
 * STREAMP implementation of Matrix Profile with multiple tasklets
 *
 */

#include <stdint.h>
#include <stdio.h>
#include <defs.h>
#include <mram.h>
#include <alloc.h>
#include <mram.h>
#include <barrier.h>
#include "common.h"

#define DOTPIP BLOCK_SIZE / sizeof(DTYPE)

__host dpu_arguments_t DPU_INPUT_ARGUMENTS;
__host dpu_result_t DPU_RESULTS[NR_TASKLETS];

// Dot product
void __attribute__ ((noinline)) dot_product(DTYPE *vectorA, DTYPE *vectorA_aux, DTYPE *vectorB, DTYPE * result) {

	for(uint32_t i = 0; i <  BLOCK_SIZE / sizeof(DTYPE); i++)
	{
		for(uint32_t j = 0; j < DOTPIP; j++)
		{
			if((j + i) > BLOCK_SIZE / sizeof(DTYPE) - 1)
			{
				result[j] += vectorA_aux[(j + i) - BLOCK_SIZE / sizeof(DTYPE)]  * vectorB[i];
			}
			else
			{
				result[j] += vectorA[j + i] * vectorB[i];
			}
		}
	}
}

BARRIER_INIT(my_barrier, NR_TASKLETS);

extern int main_kernel1(void);

int(*kernels[nr_kernels])(void) = {main_kernel1};

int main(void){
	// Kernel
	return kernels[DPU_INPUT_ARGUMENTS.kernel]();
}

// main_kernel1
int main_kernel1() {
	unsigned int tasklet_id = me();
#if PRINT
	printf("tasklet_id = %u\n", tasklet_id);
#endif
	if(tasklet_id == 0){
		mem_reset(); // Reset the heap
	}
	// Barrier
	barrier_wait(&my_barrier);

	// Input arguments
	uint32_t query_length  = DPU_INPUT_ARGUMENTS.query_length;
	DTYPE query_mean       = DPU_INPUT_ARGUMENTS.query_mean;
	DTYPE query_std        = DPU_INPUT_ARGUMENTS.query_std;
	uint32_t slice_per_dpu = DPU_INPUT_ARGUMENTS.slice_per_dpu;

	// Boundaries for current tasklet
	uint32_t myStartElem = tasklet_id  * (slice_per_dpu / (NR_TASKLETS));
	uint32_t myEndElem   = myStartElem + (slice_per_dpu / (NR_TASKLETS)) - 1;

	// Check time series limit
	if(myEndElem > slice_per_dpu - query_length) myEndElem = slice_per_dpu - query_length;

	// Starting address of the current processing block in MRAM
	uint32_t mem_offset = (uint32_t) DPU_MRAM_HEAP_POINTER;

	// Starting address of the query subsequence
	uint32_t current_mram_block_addr_query = (uint32_t)(mem_offset);
	mem_offset += query_length * sizeof(DTYPE);

	// Starting address of the time series slice
	mem_offset += myStartElem * sizeof(DTYPE);
	uint32_t starting_offset_ts = mem_offset;
	uint32_t current_mram_block_addr_TS = (uint32_t) mem_offset;

	// Starting address of the time series means
	mem_offset += (slice_per_dpu + query_length) * sizeof(DTYPE);
	uint32_t current_mram_block_addr_TSMean = (uint32_t)(mem_offset);

	// Starting address of the time series standard deviations
	mem_offset += (slice_per_dpu + query_length) * sizeof(DTYPE);
	uint32_t current_mram_block_addr_TSSigma = (uint32_t)(mem_offset);

	// Initialize local caches to store the MRAM blocks
	DTYPE *cache_TS       = (DTYPE *) mem_alloc(BLOCK_SIZE);
	DTYPE *cache_TS_aux   = (DTYPE *) mem_alloc(BLOCK_SIZE);
	DTYPE *cache_query    = (DTYPE *) mem_alloc(BLOCK_SIZE);
	DTYPE *cache_TSMean   = (DTYPE *) mem_alloc(BLOCK_SIZE);
	DTYPE *cache_TSSigma  = (DTYPE *) mem_alloc(BLOCK_SIZE);
	DTYPE *cache_dotprods = (DTYPE *) mem_alloc(BLOCK_SIZE);

	// Create result structure pointer
	dpu_result_t *result = &DPU_RESULTS[tasklet_id];

	// Auxiliary variables
	DTYPE distance;
	DTYPE min_distance = DTYPE_MAX;
	uint32_t min_index = 0;


	for(uint32_t i = myStartElem; i < myEndElem; i+= (BLOCK_SIZE / sizeof(DTYPE)))
	{
		for(uint32_t d = 0; d < DOTPIP; d++)
			cache_dotprods[d] = 0;

		current_mram_block_addr_TS    = (uint32_t) starting_offset_ts + (i - myStartElem) * sizeof(DTYPE);
		current_mram_block_addr_query = (uint32_t) DPU_MRAM_HEAP_POINTER;

		for(uint32_t j = 0; j < (query_length) / (BLOCK_SIZE / sizeof(DTYPE)); j++)
		{
			mram_read((__mram_ptr void const *) current_mram_block_addr_TS, cache_TS, BLOCK_SIZE);
			mram_read((__mram_ptr void const *) current_mram_block_addr_TS + BLOCK_SIZE, cache_TS_aux, BLOCK_SIZE);
			mram_read((__mram_ptr void const *) current_mram_block_addr_query, cache_query, BLOCK_SIZE);

			current_mram_block_addr_TS    += BLOCK_SIZE;
			current_mram_block_addr_query += BLOCK_SIZE;
			dot_product(cache_TS, cache_TS_aux, cache_query, cache_dotprods);
		}


		mram_read((__mram_ptr void const *) current_mram_block_addr_TSMean, cache_TSMean, BLOCK_SIZE);
		mram_read((__mram_ptr void const *) current_mram_block_addr_TSSigma, cache_TSSigma, BLOCK_SIZE);
		current_mram_block_addr_TSMean  += BLOCK_SIZE;
		current_mram_block_addr_TSSigma += BLOCK_SIZE;

		for (uint32_t k = 0; k < (BLOCK_SIZE / sizeof(DTYPE)); k++)
		{
			distance = 2 * ((DTYPE) query_length - (cache_dotprods[k] - (DTYPE) query_length * cache_TSMean[k]
						* query_mean) / (cache_TSSigma[k] * query_std));

			if(distance < min_distance)
			{
				min_distance =  distance;
				min_index    =  i + k;
			}
		}
	}

	// Save the result
	result->minValue = min_distance;
	result->minIndex = min_index;

	return 0;
}
