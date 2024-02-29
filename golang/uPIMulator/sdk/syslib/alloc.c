/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <stddef.h>
#include <dpuruntime.h>
#include <atomic_bit.h>
#include <attributes.h>

volatile unsigned int __sys_heap_pointer = (unsigned int)(&__sys_heap_pointer_reset);

ATOMIC_BIT_INIT(__heap_pointer);

/* noinline, because part of grind tracked functions
 * Also used by seqread.inc
 */
void *__noinline
mem_alloc_nolock(size_t size)
{
    unsigned int pointer = __HEAP_POINTER;

    if (size != 0) {
        pointer = (pointer + 7) & ~7;

        unsigned int new_heap_pointer, dummy;

        __asm__ volatile("\tadd %[nhp], %[ptr], %[sz], nc, . + 2\n"
                         "\tfault " __STR(__FAULT_ALLOC_HEAP_FULL__) "\n"
                                                                     "\tlbu %[dumb], %[nhp], -1\n"
                         : [nhp] "=r"(new_heap_pointer), [dumb] "=r"(dummy)
                         : [ptr] "r"(pointer), [sz] "r"(size));

        __HEAP_POINTER = new_heap_pointer;
    }
    return (void *)pointer;
}

void *
mem_alloc(size_t size)
{
    ATOMIC_BIT_ACQUIRE(__heap_pointer);
    void *pointer = mem_alloc_nolock(size);
    ATOMIC_BIT_RELEASE(__heap_pointer);
    return pointer;
}

void *
mem_reset()
{
    ATOMIC_BIT_ACQUIRE(__heap_pointer);

    void *initial = &__sys_heap_pointer_reset;

    __sys_heap_pointer = (unsigned int)initial;

    ATOMIC_BIT_RELEASE(__heap_pointer);

    return (void *)initial;
}
