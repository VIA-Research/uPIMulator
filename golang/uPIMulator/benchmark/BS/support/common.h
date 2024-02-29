#ifndef _COMMON_H_
#define _COMMON_H_

#ifdef TL
#define TASKLETS_INITIALIZER TASKLETS(TL, main, 2048, 2)
#define NB_OF_TASKLETS_PER_DPU TL
#else
#define TASKLETS_INITIALIZER TASKLETS(16, main, 2048, 2)
#define NB_OF_TASKLETS_PER_DPU 16
#endif

// Transfer size between MRAM and WRAM
#ifdef BL
#define BLOCK_SIZE_LOG2 BL
#define BLOCK_SIZE (1 << BLOCK_SIZE_LOG2)
#else
#define BLOCK_SIZE_LOG2 8
#define BLOCK_SIZE (1 << BLOCK_SIZE_LOG2)
#endif

// Data type
#define DTYPE int64_t

// Vector size
#define INPUT_SIZE 2048576

typedef struct {
	uint64_t input_size;
	uint64_t slice_per_dpu;
	enum kernels {
		kernel1 = 0,
		nr_kernels = 1,
	} kernel;
} dpu_arguments_t;

// Structures used by both the host and the dpu to communicate information
typedef struct {
    DTYPE found;
} dpu_results_t;

#ifndef ENERGY
#define ENERGY 0
#endif
#define PRINT 0

#define ANSI_COLOR_RED     "\x1b[31m"
#define ANSI_COLOR_GREEN   "\x1b[32m"
#define ANSI_COLOR_RESET   "\x1b[0m"
#endif
