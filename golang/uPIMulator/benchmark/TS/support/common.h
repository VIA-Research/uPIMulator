#ifndef _COMMON_H_
#define _COMMON_H_

// Transfer size between MRAM and WRAM
#ifdef BL
#define BLOCK_SIZE_LOG2 BL
#define BLOCK_SIZE (1 << BLOCK_SIZE_LOG2)
#else
#define BLOCK_SIZE_LOG2 8
#define BLOCK_SIZE (1 << BLOCK_SIZE_LOG2)
#endif

// Data type
#define DTYPE int32_t
#define DTYPE_MAX INT32_MAX

typedef struct  {
	uint32_t ts_length;
    uint32_t query_length;
    DTYPE query_mean;
    DTYPE query_std;
    uint32_t slice_per_dpu;
    int32_t exclusion_zone;
    enum kernels {
		kernel1 = 0,
		nr_kernels = 1,
	} kernel;
}dpu_arguments_t;

typedef struct  {
    DTYPE minValue;
    uint32_t minIndex;
    DTYPE maxValue;
    uint32_t maxIndex;
}dpu_result_t;

#ifndef ENERGY
#define ENERGY 0
#endif
#define PRINT 0 

#define ANSI_COLOR_RED     "\x1b[31m"
#define ANSI_COLOR_GREEN   "\x1b[32m"
#define ANSI_COLOR_RESET   "\x1b[0m"
#endif
