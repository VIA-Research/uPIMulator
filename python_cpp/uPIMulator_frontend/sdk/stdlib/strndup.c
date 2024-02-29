/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include "string.h"
#include "buddy_alloc.h"

char *
strndup(const char *string, size_t n)
{
    size_t length = strnlen(string, n);
    char *result = buddy_alloc(length + 1);

    if (result != NULL) {
        memcpy(result, string, length);
        ((char *)result)[length + 1] = '\0';
    }

    return result;
}