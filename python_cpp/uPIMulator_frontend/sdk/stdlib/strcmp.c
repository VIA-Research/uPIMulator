/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <string.h>

int
strcmp(const char *string1, const char *string2)
{
    unsigned char c1 = *string1;
    unsigned char c2 = *string2;

    while (c1 != '\0') {
        if (c1 - c2 != 0) {
            return c1 - c2;
        }

        string1++;
        string2++;
        c1 = *string1;
        c2 = *string2;
    }

    return c1 - c2;
}