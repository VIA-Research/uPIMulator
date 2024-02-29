
#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

typedef unsigned int uint32_t;

struct COOMatrix {
    uint32_t numRows;
    uint32_t numCols;
    uint32_t numNonzeros;
    uint32_t* rowIdxs;
    uint32_t* colIdxs;
};

static struct COOMatrix readCOOMatrix(const char* fileName) {

    struct COOMatrix cooMatrix;

    // Initialize fields
    FILE* fp = fopen(fileName, "r");
    assert(fscanf(fp, "%u", &cooMatrix.numRows));
    assert(fscanf(fp, "%u", &cooMatrix.numCols));
    assert(fscanf(fp, "%u", &cooMatrix.numNonzeros));
    cooMatrix.rowIdxs = (uint32_t*) malloc(cooMatrix.numNonzeros*sizeof(uint32_t));
    cooMatrix.colIdxs = (uint32_t*) malloc(cooMatrix.numNonzeros*sizeof(uint32_t));

    // Read the nonzeros
    for(uint32_t i = 0; i < cooMatrix.numNonzeros; ++i) {
        uint32_t rowIdx;
        assert(fscanf(fp, "%u", &rowIdx));
        cooMatrix.rowIdxs[i] = rowIdx - 1; // File format indexes begin at 1
        uint32_t colIdx;
        assert(fscanf(fp, "%u", &colIdx));
        cooMatrix.colIdxs[i] = colIdx - 1; // File format indexes begin at 1
    }

    fclose(fp);

    return cooMatrix;

}

static void freeCOOMatrix(struct COOMatrix cooMatrix) {
    free(cooMatrix.rowIdxs);
    free(cooMatrix.colIdxs);
}

int main(int argc, char** argv) {

    const char* fileName = (argc > 1)?argv[1]:"bcsstk30.mtx";
    unsigned int replicationFactor = (argc > 2)?atoi(argv[2]):4;
    const char* outFileName = (argc > 3)?argv[3]:"out.mtx";

    struct COOMatrix cooMatrix = readCOOMatrix(fileName);

    FILE* fp = fopen(outFileName, "w");
    fprintf(fp, "%u  %u  %u\n", cooMatrix.numRows*replicationFactor, cooMatrix.numCols, cooMatrix.numNonzeros*replicationFactor);
    for(unsigned int i = 0; i < cooMatrix.numNonzeros; ++i) {
        unsigned int row = cooMatrix.rowIdxs[i];
        unsigned int col = cooMatrix.colIdxs[i];
        for(unsigned int r = 0; r < replicationFactor; ++r) {
            fprintf(fp, "%u %u\n", row + 1 + r*cooMatrix.numRows, col + 1);
        }
    }
    fclose(fp);

    freeCOOMatrix(cooMatrix);

    return 0;

}
