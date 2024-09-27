/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <stddef.h>

char *
strncpy(char *destination, const char *source, size_t size)
{
    char *ptr = destination;
    char c = *source;
    size_t each_byte;

    for (each_byte = 0; each_byte < size; ++each_byte) {
        if (c == '\0') {
            for (; each_byte < size; ++each_byte) {
                *ptr = '\0';
                ptr++;
            }

            return destination;
        }

        *ptr = c;
        ptr++;
        source++;
        c = *source;
    }

    return destination;
}