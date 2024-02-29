
#ifndef _MRAM_MANAGEMENT_H_
#define _MRAM_MANAGEMENT_H_

#include "../support/common.h"
#include "../support/utils.h"

#define DPU_CAPACITY (64 << 20) // A DPU's capacity is 64 MiB

struct mram_heap_allocator_t {
    uint32_t totalAllocated;
};

static void init_allocator(struct mram_heap_allocator_t* allocator) {
    allocator->totalAllocated = 0;
}

static uint32_t mram_heap_alloc(struct mram_heap_allocator_t* allocator, uint32_t size) {
    uint32_t ret = allocator->totalAllocated;
    allocator->totalAllocated += ROUND_UP_TO_MULTIPLE_OF_8(size);
    if(allocator->totalAllocated > DPU_CAPACITY) {
        PRINT_ERROR("        Total memory allocated is %d bytes which exceeds the DPU capacity (%d bytes)!", allocator->totalAllocated, DPU_CAPACITY);
        exit(0);
    }
    return ret;
}

static void copyToDPU(struct dpu_set_t dpu, uint8_t* hostPtr, uint32_t mramIdx, uint32_t size) {
    DPU_ASSERT(dpu_copy_to(dpu, DPU_MRAM_HEAP_POINTER_NAME, mramIdx, hostPtr, ROUND_UP_TO_MULTIPLE_OF_8(size)));
}

static void copyFromDPU(struct dpu_set_t dpu, uint32_t mramIdx, uint8_t* hostPtr, uint32_t size) {
    DPU_ASSERT(dpu_copy_from(dpu, DPU_MRAM_HEAP_POINTER_NAME, mramIdx, hostPtr, ROUND_UP_TO_MULTIPLE_OF_8(size)));
}

#endif

