/*
* Vector addition with multiple tasklets
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

// Barrier
BARRIER_INIT(my_barrier, NR_TASKLETS);
// main_kernel
int main(void) {
    unsigned int tasklet_id = me();
    if (tasklet_id == 0){ // Initialize once the cycle counter
        mem_reset(); // Reset the heap
    }
    // Barrier
    barrier_wait(&my_barrier);
    uint64_t result = 0;
    uint64_t task_num = -1;
    mram_read(DPU_MRAM_HEAP_POINTER, &task_num, 8);
    for (uint64_t i = 0; i < task_num; i++) {
        result += task_num;
    }
    mram_write(&result, DPU_MRAM_HEAP_POINTER + 64, 8);

    return 0;
}
