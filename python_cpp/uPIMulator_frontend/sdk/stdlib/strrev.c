/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <string.h>

char *
strrev(char *string)
{
    size_t length = strlen(string);

    for (size_t each_char = 0; each_char < length / 2; ++each_char) {
        char c = string[each_char];
        string[each_char] = string[length - each_char - 1];
        string[length - each_char - 1] = c;
    }

    return string;
}
