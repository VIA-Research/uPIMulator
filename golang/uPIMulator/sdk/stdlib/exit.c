/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */
#include "defs.h"
#include "stdlib.h"

#define unreachable() __builtin_unreachable()

void
exit(int __attribute__((unused)) status)
{
    __asm__ volatile("stop true, __sys_end");
    unreachable();
}
