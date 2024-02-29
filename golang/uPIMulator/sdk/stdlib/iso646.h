/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef _DPUSYSCORE_ISO646_H_
#define _DPUSYSCORE_ISO646_H_

/**
 * @file iso646.h
 * @brief Alternative spellings for operators not supported by the ISO646 standard character set.
 */

/**
 * @def and
 * @brief Logical AND.
 */
#define and &&
/**
 * @def and_eq
 * @brief Bitwise AND accumulation.
 */
#define and_eq &=
/**
 * @def bitand
 * @brief Bitwise AND.
 */
#define bitand &
/**
 * @def bitor
 * @brief Bitwise OR.
 */
#define bitor |
/**
 * @def compl
 * @brief Bitwise NOT.
 */
#define compl ~
/**
 * @def not
 * @brief Logical NOT.
 */
#define not !
/**
 * @def not_eq
 * @brief Difference.
 */
#define not_eq !=
/**
 * @def or
 * @brief Logical OR.
 */
#define or ||
/**
 * @def or_eq
 * @brief Bitwise OR accumulation.
 */
#define or_eq |=
/**
 * @def xor
 * @brief Bitwise XOR.
 */
#define xor ^
/**
 * @def xor_eq
 * @brief Bitwise XOR accumulation.
 */
#define xor_eq ^=

#endif /* _DPUSYSCORE_ISO646_H_ */
