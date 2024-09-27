/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <stddef.h>

void *
memchr(const void *area, int character, size_t size)
{
    const char *ptr = (const char *)area;

    for (size_t each_byte = 0; each_byte < size; ++each_byte) {
        if (ptr[each_byte] == character) {
            return (void *)(ptr + each_byte);
        }
    }

    return NULL;
}