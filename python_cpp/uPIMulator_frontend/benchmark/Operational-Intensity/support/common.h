#ifndef _COMMON_H_
#define _COMMON_H_

// Transfer size between MRAM and WRAM
#ifdef BL
#define BLOCK_SIZE_LOG2 BL
#define BLOCK_SIZE (1 << BLOCK_SIZE_LOG2)
#else
#define BLOCK_SIZE_LOG2 8
#define BLOCK_SIZE (1 << BLOCK_SIZE_LOG2)
#define BL BLOCK_SIZE_LOG2
#endif

// Data type
#ifdef UINT32
#define T uint32_t
#elif UINT64
#define T uint64_t
#elif INT32
#define T int32_t
#elif INT64
#define T int64_t
#elif FLOAT
#define T float
#elif DOUBLE
#define T double
#elif CHAR
#define T char
#elif SHORT
#define T short
#endif

// Structures used by both the host and the dpu to communicate information
typedef struct {
    uint32_t size;
    enum kernels {
        kernel1 = 0,
        nr_kernels = 1,
    } kernel;
    uint32_t repetitions;
    uint32_t stride;
    T scalar;
} dpu_arguments_t;

typedef struct {
    uint64_t cycles;
} dpu_results_t;

#define PERF 1 // Use perfcounters?
#define PRINT 0

#define ANSI_COLOR_RED     "\x1b[31m"
#define ANSI_COLOR_GREEN   "\x1b[32m"
#define ANSI_COLOR_RESET   "\x1b[0m"
#endif
