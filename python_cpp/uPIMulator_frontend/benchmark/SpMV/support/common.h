
/* Common data structures between host and DPUs */

#ifndef _COMMON_H_
#define _COMMON_H_

#define ROUND_UP_TO_MULTIPLE_OF_2(x)    ((((x) + 1)/2)*2)
#define ROUND_UP_TO_MULTIPLE_OF_8(x)    ((((x) + 7)/8)*8)

struct DPUParams {
    uint32_t dpuNumRows; /* Number of rows assigned to the DPU */
    uint32_t dpuRowPtrsOffset; /* Offset of the row pointers */
    uint32_t dpuRowPtrs_m;
    uint32_t dpuNonzeros_m;
    uint32_t dpuInVector_m;
    uint32_t dpuOutVector_m;
};

struct Nonzero {
    uint32_t col;
    float value;
};

#endif

