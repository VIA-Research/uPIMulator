/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_MUTEX_H
#define DPUSYSCORE_MUTEX_H

/**
 * @file mutex.h
 * @brief Mutual exclusions.
 *
 * A mutex ensures mutual exclusion between threads: only one runtime can have the mutex at a time, blocking all the
 * other threads trying to take the mutex.
 *
 * @internal All the mutexes are stored in a table in WRAM. In this table, each byte represents a mutex,
 *           and can be accessed directly by taking the base address of the table and adding it the sysname of the mutex
 *           we want (thus, the sysname should be an integer in the range [0; NB_MUTEX -1]). The result of
 *           this addition is what a mutex_get will return.
 *           A lock is made by using an lb_a instruction on the address of the mutex given as a parameter.
 *           An unlock is made by using an sb_r instruction on the address of the mutex given as a parameter.
 *           The id of the runtime doing the unlock is what is currently stored at the address of the mutex.
 *           The base address of this table is associated with the pointer defined by __MUTEX_TABLE__.
 */

#include <stdint.h>
#include <sysdef.h>
#include <stdbool.h>
#include <atomic_bit.h>

/**
 * @typedef mutex_id_t
 * @brief A mutex object reference, as declared by MUTEX_INIT.
 */
typedef uint8_t *mutex_id_t;

/**
 * @def MUTEX_GET
 * @hideinitializer
 * @brief Return the symbol to use when using the mutex associated to the given name.
 */
#define MUTEX_GET(_name) _name

/**
 * @def MUTEX_INIT
 * @hideinitializer
 * @brief Declare and initialize a mutex associated to the given name.
 */
#define MUTEX_INIT(_name)                                                                                                        \
    ATOMIC_BIT_INIT(__CONCAT(mutex_, _name));                                                                                    \
    const mutex_id_t MUTEX_GET(_name) = &ATOMIC_BIT_GET(__CONCAT(mutex_, _name))

/**
 * @fn mutex_lock
 * @brief Takes the lock on the given mutex.
 * @param mutex the mutex we want to lock
 */
static inline void
mutex_lock(mutex_id_t mutex)
{
    __asm__ volatile("acquire %[mtx], 0, nz, ." : : [mtx] "r"(mutex) :);
}

/**
 * @fn mutex_trylock
 * @brief Tries to take the lock on the given mutex. If the lock is already taken, returns immediately.
 * @param mutex the mutex we want to lock
 * @return Whether the mutex has been successfully locked.
 */
static inline bool
mutex_trylock(mutex_id_t mutex)
{
    bool result = true;
    __asm__ volatile("acquire %[mtx], 0, z, .+2; move %[res], 0" : [res] "+r"(result) : [mtx] "r"(mutex) :);
    return result;
}

/**
 * @fn mutex_unlock
 * @brief Releases the lock on the given mutex.
 * @param mutex the mutex we want to unlock
 */
static inline void
mutex_unlock(mutex_id_t mutex)
{
    __asm__ volatile("release %[mtx], 0, nz, .+1" : : [mtx] "r"(mutex) :);
}

#endif /* DPUSYSCORE_MUTEX_H */
