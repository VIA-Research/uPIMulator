/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_MRAM_H
#define DPUSYSCORE_MRAM_H

#include <stdint.h>
#include <attributes.h>

/**
 * @file mram.h
 * @brief MRAM Transfer Management.
 */

#define DPU_MRAM_HEAP_POINTER ((__mram_ptr void *)(&__sys_used_mram_end))
extern __mram_ptr __dma_aligned uint8_t __sys_used_mram_end[0];

/**
 * @fn mram_read
 * @brief Stores the specified number of bytes from MRAM to WRAM.
 * The number of bytes must be:
 *  - at least 8
 *  - at most 2048
 *  - a multiple of 8
 *
 * @param from source address in MRAM
 * @param to destination address in WRAM
 * @param nb_of_bytes number of bytes to transfer
 */
static inline void
mram_read(const __mram_ptr void *from, void *to, unsigned int nb_of_bytes)
{
    __builtin_dpu_ldma(to, from, nb_of_bytes);
}

/**
 * @fn mram_write
 * @brief Stores the specified number of bytes from WRAM to MRAM.
 * The number of bytes must be:
 *  - at least 8
 *  - at most 2048
 *  - a multiple of 8
 *
 * @param from source address in WRAM
 * @param to destination address in MRAM
 * @param nb_of_bytes number of bytes to transfer
 */
static inline void
mram_write(const void *from, __mram_ptr void *to, unsigned int nb_of_bytes)
{
    __builtin_dpu_sdma(from, to, nb_of_bytes);
}

#endif /* DPUSYSCORE_MRAM_H */
