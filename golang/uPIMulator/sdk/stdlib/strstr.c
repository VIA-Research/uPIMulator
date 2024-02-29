/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include "stdbool.h"
#include "string.h"
#include "buddy_alloc.h"

char *
strstr(const char *haystack, const char *needle)
{
    char *current_needle = (char *)needle;
    char *start_haystack = (char *)haystack;
    char *current_haystack = start_haystack;

    while (true) {
        char needle_char = *current_needle;

        if (needle_char == '\0') {
            return start_haystack;
        }

        char haystack_char = *current_haystack;

        if (haystack_char == needle_char) {
            current_haystack++;
            current_needle++;
        } else if (haystack_char == '\0') {
            return NULL;
        } else {
            current_needle = (char *)needle;
            current_haystack = ++start_haystack;
        }
    }
}
