/*
* Random Access (GUPS) with multiple tasklets
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


T ran[128]; // Current random numbers

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
        // perfcounter_config(COUNT_CYCLES, true);
#endif
    }
    // Barrier
    barrier_wait(&my_barrier);
#if PERF
    // perfcounter_cycles cycles;
    // timer_start(&cycles); // START TIMER
    dpu_results_t *result = &DPU_RESULTS[tasklet_id];
    // result->cycles = 0;
#endif	

    uint32_t input_size_dpu = DPU_INPUT_ARGUMENTS.size / sizeof(T);

    // Number of updates to table (suggested: 16x number of table entries)
    int NUPDATE = 16 * input_size_dpu;

    for(int j = tasklet_id; j < 128; j += NR_TASKLETS){
        ran[j] = HPCC_starts((NUPDATE/128) * j);
    }
    // Barrier
    barrier_wait(&my_barrier);

    // Address of the current processing block in MRAM
    uint32_t mram_base_addr_A = (uint32_t)DPU_MRAM_HEAP_POINTER;

    // Initialize a local cache to store the MRAM block
    T *cache_A = (T *) mem_alloc(sizeof(T));
	
    for (int i = 0; i < NUPDATE/128; i++){
        for (int j = tasklet_id; j < 128; j += NR_TASKLETS){

            ran[j] = (ran[j] << 1) ^ ((S) ran[j] < 0 ? POLY : 0);

            // Table[ran[j] & (TableSize-1)] ^= ran[j]; is computed as follows (3 steps)
            // 1. Load cache 
            mram_read((__mram_ptr void const*)(mram_base_addr_A + (ran[j] & (input_size_dpu - 1)) * sizeof(T)), cache_A, sizeof(T));

            // 2. Update
            //*cache_A ^= ran[j];
            *cache_A = ran[j] & (input_size_dpu - 1);

            // 3. Write cache 
            mram_write(cache_A, (__mram_ptr void*)(mram_base_addr_A + (ran[j] & (input_size_dpu - 1)) * sizeof(T)), sizeof(T));
        }
    }

#if PERF
    // result->cycles += timer_stop(&cycles); // STOP TIMER
#endif
	
    return 0;
}
