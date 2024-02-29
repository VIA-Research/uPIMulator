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
#define DIV 2 // Shift right to divide by sizeof(T)
#elif UINT64
#define T uint64_t
#define DIV 3 // Shift right to divide by sizeof(T)
#elif INT32
#define T int32_t
#define DIV 2 // Shift right to divide by sizeof(T)
#elif INT64
#define T int64_t
#define DIV 3 // Shift right to divide by sizeof(T)
#elif FLOAT
#define T float
#define DIV 2 // Shift right to divide by sizeof(T)
#elif DOUBLE
#define T double
#define DIV 3 // Shift right to divide by sizeof(T)
#elif CHAR
#define T char
#define DIV 0 // Shift right to divide by sizeof(T)
#elif SHORT
#define T short
#define DIV 1 // Shift right to divide by sizeof(T)
#endif

#define REGS (BLOCK_SIZE >> DIV)

// Structures used by both the host and the dpu to communicate information
typedef struct {
    uint32_t size;
	enum kernels {
	    kernel1 = 0,
	    kernel2 = 1,
	    nr_kernels = 2,
	} kernel;
    T t_count;
} dpu_arguments_t;

typedef struct {
    T t_count;
} dpu_results_t;

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
