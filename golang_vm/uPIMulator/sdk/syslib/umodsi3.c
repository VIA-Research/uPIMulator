/* ===-- umodsi3.c - Implement __umodsi3 -----------------------------------===
 *
 *                     The LLVM Compiler Infrastructure
 *
 * This file is dual licensed under the MIT and the University of Illinois Open
 * Source Licenses. See LICENSE_LLVM.TXT for details.
 *
 * ===----------------------------------------------------------------------===
 *
 * This file implements __umodsi3 for the compiler_rt library.
 *
 * ===----------------------------------------------------------------------===
 */

#include "int_lib.h"

/* Returns: a % b */

extern unsigned long
__udiv32(unsigned int, unsigned int);

COMPILER_RT_ABI su_int
__umodsi3(su_int a, su_int b)
{
    unsigned long res = __udiv32(a, b);
    return (unsigned int)res;
}
