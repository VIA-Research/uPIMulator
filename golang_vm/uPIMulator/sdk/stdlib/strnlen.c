/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <stddef.h>

size_t
strnlen(const char *string, size_t max_len)
{
    size_t len = 0;

    while ((string[len] != '\0') && len < max_len) {
        len++;
    }

    return len;
}
