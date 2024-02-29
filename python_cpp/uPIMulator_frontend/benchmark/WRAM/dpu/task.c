/*
* WRAM Access
*
*/
#include <stdint.h>
#include <stdio.h>
#include <defs.h>
#include <mram.h>
#include <alloc.h>
#include <perfcounter.h>
#include <barrier.h>

#include "../support/common.h"
#include "../support/cyclecount.h"

__host dpu_arguments_t DPU_INPUT_ARGUMENTS;
__host dpu_results_t DPU_RESULTS[NR_TASKLETS];

// Copy
void copy_pattern_dpu(T *bufferC, T *bufferB, uint32_t *bufferA) {

    #pragma unroll
    for (unsigned int i = 0; i < (BLOCK_SIZE >> DIV); i++){
		
        uint32_t address = bufferA[i];
        bufferC[address] = bufferB[address];

    }

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

        perfcounter_config(COUNT_CYCLES, true);
    }
    perfcounter_cycles cycles;
    // Barrier
    barrier_wait(&my_barrier);
#ifndef WRAM
    timer_start(&cycles); // START TIMER
#endif

    uint32_t input_size_dpu = DPU_INPUT_ARGUMENTS.size;

    dpu_results_t *result = &DPU_RESULTS[tasklet_id];
    result->cycles = 0;

    const uint32_t A_SIZE = (BLOCK_SIZE >> DIV) << 2;
    // Address of the current processing block in MRAM
    uint32_t mram_base_addr_A = (uint32_t)(DPU_MRAM_HEAP_POINTER + (tasklet_id * A_SIZE));
    uint32_t mram_base_addr_B = (uint32_t)(DPU_MRAM_HEAP_POINTER + (tasklet_id << BLOCK_SIZE_LOG2) + input_size_dpu * sizeof(uint32_t));
    uint32_t mram_base_addr_C = (uint32_t)(DPU_MRAM_HEAP_POINTER + (tasklet_id << BLOCK_SIZE_LOG2) + input_size_dpu * (sizeof(uint32_t) + sizeof(T)));

    // Initialize a local cache to store the MRAM block
    uint32_t *cache_A = (uint32_t *) mem_alloc(A_SIZE);
    T *cache_B = (T *) mem_alloc(BLOCK_SIZE);
    T *cache_C = (T *) mem_alloc(BLOCK_SIZE);

    uint32_t A_byte_index = 0; 
    for(unsigned int byte_index = 0; byte_index < (input_size_dpu << DIV); byte_index += BLOCK_SIZE * NR_TASKLETS){

        // Load cache with current MRAM block
        mram_read((__mram_ptr void const*)(mram_base_addr_A + A_byte_index), cache_A, A_SIZE);
        mram_read((__mram_ptr void const*)(mram_base_addr_B + byte_index), cache_B, BLOCK_SIZE);
        mram_read((__mram_ptr void const*)(mram_base_addr_C + byte_index), cache_C, BLOCK_SIZE); // Clean cache_C

#ifdef WRAM
        // Barrier
        barrier_wait(&my_barrier);
        timer_start(&cycles); // START TIMER
#endif

        // Copy
        copy_pattern_dpu(cache_C, cache_B, cache_A);

#ifdef WRAM
        result->cycles += timer_stop(&cycles); // STOP TIMER
        // Barrier
        barrier_wait(&my_barrier);
#endif

        // Write cache to current MRAM block
        mram_write(cache_C, (__mram_ptr void*)(mram_base_addr_C + byte_index), BLOCK_SIZE);

        A_byte_index += A_SIZE * NR_TASKLETS;
    }

#ifndef WRAM
    result->cycles = timer_stop(&cycles); // STOP TIMER
#endif
    return 0;
}
