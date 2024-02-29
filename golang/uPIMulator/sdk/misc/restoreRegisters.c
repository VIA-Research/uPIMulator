/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/*
 * The "restore registers" program, is used by debugging processes to restore every registers of every runtime.
 * The program should be booted once on runtime 0.
 */

#include "restore_carry_and_zero_flag.h"

void __attribute__((naked, used, section(".text.__bootstrap"))) __bootstrap()
{
    /* clang-format off */
    __asm__ volatile(
        "  jeq id, " __STR(NR_THREADS) " - 1, .+2\n"
        "  boot id, 1\n"
        "  ld d2,  id8, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 *  1)\n"
        "  ld d4,  id8, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 *  2)\n"
        "  ld d6,  id8, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 *  3)\n"
        "  ld d8,  id8, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 *  4)\n"
        "  ld d10, id8, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 *  5)\n"
        "  ld d12, id8, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 *  6)\n"
        "  ld d14, id8, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 *  7)\n"
        "  ld d16, id8, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 *  8)\n"
        "  ld d18, id8, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 *  9)\n"
        "  ld d20, id8, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 * 10)\n"
        "  ld d22, id8, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 * 11)\n"
        "  jnz id, atomic_done\n"
        "  move r0, " __STR(NR_ATOMIC_BITS) " - 1\n"
        "atomic_loop:\n"
        "  lbu r1, r0, 0\n"
        "  jz r1, atomic_release\n"
        "  acquire r0, 0, true, atomic_next\n"
        "atomic_release:\n"
        "  release r0, 0, nz, atomic_next\n"
        "atomic_next:\n"
        "  add r0, r0, -1, pl, atomic_loop\n"
        "atomic_done:\n"
        RESTORE_CARRY_AND_ZERO_FLAG
    );
    /* clang-format on */
}
