/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include "string.h"

char *
strtok_r(char *str, const char *delim, char **saveptr)
{
    char *end;

    if (str == NULL) {
        str = *saveptr;
    }

    if (*str == '\0') {
        *saveptr = str;
        return NULL;
    }

    str += strspn(str, delim);

    if (*str == '\0') {
        *saveptr = str;
        return NULL;
    }

    end = str + strcspn(str, delim);
    if (*end == '\0') {
        *saveptr = end;
        return str;
    }

    *end = '\0';
    *saveptr = end + 1;
    return str;
}
