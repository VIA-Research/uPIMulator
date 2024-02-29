/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_MACRO_UTILS_H
#define DPUSYSCORE_MACRO_UTILS_H

/**
 * @file macro_utils.h
 * @brief Provide utility macros.
 */

#define __STR(x) __STR_AGAIN(x)
#define __STR_AGAIN(x) #x

#define __CONCAT(x, y) __CONCAT_AGAIN(x, y)
#define __CONCAT_AGAIN(x, y) x##y

#define __REPEAT_0(x)
#define __REPEAT_1(x) x(0) __REPEAT_0(x)
#define __REPEAT_2(x) x(1) __REPEAT_1(x)
#define __REPEAT_3(x) x(2) __REPEAT_2(x)
#define __REPEAT_4(x) x(3) __REPEAT_3(x)
#define __REPEAT_5(x) x(4) __REPEAT_4(x)
#define __REPEAT_6(x) x(5) __REPEAT_5(x)
#define __REPEAT_7(x) x(6) __REPEAT_6(x)
#define __REPEAT_8(x) x(7) __REPEAT_7(x)
#define __REPEAT_9(x) x(8) __REPEAT_8(x)
#define __REPEAT_10(x) x(9) __REPEAT_9(x)
#define __REPEAT_11(x) x(10) __REPEAT_10(x)
#define __REPEAT_12(x) x(11) __REPEAT_11(x)
#define __REPEAT_13(x) x(12) __REPEAT_12(x)
#define __REPEAT_14(x) x(13) __REPEAT_13(x)
#define __REPEAT_15(x) x(14) __REPEAT_14(x)
#define __REPEAT_16(x) x(15) __REPEAT_15(x)
#define __REPEAT_17(x) x(16) __REPEAT_16(x)
#define __REPEAT_18(x) x(17) __REPEAT_17(x)
#define __REPEAT_19(x) x(18) __REPEAT_18(x)
#define __REPEAT_20(x) x(19) __REPEAT_19(x)
#define __REPEAT_21(x) x(20) __REPEAT_20(x)
#define __REPEAT_22(x) x(21) __REPEAT_21(x)
#define __REPEAT_23(x) x(22) __REPEAT_22(x)
#define __REPEAT_24(x) x(23) __REPEAT_23(x)
#define __FOR_EACH_THREAD(x) __CONCAT(__REPEAT_, NR_THREADS)(x)

#endif /* DPUSYSCORE_MACRO_UTILS_H */
