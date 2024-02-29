/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_ATTRIBUTES_H
#define DPUSYSCORE_ATTRIBUTES_H

/**
 * @file attributes.h
 * @brief Provides common useful compiler attributes.
 */

#define DEPRECATED __attribute__((deprecated))

#if __STDC_VERSION__ >= 201112L
#define __NO_RETURN _Noreturn
#else
#define __NO_RETURN
#endif /* __STDC_VERSION__ */

#define __weak __attribute__((weak))

#define __section(s) __attribute__((section(s)))

#define __aligned(a) __attribute__((aligned(a)))

#define __used __attribute__((used))

#define __noinline __attribute__((noinline))

#define __atomic_bit __section(".atomic")

#define __dma_aligned __aligned(8)

#define __keep __used __section(".data.__sys_keep")

#define __host __aligned(8) __used __section(".dpu_host")

// Use this macro at variable definition to place this variable into the section
// .data.immediate_memory and then makes it possible to use this variable
// directly as an immediate into load store instructions (and then avoids the need
// to move the address into a register before): immediate values are 12 signed bits
// large.
#define __lower_data(name) __attribute__((used, section(".data.immediate_memory." name)))

/**
 * @def __mram_ptr
 * @brief An attribute declaring that a pointer is an address in MRAM.
 *
 * A typical usage is: ``unsigned int __mram_ptr * array32 = (unsigned int __mram_ptr *) 0xf000;``
 *
 * Performing a cast between a pointer in MRAM and a pointer in WRAM is not allowed by the compiler.
 *
 */
#define __mram_ptr __attribute__((address_space(255)))

#define __mram __mram_ptr __section(".mram") __dma_aligned __used

#define __mram_noinit __mram_ptr __section(".mram.noinit") __dma_aligned __used

#define __mram_keep __mram_ptr __section(".mram.keep") __dma_aligned __used

#define __mram_noinit_keep __mram_ptr __section(".mram.noinit.keep") __dma_aligned __used

#endif /* DPUSYSCORE_ATTRIBUTES_H */
