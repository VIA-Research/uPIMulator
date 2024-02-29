/*===-- udivmodsi4.c - Implement __udivmodsi4 ------------------------------===
 *
 *                    The LLVM Compiler Infrastructure
 *
 * This file is dual licensed under the MIT and the University of Illinois Open
 * Source Licenses. See LICENSE_LLVM.TXT for details.
 *
 * ===----------------------------------------------------------------------===
 *
 * This file implements __udivmodsi4 for the compiler_rt library.
 *
 * ===----------------------------------------------------------------------===
 */

#include "int_lib.h"

/* Returns: a / b, *rem = a % b  */

extern unsigned long
__udiv32(unsigned int, unsigned int);

COMPILER_RT_ABI su_int
__udivmodsi4(su_int a, su_int b, su_int *rem)
{
    unsigned long res = __udiv32(a, b);
    *rem = (unsigned int)res;
    return (unsigned int)(res >> 32);
}
