/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

char *
stpcpy(char *destination, const char *source)
{
    char c = *source;

    while (c != '\0') {
        *destination = c;
        destination++;
        source++;
        c = *source;
    }

    *destination = c;

    return destination;
}