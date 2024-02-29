/*
* Reduction with multiple tasklets
*
*/
#include <stdint.h>
#include <stdio.h>
#include <defs.h>
#include <mram.h>
#include <alloc.h>
#include <perfcounter.h>
#include <handshake.h>
#include <barrier.h>

#include "../support/common.h"
#include "../support/cyclecount.h"

__host dpu_arguments_t DPU_INPUT_ARGUMENTS;
__host dpu_results_t DPU_RESULTS[NR_TASKLETS];

// Array for communication between adjacent tasklets
T message[NR_TASKLETS];

// Reduction in each tasklet
T __attribute__ ((noinline)) reduction(T *input, unsigned int l_size){
    T output = 0;
    for (unsigned int j = 0; j < l_size; j++){
        output += input[j];
    }
    return output;
}

// Barrier
BARRIER_INIT(my_barrier, NR_TASKLETS);

extern int main_kernel1(void);

int (*kernels[nr_kernels])(void) = {main_kernel1};

int main(void) { 
    // Kernel
    return kernels[DPU_INPUT_ARGUMENTS.kernel](); 
}

// main_kernel1
int main_kernel1() {
    unsigned int tasklet_id = me();
#if PRINT
    printf("tasklet_id = %u\n", tasklet_id);
#endif
    if (tasklet_id == 0){ // Initialize once the cycle counter
        mem_reset(); // Reset the heap
#if PERF
        perfcounter_config(COUNT_CYCLES, true);
#endif
    }
    // Barrier
    barrier_wait(&my_barrier);

    dpu_results_t *result = &DPU_RESULTS[tasklet_id];
#if PERF && !PERF_SYNC
    result->cycles = 0;
    perfcounter_cycles cycles;
    timer_start(&cycles); // START TIMER
#endif

    uint32_t input_size_dpu_bytes = DPU_INPUT_ARGUMENTS.size; // Input size per DPU in bytes

    // Address of the current processing block in MRAM
    uint32_t base_tasklet = tasklet_id << BLOCK_SIZE_LOG2;
    uint32_t mram_base_addr_A = (uint32_t)DPU_MRAM_HEAP_POINTER;

    // Initialize a local cache to store the MRAM block
    T *cache_A = (T *) mem_alloc(BLOCK_SIZE);
	
    // Local count
    T l_count = 0;

#if !PERF_SYNC // COMMENT OUT TO COMPARE SYNC PRIMITIVES (Experiment in Appendix)
    for(unsigned int byte_index = base_tasklet; byte_index < input_size_dpu_bytes; byte_index += BLOCK_SIZE * NR_TASKLETS){

        // Bound checking
        uint32_t l_size_bytes = (byte_index + BLOCK_SIZE >= input_size_dpu_bytes) ? (input_size_dpu_bytes - byte_index) : BLOCK_SIZE;

        // Load cache with current MRAM block
        mram_read((__mram_ptr void const*)(mram_base_addr_A + byte_index), cache_A, l_size_bytes);
		
        // Reduction in each tasklet
        l_count += reduction(cache_A, l_size_bytes >> DIV);

    }
#endif

    // Reduce local counts
    message[tasklet_id] = l_count;

#if PERF && PERF_SYNC // TIMER FOR SYNC PRIMITIVES
    result->cycles = 0;
    perfcounter_cycles cycles;
    timer_start(&cycles); // START TIMER
#endif
#ifdef TREE // Tree-based reduction
#ifdef BARRIER
    // Barrier
    barrier_wait(&my_barrier);
#endif

    #pragma unroll
    for (unsigned int offset = 1; offset < NR_TASKLETS; offset <<= 1){

        if((tasklet_id & (2*offset - 1)) == 0){
#ifndef BARRIER
            // Wait
            handshake_wait_for(tasklet_id + offset);
#endif
            message[tasklet_id] += message[tasklet_id + offset];
        }

#ifdef BARRIER
        // Barrier
        barrier_wait(&my_barrier);
#else
        else if ((tasklet_id & (offset - 1)) == 0){ // Ensure that wait and notify are in pair
            // Notify
            handshake_notify();
        }
#endif

    }

#else  // Single-thread reduction
    // Barrier
    barrier_wait(&my_barrier);
    if(tasklet_id == 0)
        #pragma unroll
        for (unsigned int each_tasklet = 1; each_tasklet < NR_TASKLETS; each_tasklet++){
            message[0] += message[each_tasklet];
        }
#endif
#if PERF && PERF_SYNC // TIMER FOR SYNC PRIMITIVES
    result->cycles = timer_stop(&cycles); // STOP TIMER
#endif

    // Total count in this DPU
    if(tasklet_id == 0){
        result->t_count = message[tasklet_id];
    }

#if PERF && !PERF_SYNC
    result->cycles = timer_stop(&cycles); // STOP TIMER
#endif

    return 0;
}
