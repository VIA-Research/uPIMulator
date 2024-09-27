/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <string.h>

char *
strcat(char *destination, const char *source)
{
    size_t length = strlen(destination);
    unsigned int i;

    for (i = 0; source[i] != '\0'; i++) {
        destination[length + i] = source[i];
    }

    destination[length + i] = '\0';

    return destination;
}