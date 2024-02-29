/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <seqread.h>
#include <mram.h>
#include <dpuconst.h>
#include <dpuruntime.h>
#include <atomic_bit.h>
#include <stddef.h>

#define PAGE_SIZE (SEQREAD_CACHE_SIZE)
#define PAGE_ALLOC_SIZE (2 * PAGE_SIZE)
#define PAGE_OFF_MASK (PAGE_SIZE - 1)
#define PAGE_IDX_MASK (~PAGE_OFF_MASK)

#define MRAM_READ_PAGE(from, to) mram_read((__mram_ptr void *)(from), (void *)(to), PAGE_ALLOC_SIZE)

extern void *
mem_alloc_nolock(size_t size);
ATOMIC_BIT_EXTERN(__heap_pointer);

seqreader_buffer_t
__SEQREAD_ALLOC(void)
{
    ATOMIC_BIT_ACQUIRE(__heap_pointer);

    unsigned int heap_pointer = __HEAP_POINTER;
    seqreader_buffer_t pointer = (seqreader_buffer_t)((heap_pointer + PAGE_OFF_MASK) & PAGE_IDX_MASK);
    size_t size = pointer + PAGE_ALLOC_SIZE - heap_pointer;
    /* We already compute the return pointer
     * mem_alloc_nolock is only used to really reserve the memory area in the heap */
    mem_alloc_nolock(size);

    ATOMIC_BIT_RELEASE(__heap_pointer);

    return pointer;
}

void *
__SEQREAD_INIT(seqreader_buffer_t cache, __mram_ptr void *mram_addr, seqreader_t *reader)
{
    reader->wram_cache = cache;
    reader->mram_addr = (uintptr_t)(1 << __DPU_MRAM_SIZE_LOG2);

    return __SEQREAD_SEEK(mram_addr, reader);
}

void *
__SEQREAD_GET(void *ptr, uint32_t inc, seqreader_t *reader)
{
    return (void *)__builtin_dpu_seqread_get((uintptr_t)ptr, inc, reader, PAGE_SIZE);
}

void *
__SEQREAD_SEEK(__mram_ptr void *mram_addr, seqreader_t *reader)
{
    uintptr_t target_addr = (uintptr_t)mram_addr;
    uintptr_t current_addr = (uintptr_t)reader->mram_addr;
    uintptr_t wram_cache = (uintptr_t)reader->wram_cache;
    uintptr_t mram_offset = target_addr - current_addr;
    if ((mram_offset & PAGE_IDX_MASK) != 0) {
        uintptr_t target_addr_idx_page = target_addr & PAGE_IDX_MASK;
        MRAM_READ_PAGE(target_addr_idx_page, wram_cache);
        mram_offset = target_addr & PAGE_OFF_MASK;
        reader->mram_addr = target_addr_idx_page;
    }
    return (void *)(mram_offset + wram_cache);
}

__mram_ptr void *
__SEQREAD_TELL(void *ptr, seqreader_t *reader)
{
    return (__mram_ptr void *)((uintptr_t)reader->mram_addr + ((uintptr_t)ptr & PAGE_OFF_MASK));
}
