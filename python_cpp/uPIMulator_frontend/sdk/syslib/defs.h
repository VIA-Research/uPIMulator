/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_DEFS_H
#define DPUSYSCORE_DEFS_H

#include <sysdef.h>
#include <dpufault.h>
#include <macro_utils.h>
#include <attributes.h>

/**
 * @file defs.h
 * @brief Miscellaneous system functions.
 *
 * General purpose definitions.
 */

#if __STDC_VERSION__ >= 201112L
#define __ATTRIBUTE_NO_RETURN__ _Noreturn
#else
#define __ATTRIBUTE_NO_RETURN__
#endif /* __STDC_VERSION__ */

/**
 * @fn me
 * @internal This just returns the value of the special register id.
 * @return The current tasklet's sysname.
 */
static inline sysname_t
me()
{
    return __builtin_dpu_tid();
}

/**
 * @fn halt
 * @brief Halts the DPU.
 * @throws FAULT_HALT always
 */
__ATTRIBUTE_NO_RETURN__ static inline void
halt()
{
    __builtin_trap();
    __builtin_unreachable();
}

/**
 * @fn check_stack
 * @return the number of unused 32-bits words in the current runtime's stack.
 *         If the number is negative, it indicates by how much 32-bits words the stack overflowed.
 *
 * @internal This fetches the position of the next stack in memory from the Stack Pointer Table (cf. tasklet.h)
 *           and compute the remaining bytes.
 */
int
check_stack();

#endif /* DPUSYSCORE_DEFS_H */
