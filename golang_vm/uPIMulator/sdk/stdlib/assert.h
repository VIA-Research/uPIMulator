/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef _DPUSYSCORE_ASSERT_H_
#define _DPUSYSCORE_ASSERT_H_

/**
 * @file assert.h
 * @brief Provides a way to verify assumptions with <code>assert</code>.
 */

#define static_assert _Static_assert

#ifdef NDEBUG

/**
 * @def assert
 * @hideinitializer
 * @brief When NDEBUG is defined, <code>assert</code> is not available and calling it will do nothing.
 */
#define assert(ignore) ((void)0)

#else

#include <dpufault.h>
#include <macro_utils.h>

/**
 * @def assert
 * @hideinitializer
 * @brief Verify the assumption of the specified expression, resulting in a fault if it fails.
 *
 * @param expression the assumption to verify
 * @throws FAULT_ASSERT_FAILED when the assertion failed
 * @todo add a diagnostic message to the log, if it exists, when the assertion fails
 */
#define assert(expression)                                                                                                       \
    do {                                                                                                                         \
        if (!(expression)) {                                                                                                     \
            __asm__ volatile("fault " __STR(__FAULT_ASSERT_FAILED__));                                                           \
        }                                                                                                                        \
    } while (0)

#endif /* NDEBUG */

#endif /* _DPUSYSCORE_ASSERT_H_ */
