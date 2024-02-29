/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <string.h>
#include <stdbool.h>

// TODO Possible optimization:
//  use of the table of indexation that indicates if the character should or should not be accepted/rejected
//  in this case we will need 128 bits (as there are 128 ascii characters) ( = 16 bytes = 4 words ) per runtime
// => 4x24 = 96 words of 32 bits.
//
//  TODO Another solution would be to stock this table only temporarily with the allocation function, but, currently, it's not an
//  option.

size_t
strspn(const char *string, const char *accept)
{
    size_t prefix_length;

    for (prefix_length = 0; string[prefix_length] != '\0'; ++prefix_length) {
        char c = string[prefix_length];

        unsigned int accept_index = 0;
        while (true) {
            char a = accept[accept_index];

            if (c == a) {
                break;
            }

            if (a == '\0') {
                return prefix_length;
            }

            accept_index++;
        }
    }

    return prefix_length;
}