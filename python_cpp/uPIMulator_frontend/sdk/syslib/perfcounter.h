/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_PERFCOUNTER_H
#define DPUSYSCORE_PERFCOUNTER_H

#include <stdint.h>
#include <stdbool.h>

/**
 * @file perfcounter.h
 * @brief Utilities concerning the performance counter register.
 *
 */

/**
 * @typedef perfcounter_t
 * @brief A value which can be stored by the performance counter.
 */
typedef uint64_t perfcounter_t;

/**
 * @enum perfcounter_config_t
 * @brief A configuration for the performance counter, defining what should be counted.
 *
 * @var COUNT_SAME          keep the previous configuration
 * @var COUNT_CYCLES        switch to counting clock cycles
 * @var COUNT_INSTRUCTIONS  switch to counting executed instructions
 * @var COUNT_NOTHING       does not count anything
 */
typedef enum _perfcounter_config_t {
    COUNT_SAME = 0,
    COUNT_CYCLES = 1,
    COUNT_INSTRUCTIONS = 2,
    COUNT_NOTHING = 3,
} perfcounter_config_t;

/**
 * @def CLOCKS_PER_SEC
 * @hideinitializer
 * @brief A number used to convert the value returned by the perfcounter_get and perfcounter_config functions into seconds,
 *        when counting clock cycles.
 */
extern const volatile uint32_t CLOCKS_PER_SEC;

/**
 * @fn perfcounter_get
 * @brief Fetch the value of the performance counter register.
 *
 * @return The current value of the performance counter register, or undefined if perfcounter_config has not been called before.
 */
perfcounter_t
perfcounter_get(void);

#ifndef DPU_PROFILING
/**
 * @fn perfcounter_config
 * @brief Configure the performance counter behavior.
 *
 * This function cannot be used when profiling an application.
 *
 * @param config        The new behavior for the performance counter register
 * @param reset_value   Whether the performance counter register should be set to 0
 *
 * @return The current value of the performance counter register, or undefined if perfcounter_config has not been called before.
 */
perfcounter_t
perfcounter_config(perfcounter_config_t config, bool reset_value);
#else
#define perfcounter_config(config, reset_value)                                                                                  \
    do {                                                                                                                         \
    } while (0)
#endif /* !DPU_PROFILING */

#endif /* DPUSYSCORE_PERFCOUNTER_H */
