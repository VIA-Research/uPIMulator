/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/*
 * 64x64 signed division.
 *
 * This is the actual libcall implementation, as requested by the compiler.
 */
#include <stdint.h>
extern uint64_t
__udiv64(uint64_t dividend, uint64_t divider, int ask_remainder);

int64_t
__moddi3(int64_t dividend, int64_t divider)
{
    if (dividend >= 0) {
        if (divider >= 0) {
            return __udiv64(dividend, divider, 1);
        } else {
            return __udiv64(dividend, -divider, 1);
        }
    } else if (divider >= 0) {
        // Negative dividend, positive divider
        return -__udiv64(-dividend, divider, 1);
    } else {
        // Negative dividend, negative divider
        return -__udiv64(-dividend, -divider, 1);
    }
}
