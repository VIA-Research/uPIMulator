/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef _DPUSYSCORE_LIMITS_H_
#define _DPUSYSCORE_LIMITS_H_

#define SCHAR_MAX (0x0000007f)
#define SHRT_MAX (0x00007fff)
#define INT_MAX (0x7fffffff)
#define LONG_MAX (0x7fffffffffffffffl)
#define LLONG_MAX (0x7fffffffffffffffl)

#define SCHAR_MIN (-SCHAR_MAX - 1)
#define SHRT_MIN (-SHRT_MAX - 1)
#define INT_MIN (-INT_MAX - 1)
#define LONG_MIN (-LONG_MAX - 1)
#define LLONG_MIN (-LLONG_MAX - 1)

#define UCHAR_MAX (SCHAR_MAX * 2 + 1)
#define USHRT_MAX (SHRT_MAX * 2 + 1)
#define UINT_MAX (INT_MAX * 2U + 1U)
#define ULONG_MAX (LONG_MAX * 2UL + 1UL)
#define ULLONG_MAX (LLONG_MAX * 2UL + 1UL)

#ifdef __CHAR_UNSIGNED__ /* -funsigned-char */
#define CHAR_MIN 0
#define CHAR_MAX UCHAR_MAX
#else
#define CHAR_MIN SCHAR_MIN
#define CHAR_MAX SCHAR_MAX
#endif

/* The maximum number of bytes in a multi-byte character.  */
#define MB_LEN_MAX 16

/* Limits of integral types */

/**
 * @def CHAR_BIT
 * @hideinitializer
 * @brief The number of bits in a char type.
 */
#define CHAR_BIT (8)

/**
 * @def WORD_BIT
 * @hideinitializer
 * @brief The number of bits in a word type.
 */
#define WORD_BIT (32)

/**
 * @def LONG_BIT
 * @hideinitializer
 * @brief The number of bits in a pseudo-long type.
 */
#define LONG_BIT (32)

/* Minimum of signed integral types */

/**
 * @def INT8_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>int8_t</code>.
 */
#define INT8_MIN (-0x7f - 1)

/**
 * @def INT16_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>int16_t</code>.
 */
#define INT16_MIN (-0x7fff - 1)

/**
 * @def INT32_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>int32_t</code>.
 */
#define INT32_MIN (-0x7fffffff - 1)

/**
 * @def INT64_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>int64_t</code>.
 */
#define INT64_MIN (-0x7fffffffffffffffL - 1L)

/* Maximum of signed integral types */

/**
 * @def INT8_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>int8_t</code>.
 */
#define INT8_MAX (0x7f)
/**
 * @def INT16_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>int16_t</code>.
 */
#define INT16_MAX (0x7fff)
/**
 * @def INT32_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>int32_t</code>.
 */
#define INT32_MAX (0x7fffffff)
/**
 * @def INT64_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>int64_t</code>.
 */
#define INT64_MAX (0x7fffffffffffffffL)

/* Maximum of unsigned integral types */

/**
 * @def UINT8_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>uint8_t</code>.
 */
#define UINT8_MAX (0xff)
/**
 * @def UINT16_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>uint16_t</code>.
 */
#define UINT16_MAX (0xffff)
/**
 * @def UINT32_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>uint32_t</code>.
 */
#define UINT32_MAX (0xffffffff)
/**
 * @def UINT64_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>uint64_t</code>.
 */
#define UINT64_MAX (0xffffffffffffffffUL)

/* Minimum of signed integral types having a minimum size */

/**
 * @def INT_LEAST8_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>int_least8_t</code>.
 */
#define INT_LEAST8_MIN (-0x7f - 1)
/**
 * @def INT_LEAST16_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>int_least16_t</code>.
 */
#define INT_LEAST16_MIN (-0x7fff - 1)
/**
 * @def INT_LEAST32_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>int_least32_t</code>.
 */
#define INT_LEAST32_MIN (-0x7fffffff - 1)
/**
 * @def INT_LEAST64_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>int_least64_t</code>.
 */
#define INT_LEAST64_MIN (-0x7fffffffffffffffL - 1L)

/* Maximum of signed integral types having a minimum size */

/**
 * @def INT_LEAST8_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>int_least8_t</code>.
 */
#define INT_LEAST8_MAX (0x7f)
/**
 * @def INT_LEAST16_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>int_least16_t</code>.
 */
#define INT_LEAST16_MAX (0x7fff)
/**
 * @def INT_LEAST32_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>int_least32_t</code>.
 */
#define INT_LEAST32_MAX (0x7fffffff)
/**
 * @def INT_LEAST64_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>int_least64_t</code>.
 */
#define INT_LEAST64_MAX (0x7fffffffffffffffL)

/* Maximum of unsigned integral types having a minimum size */

/**
 * @def UINT_LEAST8_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>uint_least8_t</code>.
 */
#define UINT_LEAST8_MAX (0xff)
/**
 * @def UINT_LEAST16_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>uint_least16_t</code>.
 */
#define UINT_LEAST16_MAX (0xffff)
/**
 * @def UINT_LEAST32_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>uint_least32_t</code>.
 */
#define UINT_LEAST32_MAX (0xffffffff)
/**
 * @def UINT_LEAST64_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>uint_least64_t</code>.
 */
#define UINT_LEAST64_MAX (0xffffffffffffffffUL)

/* Minimum of fast signed integral types having a minimum size */

/**
 * @def INT_FAST8_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>int_fast8_t</code>.
 */
#define INT_FAST8_MIN (-0x7f - 1)
/**
 * @def INT_FAST16_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>int_fast16_t</code>.
 */
#define INT_FAST16_MIN (-0x7fffffff - 1)
/**
 * @def INT_FAST32_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>int_fast32_t</code>.
 */
#define INT_FAST32_MIN (-0x7fffffff - 1)
/**
 * @def INT_FAST64_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>int_fast64_t</code>.
 */
#define INT_FAST64_MIN (-0x7fffffffffffffffL - 1L)

/* Maximum of fast signed integral types having a minimum size */

/**
 * @def INT_FAST8_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>int_fast8_t</code>.
 */
#define INT_FAST8_MAX (0x7f)
/**
 * @def INT_FAST16_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>int_fast16_t</code>.
 */
#define INT_FAST16_MAX (0x7fffffff)
/**
 * @def INT_FAST32_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>int_fast32_t</code>.
 */
#define INT_FAST32_MAX (0x7fffffff)
/**
 * @def INT_FAST64_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>int_fast64_t</code>.
 */
#define INT_FAST64_MAX (0x7fffffffffffffffL)

/* Maximum of fast unsigned integral types having a minimum size */

/**
 * @def UINT_FAST8_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>uint_fast8_t</code>.
 */
#define UINT_FAST8_MAX (0xff)
/**
 * @def UINT_FAST16_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>uint_fast16_t</code>.
 */
#define UINT_FAST16_MAX (0xffffffffU)
/**
 * @def UINT_FAST32_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>uint_fast32_t</code>.
 */
#define UINT_FAST32_MAX (0xffffffffU)
/**
 * @def UINT_FAST64_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>uint_fast64_t</code>.
 */
#define UINT_FAST64_MAX (0xffffffffffffffffUL)

/* Limits for integral types holding void* pointers */

/**
 * @def INTPTR_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>intptr_t</code>.
 */
#define INTPTR_MIN (-0x7fffffff - 1)
/**
 * @def INTPTR_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>intptr_t</code>.
 */
#define INTPTR_MAX (0x7fffffff)
/**
 * @def UINTPTR_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>uintptr_t</code>.
 */
#define UINTPTR_MAX (0xffffffffU)

/* Limits of greatest-width integer types */

/**
 * @def INTMAX_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>intmax_t</code>.
 */
#define INTMAX_MIN (-0x7fffffffffffffffLL - 1)
/**
 * @def INTMAX_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>intmax_t</code>.
 */
#define INTMAX_MAX (0x7fffffffffffffffLL)
/**
 * @def UINTMAX_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>uintmax_t</code>.
 */
#define UINTMAX_MAX (0xffffffffffffffffULL)

/* Limits of others integer types */

/**
 * @def PTRDIFF_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>ptrdiff_t</code>.
 * @see ptrdiff_t
 */
#define PTRDIFF_MIN (-0x7fffffff - 1)
/**
 * @def PTRDIFF_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>ptrdiff_t</code>.
 * @see ptrdiff_t
 */
#define PTRDIFF_MAX (0x7fffffff)

/**
 * @def SIZE_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>size_t</code>.
 * @see size_t
 */
#define SIZE_MAX (0xffffffffU)

/**
 * @def WCHAR_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>wchar_t</code>.
 * @see wchar_t
 */
#define WCHAR_MIN (-0x7fffffff - 1)
/**
 * @def WCHAR_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>wchar_t</code>.
 * @see wchar_t
 */
#define WCHAR_MAX (0x7fffffff)

/**
 * @def WINT_MIN
 * @hideinitializer
 * @brief The minimum value for a value of type <code>wint_t</code>.
 */
#define WINT_MIN (0u)
/**
 * @def WINT_MAX
 * @hideinitializer
 * @brief The maximum value for a value of type <code>wint_t</code>.
 */
#define WINT_MAX (0xffffffffu)

/* Macros for integer constant expressions */

/* Signed */

/**
 * @def INT8_C
 * @hideinitializer
 * @brief Expands the value to an expression corresponding to the type <code>int_least8_t</code>
 */
#define INT8_C(value) value
/**
 * @def INT16_C
 * @hideinitializer
 * @brief Expands the value to an expression corresponding to the type <code>int_least16_t</code>
 */
#define INT16_C(value) value
/**
 * @def INT32_C
 * @hideinitializer
 * @brief Expands the value to an expression corresponding to the type <code>int_least32_t</code>
 */
#define INT32_C(value) value
/**
 * @def INT64_C
 * @hideinitializer
 * @brief Expands the value to an expression corresponding to the type <code>int_least64_t</code>
 */
#define INT64_C(value) value##LL

/* Unsigned */

/**
 * @def UINT8_C
 * @hideinitializer
 * @brief Expands the value to an expression corresponding to the type <code>uint_least8_t</code>
 */
#define UINT8_C(value) value##U
/**
 * @def UINT16_C
 * @hideinitializer
 * @brief Expands the value to an expression corresponding to the type <code>uint_least16_t</code>
 */
#define UINT16_C(value) value##U
/**
 * @def UINT32_C
 * @hideinitializer
 * @brief Expands the value to an expression corresponding to the type <code>uint_least32_t</code>
 */
#define UINT32_C(value) value##U
/**
 * @def UINT64_C
 * @hideinitializer
 * @brief Expands the value to an expression corresponding to the type <code>uint_least64_t</code>
 */
#define UINT64_C(value) value##ULL

/* Maximum types */

/**
 * @def INTMAX_C
 * @hideinitializer
 * @brief Expands the value to an expression corresponding to the type <code>intmax_t</code>
 */
#define INTMAX_C(value) value##LL
/**
 * @def UINTMAX_C
 * @hideinitializer
 * @brief Expands the value to an expression corresponding to the type <code>uintmax_t</code>
 */
#define UINTMAX_C(value) value##ULL

#endif /* _DPUSYSCORE_LIMITS_H_ */
