/*
* SpMV with multiple tasklets
*
*/
#include <stdio.h>

#include <alloc.h>
#include <barrier.h>
#include <defs.h>
#include <mram.h>
#include <perfcounter.h>
#include <seqread.h>

#include "../support/common.h"

#define PRINT_ERROR(fmt, ...) printf("\033[0;31mERROR:\033[0m   "fmt"\n", ##__VA_ARGS__)

#define MIN(x, y)   (((x) < (y))?(x):(y))

BARRIER_INIT(my_barrier, NR_TASKLETS);

// main
int main() {

    if(me() == 0) {
        mem_reset(); // Reset the heap
    }
    // Barrier
    barrier_wait(&my_barrier);

    // Load parameters
    uint32_t params_m = (uint32_t) DPU_MRAM_HEAP_POINTER;
    struct DPUParams* params_w = (struct DPUParams*) mem_alloc(ROUND_UP_TO_MULTIPLE_OF_8(sizeof(struct DPUParams)));
    mram_read((__mram_ptr void const*)params_m, params_w, ROUND_UP_TO_MULTIPLE_OF_8(sizeof(struct DPUParams)));
    uint32_t numRows = params_w->dpuNumRows;

    // Sanity check
    if(me() == 0) {
        if(numRows%2 != 0) {
            // The number of rows assigned to the DPU must be a multiple of two to ensure that writes to the output vector are aligned to 8 bytes
            PRINT_ERROR("The number of rows is not a multiple of two!");
        }
    }

    // Identify tasklet's rows
    uint32_t numRowsPerTasklet = ROUND_UP_TO_MULTIPLE_OF_2((numRows - 1)/NR_TASKLETS + 1); // Multiple of two to ensure that access to rowPtrs and outVector is 8-byte aligned
    uint32_t taskletRowsStart = me()*numRowsPerTasklet;
    uint32_t taskletNumRows;
    if(taskletRowsStart > numRows) {
        taskletNumRows = 0;
    } else if(taskletRowsStart + numRowsPerTasklet > numRows) {
        taskletNumRows = numRows - taskletRowsStart;
    } else {
        taskletNumRows = numRowsPerTasklet;
    }

    // Only process tasklets with nonzero number of rows
    if(taskletNumRows > 0) {

        // Extract parameters
        uint32_t rowPtrsOffset = params_w->dpuRowPtrsOffset;
        uint32_t rowPtrs_m = ((uint32_t)DPU_MRAM_HEAP_POINTER) + params_w->dpuRowPtrs_m;
        uint32_t nonzeros_m = ((uint32_t)DPU_MRAM_HEAP_POINTER) + params_w->dpuNonzeros_m;
        uint32_t inVector_m = ((uint32_t)DPU_MRAM_HEAP_POINTER) + params_w->dpuInVector_m;
        uint32_t outVector_m = ((uint32_t)DPU_MRAM_HEAP_POINTER) + params_w->dpuOutVector_m;

        // Initialize row pointer sequential reader
        uint32_t taskletRowPtrs_m = rowPtrs_m + taskletRowsStart*sizeof(uint32_t);
        seqreader_t rowPtrReader;
        uint32_t* taskletRowPtrs_w = seqread_init(seqread_alloc(), (__mram_ptr void*)taskletRowPtrs_m, &rowPtrReader);
        uint32_t firstRowPtr = *taskletRowPtrs_w;

        // Initialize nonzeros sequential reader
        uint32_t taskletNonzerosStart = firstRowPtr - rowPtrsOffset;
        uint32_t taskletNonzeros_m = nonzeros_m + taskletNonzerosStart*sizeof(struct Nonzero); // 8-byte aligned because Nonzero is 8 bytes
        seqreader_t nonzerosReader;
        struct Nonzero* taskletNonzeros_w = seqread_init(seqread_alloc(), (__mram_ptr void*)taskletNonzeros_m, &nonzerosReader);

        // Initialize input vector cache
        uint32_t inVectorTileSize = 64;
        float* inVectorTile_w = mem_alloc(inVectorTileSize*sizeof(float));
        mram_read((__mram_ptr void const*)inVector_m, inVectorTile_w, 256);
        uint32_t currInVectorTileIdx = 0;

        // Initialize output vector cache
        uint32_t taskletOutVector_m = outVector_m + taskletRowsStart*sizeof(float);
        uint32_t outVectorTileSize = 64;
        float* outVectorTile_w = mem_alloc(outVectorTileSize*sizeof(float));

        // SpMV
        uint32_t nextRowPtr = firstRowPtr;
        for(uint32_t row = 0; row < taskletNumRows; ++row) {

            // Find row nonzeros
            taskletRowPtrs_w = seqread_get(taskletRowPtrs_w, sizeof(uint32_t), &rowPtrReader);
            uint32_t rowPtr = nextRowPtr;
            nextRowPtr = *taskletRowPtrs_w;
            uint32_t taskletNNZ = nextRowPtr - rowPtr;

            // Multiply row with vector
            float outValue = 0.0f;
            for(uint32_t nzIdx = 0; nzIdx < taskletNNZ; ++nzIdx) {

                // Get matrix value
                float matValue = taskletNonzeros_w->value;

                // Get input vector value
                uint32_t col = taskletNonzeros_w->col;
                uint32_t inVectorTileIdx = col/inVectorTileSize;
                uint32_t inVectorTileOffset = col%inVectorTileSize;
                if(inVectorTileIdx != currInVectorTileIdx) {
                    mram_read((__mram_ptr void const*)(inVector_m + inVectorTileIdx*inVectorTileSize*sizeof(float)), inVectorTile_w, 256);
                    currInVectorTileIdx = inVectorTileIdx;
                }
                float inValue = inVectorTile_w[inVectorTileOffset];

                // Multiply and add
                outValue += matValue*inValue;

                // Read next nonzero
                taskletNonzeros_w = seqread_get(taskletNonzeros_w, sizeof(struct Nonzero), &nonzerosReader); // Last read will be out of bounds and unused

            }

            // Store output
            uint32_t outVectorTileIdx = row/outVectorTileSize;
            uint32_t outVectorTileOffset = row%outVectorTileSize;
            outVectorTile_w[outVectorTileOffset] = outValue;
            if(outVectorTileOffset == outVectorTileSize - 1) { // Last element in tile
                mram_write(outVectorTile_w, (__mram_ptr void*)(taskletOutVector_m + outVectorTileIdx*outVectorTileSize*sizeof(float)), 256);
            } else if(row == taskletNumRows - 1) { // Last row for tasklet
                mram_write(outVectorTile_w, (__mram_ptr void*)(taskletOutVector_m + outVectorTileIdx*outVectorTileSize*sizeof(float)), (taskletNumRows%outVectorTileSize)*sizeof(float));
            }

        }
    }

    return 0;
}
