/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <stddef.h>

size_t
strlen(const char *string)
{
    const char *ptr = string;

    while (*ptr != '\0') {
        ptr++;
    }

    return ptr - string;
}
