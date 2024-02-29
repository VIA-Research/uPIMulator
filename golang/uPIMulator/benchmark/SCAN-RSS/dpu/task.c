/*
* Scan with multiple tasklets (Reduce-scan-scan)
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

__host dpu_arguments_t DPU_INPUT_ARGUMENTS;
__host dpu_results_t DPU_RESULTS[NR_TASKLETS];

// Array for communication between adjacent tasklets
T message[NR_TASKLETS];
T message_partial_count;

// Reduction in each tasklet
T __attribute__ ((noinline)) reduction(T *input){
    T output = 0;
    #pragma unroll
    for(unsigned int j = 0; j < REGS; j++) {
        output += input[j];
    }
    return output;
}
// Scan in each tasklet
T __attribute__ ((noinline)) scan(T *output, T *input){
    output[0] = input[0];
    #pragma unroll
    for(unsigned int j = 1; j < REGS; j++) {
        output[j] = output[j - 1] + input[j];
    }
    return output[REGS - 1];
}
// Handshake with adjacent tasklets
T __attribute__ ((noinline)) handshake_sync(T l_count, unsigned int tasklet_id){
    T p_count;
    // Wait and read message
    if(tasklet_id != 0){
        handshake_wait_for(tasklet_id - 1);
        p_count = message[tasklet_id];
    }
    else
        p_count = 0;
    // Write message and notify
    if(tasklet_id < NR_TASKLETS - 1){
        message[tasklet_id + 1] = p_count + l_count;
        handshake_notify();
    }
    return p_count;
}

// Barrier
BARRIER_INIT(my_barrier, NR_TASKLETS);

// Add in each tasklet
void __attribute__ ((noinline)) add_(T *output, T p_count){
    #pragma unroll
    for(unsigned int j = 0; j < REGS; j++) {
        output[j] += p_count;
    }
}

extern int main_kernel1(void);
extern int main_kernel2(void);

int (*kernels[nr_kernels])(void) = {main_kernel1, main_kernel2};

int main(void) { 
    // Kernel
    return kernels[DPU_INPUT_ARGUMENTS.kernel](); 
}

// Reduction
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

    dpu_results_t *result = &DPU_RESULTS[tasklet_id];

    uint32_t input_size_dpu_bytes = DPU_INPUT_ARGUMENTS.size; // Input size per DPU in bytes

    // Address of the current processing block in MRAM
    uint32_t base_tasklet = tasklet_id << BLOCK_SIZE_LOG2;
    uint32_t mram_base_addr_A = (uint32_t)DPU_MRAM_HEAP_POINTER;

    // Initialize a local cache to store the MRAM block
    T *cache_A = (T *) mem_alloc(BLOCK_SIZE);
	
    // Local count
    T l_count = 0;

    for(unsigned int byte_index = base_tasklet; byte_index < input_size_dpu_bytes; byte_index += BLOCK_SIZE * NR_TASKLETS){

        // Load cache with current MRAM block
        mram_read((const __mram_ptr void*)(mram_base_addr_A + byte_index), cache_A, BLOCK_SIZE);

        // Reduction in each tasklet
        l_count += reduction(cache_A);

    }

    // Reduce local counts
    message[tasklet_id] = l_count;

    // Single-runtime reduction
    // Barrier
    barrier_wait(&my_barrier);
    if(tasklet_id == 0){
        for (unsigned int each_tasklet = 1; each_tasklet < NR_TASKLETS; each_tasklet++){
            message[0] += message[each_tasklet];
        }
        // Total count in this DPU
        result->t_count = message[0];
    }

    return 0;
}

// Scan
int main_kernel2() {
    unsigned int tasklet_id = me();
#if PRINT
    printf("tasklet_id = %u\n", tasklet_id);
#endif
    if (tasklet_id == 0){ // Initialize once the cycle counter
        mem_reset(); // Reset the heap
    }
    // Barrier
    barrier_wait(&my_barrier);

    dpu_results_t *result = &DPU_RESULTS[tasklet_id];

    uint32_t input_size_dpu_bytes = DPU_INPUT_ARGUMENTS.size; // Input size per DPU in bytes

    // Address of the current processing block in MRAM
    uint32_t base_tasklet = tasklet_id << BLOCK_SIZE_LOG2;
    uint32_t mram_base_addr_A = (uint32_t)DPU_MRAM_HEAP_POINTER;
    uint32_t mram_base_addr_B = (uint32_t)(DPU_MRAM_HEAP_POINTER + input_size_dpu_bytes);

    // Initialize a local cache to store the MRAM block
    T *cache_A = (T *) mem_alloc(BLOCK_SIZE);
    T *cache_B = (T *) mem_alloc(BLOCK_SIZE);
	
    // Initialize shared variable
    if(tasklet_id == NR_TASKLETS - 1)
        message_partial_count = DPU_INPUT_ARGUMENTS.t_count;
    // Barrier
    barrier_wait(&my_barrier);

    for(unsigned int byte_index = base_tasklet; byte_index < input_size_dpu_bytes; byte_index += BLOCK_SIZE * NR_TASKLETS){

        // Load cache with current MRAM block
        mram_read((const __mram_ptr void*)(mram_base_addr_A + byte_index), cache_A, BLOCK_SIZE);

        // Scan in each tasklet
        T l_count = scan(cache_B, cache_A); 

        // Sync with adjacent tasklets
        T p_count = handshake_sync(l_count, tasklet_id);

        // Barrier
        barrier_wait(&my_barrier);

        // Add in each tasklet
        add_(cache_B, message_partial_count + p_count);

        // Write cache to current MRAM block
        mram_write(cache_B, (__mram_ptr void*)(mram_base_addr_B + byte_index), BLOCK_SIZE);

        // Total count in this DPU
        if(tasklet_id == NR_TASKLETS - 1){
            result->t_count = message_partial_count + p_count + l_count;
            message_partial_count = result->t_count;
        }
	}
		
    return 0;
}
