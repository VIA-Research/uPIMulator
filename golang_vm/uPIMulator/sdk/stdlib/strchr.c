/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <string.h>

char *
strchr(const char *string, int character)
{
    char *str = (char *)string;
    unsigned char c = *str;

    while (1) {
        if (c == character) {
            return str;
        }
        if (c == '\0') {
            return NULL;
        }

        str++;
        c = *str;
    }
}
