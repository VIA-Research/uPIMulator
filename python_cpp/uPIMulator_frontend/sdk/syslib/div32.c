/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <dpuruntime.h>

void __attribute__((naked, noinline, no_instrument_function)) __udiv32(void)
{
    __asm__ volatile("  "__ADD_PROFILING_ENTRY__
                     "  clz r3, r1, max, __udiv32_division_by_zero\n" // r3 = by how many the divider can be shifted on 32-bit
                     "  clz r4, r0\n" // r4 = number of useless bits of the dividend
                     "  sub r3, r4, r3, gtu, __udiv32_result_0\n" // r3 = the maximal shift to be done
                     "  move r4, r1\n"
                     "  move.u d0, r0\n"
                     "  jump r3, __udiv32_base\n" // As we will jump backward relatively to __udiv32_base
                     "  div_step d0, r4, d0, 31\n"
                     "  div_step d0, r4, d0, 30\n"
                     "  div_step d0, r4, d0, 29\n"
                     "  div_step d0, r4, d0, 28\n"
                     "  div_step d0, r4, d0, 27\n"
                     "  div_step d0, r4, d0, 26\n"
                     "  div_step d0, r4, d0, 25\n"
                     "  div_step d0, r4, d0, 24\n"
                     "  div_step d0, r4, d0, 23\n"
                     "  div_step d0, r4, d0, 22\n"
                     "  div_step d0, r4, d0, 21\n"
                     "  div_step d0, r4, d0, 20\n"
                     "  div_step d0, r4, d0, 19\n"
                     "  div_step d0, r4, d0, 18\n"
                     "  div_step d0, r4, d0, 17\n"
                     "  div_step d0, r4, d0, 16\n"
                     "  div_step d0, r4, d0, 15\n"
                     "  div_step d0, r4, d0, 14\n"
                     "  div_step d0, r4, d0, 13\n"
                     "  div_step d0, r4, d0, 12\n"
                     "  div_step d0, r4, d0, 11\n"
                     "  div_step d0, r4, d0, 10\n"
                     "  div_step d0, r4, d0, 9\n"
                     "  div_step d0, r4, d0, 8\n"
                     "  div_step d0, r4, d0, 7\n"
                     "  div_step d0, r4, d0, 6\n"
                     "  div_step d0, r4, d0, 5\n"
                     "  div_step d0, r4, d0, 4\n"
                     "  div_step d0, r4, d0, 3\n"
                     "  div_step d0, r4, d0, 2\n"
                     "  div_step d0, r4, d0, 1\n"
                     "__udiv32_base:\n"
                     "  div_step d0, r4, d0, 0\n"
                     "__udiv32_exit:\n"
                     "  jump r23\n"
                     "__udiv32_result_0:\n"
                     "  move.u d0, r0, true, __udiv32_exit\n"
                     "__udiv32_division_by_zero:\n"
                     "  fault "__STR(__FAULT_DIVISION_BY_ZERO__));
}

void __attribute__((naked, noinline, no_instrument_function)) __div32(void)
{
    __asm__ volatile("  "__ADD_PROFILING_ENTRY__
                     "sd r22, 0, d22\n"
                     "add r22, r22, 8\n"
                     // The quotient's sign depends on the sign of the dividend and divider... After few tries it sounds
                     // like the quickest way to select the operators is to branch according to the cases.
                     "  clo r3, r0, z, __div32_pos_dividend\n"
                     "  clo r3, r1, z, __div32_neg_dividend_pos_divider\n"
                     "__div32_neg_dividend_neg_divider:\n" // As a result, the quotient is positive and the remainder negative
                     "  neg r0, r0\n"
                     "  neg r1, r1\n"
                     "  call r23, __udiv32\n"
                     "  neg r1, r1, true, __div32_exit\n"
                     "__div32_neg_dividend_pos_divider:\n" // As a result, the quotient is negative and the remainder negative
                     "  neg r0, r0\n"
                     "  call r23, __udiv32\n"
                     "  neg r1, r1\n"
                     "  neg r0, r0, true, __div32_exit\n"
                     "__div32_pos_dividend:\n"
                     "  clo r3, r1, z, __div32_pos_dividend_pos_divider\n"
                     "  neg r1, r1\n" // As a result, the quotient is negative and the remainder positive
                     "  call r23, __udiv32\n"
                     "  neg r0, r0, true, __div32_exit\n"
                     "__div32_pos_dividend_pos_divider:\n" // The dividend and divider are both positive
                     "  call r23, __udiv32\n"
                     "__div32_exit:\n"
                     "  ld d22, r22, -8\n"
                     "  jump r23\n");
}
