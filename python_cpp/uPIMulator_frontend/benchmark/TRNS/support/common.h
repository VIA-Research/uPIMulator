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
#define T int64_t

// Structures used by both the host and the dpu to communicate information 
typedef struct {
    uint32_t m;
    uint32_t n;
    uint32_t M_;
	enum kernels {
	    kernel1 = 0,
	    kernel2 = 1,
	    nr_kernels = 2,
	} kernel;
} dpu_arguments_t;

#ifndef ENERGY
#define ENERGY 0
#endif
#define PRINT 0 

#define ANSI_COLOR_RED     "\x1b[31m"
#define ANSI_COLOR_GREEN   "\x1b[32m"
#define ANSI_COLOR_RESET   "\x1b[0m"

#define divceil(n, m) (((n)-1) / (m) + 1)
#define roundup(n, m) ((n / m) * m + m)
#endif
