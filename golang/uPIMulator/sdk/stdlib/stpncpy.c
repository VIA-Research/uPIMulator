/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <stddef.h>

char *
stpncpy(char *destination, const char *source, size_t size)
{
    char c = *source;
    size_t each_byte;

    for (each_byte = 0; each_byte < size; ++each_byte) {
        if (c == '\0') {
            char *null_char_ptr = destination;

            for (; each_byte < size; ++each_byte) {
                *destination = '\0';
                destination++;
            }

            return null_char_ptr;
        }

        *destination = c;
        destination++;
        source++;
        c = *source;
    }

    return destination;
}