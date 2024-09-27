/*
* Unique with multiple tasklets
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
uint32_t message[NR_TASKLETS];
T        message_value[NR_TASKLETS];
uint32_t message_offset[NR_TASKLETS];
uint32_t message_partial_count;
T        message_last_from_last;

// UNI in each tasklet
unsigned int __attribute__ ((noinline)) unique(T *output, T *input){
    unsigned int pos = 0;
    output[pos] = input[pos];
    pos++;
    #pragma unroll
    for(unsigned int j = 1; j < REGS; j++) {
        if(input[j] != input[j - 1]) {
            output[pos] = input[j];
            pos++;
        }
    }
    return pos;
}

// Handshake with adjacent tasklets
uint3 __attribute__ ((noinline)) handshake_sync(T *output, unsigned int l_count, unsigned int tasklet_id){
    unsigned int p_count, o_count, offset;
    // Wait and read message
    if(tasklet_id != 0){
        handshake_wait_for(tasklet_id - 1);
        p_count = message[tasklet_id];
        offset = (message_value[tasklet_id] == output[0])?1:0;
        o_count = message_offset[tasklet_id];
    }
    else{
        p_count = 0;
        offset = (message_last_from_last == output[0])?1:0;
        o_count = 0;
    }
    // Write message and notify
    if(tasklet_id < NR_TASKLETS - 1){
        message[tasklet_id + 1] = p_count + l_count;
        message_value[tasklet_id + 1] = output[l_count - 1];
        message_offset[tasklet_id + 1] = o_count + offset;
        handshake_notify();
    }
    uint3 result = {p_count, o_count, offset}; 
    return result;
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
    if(tasklet_id == NR_TASKLETS - 1){
        message_partial_count = 0;
        message_last_from_last = 0xFFFFFFFF; // A value that is not in the input array
    }
    // Barrier
    barrier_wait(&my_barrier);

    unsigned int i = 0; // Iteration count
    for(unsigned int byte_index = base_tasklet; byte_index < input_size_dpu_bytes; byte_index += BLOCK_SIZE * NR_TASKLETS){

        // Load cache with current MRAM block
        mram_read((__mram_ptr void const*)(mram_base_addr_A + byte_index), cache_A, BLOCK_SIZE);

        // UNI in each tasklet
        unsigned int l_count = unique(cache_B, cache_A); // In-place or out-of-place?

        // Sync with adjacent tasklets
        uint3 po_count = handshake_sync(cache_B, l_count, tasklet_id);

        // Write cache to current MRAM block
        mram_write(&cache_B[po_count.z], (__mram_ptr void*)(mram_base_addr_B + (message_partial_count + po_count.x - po_count.y) * sizeof(T)), l_count * sizeof(T));

        // First
        if(tasklet_id == 0 && i == 0){
            result->first = cache_B[0];
        }
        
        // Total count in this DPU
        if(tasklet_id == NR_TASKLETS - 1){
            message_last_from_last = cache_B[l_count - 1];
            result->last = cache_B[l_count - 1];
            result->t_count = message_partial_count + po_count.x + l_count - po_count.y - po_count.z;
            message_partial_count = result->t_count;
        }

        // Barrier
        barrier_wait(&my_barrier);

        i++;
    }

    return 0;
}
