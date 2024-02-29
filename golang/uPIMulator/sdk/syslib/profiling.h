/* Copyright 2021 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_PROFILING_H
#define DPUSYSCORE_PROFILING_H

/**
 * @file profiling.h
 * @brief Code section profiling management.
 */

#include <attributes.h>
#include <limits.h>
#include <profiling_internals.h>

#define PROFILING_RESET_VALUE (UINT32_MAX)

/**
 * @def PROFILING_INIT
 * @hideinitializer
 * @brief Declare and initialize a profiling context associated to the given name.
 */
#define PROFILING_INIT(_name)                                                                                                    \
    __section(".dpu_profiling") dpu_profiling_t _name = {                                                                        \
        .start = { [0 ...(NR_THREADS - 1)] = PROFILING_RESET_VALUE },                                                            \
        .count = { [0 ...(NR_THREADS - 1)] = 0 },                                                                                \
    }

#ifdef DPU_PROFILING
/**
 * @fn profiling_start
 * @brief Start profiling a code section.
 *
 * This function saves the perfcounter current value in the profiling context.
 *
 * @param context the profiling context to use.
 */
void
profiling_start(dpu_profiling_t *context);

/**
 * @fn profiling_stop
 * @brief Stop profiling a code section.
 *
 * This function gets the perfcounter current value and computes the number of cyles spent in the code section.
 * The profiling_start function must be called beforehand.
 *
 * @param context the profiling context to use.
 */
void
profiling_stop(dpu_profiling_t *context);
#else
#define profiling_start(context)                                                                                                 \
    do {                                                                                                                         \
    } while (0)
#define profiling_stop(context)                                                                                                  \
    do {                                                                                                                         \
    } while (0)
#endif /* DPU_PROFILING */

#endif /* DPUSYSCORE_PROFILING_H */
