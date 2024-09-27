/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include "macro_utils.h"

/* clang-format off */
#define __RESTORE_CARRY_AND_ZERO_FLAG(x) \
        "add r0, zero, 0x00000001; ld d0, id8, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 * 0); stop true, 0\n" /* ... restore Z = 0, C = 0 */ \
        "add r0, mneg, 0x80000001; ld d0, id8, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 * 0); stop true, 0\n" /* ... restore Z = 0, C = 1 */ \
        "add r0, zero, 0x00000000; ld d0, id8, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 * 0); stop true, 0\n" /* ... restore Z = 1, C = 0 */ \
        "add r0, mneg, 0x80000000; ld d0, id8, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 * 0); stop true, 0\n" // ... restore Z = 1, C = 1 */

#define RESTORE_CARRY_AND_ZERO_FLAG \
        "  lw r0,  id4, " __STR(NR_ATOMIC_BITS) " + (" __STR(NR_THREADS) " * 8 * 12)\n" \
        "  add r1, r0,  r0\n" \
        "  add r0, r0,  r1\n" /* r0 =  3 * r0 (each line of ending_routines is 3 instructions) */ \
        "  or  r1, id8, 0 \n" \
        "  add r1, id4, r1\n" /* r1 = 12 * id (there are 12 instructions per runtime in ending_routines) */ \
        "  add r0, r0,  r1\n" /* r0 = r0 + r1 (compute the offset (in number of instructions) to jump to) */ \
        "  call zero, r0, ending_routines\n" \
        "ending_routines:\n" \
        __FOR_EACH_THREAD(__RESTORE_CARRY_AND_ZERO_FLAG)

/* clang-format on */
