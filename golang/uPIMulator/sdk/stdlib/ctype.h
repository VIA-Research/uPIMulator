/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef _DPUSYSCORE_CTYPE_H_
#define _DPUSYSCORE_CTYPE_H_

/**
 * @file ctype.h
 * @brief Provides useful functions for testing and mapping characters.
 */

/**
 * @brief Checks whether the specified character is a digit.

 * @param c an unsigned char or EOF
 * @return Whether the character is a digit (using 0 as false and anything else as true).
 */
static inline int
isdigit(int c)
{
    return (c >= '0') && (c <= '9');
}

/**
 * @brief Checks whether the specified character is a lowercase letter.

 * @param c an unsigned char or EOF
 * @return Whether the character is a digit (using 0 as false and anything else as true).
 */
static inline int islower(c) { return (c >= 'a') && (c <= 'z'); }

/**
 * @brief Checks whether the specified character is an uppercase letter.

 * @param c an unsigned char or EOF
 * @return Whether the character is an uppercase letter (using 0 as false and anything else as true).
 */
static inline int isupper(c) { return (c >= 'A') && (c <= 'Z'); }

/**
 * @brief Checks whether the specified character is a letter.

 * @param c an unsigned char or EOF
 * @return Whether the character is a letter (using 0 as false and anything else as true).
 */
static inline int
isalpha(int c)
{
    return islower(c) || isupper(c);
}

/**
 * @brief Checks whether the specified character is a letter or a digit.

 * @param c an unsigned char or EOF
 * @return Whether the character is a letter or a digit (using 0 as false and anything else as true).
 */
static inline int
isalnum(int c)
{
    return isalpha(c) || isdigit(c);
}

/**
 * @brief Checks whether the specified character is a control character.

 * @param c an unsigned char or EOF
 * @return Whether the character is a control character (using 0 as false and anything else as true).
 */
static inline int
iscntrl(int c)
{
    return (c <= 0x1f) || (c == 0x7f);
}

/**
 * @brief Checks whether the specified character is printable.

 * @param c an unsigned char or EOF
 * @return Whether the character is printable (using 0 as false and anything else as true).
 */
static inline int
isprint(int c)
{
    return !iscntrl(c);
}

/**
 * @brief Checks whether the specified character has graphical representation using locale.

 * @param c an unsigned char or EOF
 * @return Whether the character has graphical representation using locale (using 0 as false and anything else as true).
 */
static inline int
isgraph(int c)
{
    return isprint(c) && (c != ' ');
}

/**
 * @brief Checks whether the specified character is a punctuation character.

 * @param c an unsigned char or EOF
 * @return Whether the character is a punctuation character (using 0 as false and anything else as true).
 */
static inline int
ispunct(int c)
{
    return (c >= '!' && c <= '/') || (c >= ':' && c <= '@') || (c >= '[' && c <= '`') || (c >= '{' && c <= '~');
}

/**
 * @brief Checks whether the specified character is a white-space.

 * @param c an unsigned char or EOF
 * @return Whether the character is a white-space (using 0 as false and anything else as true).
 */
static inline int
isspace(int c)
{
    return (c >= 0x9 && c <= 0xd) || (c == ' ');
}

/**
 * @brief Checks whether the specified character is a hexadecimal digit.

 * @param c an unsigned char or EOF
 * @return Whether the character is a hexadecimal digit (using 0 as false and anything else as true).
 */
static inline int
isxdigit(int c)
{
    return isdigit(c) || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f');
}

/**
 * @brief Checks whether the specified character is a blank character.

 * @param c an unsigned char or EOF
 * @return Whether the character is a blank character (using 0 as false and anything else as true).
 */
static inline int
isblank(int c)
{
    return c == ' ' || c == '\t';
}

/**
 * @brief Converts the specified character to a lowercase letter if it is a letter;

 * @param c an unsigned char or EOF
 * @return The lowercase letter corresponding to the character, if it is a letter. The initial character otherwise.
 */
static inline int
tolower(int c)
{
    return isupper(c) ? (c + 0x20) : c;
}

/**
 * @brief Converts the specified character to a uppercase letter if it is a letter;

 * @param c an unsigned char or EOF
 * @return The uppercase letter corresponding to the character, if it is a letter. The initial character otherwise.
 */
static inline int
toupper(int c)
{
    return islower(c) ? (c - 0x20) : c;
}

#endif /* _DPUSYSCORE_CTYPE_H_ */
