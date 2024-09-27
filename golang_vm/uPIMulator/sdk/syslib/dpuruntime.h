/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_DPURUNTIME_H
#define DPUSYSCORE_DPURUNTIME_H

#include <built_ins.h>
#include <dpuconst.h>
#include <dpufault.h>
#include <macro_utils.h>
#include <stdint.h>

// todo fix: This file should not be included by another syslib header file, only by source files (conflicting definitions).

#define __INITIAL_HEAP_POINTER __sys_heap_pointer_reset
#define __HEAP_POINTER __sys_heap_pointer
#define __WAIT_QUEUE_TABLE __sys_wq_table
#define __SP_TABLE__ __sys_thread_stack_table_ptr
#define __STDOUT_BUFFER_STATE __stdout_buffer_state

/* The order needs to match the __bootstrap function expectation */
typedef struct {
    uint32_t stack_size;
    uint32_t stack_ptr;
} thread_stack_t;

extern unsigned int __INITIAL_HEAP_POINTER;
extern volatile unsigned int __HEAP_POINTER;
extern unsigned char __WAIT_QUEUE_TABLE[];
extern thread_stack_t __SP_TABLE__[];

#define __acquire(base, off) __builtin_acquire_rici(base, off, "nz", __AT_THIS_INSTRUCTION)
#define __release(base, off, at) __builtin_release_rici(base, off, "nz", at)

#define __resume(base, off) __builtin_resume_rici(base, off, "nz", __AT_THIS_INSTRUCTION)
#define __stop() __builtin_stop_ci("false", "0")
#define __stop_at(label) __builtin_stop_ci("true", label)

#define likely(x) __builtin_expect((x), 1)
#define unlikely(x) __builtin_expect((x), 0)
#define unreachable() __builtin_unreachable()

#define count_leading_zeros(x) __builtin_clz(x)
#define count_population(x) __builtin_popcount(x)

#define __EMPTY_WAIT_QUEUE 0xFF

#define __AT_THIS_INSTRUCTION ".+0"
#define __AT_NEXT_INSTRUCTION ".+1"

// Use this macro at the beginning of an assembly function in order to get profiled.
#ifdef DPU_PROFILING
#define __ADD_PROFILING_ENTRY__ "call r23, mcount\n"
#else
#define __ADD_PROFILING_ENTRY__ "\n"
#endif

#ifdef DPU_PROFILING
/* Reset counter + count cycles */
#define __CONFIG_PERFCOUNTER_ENTRY__                                                                                             \
    "  move r23, 3\n"                                                                                                            \
    "  time_cfg zero, r23\n"
#define __SAVE_PERFCOUNTER_ENTRY__                                                                                               \
    "  time r23\n"                                                                                                               \
    "  sw zero, perfcounter_end_value, r23\n"
#else
#define __CONFIG_PERFCOUNTER_ENTRY__ "\n"
#define __SAVE_PERFCOUNTER_ENTRY__ "\n"
#endif

#endif /* DPUSYSCORE_DPURUNTIME_H */
