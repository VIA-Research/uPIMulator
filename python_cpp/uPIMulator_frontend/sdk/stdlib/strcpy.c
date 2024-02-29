/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <string.h>

char *
strcpy(char *destination, const char *source)
{
    char *ptr = destination;
    char c = *source;

    while (c != '\0') {
        *ptr = c;
        ptr++;
        source++;
        c = *source;
    }

    *ptr = c;

    return destination;
}