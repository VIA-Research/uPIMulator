/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <string.h>

char *
strncat(char *destination, const char *source, size_t size)
{
    size_t length = strlen(destination);
    size_t i;

    for (i = 0; (i < size) && (source[i] != '\0'); i++) {
        destination[length + i] = source[i];
    }

    destination[length + i] = '\0';

    return destination;
}