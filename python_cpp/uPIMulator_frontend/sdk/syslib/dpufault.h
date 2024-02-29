/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_DPUFAULT_H
#define DPUSYSCORE_DPUFAULT_H

// A list of "fault codes"
#define __FAULT_ALLOC_HEAP_FULL__ 1
#define __FAULT_DIVISION_BY_ZERO__ 2
#define __FAULT_ASSERT_FAILED__ 3
// Used in the compiler to implement a trap
#define __FAULT_HALT__ 4
#define __FAULT_PRINTF_OVERFLOW__ 5
#define __FAULT_ALREADY_PROFILING__ 6
#define __FAULT_NOT_PROFILING__ 7

#endif
