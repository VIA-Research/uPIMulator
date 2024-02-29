/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_BARRIER_H
#define DPUSYSCORE_BARRIER_H

/**
 * @file barrier.h
 * @brief Synchronization with barriers.
 *
 * This synchronization mechanism allows to suspend a fixed number of tasklets until the expected number of subscribers is
 * present. When the required number of tasklets reached the barrier, the counter of the barrier will be reinitialised to the
 * original value.
 *
 * @internal The barriers are represented by a static value, defining the number of expected tasklets, a counter for the
 *           current number of tasklets suspended by this barrier and a wait queue entry.
 *           Whenever a new tasklet reaches the barrier (barrier_wait), the counter is  decremented and the tasklet
 *           is put into the wait queue.
 *           If the counter is reduced to 0, all the tasklets that were suspended by this barrier will be resumed and
 *           the counter will be reinitialised to its initial value.
 */

#include <attributes.h>
#include <atomic_bit.h>
#include <stdint.h>

/**
 * @typedef barrier_t
 * @brief A barrier object, as declared by BARRIER_INIT.
 */
typedef struct barrier_t {
    uint8_t wait_queue;
    uint8_t count;
    uint8_t initial_count;
    uint8_t lock;
} barrier_t;

/**
 * @def BARRIER_INIT
 * @hideinitializer
 * @brief Declare and initialize a barrier associated to the given name.
 */
/* clang-format off */
#define BARRIER_INIT(_name, _counter)                                                                                            \
    _Static_assert((_counter < 128) && (_counter >= -127), "barrier counter must be encoded on a byte");                         \
    ATOMIC_BIT_INIT(__CONCAT(barrier_, _name));                                                                                  \
    extern barrier_t _name;                                                                                                      \
    __asm__(".section .data." __STR(_name) "\n"                                                                                  \
            ".type " __STR( _name) ",@object\n"                                                                                  \
            ".globl " __STR( _name) "\n"                                                                                         \
            ".p2align 2\n" __STR(_name) ":\n"                                                                                    \
            ".byte 0xFF\n"                                                                                                       \
            ".byte " __STR(_counter) "\n"                                                                                        \
            ".byte " __STR(_counter) "\n"                                                                                        \
            ".byte " __STR(ATOMIC_BIT_GET(__CONCAT(barrier_,_name))) "\n"                                                        \
            ".size " __STR(_name) ", 4\n"                                                                                        \
            ".text");
/* clang-format on */

/**
 * @fn barrier_wait
 * @brief Decrements the counter associated to the barrier and suspends the invoking tasklet.
 *
 * The counter of the barrier is decremented and the invoking tasklet is suspended until
 * the counter associated to the barrier is reduced to 0.
 *
 * @param barrier the barrier the tasklet will be associated to.
 */
void
barrier_wait(barrier_t *barrier);

#endif /* DPUSYSCORE_BARRIER_H */
