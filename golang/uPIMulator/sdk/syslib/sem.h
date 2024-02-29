/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_SEM_H
#define DPUSYSCORE_SEM_H

/**
 * @file sem.h
 * @brief Synchronization with semaphores.
 *
 * A semaphore is characterized by a counter and a wait queue. It provides two functions:
 *
 *   - Take: the counter is decremented by 1. If the counter is negative, the runtime is blocked (stop) and placed in the
 *     semaphore's wait queue, waiting to be resume by another runtime.
 *   - Give: the counter is incremented by 1. If the counter was negative before the increment, the runtime resumes the execution
 *     of the first runtime waiting in the waiting queue. In all the cases, the runtime continues its own execution.
 *
 */

#include <attributes.h>
#include <atomic_bit.h>
#include <stdint.h>

/**
 * @typedef sem_t
 * @brief A semaphore object, as declared by SEMAPHORE_INIT.
 */
typedef struct sem_t {
    uint8_t wait_queue;
    uint8_t count;
    uint8_t initial_count;
    uint8_t lock;
} sem_t;

/**
 * @def SEMAPHORE_INIT
 * @hideinitializer
 * @brief Declare and initialize a semaphore associated to the given name.
 */
/* clang-format off */
#define SEMAPHORE_INIT(_name, _counter)                                                                                          \
    _Static_assert((_counter < 128) && (_counter >= -127), "semaphore counter must be encoded on a byte");                       \
    ATOMIC_BIT_INIT(__CONCAT(semaphore_, _name));                                                                                \
    extern sem_t (_name);                                                                                                        \
    __asm__(".section .data." __STR(_name) "\n"                                                                                  \
            ".type " __STR(_name) ",@object\n"                                                                                   \
            ".globl " __STR(_name) "\n"                                                                                          \
            ".p2align 2\n" __STR(_name) ":\n"                                                                                    \
            ".byte 0xFF\n"                                                                                                       \
            ".byte " __STR(_counter) "\n"                                                                                        \
            ".byte " __STR(_counter) "\n"                                                                                        \
            ".byte " __STR(ATOMIC_BIT_GET(__CONCAT(semaphore_, _name))) "\n"                                                     \
            ".size " __STR(_name) ", 4\n"                                                                                        \
            ".text");
/* clang-format on */

/**
 * @fn sem_take
 * @brief Takes one unit in the given semaphore (cf Take definition).
 * @param sem the semaphore we want to take
 */
void
sem_take(sem_t *sem);

/**
 * @fn sem_give
 * @brief Gives on unit in the given semaphore (cf Give definition).
 * @param sem the semaphore we want to give
 */
void
sem_give(sem_t *sem);

#endif /* DPUSYSCORE_SEM_H */
