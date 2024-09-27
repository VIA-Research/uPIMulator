/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <string.h>

int
strncmp(const char *string1, const char *string2, size_t size)
{
    for (size_t len = 0; len < size; ++len) {
        unsigned char c1 = string1[len];
        unsigned char c2 = string2[len];

        if (((c1 - c2) != 0) || (c1 == '\0')) {
            return c1 - c2;
        }
    }

    return 0;
}
