/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_STDLIB_H
#define DPUSYSCORE_STDLIB_H

/**
 * @file stdlib.h
 * @brief Elementary standard C functions: calls the system function halt.
 */

#include <stddef.h>
#include <attributes.h>

/**
 * @def EXIT_FAILURE
 * @hideinitializer
 * @brief Unsuccessful termination for exit().
 */
#define EXIT_FAILURE 1

/**
 * @def EXIT_SUCCESS
 * @hideinitializer
 * @brief Successful termination for exit().
 */
#define EXIT_SUCCESS 0

/**
 * @brief Aborts the DPU execution triggering a processor fault.
 */
__NO_RETURN void
abort(void);

/**
 * @brief Terminates the invoking tasklet, returning the specified status.
 */
__NO_RETURN void
exit(int status);

/**
 * @brief Get an environment variable, or NULL. In the DPU case, always NULL.
 */
static inline char *
getenv(__attribute__((unused)) const char *name)
{
    return NULL;
}

/**
 * @brief Returns the absolute value of the argument.
 */
static inline int
abs(int x)
{
    return (x < 0) ? -x : x;
}

/**
 * @brief Returns the absolute value of the argument.
 */
static inline long int
labs(long int x)
{
    return (x < 0) ? -x : x;
}

/**
 * @brief Returns the absolute value of the argument.
 */
static inline long long int
llabs(long long int x)
{
    return (x < 0) ? -x : x;
}

typedef struct {
    int quot;
    int rem;
} div_t;

typedef struct {
    long int quot;
    long int rem;
} ldiv_t;

typedef struct {
    long long int quot;
    long long int rem;
} lldiv_t;

static inline div_t
div(int numer, int denom)
{
    div_t result = { numer / denom, numer % denom };
    return result;
}

static inline ldiv_t
ldiv(long int numer, long int denom)
{
    ldiv_t result = { numer / denom, numer % denom };
    return result;
}

static inline lldiv_t
lldiv(long long int numer, long long int denom)
{
    lldiv_t result = { numer / denom, numer % denom };
    return result;
}

/**
 * @brief Converts a string to an integer
 *
 * Function converts the initial part of the string in nptr to an integer value. The string may begin
 * with an arbitrary amount of white space followed by a single optional '+' or '-' sign.
 *
 * Conversion stops at the first character not representing a digit. If an underflow occurs, atoi()
 * returns INT_MIN. If an overflow occurs, atoi() returns INT_MAX. In both cases errno is set to ERANGE.
 *
 *
 * @param nptr string that contains an integer in a string format
 * @return the result of conversion unless the value would overflow or underflow.
 */
int
atoi(const char *nptr);

/**
 * @brief Converts a string to a long integer (64 bits)
 *
 * Function converts the initial part of the string in nptr to a long value. The string may begin
 * with an arbitrary amount of white space followed by a single optional '+' or '-' sign.
 *
 * Conversion stops at the first character not representing a digit. If an underflow occurs, atol()
 * returns LONG_MIN. If an overflow occurs, atol() returns LONG_MAX. In both cases errno is set to ERANGE.
 *
 *
 * @param nptr string that contains an integer in a string format
 * @return the result of conversion unless the value would overflow or underflow.
 */
long
atol(const char *nptr);

///**
// * @brief Converts a string to a long integer (64 bits) according to the given base
// * between 2 and 36 inclusive, or be the special value 0
// *
// * TODO : If the given base is oustide of the range [2...36], then errno is set to EINVAL
// *
// * Function converts the initial part of the string in nptr to a long value. The string may begin
// * with an arbitrary amount of white space followed by a single optional '+' or '-' sign.
// *
// * If base is zero or 16, the string may then include a "0x" prefix, and the number will be read
// * in base 16; otherwise, a zero base is taken as 10 (decimal) unless the next character is '0',
// * in which case it is taken as 8 (octal).
// *
// * Conversion stops at the first character not representing a digit in the given base.
// * Accepted digits are : in bases above 10, the letter 'A' in either uppercase or lowercase
// * represents 10, 'B' represents 11, and so forth, with 'Z' representing 35.
// *
// * If endptr is not NULL, strtol() stores the address of the first invalid character in *endptr.
// * If there were no digits at all, strtol() stores the original value of nptr in *endptr (and
// * returns 0). In particular, if *nptr is not '\0' but **endptr is '\0' on return, the entire
// * string is valid.
// *
// * If an underflow occurs, atol() returns LONG_MIN. If an overflow occurs, atol()
// * returns LONG_MAX. TODO!!! In both cases errno is set to ERANGE.
// *
// * @param nptr string that contains an integer in a string format
// * @param endptr
// * @param base
//
// * @return the result of conversion unless the value would overflow or underflow.
//*/
// long int strtol(const char *nptr, char **endptr, int base);
// TODO : strtol() doesn't work : has to be written in assembly

#endif /* DPUSYSCORE_STDLIB_H */
