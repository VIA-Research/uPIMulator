/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/*
 * 64x64 unsigned remainder.
 *
 * This is the actual libcall implementation, as requested by the compiler.
 */
#include <stdint.h>
extern uint64_t
__udiv64(uint64_t dividend, uint64_t divider, int ask_remainder);

uint64_t
__umoddi3(uint64_t dividend, uint64_t divider)
{
    return __udiv64(dividend, divider, 1);
}
