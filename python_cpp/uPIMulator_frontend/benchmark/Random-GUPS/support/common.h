#ifndef _COMMON_H_
#define _COMMON_H_

// Structures used by both the host and the dpu to communicate information 
typedef struct {
    uint32_t size;
	enum kernels {
	    kernel1 = 0,
	    nr_kernels = 1,
	} kernel;
} dpu_arguments_t;

typedef struct {
    uint64_t cycles;
} dpu_results_t;

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
#define T uint64_t
#define S int64_t

// HPCC_starts function from HPCC benchmark (https://github.com/icl-utk-edu/hpcc/blob/main/RandomAccess/utility.c)
#define POLY 7ULL
#define PERIOD 1317624576693539401LL

T HPCC_starts(S n) {
    int i, j;
    T m2[64];
    T temp, ran;

    while (n < 0) n += PERIOD;
    while (n > PERIOD) n -= PERIOD;
    if (n == 0) return 0x1;

    temp = 0x1;
    for (i=0; i<64; i++) {
        m2[i] = temp;
        temp = (temp << 1) ^ ((S) temp < 0 ? POLY : 0);
        temp = (temp << 1) ^ ((S) temp < 0 ? POLY : 0);
    }

    for (i=62; i>=0; i--)
        if ((n >> i) & 1)
            break;

    ran = 0x2;
    while (i > 0) {
        temp = 0;
        for (j=0; j<64; j++)
            if ((ran >> j) & 1)
                temp ^= m2[j];
        ran = temp;
        i -= 1;
        if ((n >> i) & 1)
            ran = (ran << 1) ^ ((S) ran < 0 ? POLY : 0);
    }

    return ran;
}

#define PERF 1 // Use perfcounters?
#define PRINT 0

#define ANSI_COLOR_RED     "\x1b[31m"
#define ANSI_COLOR_GREEN   "\x1b[32m"
#define ANSI_COLOR_RESET   "\x1b[0m"
#endif
