/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <stddef.h>
#include <stdbool.h>
#include <ctype.h>

int
atoi(const char *nptr)
{
    int result = 0;
    bool is_positive = true;

    if (nptr == NULL) {
        return result;
    }

    while (isspace(*nptr)) {
        nptr++;
    }

    if (*nptr == '-') {
        is_positive = false;
        nptr++;
    } else if (*nptr == '+') {
        nptr++;
    }

    for (;; nptr++) {
        unsigned int digit = *nptr - '0';

        if (digit > 9) {
            break;
        }

        result = (10 * result) + digit;
    }

    return is_positive ? result : -result;
}