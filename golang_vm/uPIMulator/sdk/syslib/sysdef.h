/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_SYSDEF_H
#define DPUSYSCORE_SYSDEF_H

/**
 * @file sysdef.h
 * @brief Provides useful system abstractions.
 */

/**
 * @typedef thread_id_t
 * @brief A unique runtime number.
 */
typedef unsigned int thread_id_t;

/**
 * @typedef sysname_t
 * @brief A system name.
 *
 * Used to name system structures, like mutexes, semaphores, meetpoints, etc... In practice, system names
 * are integers, representing a unique identifier for the given type of structure.
 */
typedef unsigned int sysname_t;

#endif /* DPUSYSCORE_SYSDEF_H */
