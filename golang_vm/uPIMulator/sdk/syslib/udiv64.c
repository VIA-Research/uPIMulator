/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/*
 * 64x64 multiplication unsigned division.
 */
#include <stdint.h>
#include <dpuruntime.h>

static unsigned int
__clz__(uint64_t x)
{
    return __builtin_clzl(x);
}

uint64_t
__udiv64(uint64_t dividend, uint64_t divider, int ask_remainder)
{
    uint64_t dxo = dividend, dxe = 0;

    if (divider == 0)
        goto division_by_zero;
    if (divider > dividend) {
        if (ask_remainder == 0)
            return 0;
        else
            return dividend;
    }

    // Mimic the div_step.
    /// div_step functionality:
    //   if (Dxo >= (Ra<< #u5)) {
    //     Dxo = Dxo - (Ra<< #u5);
    //     Dxe = (Dxe << 1) | 1;
    //   } else {
    //     Dxe =  Dxe << 1;
    //   }
    int dividerl0 = __clz__(divider), dividendl0 = __clz__(dividend);

    int i = dividerl0 - dividendl0;

    for (; i >= 0; i--) {
        uint64_t pivot = ((uint64_t)divider << i);
        if (dxo >= pivot) {
            dxo = dxo - pivot;
            dxe = ((uint64_t)dxe << 1) | 1L;
        } else {
            dxe = (uint64_t)dxe << 1;
        }
    }
    if (ask_remainder == 1)
        return dxo;
    else
        return dxe;

division_by_zero:
    __asm__ volatile("fault " __STR(__FAULT_DIVISION_BY_ZERO__));
    unreachable();
}
