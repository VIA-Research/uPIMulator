/*
* Strided access with multiple tasklets
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
    timer_start(&cycles); // START TIMER	
    dpu_results_t *result = &DPU_RESULTS[tasklet_id];
    result->cycles = 0;

    uint32_t input_size_dpu = DPU_INPUT_ARGUMENTS.size / sizeof(T);
    uint32_t s = DPU_INPUT_ARGUMENTS.stride;
	
    // Address of the current processing block in MRAM
    uint32_t mram_base_addr_A = (uint32_t)(DPU_MRAM_HEAP_POINTER + (tasklet_id * (input_size_dpu * sizeof(T) / NR_TASKLETS)));
    uint32_t mram_base_addr_B = (uint32_t)(DPU_MRAM_HEAP_POINTER + (tasklet_id * (input_size_dpu * sizeof(T) / NR_TASKLETS)) + input_size_dpu * sizeof(T));

#ifdef COARSECOARSE	
    // BLOCK SIZE
    uint32_t B_SIZE = BLOCK_SIZE / sizeof(T);
    uint32_t ADDR = (input_size_dpu/NR_TASKLETS) * tasklet_id;
    uint32_t j = 0;

    // Initialize a local cache to store the MRAM block
    T *cache_A = (T *) mem_alloc(BLOCK_SIZE);
    T *cache_B = (T *) mem_alloc(BLOCK_SIZE);

    for(unsigned int byte_index = 0; byte_index < input_size_dpu * sizeof(T) / NR_TASKLETS; byte_index += BLOCK_SIZE){

        // Load cache with current MRAM block
        mram_read((__mram_ptr void const*)(mram_base_addr_A + byte_index), cache_A, BLOCK_SIZE);
        mram_read((__mram_ptr void const*)(mram_base_addr_B + byte_index), cache_B, BLOCK_SIZE);

        // Copy
        if(((ADDR + j * B_SIZE) & (s - 1)) == 0){

            for(unsigned int i = 0; i < B_SIZE; i += s){
                cache_B[i] = cache_A[i];
            }

        }

        // Write cache to current MRAM block
        mram_write(cache_B, (__mram_ptr void*)(mram_base_addr_B + byte_index), BLOCK_SIZE);
        j++;
    }
#else // FINEFINE
    // Initialize a local cache to store the MRAM block
    T *cache_A = (T *) mem_alloc(sizeof(T));
    uint32_t stride = (uint32_t)(s * sizeof(T));

    for(unsigned int byte_index = 0; byte_index < input_size_dpu * sizeof(T)  / NR_TASKLETS; byte_index += stride){

        // Load cache with current MRAM block
        mram_read((__mram_ptr void const*)(mram_base_addr_A + byte_index), cache_A, sizeof(T));

        // Write cache to current MRAM block
        mram_write(cache_A, (__mram_ptr void*)(mram_base_addr_B + byte_index), sizeof(T));
    }	
#endif

    result->cycles = timer_stop(&cycles); // STOP TIMER
	
    return 0;
}
