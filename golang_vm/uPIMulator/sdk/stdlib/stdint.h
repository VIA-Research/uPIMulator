/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef _DPUSYSCORE_STDINT_H_
#define _DPUSYSCORE_STDINT_H_

/**
 * @file stdint.h
 * @brief Provides abstraction over machine types.
 */

/* Exact integer types */

/* Signed */

/**
 * @brief A signed 8-bit value.
 */
typedef signed char int8_t;
/**
 * @brief A signed 16-bit value.
 */
typedef short int int16_t;
/**
 * @brief A signed 32-bit value.
 */
typedef int int32_t;
/**
 * @brief A signed 64-bit value.
 */
typedef long long int int64_t;

/* Unsigned */

/**
 * @brief An unsigned 8-bit value.
 */
typedef unsigned char uint8_t;
/**
 * @brief An unsigned 16-bit value.
 */
typedef unsigned short int uint16_t;
/**
 * @brief An unsigned 32-bit value.
 */
typedef unsigned int uint32_t;

/**
 * @brief An unsigned 64-bit value.
 */
typedef unsigned long int uint64_t;

/* Small types */

/* Signed */

/**
 * @brief A signed value on at least 8 bits.
 */
typedef signed char int_least8_t;
/**
 * @brief A signed value on at least 16 bits.
 */
typedef short int int_least16_t;
/**
 * @brief A signed value on at least 32 bits.
 */
typedef int int_least32_t;
/**
 * @brief A signed value on at least 64 bits.
 */
typedef long int int_least64_t;

/* Unsigned */

/**
 * @brief An unsigned value on at least 8 bits.
 */
typedef unsigned char uint_least8_t;
/**
 * @brief An unsigned value on at least 16 bits.
 */
typedef unsigned short int uint_least16_t;
/**
 * @brief An unsigned value on at least 32 bits.
 */
typedef unsigned int uint_least32_t;
/**
 * @brief An unsigned value on at least 64 bits.
 */
typedef unsigned long int uint_least64_t;

/* Fast types */

/* Signed */

/**
 * @brief A signed value on at least 8 bits, optimized for that length.
 */
typedef signed char int_fast8_t;
/**
 * @brief A signed value on at least 16 bits, optimized for that length.
 */
typedef int int_fast16_t;
/**
 * @brief A signed value on at least 32 bits, optimized for that length.
 */
typedef int int_fast32_t;
/**
 * @brief A signed value on at least 64 bits, optimized for that length.
 */
typedef long int int_fast64_t;

/* Unsigned */

/**
 * @brief An unsigned value on at least 8 bits, optimized for that length.
 */
typedef unsigned char uint_fast8_t;
/**
 * @brief An unsigned value on at least 16 bits, optimized for that length.
 */
typedef unsigned int uint_fast16_t;
/**
 * @brief An unsigned value on at least 32 bits, optimized for that length.
 */
typedef unsigned int uint_fast32_t;
/**
 * @brief An unsigned value on at least 64 bits, optimized for that length.
 */
typedef unsigned long int uint_fast64_t;

/* Types for void* pointers */

/**
 * @brief A signed value which can contain a pointer value.
 */
typedef int intptr_t;
/**
 * @brief An unsigned value which can contain a pointer value.
 */
typedef unsigned int uintptr_t;

/* Greatest-width integer types */

/**
 * @brief A signed value which can contain all signed values.
 */
typedef long long int intmax_t;
/**
 * @brief An unsigned value which can contain all unsigned values.
 */
typedef unsigned long long int uintmax_t;

#include <limits.h>

#endif /* _DPUSYSCORE_STDINT_H_ */
