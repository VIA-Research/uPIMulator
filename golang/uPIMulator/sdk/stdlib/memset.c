/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <memmram_utils.h>
#include <mram.h>
#include <stddef.h>
#include <stdint.h>

static void *
__memset_wram_1align(void *dest, int c, size_t len)
{
    uint8_t *dest8 = (uint8_t *)(dest);
    for (uint32_t i = 0; i < (len); ++i) {
        dest8[i] = (c);
    }
    return dest;
}

typedef uint32_t memset_wram_t;
/* Requisite:
 *  - dest: align on 4 bytes
 *  - len: mutiple of 4 bytes
 */
void *__attribute__((used)) __memset_wram_4align(void *dest, int c, size_t len)
{
    uint32_t cccc;
    memset_wram_t *dest32 = (memset_wram_t *)dest;

    c &= 0xff; /* Clear upper bits before ORing below */
    cccc = c | (c << 8) | (c << 16) | (c << 24);

    for (uint32_t i = 0; i < len / sizeof(memset_wram_t); ++i) {
        dest32[i] = cccc;
    }

    return dest;
}

void *
memset(void *dest, int c, size_t len)
{
    const uint32_t align = sizeof(memset_wram_t);
    const uint32_t align_off_mask = (align - 1);
    uint32_t align_offset = ((uintptr_t)dest) & align_off_mask;
    uint8_t *d = (uint8_t *)dest;

    /* memset head */
    if (align_offset != 0) {
        size_t head_len = align - align_offset;
        if (head_len > len) {
            head_len = len;
        }

        __memset_wram_1align(d, c, head_len);

        len -= head_len;
        d += head_len;
    }

    /* memset body */
    if (len >= align) {
        size_t body_len = len & (~align_off_mask);

        __memset_wram_4align(d, c, body_len);

        len -= body_len;
        d += body_len;
    }

    /* memset tail */
    if (len > 0) {
        __memset_wram_1align(d, c, len);
    }

    return dest;
}

#define MEMSET_MRAM_CACHE_SIZE (8)
/* Requisite:
 *  - dest: align on 8 bytes
 *  - len: mutiple of 8 bytes
 */
__attribute__((used)) __mram_ptr void *
__memset_mram_8align(__mram_ptr void *dest, int c, size_t len)
{
    __dma_aligned uint8_t cache64[MEMSET_MRAM_CACHE_SIZE];
    void *cache = (void *)cache64;

    __memset_wram_4align(cache, c, MEMSET_MRAM_CACHE_SIZE);

    for (uint32_t idx = 0; idx < len; idx += MEMSET_MRAM_CACHE_SIZE) {
        mram_write(cache, dest + idx, MEMSET_MRAM_CACHE_SIZE);
    }

    return dest;
}

__attribute__((used)) __mram_ptr void *
__memset_mram(__mram_ptr void *dest, int c, size_t len)
{
    __dma_aligned uint8_t cache64[MEMSET_MRAM_CACHE_SIZE];
    void *cache = (void *)cache64;
    __mram_ptr uint8_t *d = (__mram_ptr uint8_t *)((uintptr_t)dest & (~DMA_OFF_MASK));
    uint32_t align_offset = ((uintptr_t)dest) & DMA_OFF_MASK;

    /* memset head */
    if (align_offset != 0) {
        size_t head_len = MEMSET_MRAM_CACHE_SIZE - align_offset;
        if (head_len > len) {
            head_len = len;
        }

        mram_read(d, cache, MEMSET_MRAM_CACHE_SIZE);
        __memset_wram_1align(cache + align_offset, c, head_len);
        mram_write(cache, d, MEMSET_MRAM_CACHE_SIZE);

        len -= head_len;
        d += MEMSET_MRAM_CACHE_SIZE;
    }

    /* memset body */
    if (len >= MRAM_CACHE_SIZE) {
        size_t body_len = len & (~(MEMSET_MRAM_CACHE_SIZE - 1));

        __memset_mram_8align(d, c, body_len);

        len -= body_len;
        d += body_len;
    }

    /* memset tail */
    if (len > 0) {
        mram_read(d, cache, MEMSET_MRAM_CACHE_SIZE);
        __memset_wram_1align(cache, c, len);
        mram_write(cache, d, MEMSET_MRAM_CACHE_SIZE);
    }

    return dest;
}
