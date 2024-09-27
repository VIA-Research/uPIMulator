/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_DEVPRIVATE_H
#define DPUSYSCORE_DEVPRIVATE_H

/**
 * @file devprivate.h
 * @brief Reserved for internal use ... please do not use those functions unless you know exactly what you do.
 */

/**
 * @def tell
 * @brief On a simulator, injects a tell instruction to print out developer debug info.
 * @nolink
 *
 * @warning This function will not work on a target different from simulator.
 *
 * @param reg a register
 * @param val a constant value
 */
#define tell(reg, val) __asm__("tell %[r], " val : : [r] "r"(reg) :)

#endif /* DPUSYSCORE_DEVPRIVATE_H */
