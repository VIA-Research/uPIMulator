/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <stddef.h>

char *
strrchr(char *string, int character)
{
    char *pos = NULL;
    char *ptr = string;
    unsigned char c;

    do {
        c = *ptr;
        if (c == character) {
            pos = ptr;
        }
        ptr++;
    } while (c != '\0');

    return pos;
}
