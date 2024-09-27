/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <defs.h>
#include <sysdef.h>
#include <dpuruntime.h>

int
check_stack()
{
    unsigned int stack_base, stack_size;
    int stack_limit, remaining;
    thread_id_t tid = me();

    stack_base = __SP_TABLE__[tid].stack_ptr;
    stack_size = __SP_TABLE__[tid].stack_size;
    stack_limit = (int)(stack_base + stack_size);
    __asm__ volatile("sub %[r], %[l], r22" : [r] "=r"(remaining) : [l] "r"(stack_limit));

    return remaining;
}
