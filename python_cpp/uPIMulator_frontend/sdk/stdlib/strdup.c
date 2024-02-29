/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include "string.h"
#include "buddy_alloc.h"

char *
strdup(const char *string)
{
    size_t length = strlen(string) + 1; // we get the length of the string for memory allocation

    char *result = buddy_alloc(length); // we allocate length+1 bytes for the duplicate

    if (result != NULL) {
        memcpy(result, string, length); // we copy length bytes from string to the duplicate
    }

    return result;
}