/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <string.h>

char *
strsep(char **stringp, const char *delim)
{
    if (*stringp == NULL) {
        return NULL;
    }

    char *original = *stringp;
    char *delim_ptr = strpbrk(*stringp, delim);

    if (delim_ptr == NULL) {
        *stringp = NULL;
    } else {
        *delim_ptr = '\0';
        *stringp = delim_ptr + 1;
    }

    return original;
}