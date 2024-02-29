
#ifndef _MATRIX_H_
#define _MATRIX_H_

#include <assert.h>
#include <stdio.h>

#include "common.h"
#include "utils.h"

struct COOMatrix {
    uint32_t numRows;
    uint32_t numCols;
    uint32_t numNonzeros;
    uint32_t* rowIdxs;
    struct Nonzero* nonzeros;
};

struct CSRMatrix {
    uint32_t numRows;
    uint32_t numCols;
    uint32_t numNonzeros;
    uint32_t* rowPtrs;
    struct Nonzero* nonzeros;
};

static struct COOMatrix readCOOMatrix(const char* fileName) {

    struct COOMatrix cooMatrix;

    // Initialize fields
    FILE* fp = fopen(fileName, "r");
    assert(fscanf(fp, "%u", &cooMatrix.numRows));
    if(cooMatrix.numRows%2 == 1) {
        PRINT_WARNING("Reading matrix %s: number of rows must be even. Padding with an extra row.", fileName);
        cooMatrix.numRows++;
    }
    assert(fscanf(fp, "%u", &cooMatrix.numCols));
    assert(fscanf(fp, "%u", &cooMatrix.numNonzeros));
    cooMatrix.rowIdxs = (uint32_t*) malloc(ROUND_UP_TO_MULTIPLE_OF_8(cooMatrix.numNonzeros*sizeof(uint32_t)));
    cooMatrix.nonzeros = (struct Nonzero*) malloc(ROUND_UP_TO_MULTIPLE_OF_8(cooMatrix.numNonzeros*sizeof(struct Nonzero)));

    // Read the nonzeros
    for(uint32_t i = 0; i < cooMatrix.numNonzeros; ++i) {
        uint32_t rowIdx;
        assert(fscanf(fp, "%u", &rowIdx));
        cooMatrix.rowIdxs[i] = rowIdx - 1; // File format indexes begin at 1
        uint32_t colIdx;
        assert(fscanf(fp, "%u", &colIdx));
        cooMatrix.nonzeros[i].col = colIdx - 1; // File format indexes begin at 1
        cooMatrix.nonzeros[i].value = 1.0f;
    }

    return cooMatrix;

}

static void freeCOOMatrix(struct COOMatrix cooMatrix) {
    free(cooMatrix.rowIdxs);
    free(cooMatrix.nonzeros);
}

static struct CSRMatrix coo2csr(struct COOMatrix cooMatrix) {

    struct CSRMatrix csrMatrix;

    // Initialize fields
    csrMatrix.numRows = cooMatrix.numRows;
    csrMatrix.numCols = cooMatrix.numCols;
    csrMatrix.numNonzeros = cooMatrix.numNonzeros;
    csrMatrix.rowPtrs = (uint32_t*) malloc(ROUND_UP_TO_MULTIPLE_OF_8((csrMatrix.numRows + 1)*sizeof(uint32_t)));
    csrMatrix.nonzeros = (struct Nonzero*) malloc(ROUND_UP_TO_MULTIPLE_OF_8(csrMatrix.numNonzeros*sizeof(struct Nonzero)));

    // Histogram rowIdxs
    memset(csrMatrix.rowPtrs, 0, (csrMatrix.numRows + 1)*sizeof(uint32_t));
    for(uint32_t i = 0; i < cooMatrix.numNonzeros; ++i) {
        uint32_t rowIdx = cooMatrix.rowIdxs[i];
        csrMatrix.rowPtrs[rowIdx]++;
    }

    // Prefix sum rowPtrs
    uint32_t sumBeforeNextRow = 0;
    for(uint32_t rowIdx = 0; rowIdx < csrMatrix.numRows; ++rowIdx) {
        uint32_t sumBeforeRow = sumBeforeNextRow;
        sumBeforeNextRow += csrMatrix.rowPtrs[rowIdx];
        csrMatrix.rowPtrs[rowIdx] = sumBeforeRow;
    }
    csrMatrix.rowPtrs[csrMatrix.numRows] = sumBeforeNextRow;

    // Bin the nonzeros
    for(uint32_t i = 0; i < cooMatrix.numNonzeros; ++i) {
        uint32_t rowIdx = cooMatrix.rowIdxs[i];
        uint32_t nnzIdx = csrMatrix.rowPtrs[rowIdx]++;
        csrMatrix.nonzeros[nnzIdx] = cooMatrix.nonzeros[i];
    }

    // Restore rowPtrs
    for(uint32_t rowIdx = csrMatrix.numRows - 1; rowIdx > 0; --rowIdx) {
        csrMatrix.rowPtrs[rowIdx] = csrMatrix.rowPtrs[rowIdx - 1];
    }
    csrMatrix.rowPtrs[0] = 0;

    return csrMatrix;

}

static void freeCSRMatrix(struct CSRMatrix csrMatrix) {
    free(csrMatrix.rowPtrs);
    free(csrMatrix.nonzeros);
}

static void initVector(float* vec, uint32_t size) {
    for(uint32_t i = 0; i < size; ++i) {
        vec[i] = 1.0f;
    }
}

#endif

