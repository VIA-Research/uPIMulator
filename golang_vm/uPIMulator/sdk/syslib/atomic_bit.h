/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_ATOMIC_BIT_H
#define DPUSYSCORE_ATOMIC_BIT_H

/**
 * @file atomic_bit.h
 * @brief Provides direct access to the atomic bits.
 */

#include <stdint.h>
#include <macro_utils.h>
#include <attributes.h>

#define ATOMIC_BIT_GET(_name) __CONCAT(__atomic_bit_, _name)
#define ATOMIC_BIT_INIT(_name) uint8_t __atomic_bit ATOMIC_BIT_GET(_name)
#define ATOMIC_BIT_EXTERN(_name) extern ATOMIC_BIT_INIT(_name)

extern uint8_t __atomic_start_addr;
#define ATOMIC_BIT_INDEX(_name) (&ATOMIC_BIT_GET(_name) - &__atomic_start_addr)

#define __ATOMIC_BIT_ACQUIRE(_reg, _bit)                                                                                         \
    __asm__ volatile("acquire %[areg], %[abit], nz, ." : : [areg] "r"(_reg), [abit] "i"(_bit))

#define __ATOMIC_BIT_RELEASE(_reg, _bit)                                                                                         \
    __asm__ volatile("release %[areg], %[abit], nz, .+1" : : [areg] "r"(_reg), [abit] "i"(_bit))

#define ATOMIC_BIT_ACQUIRE(_name) __asm__ volatile("acquire zero, %[abit], nz, ." : : [abit] "i"(&ATOMIC_BIT_GET(_name)))

#define ATOMIC_BIT_RELEASE(_name) __asm__ volatile("release zero, %[abit], nz, .+1" : : [abit] "i"(&ATOMIC_BIT_GET(_name)))

#endif /* DPUSYSCORE_ATOMIC_BIT_H */
