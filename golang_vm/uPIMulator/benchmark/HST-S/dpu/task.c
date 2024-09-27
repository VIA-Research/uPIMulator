/*
* Histogram (HST-S) with multiple tasklets
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

__host dpu_arguments_t DPU_INPUT_ARGUMENTS;

// Array for communication between adjacent tasklets
uint32_t* message[NR_TASKLETS];
// DPU histogram
uint32_t* histo_dpu;

// Barrier
BARRIER_INIT(my_barrier, NR_TASKLETS);

// Histogram in each tasklet
void __attribute__ ((noinline)) histogram(uint32_t* histo, uint32_t bins, T *input, unsigned int l_size){
    for(unsigned int j = 0; j < l_size; j++) {
        T d = input[j];
        histo[(d * bins) >> DEPTH] += 1;
    }
}

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
    }
    // Barrier
    barrier_wait(&my_barrier);

    uint32_t input_size_dpu_bytes = DPU_INPUT_ARGUMENTS.size;
    uint32_t input_size_dpu_bytes_transfer = DPU_INPUT_ARGUMENTS.transfer_size; // Transfer input size per DPU in bytes
    uint32_t bins = DPU_INPUT_ARGUMENTS.bins;

    // Address of the current processing block in MRAM
    uint32_t base_tasklet = tasklet_id << BLOCK_SIZE_LOG2;
    uint32_t mram_base_addr_A = (uint32_t)DPU_MRAM_HEAP_POINTER;
    uint32_t mram_base_addr_histo = (uint32_t)(DPU_MRAM_HEAP_POINTER + input_size_dpu_bytes_transfer);

    // Initialize a local cache to store the MRAM block
    T *cache_A = (T *) mem_alloc(BLOCK_SIZE);
	
    // Local histogram
    uint32_t *histo = (uint32_t *) mem_alloc(bins * sizeof(uint32_t));

    // Initialize local histogram
    for(unsigned int i = 0; i < bins; i++){
        histo[i] = 0;
    }

    // Compute histogram
    for(unsigned int byte_index = base_tasklet; byte_index < input_size_dpu_bytes; byte_index += BLOCK_SIZE * NR_TASKLETS){

        // Bound checking
        uint32_t l_size_bytes = (byte_index + BLOCK_SIZE >= input_size_dpu_bytes) ? (input_size_dpu_bytes - byte_index) : BLOCK_SIZE;

        // Load cache with current MRAM block
        mram_read((const __mram_ptr void*)(mram_base_addr_A + byte_index), cache_A, l_size_bytes);

        // Histogram in each tasklet
        histogram(histo, bins, cache_A, l_size_bytes >> DIV);

    }
    message[tasklet_id] = histo;

    // Barrier
    barrier_wait(&my_barrier);

    uint32_t *histo_dpu = message[0];

    for (unsigned int i = tasklet_id; i < bins; i += NR_TASKLETS){
        uint32_t b = 0;		
        for (unsigned int j = 0; j < NR_TASKLETS; j++){			
            b += *(message[j] + i);
        }
        histo_dpu[i] = b;
    }

    // Barrier
    barrier_wait(&my_barrier);

    // Write dpu histogram to current MRAM block
    if(tasklet_id == 0){
        if(bins * sizeof(uint32_t) <= 2048)
            mram_write(histo_dpu, (__mram_ptr void*)(mram_base_addr_histo), bins * sizeof(uint32_t));
        else 
            for(unsigned int offset = 0; offset < ((bins * sizeof(uint32_t)) >> 11); offset++){
                mram_write(histo_dpu + (offset << 9), (__mram_ptr void*)(mram_base_addr_histo + (offset << 11)), 2048);
            }
    }

    return 0;
}
