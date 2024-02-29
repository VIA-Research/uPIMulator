/* Copyright 2021 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */
#include <profiling.h>
#include <defs.h>
#include <dpufault.h>
#include <dpuruntime.h>
#include <macro_utils.h>
#include <stdint.h>
#include <sysdef.h>

#ifdef DPU_PROFILING
void __attribute__((no_instrument_function)) profiling_start(dpu_profiling_t *context)
{
    thread_id_t tid = me();
    uint32_t perfcounter_value;

    if (unlikely(context->start[tid] != PROFILING_RESET_VALUE)) {
        __asm__("fault " __STR(__FAULT_ALREADY_PROFILING__));
        unreachable();
    }

    __asm__ volatile("time %[r]" : [r] "=r"(perfcounter_value));
    context->start[tid] = perfcounter_value;
}

void __attribute__((no_instrument_function)) profiling_stop(dpu_profiling_t *context)
{
    thread_id_t tid = me();
    uint32_t perfcounter_value;

    if (unlikely(context->start[tid] == PROFILING_RESET_VALUE)) {
        __asm__("fault " __STR(__FAULT_NOT_PROFILING__));
        unreachable();
    }

    __asm__ volatile("time %[r]" : [r] "=r"(perfcounter_value));
    context->count[tid] += perfcounter_value - context->start[tid];
    context->start[tid] = PROFILING_RESET_VALUE;
}
#endif /* DPU_PROFILING */
