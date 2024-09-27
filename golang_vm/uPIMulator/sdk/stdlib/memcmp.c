/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <stddef.h>

int
memcmp(const void *area1, const void *area2, size_t size)
{
    const unsigned char *ptr1 = (const unsigned char *)area1;
    const unsigned char *ptr2 = (const unsigned char *)area2;

    for (size_t each_byte = 0; each_byte < size; ++each_byte) {
        int diff = ptr1[each_byte] - ptr2[each_byte];
        if (diff != 0) {
            return diff;
        }
    }

    return 0;
}