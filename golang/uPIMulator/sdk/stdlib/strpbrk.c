/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <string.h>
#include <stdbool.h>

char *
strpbrk(const char *string, const char *accept)
{
    string += strcspn(string, accept);
    return *string ? (char *)string : NULL;
}