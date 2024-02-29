/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <dpuruntime.h>

int __attribute__((noinline)) __mulsi3(int a, int b)
{
    int dest;
    __asm__ volatile("  jgtu %2, %1, __mulsi3_swap\n"
                     "  move r2, %1\n"
                     "  move r0, %2, true, __mulsi3_start\n"
                     "__mulsi3_swap:\n"
                     "  move r2, %2\n"
                     "  move r0, %1\n"
                     "__mulsi3_start:\n"
                     "  move r1, zero\n"
                     "  mul_step d0, r2, d0, 0 , z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 1 , z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 2 , z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 3 , z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 4 , z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 5 , z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 6 , z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 7 , z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 8 , z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 9 , z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 10, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 11, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 12, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 13, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 14, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 15, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 16, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 17, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 18, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 19, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 20, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 21, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 22, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 23, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 24, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 25, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 26, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 27, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 28, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 29, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 30, z, __mulsi3_exit\n"
                     "  mul_step d0, r2, d0, 31, z, __mulsi3_exit\n"
                     "__mulsi3_exit:\n"
                     "  move %0, r1\n"
                     : "=r"(dest)
                     : "r"(a), "r"(b));
    return dest;
}
