/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <fsb_allocator.h>
#include <alloc.h>
#include <stddef.h>
#include <dpuruntime.h>
#include <attributes.h>
#include <atomic_bit.h>

ATOMIC_BIT_INIT(__fsb_lock);

// noinline, because part of grind tracked functions
fsb_allocator_t __noinline
fsb_alloc(unsigned int block_size, unsigned int nb_of_blocks)
{
    if (block_size > 0xFFFFFFF8) {
        __asm__ volatile("fault " __STR(__FAULT_ALLOC_HEAP_FULL__));
        unreachable();
    }

    block_size = (block_size == 0) ? 8 : (block_size + 7) & ~7;

    unsigned int memory_space_to_allocate = block_size * nb_of_blocks + 4;
    void *memory = mem_alloc(memory_space_to_allocate);

    unsigned int first_block = (unsigned int)memory;

    for (unsigned int each_block = 0; each_block < nb_of_blocks - 1; ++each_block) {
        unsigned int next_block_address = (unsigned int)(memory + block_size);
        *((unsigned int *)memory) = next_block_address;
        memory = (void *)next_block_address;
    }

    *((unsigned int *)memory) = 0;
    memory += block_size;

    void *free_ptr = memory;
    *((unsigned int *)free_ptr) = first_block;

    return (fsb_allocator_t)free_ptr;
}

// noinline, because part of grind tracked functions
void *__noinline
fsb_get(fsb_allocator_t allocator)
{
    void **result;
    ATOMIC_BIT_ACQUIRE(__fsb_lock);
    __asm__ volatile("lw %[res], %[alloc], 0" : [res] "=r"(result) : [alloc] "r"(allocator));

    if (result == NULL) {
        ATOMIC_BIT_RELEASE(__fsb_lock);
        return NULL;
    }

    void *next = *result;

    __asm__ volatile("sw %[alloc], 0, %[res]" : : [res] "r"(next), [alloc] "r"(allocator));
    ATOMIC_BIT_RELEASE(__fsb_lock);

    return (void *)result;
}

// noinline, because part of grind tracked functions
void __noinline
fsb_free(fsb_allocator_t allocator, void *ptr)
{
    void *next_free;
    ATOMIC_BIT_ACQUIRE(__fsb_lock);
    __asm__ volatile("lw %[res], %[alloc], 0" : [res] "=r"(next_free) : [alloc] "r"(allocator));

    *((void **)ptr) = next_free;

    __asm__ volatile("sw %[alloc], 0, %[res]" : : [res] "r"(ptr), [alloc] "r"(allocator));
    ATOMIC_BIT_RELEASE(__fsb_lock);
}
