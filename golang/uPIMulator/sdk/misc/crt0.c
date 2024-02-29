/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <macro_utils.h>
#include <attributes.h>
#include <stdint.h>
#include <dpuruntime.h>

__mram_noinit uint8_t __sys_used_mram_end[0];

uint64_t __sys_null_pointer __attribute__((used, section(".data.__sys_zero"))) = 0ULL;

uint32_t __lower_data("thread_profiling") thread_profiling[NR_THREADS];
uint32_t perfcounter_end_value;

#define DECLARE_STACK(x)                                                                                                         \
    extern uint32_t __sys_stack_thread_##x;                                                                                      \
    extern uint32_t STACK_SIZE_TASKLET_##x;
#define SET_STACK_TABLE_PTR(x)                                                                                                   \
    [x] = { .stack_ptr = (uint32_t)&__sys_stack_thread_##x, .stack_size = (uint32_t)&STACK_SIZE_TASKLET_##x },

__FOR_EACH_THREAD(DECLARE_STACK);

thread_stack_t __keep __dma_aligned __SP_TABLE__[NR_THREADS] = { __FOR_EACH_THREAD(SET_STACK_TABLE_PTR) };

__host const volatile uint32_t CLOCKS_PER_SEC;

__host const volatile uint32_t error_storage;

void __attribute__((naked, used, section(".text.__bootstrap"), no_instrument_function)) __bootstrap()
{
    /* clang-format off */
    __asm__ volatile(
        "  jnz id, __sys_start_thread\n"
        __CONFIG_PERFCOUNTER_ENTRY__
        "  sd zero, " __STR(__STDOUT_BUFFER_STATE) ", 0\n"
        "  move r23, " __STR(__atomic_end_addr) "\n"
        "__sys_atomic_bit_clear:\n"
        "  jeq r23, " __STR(__atomic_used_addr) ", __sys_start_thread\n"
        "  release r23, 0, nz, . + 1\n"
        "  add r23, r23, -1, true, __sys_atomic_bit_clear\n"
        "__sys_start_thread:\n"
        "  jeq id, NR_TASKLETS - 1, . + 2\n"
        "  boot id, 1\n"
        "  ld d22, id8, " __STR(__SP_TABLE__) "\n"
        "  call r23, main\n"
        ".globl __sys_end\n"
        "__sys_end:\n"
        __SAVE_PERFCOUNTER_ENTRY__
        "  stop true, __sys_end");
    /* clang-format on */
}
