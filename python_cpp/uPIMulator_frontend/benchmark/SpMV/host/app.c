/**
* app.c
* SpMV Host Application Source File
*
*/
#include <dpu.h>
#include <dpu_log.h>

#include <assert.h>
#include <getopt.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "mram-management.h"
#include "../support/common.h"
#include "../support/matrix.h"
#include "../support/params.h"
#include "../support/timer.h"
#include "../support/utils.h"

#define DPU_BINARY "./bin/dpu_code"

#ifndef ENERGY
#define ENERGY 0
#endif
#if ENERGY
#include <dpu_probe.h>
#endif

// Main of the Host Application
int main(int argc, char** argv) {

    // Process parameters
    struct Params p = input_params(argc, argv);

    // Timing and profiling
    Timer timer;
    float loadTime = 0.0f, dpuTime = 0.0f, retrieveTime = 0.0f;
    #if ENERGY
    struct dpu_probe_t probe;
    DPU_ASSERT(dpu_probe_init("energy_probe", &probe));
    #endif

    // Allocate DPUs and load binary
    struct dpu_set_t dpu_set, dpu;
    uint32_t numDPUs;
    DPU_ASSERT(dpu_alloc(NR_DPUS, NULL, &dpu_set));
    DPU_ASSERT(dpu_load(dpu_set, DPU_BINARY, NULL));
    DPU_ASSERT(dpu_get_nr_dpus(dpu_set, &numDPUs));
    PRINT_INFO(p.verbosity >= 1, "Allocated %d DPU(s)", numDPUs);

    // Initialize SpMV data structures
    PRINT_INFO(p.verbosity >= 1, "Reading matrix %s", p.fileName);
    struct COOMatrix cooMatrix = readCOOMatrix(p.fileName);
    PRINT_INFO(p.verbosity >= 1, "    %u rows, %u columns, %u nonzeros", cooMatrix.numRows, cooMatrix.numCols, cooMatrix.numNonzeros);
    struct CSRMatrix csrMatrix = coo2csr(cooMatrix);
    uint32_t numRows = csrMatrix.numRows;
    uint32_t numCols = csrMatrix.numCols;
    uint32_t* rowPtrs = csrMatrix.rowPtrs;
    struct Nonzero* nonzeros = csrMatrix.nonzeros;
    float* inVector = malloc(ROUND_UP_TO_MULTIPLE_OF_8(numCols*sizeof(float)));
    initVector(inVector, numCols);
    float* outVector = malloc(ROUND_UP_TO_MULTIPLE_OF_8(numRows*sizeof(float)));

    // Partition data structure across DPUs
    uint32_t numRowsPerDPU = ROUND_UP_TO_MULTIPLE_OF_2((numRows - 1)/numDPUs + 1);
    PRINT_INFO(p.verbosity >= 1, "Assigning %u rows per DPU", numRowsPerDPU);
    struct DPUParams dpuParams[numDPUs];
    unsigned int dpuIdx = 0;
    PRINT_INFO(p.verbosity == 1, "Copying data to DPUs");
    DPU_FOREACH (dpu_set, dpu) {

        // Allocate parameters
        struct mram_heap_allocator_t allocator;
        init_allocator(&allocator);
        uint32_t dpuParams_m = mram_heap_alloc(&allocator, sizeof(struct DPUParams));

        // Find DPU's rows
        uint32_t dpuStartRowIdx = dpuIdx*numRowsPerDPU;
        uint32_t dpuNumRows;
        if(dpuStartRowIdx > numRows) {
            dpuNumRows = 0;
        } else if(dpuStartRowIdx + numRowsPerDPU > numRows) {
            dpuNumRows = numRows - dpuStartRowIdx;
        } else {
            dpuNumRows = numRowsPerDPU;
        }
        dpuParams[dpuIdx].dpuNumRows = dpuNumRows;
        PRINT_INFO(p.verbosity >= 2, "    DPU %u:", dpuIdx);
        PRINT_INFO(p.verbosity >= 2, "        Receives %u rows", dpuNumRows);

        // Partition nonzeros and copy data
        if(dpuNumRows > 0) {

            // Find DPU's CSR matrix partition
            uint32_t* dpuRowPtrs_h = &rowPtrs[dpuStartRowIdx];
            uint32_t dpuRowPtrsOffset = dpuRowPtrs_h[0];
            struct Nonzero* dpuNonzeros_h = &nonzeros[dpuRowPtrsOffset];
            uint32_t dpuNumNonzeros = dpuRowPtrs_h[dpuNumRows] - dpuRowPtrsOffset;

            // Allocate MRAM
            uint32_t dpuRowPtrs_m = mram_heap_alloc(&allocator, (dpuNumRows + 1)*sizeof(uint32_t));
            uint32_t dpuNonzeros_m = mram_heap_alloc(&allocator, dpuNumNonzeros*sizeof(struct Nonzero));
            uint32_t dpuInVector_m = mram_heap_alloc(&allocator, numCols*sizeof(float));
            uint32_t dpuOutVector_m = mram_heap_alloc(&allocator, dpuNumRows*sizeof(float));
            assert((dpuNumRows*sizeof(float))%8 == 0 && "Output sub-vector must be a multiple of 8 bytes!");
            PRINT_INFO(p.verbosity >= 2, "        Total memory allocated is %d bytes", allocator.totalAllocated);

            // Set up DPU parameters
            dpuParams[dpuIdx].dpuRowPtrsOffset = dpuRowPtrsOffset;
            dpuParams[dpuIdx].dpuRowPtrs_m = dpuRowPtrs_m;
            dpuParams[dpuIdx].dpuNonzeros_m = dpuNonzeros_m;
            dpuParams[dpuIdx].dpuInVector_m = dpuInVector_m;
            dpuParams[dpuIdx].dpuOutVector_m = dpuOutVector_m;

            // Send data to DPU
            PRINT_INFO(p.verbosity >= 2, "        Copying data to DPU");
            startTimer(&timer);
            copyToDPU(dpu, (uint8_t*)dpuRowPtrs_h, dpuRowPtrs_m, (dpuNumRows + 1)*sizeof(uint32_t));
            copyToDPU(dpu, (uint8_t*)dpuNonzeros_h, dpuNonzeros_m, dpuNumNonzeros*sizeof(struct Nonzero));
            copyToDPU(dpu, (uint8_t*)inVector, dpuInVector_m, numCols*sizeof(float));
            stopTimer(&timer);
            loadTime += getElapsedTime(timer);

        }

        // Send parameters to DPU
        PRINT_INFO(p.verbosity >= 2, "        Copying parameters to DPU");
        startTimer(&timer);
        copyToDPU(dpu, (uint8_t*)&dpuParams[dpuIdx], dpuParams_m, sizeof(struct DPUParams));
        stopTimer(&timer);
        loadTime += getElapsedTime(timer);

        ++dpuIdx;

    }
    PRINT_INFO(p.verbosity >= 1, "    CPU-DPU Time: %f ms", loadTime*1e3);

    // Run all DPUs
    PRINT_INFO(p.verbosity >= 1, "Booting DPUs");
    startTimer(&timer);
    #if ENERGY
    DPU_ASSERT(dpu_probe_start(&probe));
    #endif
    DPU_ASSERT(dpu_launch(dpu_set, DPU_SYNCHRONOUS));
    #if ENERGY
    DPU_ASSERT(dpu_probe_stop(&probe));
    double energy;
    DPU_ASSERT(dpu_probe_get(&probe, DPU_ENERGY, DPU_AVERAGE, &energy));
    PRINT_INFO(p.verbosity >= 1, "    DPU Energy: %f J", energy);
    #endif
    stopTimer(&timer);
    dpuTime += getElapsedTime(timer);
    PRINT_INFO(p.verbosity >= 1, "    DPU Time: %f ms", dpuTime*1e3);

    // Copy back result
    PRINT_INFO(p.verbosity >= 1, "Copying back the result");
    startTimer(&timer);
    dpuIdx = 0;
    DPU_FOREACH (dpu_set, dpu) {
        unsigned int dpuNumRows = dpuParams[dpuIdx].dpuNumRows;
        if(dpuNumRows > 0) {
            uint32_t dpuStartRowIdx = dpuIdx*numRowsPerDPU;
            copyFromDPU(dpu, dpuParams[dpuIdx].dpuOutVector_m, (uint8_t*)(outVector + dpuStartRowIdx), dpuNumRows*sizeof(float));
        }
        ++dpuIdx;
    }
    stopTimer(&timer);
    retrieveTime += getElapsedTime(timer);
    PRINT_INFO(p.verbosity >= 1, "    DPU-CPU Time: %f ms", retrieveTime*1e3);
    if(p.verbosity == 0) PRINT("CPU-DPU Time(ms): %f    DPU Kernel Time (ms): %f    DPU-CPU Time (ms): %f", loadTime*1e3, dpuTime*1e3, retrieveTime*1e3);

    // Calculating result on CPU
    PRINT_INFO(p.verbosity >= 1, "Calculating result on CPU");
    float* outVectorReference = malloc(numRows*sizeof(float));
    for(uint32_t rowIdx = 0; rowIdx < numRows; ++rowIdx) {
        float sum = 0.0f;
        for(uint32_t i = rowPtrs[rowIdx]; i < rowPtrs[rowIdx + 1]; ++i) {
            uint32_t colIdx = nonzeros[i].col;
            float value = nonzeros[i].value;
            sum += inVector[colIdx]*value;
        }
        outVectorReference[rowIdx] = sum;
    }

    // Verify the result
    PRINT_INFO(p.verbosity >= 1, "Verifying the result");
    for(uint32_t rowIdx = 0; rowIdx < numRows; ++rowIdx) {
        float diff = (outVectorReference[rowIdx] - outVector[rowIdx])/outVectorReference[rowIdx];
        const float tolerance = 0.00001;
        if(diff > tolerance || diff < -tolerance) {
            PRINT_ERROR("Mismatch at index %u (CPU result = %f, DPU result = %f)", rowIdx, outVectorReference[rowIdx], outVector[rowIdx]);
        }
    }

    // Display DPU Logs
    if(p.verbosity >= 2) {
        PRINT_INFO(p.verbosity >= 2, "Displaying DPU Logs:");
        dpuIdx = 0;
        DPU_FOREACH (dpu_set, dpu) {
            PRINT("DPU %u:", dpuIdx);
            DPU_ASSERT(dpu_log_read(dpu, stdout));
            ++dpuIdx;
        }
    }

    // Deallocate data structures
    freeCOOMatrix(cooMatrix);
    freeCSRMatrix(csrMatrix);
    free(inVector);
    free(outVector);
    free(outVectorReference);

    return 0;
}
