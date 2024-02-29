/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <memmram_utils.h>
#include <mram.h>
#include <stddef.h>
#include <stdint.h>

__attribute__((used)) void *
__memcpy_wram_4align(void *dest, const void *src, size_t len)
{
    uint32_t *dw = (uint32_t *)dest;
    uint32_t *sw = (uint32_t *)src;

    for (uint32_t i = 0; i < (len / sizeof(uint32_t)); ++i) {
        dw[i] = sw[i];
    }
    return dest;
}

void *
memcpy(void *dest, const void *src, size_t len)
{
    uint8_t *d = (uint8_t *)dest;
    const uint8_t *s = (const uint8_t *)src;
    uint32_t *dw;
    const uint32_t *sw;
    uint8_t *head;
    uint8_t *const tail = (uint8_t *)dest + len;
    /* Set 'body' to the last word boundary */
    uint32_t *const body = (uint32_t *)((uintptr_t)tail & ~3);

    if (((uintptr_t)dest & 3) != ((uintptr_t)src & 3)) {
        /* Misaligned. no body, no tail. */
        head = tail;
    } else {
        /* Aligned */
        if ((uintptr_t)tail < (((uintptr_t)d + 3) & ~3))
            /* len is shorter than the first word boundary */
            head = tail;
        else
            /* Set 'head' to the first word boundary */
            head = (uint8_t *)(((uintptr_t)d + 3) & ~3);
    }

    /* Copy head */
    uint32_t head_len = head - d;
    if (head_len != 0) {
        for (uint32_t i = 0; i < head_len; ++i)
            d[i] = s[i];
    }

    /* Copy body */
    dw = (uint32_t *)(d + head_len);
    sw = (uint32_t *)(s + head_len);

    uint32_t body_len = (body < dw) ? 0 : body - dw;
    if (body_len != 0) {
        __memcpy_wram_4align(dw, sw, body_len * sizeof(uint32_t));
    }

    /* Copy tail */
    d = (uint8_t *)(dw + body_len);
    s = (const uint8_t *)(sw + body_len);
    uint32_t tail_len = tail - d;
    if (tail_len != 0) {
        for (uint32_t i = 0; i < tail_len; ++i)
            d[i] = s[i];
    }

    return dest;
}

__attribute__((used)) __mram_ptr void *
__memcpy_mw(__mram_ptr void *dest, const void *src, size_t len)
{
    uint64_t destCache64[MRAM_CACHE_SIZE / sizeof(uint64_t)];
    void *destCache = (void *)destCache64;

    uint32_t srcOff = ((uintptr_t)src) & DMA_OFF_MASK;
    uint32_t destOff = ((uintptr_t)dest) & DMA_OFF_MASK;
    size_t remaining = len;

    uint32_t idx = 0;

    if (destOff != 0) {
        size_t part = MIN(remaining, MRAM_CACHE_SIZE - destOff);
        mram_read(dest, destCache, MRAM_CACHE_SIZE);
        memcpy(destCache + destOff, src, part);
        mram_write(destCache, dest, MRAM_CACHE_SIZE);
        remaining -= part;
        idx += part;
    }

    if (srcOff == destOff) {
        while (remaining >= MRAM_CACHE_SIZE) {
            mram_write(src + idx, dest + idx, MRAM_CACHE_SIZE);
            remaining -= MRAM_CACHE_SIZE;
            idx += MRAM_CACHE_SIZE;
        }
    } else {
        while (remaining >= MRAM_CACHE_SIZE) {
            memcpy(destCache, src + idx, MRAM_CACHE_SIZE);
            mram_write(destCache, dest + idx, MRAM_CACHE_SIZE);
            remaining -= MRAM_CACHE_SIZE;
            idx += MRAM_CACHE_SIZE;
        }
    }

    if (remaining != 0) {
        mram_read(dest + idx, destCache, MRAM_CACHE_SIZE);
        memcpy(destCache, src + idx, remaining);
        mram_write(destCache, dest + idx, MRAM_CACHE_SIZE);
    }

    return dest;
}

__attribute__((used)) void *
__memcpy_wm(void *dest, const __mram_ptr void *src, size_t len)
{
    uint64_t srcCache64[MRAM_CACHE_SIZE / sizeof(uint64_t)];
    void *srcCache = (void *)srcCache64;

    uint32_t srcOff = ((uintptr_t)src) & DMA_OFF_MASK;
    uint32_t destOff = ((uintptr_t)dest) & DMA_OFF_MASK;
    size_t remaining = len;

    uint32_t idx = 0;
    size_t part = MIN(remaining, MRAM_CACHE_SIZE - srcOff);

    mram_read(src, srcCache, MRAM_CACHE_SIZE);
    memcpy(dest, srcCache + srcOff, part);
    remaining -= part;
    idx += part;

    if (srcOff == destOff) {
        while (remaining >= MRAM_CACHE_SIZE) {
            mram_read(src + idx, dest + idx, MRAM_CACHE_SIZE);
            remaining -= MRAM_CACHE_SIZE;
            idx += MRAM_CACHE_SIZE;
        }
    } else {
        while (remaining >= MRAM_CACHE_SIZE) {
            mram_read(src + idx, srcCache, MRAM_CACHE_SIZE);
            memcpy(dest + idx, srcCache, MRAM_CACHE_SIZE);
            remaining -= MRAM_CACHE_SIZE;
            idx += MRAM_CACHE_SIZE;
        }
    }

    if (remaining != 0) {
        mram_read(src + idx, srcCache, MRAM_CACHE_SIZE);
        memcpy(dest + idx, srcCache, remaining);
    }

    return dest;
}

__attribute__((used)) __mram_ptr void *
__memcpy_mm(__mram_ptr void *dest, const __mram_ptr void *src, size_t len)
{
    uint64_t srcCache64[MRAM_CACHE_SIZE / sizeof(uint64_t)];
    uint64_t destCache64[MRAM_CACHE_SIZE / sizeof(uint64_t)];
    void *srcCache = (void *)srcCache64;
    void *destCache = (void *)destCache64;

    uint32_t srcOff = ((uintptr_t)src) & DMA_OFF_MASK;
    uint32_t destOff = ((uintptr_t)dest) & DMA_OFF_MASK;
    size_t remaining = len;

    if (srcOff == destOff) {
        uint32_t idx = 0;

        if (destOff != 0) {
            size_t part = MIN(remaining, MRAM_CACHE_SIZE - srcOff);
            mram_read(dest, destCache, MRAM_CACHE_SIZE);
            mram_read(src, srcCache, MRAM_CACHE_SIZE);
            memcpy(destCache + destOff, srcCache + srcOff, part);
            mram_write(destCache, dest, MRAM_CACHE_SIZE);
            remaining -= part;
            idx += part;
        }

        while (remaining >= MRAM_CACHE_SIZE) {
            mram_read(src + idx, srcCache, MRAM_CACHE_SIZE);
            mram_write(srcCache, dest + idx, MRAM_CACHE_SIZE);
            remaining -= MRAM_CACHE_SIZE;
            idx += MRAM_CACHE_SIZE;
        }

        if (remaining != 0) {
            mram_read(dest + idx, destCache, MRAM_CACHE_SIZE);
            mram_read(src + idx, srcCache, MRAM_CACHE_SIZE);
            memcpy(destCache, srcCache, remaining);
            mram_write(destCache, dest + idx, MRAM_CACHE_SIZE);
        }
    } else {
        uint32_t srcIdx = 0;
        uint32_t destIdx = 0;
        size_t initLen = MIN(remaining, MRAM_CACHE_SIZE - MIN(destOff, srcOff));

        if (initLen == remaining) {
            mram_read(dest, destCache, MRAM_CACHE_SIZE);
            mram_read(src, srcCache, MRAM_CACHE_SIZE);
            memcpy(destCache + destOff, srcCache + srcOff, remaining);
            mram_write(destCache, dest, MRAM_CACHE_SIZE);
            return dest;
        }

        mram_read(src, srcCache, MRAM_CACHE_SIZE);
        srcIdx += MRAM_CACHE_SIZE;

        if (destOff != 0) {
            mram_read(dest, destCache, DMA_ALIGNED(destOff));

            if (destOff > srcOff) {
                size_t part = MRAM_CACHE_SIZE - destOff;
                memcpy(destCache + destOff, srcCache + srcOff, part);

                srcOff += part;
            } else {
                size_t part = MRAM_CACHE_SIZE - srcOff;
                memcpy(destCache + destOff, srcCache + srcOff, part);
                mram_read(src + srcIdx, srcCache, MRAM_CACHE_SIZE);
                srcIdx += MRAM_CACHE_SIZE;

                size_t part2 = srcOff - destOff;
                memcpy(destCache + destOff + part, srcCache, part2);

                srcOff = part2;
            }

            mram_write(destCache, dest + destIdx, MRAM_CACHE_SIZE);
            destIdx += MRAM_CACHE_SIZE;
            remaining -= MRAM_CACHE_SIZE - destOff;
        }

        while (remaining >= MRAM_CACHE_SIZE) {
            size_t part = MRAM_CACHE_SIZE - srcOff;
            memcpy(destCache, srcCache + srcOff, part);
            mram_read(src + srcIdx, srcCache, MRAM_CACHE_SIZE);
            srcIdx += MRAM_CACHE_SIZE;

            size_t part2 = srcOff;
            memcpy(destCache + part, srcCache, part2);
            mram_write(destCache, dest + destIdx, MRAM_CACHE_SIZE);
            remaining -= MRAM_CACHE_SIZE;
            destIdx += MRAM_CACHE_SIZE;
        }

        if (remaining != 0) {
            mram_read(dest + destIdx, destCache, MRAM_CACHE_SIZE);

            size_t part = MRAM_CACHE_SIZE - srcOff;
            memcpy(destCache, srcCache + srcOff, part);

            if (remaining > part) {
                size_t part2 = remaining - part;
                mram_read(src + srcIdx, srcCache, MRAM_CACHE_SIZE);
                memcpy(destCache + part, srcCache, part2);
            }

            mram_write(destCache, dest + destIdx, MRAM_CACHE_SIZE);
        }
    }

    return dest;
}
