/* ===-- modsi3.c - Implement __modsi3 -------------------------------------===
 *
 *                     The LLVM Compiler Infrastructure
 *
 * This file is dual licensed under the MIT and the University of Illinois Open
 * Source Licenses. See LICENSE_LLVM.TXT for details.
 *
 * ===----------------------------------------------------------------------===
 *
 * This file implements __modsi3 for the compiler_rt library.
 *
 * ===----------------------------------------------------------------------===
 */

#include "int_lib.h"

/* Returns: a % b */

extern long
__div32(int, int);

COMPILER_RT_ABI si_int
__modsi3(si_int a, si_int b)
{
    long res = __div32(a, b);
    return (int)res;
}
