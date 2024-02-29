/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <stddef.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdarg.h>
#include <string.h>
#include <atomic_bit.h>
#include <attributes.h>
#include <mram.h>
#include <dpuruntime.h>

#define DEFAULT_STDOUT_BUFFER_SIZE (1 << 20)

unsigned char __weak __mram_noinit __stdout_buffer[DEFAULT_STDOUT_BUFFER_SIZE];
unsigned int __weak __stdout_buffer_size = DEFAULT_STDOUT_BUFFER_SIZE;

/* __lower_data: needed to make sure that the structure address will be less that a signed12
 *               (sd endian:e ra off:s12 imm:s16 used in bootstrap).
 *
 * __dma_aligned: needed to make sure that the structure address will be aligned on 8 bytes (for sd in bootstrap as well).
 *
 * This structure is initialize at zero in the bootsrap
 */
__lower_data(__STR(__STDOUT_BUFFER_STATE)) __dma_aligned struct {
    uint32_t wp;
    uint32_t has_wrapped;
} __STDOUT_BUFFER_STATE;

static uint32_t __stdout_buffer_write_pointer_initial;
static uint32_t __stdout_nr_of_wrapping;

#define STDOUT_CACHE_BUFFER_SIZE 8
_Static_assert((STDOUT_CACHE_BUFFER_SIZE >= 8) && (STDOUT_CACHE_BUFFER_SIZE <= 2048) && (STDOUT_CACHE_BUFFER_SIZE % 8 == 0),
    "STDOUT_CACHE_BUFFER_SIZE needs to be a multiple of 8 in ]0; 2048]");

static char __stdout_cache_buffer[STDOUT_CACHE_BUFFER_SIZE] __dma_aligned;
static unsigned int __stdout_cache_write_index;

ATOMIC_BIT_INIT(__stdout_buffer_lock);

__attribute__((noinline)) static void
__transfer_cache_to_mram()
{
    __mram_ptr void *offset_in_mram = (__mram_ptr void *)(__STDOUT_BUFFER_STATE.wp + (uintptr_t)__stdout_buffer);

    __STDOUT_BUFFER_STATE.wp += STDOUT_CACHE_BUFFER_SIZE;
    if (__STDOUT_BUFFER_STATE.wp >= __stdout_buffer_size) {
        __STDOUT_BUFFER_STATE.wp = 0;
        __STDOUT_BUFFER_STATE.has_wrapped = true;
        __stdout_nr_of_wrapping++;
    }

    mram_write(__stdout_cache_buffer, offset_in_mram, STDOUT_CACHE_BUFFER_SIZE);
}

// Generic template that will be used everywhere: cache a byte and flush to MRAM
// when the cache is full.
__attribute__((noinline)) static void
__write_byte_and_flush_if_needed(uint8_t byte)
{
    __stdout_cache_buffer[__stdout_cache_write_index++] = byte;
    if (__stdout_cache_write_index == STDOUT_CACHE_BUFFER_SIZE) {
        __transfer_cache_to_mram();
        __stdout_cache_write_index = 0;
    }
}

__attribute__((noinline)) static void
__finalized_print_sequence()
{
    memset(__stdout_cache_buffer + __stdout_cache_write_index, 0, STDOUT_CACHE_BUFFER_SIZE - __stdout_cache_write_index);
    __transfer_cache_to_mram();

    if (__stdout_nr_of_wrapping > 1
        || (__stdout_nr_of_wrapping == 1 && __STDOUT_BUFFER_STATE.wp > __stdout_buffer_write_pointer_initial))
        __asm__("fault " __STR(__FAULT_PRINTF_OVERFLOW__)); // need to throw fault because we will not be able to print the buffer
}

__attribute__((noinline)) static void
__open_print_sequence()
{
    ATOMIC_BIT_ACQUIRE(__stdout_buffer_lock);
    __stdout_cache_write_index = 0;
    __stdout_nr_of_wrapping = 0;
    __stdout_buffer_write_pointer_initial = __STDOUT_BUFFER_STATE.wp;
}

/* Nothing else that the release instruction should be in this function in order to make sure that the print routine in complete
 * at this point*/
__attribute__((noinline)) static void
__close_print_sequence()
{
    ATOMIC_BIT_RELEASE(__stdout_buffer_lock);
}

void
printf(const char *restrict format, ...)
{
    bool insert_string_arg = true;
    bool insert_string_arg_end_character = false;
    char *current_format_char_ptr = (char *)format;

    __open_print_sequence();

    va_list args;
    va_start(args, format);

    for (; *current_format_char_ptr != '\0'; ++current_format_char_ptr) {
        if (*current_format_char_ptr == '%') {
            ++current_format_char_ptr;
            if (*current_format_char_ptr == '\0')
                break;
            if (*current_format_char_ptr == '%')
                goto standard_character_format_process;

            __write_byte_and_flush_if_needed('%');

            while (*current_format_char_ptr != '\0') {
                if (*current_format_char_ptr == 'l') {
                    __write_byte_and_flush_if_needed(*current_format_char_ptr);
                    ++current_format_char_ptr;
                    continue;
                }
                if ((*current_format_char_ptr == 'L') || (*current_format_char_ptr == 'z')) {
                    ++current_format_char_ptr;
                    continue;
                }
                if (*current_format_char_ptr == 'i') {
                    __write_byte_and_flush_if_needed('d');
                    break;
                }
                __write_byte_and_flush_if_needed(*current_format_char_ptr);

                if (((*current_format_char_ptr >= 'A') && (*current_format_char_ptr <= 'Z'))
                    || ((*current_format_char_ptr >= 'a') && (*current_format_char_ptr <= 'z')))
                    break;

                ++current_format_char_ptr;
            }

            insert_string_arg = true;

        } else {
        standard_character_format_process:
            if (insert_string_arg) {
                __write_byte_and_flush_if_needed('%');
                __write_byte_and_flush_if_needed('s');
                insert_string_arg = false;
            }
        }
    }

    __write_byte_and_flush_if_needed('\0');
    current_format_char_ptr = (char *)format;

    for (; *current_format_char_ptr != '\0'; ++current_format_char_ptr) {
        if (*current_format_char_ptr == '%') {
            ++current_format_char_ptr;

            if (*current_format_char_ptr == '\0')
                break;
            if (*current_format_char_ptr == '%')
                goto standard_character_process;

            if (insert_string_arg_end_character) {
                insert_string_arg_end_character = false;
                __write_byte_and_flush_if_needed('\0');
            }

            bool arg_is_64_bits = false;

            while (*current_format_char_ptr != '\0') {
                if ((*current_format_char_ptr == 'l') || (*current_format_char_ptr == 'L')) {
                    arg_is_64_bits = true;
                    current_format_char_ptr++;
                    continue;
                } else if (*current_format_char_ptr == 'z') {
                    current_format_char_ptr++;
                    continue;
                }

                if (((*current_format_char_ptr >= 'A') && (*current_format_char_ptr <= 'Z'))
                    || ((*current_format_char_ptr >= 'a') && (*current_format_char_ptr <= 'z')))
                    break;

                ++current_format_char_ptr;
            }

            switch (*current_format_char_ptr) {
                case 's': {
                    char *arg = (char *)va_arg(args, int);
                    while (*arg != '\0') {
                        __write_byte_and_flush_if_needed(*arg);
                        arg++;
                    }
                    __write_byte_and_flush_if_needed('\0');
                    break;
                }
                case 'c': {
                    char arg_as_char = (char)va_arg(args, int);
                    __write_byte_and_flush_if_needed(arg_as_char);
                    break;
                }
                case 'f':
                case 'e':
                case 'E':
                case 'g':
                case 'G': {
                    __asm__ volatile("nop");
                    double val = va_arg(args, double);
                    char *arg = (char *)&val;
                    for (int i = 0; i < 8; i++) {
                        char arg_byte = arg[i];
                        __write_byte_and_flush_if_needed(arg_byte);
                    }
                    break;
                }
                default: {
                    unsigned int arg_size_in_bytes;
                    long val;

                    if (arg_is_64_bits) {
                        val = va_arg(args, long);
                        arg_size_in_bytes = 8;
                    } else {
                        val = (long)va_arg(args, int);
                        arg_size_in_bytes = 4;
                    }

                    char *arg = (char *)&val;
                    for (unsigned int i = 0; i < arg_size_in_bytes; i++) {
                        char arg_byte = arg[i];
                        __write_byte_and_flush_if_needed(arg_byte);
                    }
                }
            }
        } else {
        standard_character_process:
            __write_byte_and_flush_if_needed(*current_format_char_ptr);
            insert_string_arg_end_character = true;
        }
    }

    if (insert_string_arg_end_character) {
        __write_byte_and_flush_if_needed('\0');
    }

    va_end(args);

    __finalized_print_sequence();
    __close_print_sequence();
}

void
puts(const char *str)
{
    __open_print_sequence();

    __write_byte_and_flush_if_needed('%');
    __write_byte_and_flush_if_needed('s');
    __write_byte_and_flush_if_needed('\0');

    for (char *current_char_ptr = (char *)str; *current_char_ptr != '\0'; current_char_ptr++) {
        __write_byte_and_flush_if_needed(*current_char_ptr);
    }

    __write_byte_and_flush_if_needed('\n');
    __write_byte_and_flush_if_needed('\0');

    __finalized_print_sequence();
    __close_print_sequence();
}

void
putchar(int c)
{
    __open_print_sequence();

    __write_byte_and_flush_if_needed('%');
    __write_byte_and_flush_if_needed('c');
    __write_byte_and_flush_if_needed('\0');

    char arg_as_char = (char)c;
    __write_byte_and_flush_if_needed(arg_as_char);

    __finalized_print_sequence();
    __close_print_sequence();
}
