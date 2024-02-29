
#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdint.h>

#include "../../support/matrix.h"
#include "../../support/params.h"
#include "../../support/timer.h"
#include "../../support/utils.h"

__global__ void spmv_kernel(CSRMatrix csrMatrix, float* inVector, float* outVector) {
    unsigned int row = blockIdx.x*blockDim.x + threadIdx.x;
    if(row < csrMatrix.numRows) {
        float sum = 0.0f;
        for(unsigned int i = csrMatrix.rowPtrs[row]; i < csrMatrix.rowPtrs[row + 1]; ++i) {
            struct Nonzero nonzero = csrMatrix.nonzeros[i];
            sum += inVector[nonzero.col]*nonzero.value;
        }
        outVector[row] = sum;
    }
}

int main(int argc, char** argv) {

    // Process parameters
    struct Params p = input_params(argc, argv);

    // Initialize SpMV data structures
    PRINT_INFO(p.verbosity >= 1, "Reading matrix %s", p.fileName);
    struct COOMatrix cooMatrix = readCOOMatrix(p.fileName);
    PRINT_INFO(p.verbosity >= 1, "    %u rows, %u columns, %u nonzeros", cooMatrix.numRows, cooMatrix.numCols, cooMatrix.numNonzeros);
    struct CSRMatrix csrMatrix = coo2csr(cooMatrix);
    float* inVector = (float*) malloc(csrMatrix.numCols*sizeof(float));
    float* outVector = (float*) malloc(csrMatrix.numRows*sizeof(float));
    initVector(inVector, csrMatrix.numCols);

    // Allocate data structures on GPU
    CSRMatrix csrMatrix_d;
    csrMatrix_d.numRows = csrMatrix.numRows;
    csrMatrix_d.numCols = csrMatrix.numCols;
    csrMatrix_d.numNonzeros = csrMatrix.numNonzeros;
    cudaMalloc((void**) &csrMatrix_d.rowPtrs, (csrMatrix_d.numRows + 1)*sizeof(unsigned int));
    cudaMalloc((void**) &csrMatrix_d.nonzeros, csrMatrix_d.numNonzeros*sizeof(struct Nonzero));
    float* inVector_d;
    cudaMalloc((void**) &inVector_d, csrMatrix_d.numCols*sizeof(float));
    float* outVector_d;
    cudaMalloc((void**) &outVector_d, csrMatrix_d.numRows*sizeof(float));

    // Copy data to GPU
    cudaMemcpy(csrMatrix_d.rowPtrs, csrMatrix.rowPtrs, (csrMatrix_d.numRows + 1)*sizeof(unsigned int), cudaMemcpyHostToDevice);
    cudaMemcpy(csrMatrix_d.nonzeros, csrMatrix.nonzeros, csrMatrix_d.numNonzeros*sizeof(struct Nonzero), cudaMemcpyHostToDevice);
    cudaMemcpy(inVector_d, inVector, csrMatrix_d.numCols*sizeof(float), cudaMemcpyHostToDevice);
    cudaDeviceSynchronize();

    // Calculating result on GPU
    PRINT_INFO(p.verbosity >= 1, "Calculating result on GPU");
    Timer timer;
    startTimer(&timer);
    unsigned int numThreadsPerBlock = 1024;
    unsigned int numBlocks = (csrMatrix_d.numRows + numThreadsPerBlock - 1)/numThreadsPerBlock;
    spmv_kernel <<< numBlocks, numThreadsPerBlock >>> (csrMatrix_d, inVector_d, outVector_d);
    cudaDeviceSynchronize();
    stopTimer(&timer);
    if(p.verbosity == 0) PRINT("%f", getElapsedTime(timer)*1e3);
    PRINT_INFO(p.verbosity >= 1, "    Elapsed time: %f ms", getElapsedTime(timer)*1e3);

    // Copy data from GPU
    cudaMemcpy(outVector, outVector_d, csrMatrix_d.numRows*sizeof(float), cudaMemcpyDeviceToHost);
    cudaDeviceSynchronize();

    // Calculating result on CPU
    PRINT_INFO(p.verbosity >= 1, "Calculating result on CPU");
    float* outVectorReference = (float*) malloc(csrMatrix.numRows*sizeof(float));
    for(uint32_t rowIdx = 0; rowIdx < csrMatrix.numRows; ++rowIdx) {
        float sum = 0.0f;
        for(uint32_t i = csrMatrix.rowPtrs[rowIdx]; i < csrMatrix.rowPtrs[rowIdx + 1]; ++i) {
            uint32_t colIdx = csrMatrix.nonzeros[i].col;
            float value = csrMatrix.nonzeros[i].value;
            sum += inVector[colIdx]*value;
        }
        outVectorReference[rowIdx] = sum;
    }

    // Verify the result
    PRINT_INFO(p.verbosity >= 1, "Verifying the result");
    for(uint32_t rowIdx = 0; rowIdx < csrMatrix.numRows; ++rowIdx) {
        float diff = (outVectorReference[rowIdx] - outVector[rowIdx])/outVectorReference[rowIdx];
        const float tolerance = 0.00001;
        if(diff > tolerance || diff < -tolerance) {
            PRINT_ERROR("Mismatch at index %u (CPU result = %f, DPU result = %f)", rowIdx, outVectorReference[rowIdx], outVector[rowIdx]);
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

