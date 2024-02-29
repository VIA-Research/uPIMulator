/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <memmram_utils.h>
#include <mram.h>
#include <string.h>
#include <stdint.h>

void *
memmove(void *dest, const void *src, size_t len)
{
    if ((uintptr_t)dest <= (uintptr_t)src || (uintptr_t)dest >= (uintptr_t)src + len) {
        /* Start of destination doesn't overlap source, so just use
         * memcpy(). */
        return memcpy(dest, src, len);
    } else {
        /* Need to copy from tail because there is overlap. */
        char *d = (char *)dest + len;
        const char *s = (const char *)src + len;
        uint32_t *dw;
        const uint32_t *sw;
        char *head;
        char *const tail = (char *)dest;
        /* Set 'body' to the last word boundary */
        uint32_t *const body = (uint32_t *)(((uintptr_t)tail + 3) & ~3);

        if (((uintptr_t)dest & 3) != ((uintptr_t)src & 3)) {
            /* Misaligned. no body, no tail. */
            head = tail;
        } else {
            /* Aligned */
            if ((uintptr_t)tail > ((uintptr_t)d & ~3))
                /* Shorter than the first word boundary */
                head = tail;
            else
                /* Set 'head' to the first word boundary */
                head = (char *)((uintptr_t)d & ~3);
        }

        /* Copy head */
        uint32_t head_len = d - head;
        for (int32_t i = head_len - 1; i >= 0; --i)
            d[i - head_len] = s[i - head_len];

        /* Copy body */
        dw = (uint32_t *)(d - head_len);
        sw = (uint32_t *)(s - head_len);

        uint32_t body_len = (dw < body) ? 0 : dw - body;
        for (int32_t i = body_len - 1; i >= 0; --i)
            dw[i - body_len] = sw[i - body_len];

        /* Copy tail */
        d = (char *)(dw - body_len);
        s = (const char *)(sw - body_len);

        uint32_t tail_len = d - tail;
        for (int32_t i = tail_len - 1; i >= 0; --i)
            d[i - tail_len] = s[i - tail_len];

        return dest;
    }
}

__mram_ptr void *
__memmove_mm(__mram_ptr void *dest, __mram_ptr const void *src, size_t len)
{
    if ((uintptr_t)dest <= (uintptr_t)src || (uintptr_t)dest >= (uintptr_t)src + len) {
        /* Start of destination doesn't overlap source, so just use
         * memcpy(). */
        return (__mram_ptr void *)memcpy(dest, src, len);
    } else {
        uint64_t srcCache64[MRAM_CACHE_SIZE / sizeof(uint64_t)];
        uint64_t destCache64[MRAM_CACHE_SIZE / sizeof(uint64_t)];
        void *srcCache = (void *)srcCache64;
        void *destCache = (void *)destCache64;

        __mram_ptr const void *srcIdx = src + len;
        __mram_ptr void *dstIdx = dest + len;
        uint32_t remaining = len;

        uint32_t srcOff = ((uintptr_t)srcIdx) & DMA_OFF_MASK;
        uint32_t dstOff = ((uintptr_t)dstIdx) & DMA_OFF_MASK;

        if (srcOff == dstOff) {
            size_t part = MIN(remaining, srcOff);
            uint32_t off = srcOff - part;

            if (dstOff != 0) {
                srcIdx -= srcOff;
                dstIdx -= dstOff;

                mram_read(dstIdx, destCache, MRAM_CACHE_SIZE);
                mram_read(srcIdx, srcCache, MRAM_CACHE_SIZE);
                memcpy(destCache + off, srcCache + off, part);
                mram_write(destCache, dstIdx, MRAM_CACHE_SIZE);
                remaining -= part;
            }

            srcIdx -= MRAM_CACHE_SIZE;
            dstIdx -= MRAM_CACHE_SIZE;

            while (remaining >= MRAM_CACHE_SIZE) {
                mram_read(srcIdx, srcCache, MRAM_CACHE_SIZE);
                mram_write(srcCache, dstIdx, MRAM_CACHE_SIZE);
                remaining -= MRAM_CACHE_SIZE;
                srcIdx -= MRAM_CACHE_SIZE;
                dstIdx -= MRAM_CACHE_SIZE;
            }

            if (remaining != 0) {
                uint32_t off = MRAM_CACHE_SIZE - remaining;

                mram_read(dstIdx, destCache, MRAM_CACHE_SIZE);
                mram_read(srcIdx, srcCache, MRAM_CACHE_SIZE);
                memcpy(destCache + off, srcCache + off, remaining);
                mram_write(destCache, dstIdx, MRAM_CACHE_SIZE);
            }
        } else {
            size_t initLen = MIN(remaining, MIN(dstOff, srcOff));

            if (initLen == remaining) {
                mram_read(dest, destCache, MRAM_CACHE_SIZE);
                mram_read(src, srcCache, MRAM_CACHE_SIZE);
                memcpy(destCache + dstOff - remaining, srcCache + srcOff - remaining, remaining);
                mram_write(destCache, dest, MRAM_CACHE_SIZE);
                return dest;
            }

            mram_read(srcIdx, srcCache, MRAM_CACHE_SIZE);
            srcIdx -= MRAM_CACHE_SIZE;

            if (dstOff != 0) {
                size_t part = DMA_ALIGNED(dstOff);
                mram_read(dstIdx, destCache + MRAM_CACHE_SIZE - part, part);

                if (srcOff > dstOff) {
                    part = MRAM_CACHE_SIZE - (DMA_ALIGNMENT - dstOff);
                    memcpy(destCache, srcCache + srcOff - part, part);
                    srcOff -= part;
                } else {
                    part = MRAM_CACHE_SIZE - (DMA_ALIGNMENT - srcOff);
                    memcpy(destCache + dstOff - part, srcCache, part);
                    mram_read(srcIdx, srcCache, MRAM_CACHE_SIZE);
                    srcIdx -= MRAM_CACHE_SIZE;

                    size_t part2 = dstOff - part;
                    memcpy(destCache, srcCache + MRAM_CACHE_SIZE - part2, part2);

                    srcOff = MRAM_CACHE_SIZE - part2;
                }

                mram_write(destCache, dstIdx, MRAM_CACHE_SIZE);
                dstIdx -= MRAM_CACHE_SIZE;
                remaining -= MRAM_CACHE_SIZE - (DMA_ALIGNMENT - dstOff);
            }

            while (remaining >= MRAM_CACHE_SIZE) {
                size_t part = MRAM_CACHE_SIZE - (DMA_ALIGNMENT - srcOff);
                memcpy(destCache + MRAM_CACHE_SIZE - part, srcCache, part);
                mram_read(srcIdx, srcCache, MRAM_CACHE_SIZE);
                srcIdx -= MRAM_CACHE_SIZE;

                size_t part2 = MRAM_CACHE_SIZE - part;
                memcpy(destCache, srcCache + MRAM_CACHE_SIZE - part2, part2);
                mram_write(destCache, dstIdx, MRAM_CACHE_SIZE);
                dstIdx -= MRAM_CACHE_SIZE;
                remaining -= MRAM_CACHE_SIZE;
            }

            if (remaining != 0) {
                mram_read(dstIdx, destCache, MRAM_CACHE_SIZE);

                size_t part = MRAM_CACHE_SIZE - (DMA_ALIGNMENT - srcOff);
                memcpy(destCache + MRAM_CACHE_SIZE - part, srcCache, part);

                if (remaining > part) {
                    size_t part2 = remaining - part;
                    mram_read(srcIdx, srcCache, MRAM_CACHE_SIZE);
                    memcpy(destCache + MRAM_CACHE_SIZE - remaining, srcCache + MRAM_CACHE_SIZE - part2, part2);
                }

                mram_write(destCache, dstIdx, MRAM_CACHE_SIZE);
            }
        }

        return dest;
    }
}
