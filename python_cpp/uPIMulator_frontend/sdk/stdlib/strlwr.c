/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <ctype.h>

char *
strlwr(char *string)
{
    char *ptr = string;
    char c;

    while ((c = *ptr) != '\0') {
        *ptr = tolower(c);
        ptr++;
    }

    return string;
}
