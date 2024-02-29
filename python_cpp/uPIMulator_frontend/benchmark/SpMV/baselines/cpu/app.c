
#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdint.h>

#include <omp.h>

#include "../../support/matrix.h"
#include "../../support/params.h"
#include "../../support/timer.h"
#include "../../support/utils.h"

int main(int argc, char** argv) {

    // Process parameters
    struct Params p = input_params(argc, argv);

    // Initialize SpMV data structures
    PRINT_INFO(p.verbosity >= 1, "Reading matrix %s", p.fileName);
    struct COOMatrix cooMatrix = readCOOMatrix(p.fileName);
    PRINT_INFO(p.verbosity >= 1, "    %u rows, %u columns, %u nonzeros", cooMatrix.numRows, cooMatrix.numCols, cooMatrix.numNonzeros);
    struct CSRMatrix csrMatrix = coo2csr(cooMatrix);
    float* inVector = malloc(csrMatrix.numCols*sizeof(float));
    float* outVector = malloc(csrMatrix.numRows*sizeof(float));
    initVector(inVector, csrMatrix.numCols);

    // Calculating result on CPU
    PRINT_INFO(p.verbosity >= 1, "Calculating result on CPU");
    omp_set_num_threads(4);
    Timer timer;
    startTimer(&timer);
    #pragma omp parallel for
    for(uint32_t rowIdx = 0; rowIdx < csrMatrix.numRows; ++rowIdx) {
        float sum = 0.0f;
        for(uint32_t i = csrMatrix.rowPtrs[rowIdx]; i < csrMatrix.rowPtrs[rowIdx + 1]; ++i) {
            uint32_t colIdx = csrMatrix.nonzeros[i].col;
            float value = csrMatrix.nonzeros[i].value;
            sum += inVector[colIdx]*value;
        }
        outVector[rowIdx] = sum;
    }
    stopTimer(&timer);
    if(p.verbosity == 0) PRINT("%f", getElapsedTime(timer)*1e3);
    PRINT_INFO(p.verbosity >= 1, "    Elapsed time: %f ms", getElapsedTime(timer)*1e3);

    // Deallocate data structures
    freeCOOMatrix(cooMatrix);
    freeCSRMatrix(csrMatrix);
    free(inVector);
    free(outVector);

    return 0;

}
