/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef _DPUSYSCORE_STDIO_H_
#define _DPUSYSCORE_STDIO_H_

#include <attributes.h>
#include <stddef.h>

/**
 * @file stdio.h
 * @brief Standard input/output library functions.
 */

/**
 * @def STDOUT_BUFFER_INIT
 * @hideinitializer
 * @brief Declares the stdout buffer. Should be used as when declaring a global variable.
 * @param size the size of the stdout buffer. Must be a multiple of 8, and greater than 0.
 */
#define STDOUT_BUFFER_INIT(size)                                                                                                 \
    _Static_assert((size >= 8) && (((size)&7) == 0), "stdout buffer size must be a multiple of 8 and > 0");                      \
    unsigned char __dma_aligned __mram_noinit __stdout_buffer[(size)];                                                           \
    const unsigned int __stdout_buffer_size = (size);

/**
 * @fn printf
 * @brief Writes the formatted data in the stdout buffer.
 *
 * This function has a prototype close to the one of the standard printf function.
 * However, the format string comply to the java.util.Formatter format, which is
 * similar to the printf format, but not quit exactly the same. Date formatter may
 * produce interpreted results, but they will probably be incorrect. Every format
 * specifier should reference one and only one of the variadic argument (eg. "%n"
 * is not supported).
 *
 * There is no compile-time check to verify that the format is correct: any other
 * character in the format string will not be interpreted.
 *
 * @param format how the logged data should be formatted
 * @param ... the different data to be printed
 */
void __attribute__((format(printf, 1, 2))) printf(const char *restrict format, ...);

/**
 * @fn puts
 * @brief Writes the string in the stdout buffer. A newline character is appended to the output.
 * @param str the null-terminated string to be written
 */
void
puts(const char *str);

/**
 * @fn putchar
 * @brief Writes the character in the stdout buffer.
 * @param c the character to be written
 */
void
putchar(int c);

#endif /* _DPUSYSCORE_STDIO_H_ */
