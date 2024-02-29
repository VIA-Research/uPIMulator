/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_STRING_H
#define DPUSYSCORE_STRING_H

/**
 * @file string.h
 * @brief Provides functions to manipulate arrays of characters.
 */

#include <stddef.h>

/**
 * @brief Computes the length of the given null-terminated string.
 *
 * @param string the string for which we want the length.
 * @return The length of the string, not including the null character.
 */
size_t
strlen(const char *string);

/**
 * @brief Computes the length of the given null-terminated string if this length is less than <code>max_len</code>.
 * strnlen checks at most <code>max_len</code> bytes and returns <code>max_len</code> if it has read as many bytes.
 *
 * @param string the string for which we want to find the length.
 * @param max_len maximum number of bytes to check
 * @return The length of the string, not including the null character or <code>max_len</code> if null character
 * wasn't found.
 */
size_t
strnlen(const char *string, size_t max_len);

/**
 * @brief Compares the first <code>size</code> bytes of <code>area1</code> and <code>area2</code>.
 *
 * @param area1 the pointer to the start of the first area of the comparison
 * @param area2 the pointer to the start of the second area of the comparison
 * @param size the number of bytes to compare between each area.
 * @return <code>0</code> if the areas are the same, a non-zero value otherwise.
 */
int
memcmp(const void *area1, const void *area2, size_t size);

/**
 * @brief Compares the two null-terminated strings <code>string1</code> and <code>string2</code>.
 *
 * @param string1 the pointer to the start of the first string of the comparison
 * @param string2 the pointer to the start of the second string of the comparison
 * @return <code>0</code> if the strings are the same, a non-zero value otherwise.
 */
int
strcmp(const char *string1, const char *string2);

/**
 * @brief Compares the first <code>size</code> bytes of the two null-terminated strings <code>string1</code> and
 * <code>string2</code>.
 *
 * @param string1 the pointer to the start of the first string of the comparison
 * @param string2 the pointer to the start of the second string of the comparison
 * @param size the maximum number of bytes to compare between each string.
 * @return <code>0</code> if the strings are the same, a non-zero value otherwise.
 */
int
strncmp(const char *string1, const char *string2, size_t size);

/**
 * @brief Set the first <code>size</code> bytes of <code>area</code> at <code>value</code>.
 *
 * @param area the pointer to the start of the area to set
 * @param value the value at which the area is set
 * @param size the number of bytes being set
 * @return A pointer to the start of the set area.
 */
void *
memset(void *area, int value, size_t size);

/**
 * @brief Search for the first occurrence of <code>character</code> in the first <code>size</code> bytes of <code>area</code>.
 *
 * @param area the pointer to the start of the area to search
 * @param character the value to search for
 * @param size the number of bytes to search
 * @return A pointer to the first occurrence, if it exists, <code>NULL</code> otherwise.
 */
void *
memchr(const void *area, int character, size_t size);

/**
 * @brief Concatenate the string <code>source</code> after the string <code>destination</code>.
 *
 * @param destination the pointer to the start of the first string
 * @param source the pointer to the start of the second string
 * @return A pointer to the start of the concatenated string.
 */
char *
strcat(char *destination, const char *source);

/**
 * @brief Concatenate the first <code>size</code> bytes of the string <code>source</code> after the string
 * <code>destination</code>.
 *
 * @param destination the pointer to the start of the first string
 * @param source the pointer to the start of the second string
 * @param size the maximum number of bytes to concatenate
 * @return A pointer to the start of the concatenated string.
 */
char *
strncat(char *destination, const char *source, size_t size);

/**
 * @brief Search for the first occurrence of <code>character</code> in the <code>string</code>.
 *
 * @param string the pointer to the start of the string to search
 * @param character the value to search for
 * @return A pointer to the first occurrence, if it exists, <code>NULL</code> otherwise.
 */
char *
strchr(const char *string, int character);

/**
 * @brief Search for the last occurrence of <code>character</code> in the <code>string</code>.
 *
 * @param string the pointer to the start of the string to search
 * @param character the value to search for
 * @return A pointer to the last occurrence, if it exists, <code>NULL</code> otherwise.
 */
char *
strrchr(const char *string, int character);

/**
 * @brief Copy <code>size</code> bytes from <code>source</code> into <code>destination</code>.
 *
 * @warning This function is not safe for overlapping memory blocks.
 *
 * @param destination the pointer to the start of the destination of the copy
 * @param source the pointer to the start of the area to copy
 * @param size the number of bytes to copy
 * @return A pointer to the start of the copied area.
 */
void *
memcpy(void *destination, const void *source, size_t size);

/**
 * @brief Copy <code>size</code> bytes from <code>source</code> into <code>destination</code>.
 *
 * This is a safer method than <code>memcpy</code> for overlapping memory blocks.
 *
 * @see memcpy
 * @param destination the pointer to the start of the destination of the copy
 * @param source the pointer to the start of the area to copy
 * @param size the number of bytes to copy
 * @return A pointer to the start of the copied area.
 */
void *
memmove(void *destination, const void *source, size_t size);

/**
 * @brief Copy the string <code>source</code> into the string<code>destination</code>.
 *
 * @warning This function is not safe for overlapping strings.
 *
 * @param destination the pointer to the start of the destination of the copy
 * @param source the pointer to the start of the string to copy
 * @return A pointer to the start of the copied string.
 */
char *
strcpy(char *destination, const char *source);

/**
 * @brief Copy <code>size</code> bytes from the string <code>source</code> into the string<code>destination</code>.
 *
 * @warning This function is not safe for overlapping strings.
 *
 * @param destination the pointer to the start of the destination of the copy
 * @param source the pointer to the start of the string to copy
 * @param size the maximum number of bytes to copy
 * @return A pointer to the start of the copied string.
 */
char *
strncpy(char *destination, const char *source, size_t size);

/**
 * @def strxfrm
 * @hideinitializer
 * @brief Transform the first <code>size</code> bytes of the string <code>source</code> into current locale and place them in the
 * string <code>destination</code>.
 *
 * There is no concept of "locale" in the DPU, implying that the related functions behave as native, "locale-less", functions.
 * This function is just a synonym of <code>strncpy</code>.
 *
 * @warning This function is not safe for overlapping strings.
 *
 * @see strncpy
 * @param destination the pointer to the start of the destination of the copy
 * @param source the pointer to the start of the string to copy
 * @param size the maximum number of bytes to copy
 * @return A pointer to the start of the copied string.
 */
#define strxfrm strncpy

/**
 * @def strcoll
 * @hideinitializer
 * @brief Compare two null-terminated strings using the current locale.
 *
 * There is no concept of "locale" in the DPU, implying that the related functions behave as native, "locale-less", functions.
 * This function is just a synonym of <code>strcmp</code>.
 *
 * @see strcmp
 * @param string1 the pointer to the start of the first string of the comparison
 * @param string2 the pointer to the start of the second string of the comparison
 * @return <code>0</code> if the strings are the same, a non-zero value otherwise.
 */
#define strcoll strcmp

/**
 * @brief Copy the string <code>source</code> into the string <code>destination</code>.
 *
 * @warning This function is not safe for overlapping strings.
 *
 * @param destination the pointer to the start of the destination of the copy
 * @param source the pointer to the start of the string to copy
 * @return A pointer to the end (the address of the terminating null byte)
 * of the copied string.
 */
char *
stpcpy(char *destination, const char *source);

/**
 * @brief Copy <code>size</code> bytes from the string <code>source</code> into the string<code>destination</code>.
 *
 * @warning This function is not safe for overlapping strings.
 * If <code>size</code> is less than the length of the <code>source</code>, then
 * the remaining characters in <code>destination</code> will be filled with '\0'
 *
 * @param destination the pointer to the start of the destination of the copy
 * @param source the pointer to the start of the string to copy
 * @param size the maximum number of bytes to copy
 * @return A pointer to the end (the address of the terminating null byte)
 * of the copied string.
 */
char *
stpncpy(char *destination, const char *source, size_t size);

/**
 * @brief Converts every character of a null-terminated string into lowercase
 *
 * Convertion is done in place. Only uppercase latin characters
 * will become lowercase, all other characters will remain
 * the same.
 *
 * @param string the string we want to convert to lowercase.
 * @return The same pointer as the one passed in parameters
 */
char *
strlwr(char *string);

/**
 * @brief Converts every character of a null-terminated string into lowercase
 *
 * Convertion is done in place. Only uppercase latin characters
 * will become lowercase, all other characters will remain
 * the same.
 *
 * @param string the string we want to convert to lowercase.
 * @return The same pointer as the one passed in parameters
 */
char *
strupr(char *string);

/**
 * @brief Reverses the order of characters in the string
 *
 * For example, a string "Hello" becomes "olleH". The NULL character at the end
 * of the string remains at the end.
 *
 * @param string the string we want to reverse
 * @return The same pointer as the one passed in parameters
 */
char *
strrev(char *string);

/**
 * @brief Returns a string corresponding to the error number
 *
 * Warning : the returned pointer should be duplicated if the user ever needs to modify it
 *
 * @param errnum number of error
 * @return the pointer to the message corresponding to the errnum or NULL if none was found
 */
char *
strerror(int errnum);

/**
 * @brief Returns a pointer to a new string which is a duplicate of the argument
 *
 * Warning : Memory for the new string is obtained with buddy_alloc() //TODO buddy_alloc/malloc for now?
 * and should be freed with buddy_free().
 * buddy_init() should be called before calling strerror().
 *
 * @param string string to duplicate
 * @return the pointer to the duplicate of the argument string or NULL if couldn't allocate enough memory space
 */
// char *strdup(const char *string);

/**
 * @brief Returns a pointer to a new string which is a duplicate of the argument (copies at most n bytes)
 *
 * Warning : Memory for the new string is obtained with buddy_alloc() //TODO buddy_alloc/malloc for now?
 * and should be freed with buddy_free().
 * buddy_init() should be called before calling strerror().
 *
 * @param string string to duplicate
 * @param n max number of characters to duplicate
 * @return the pointer to the duplicate of the argument string or NULL if couldn't allocate enough memory space
 */
// char *strndup(const char *string, size_t n);

/**
 * @brief Calculates the length of the longest prefix of string which consists entirely of bytes in accept.
 *
 *  The function does not sort the characters and does not delete the double occurrences in accept.
 *  If the user wishes to accelerate the function, she/he needs to do this separately.
 *  This decision was made, because the usefulness of such operations highly depends on the parameters.
 *
 * @param string the target of the function
 * @param accept key characters that the longest prefix consists of
 * @return the index of the first character in string that is not in accept
 */
size_t
strspn(const char *string, const char *accept);

/**
 * @brief Calculates the length of the longest prefix of string which consists entirely of bytes not in reject.
 *
 *  The function does not sort the characters and does not delete the double occurrences in accept.
 *  If the user wishes to accelerate the function, she/he needs to do this separately.
 *  This decision was made, because the usefulness of such operations highly depends on the parameters.
 *
 * @param string the target of the function
 * @param reject key characters that must not be in the longest prefix
 * @return the index of the first character in string that is the same as any of the ones in reject
 */
size_t
strcspn(const char *string, const char *reject);

/**
 * @brief Locates the first occurrence in the target string of any of the bytes in the string accept.
 *
 *  The function does not sort the characters and does not delete the double occurrences in accept.
 *  If the user wishes to accelerate the function, she/he needs to do this separately.
 *  This decision was made, because the usefulness of such operations highly depends on the parameters.
 *
 * @param string the target of the function
 * @param accept key characters
 * @return a pointer to the byte in string that matches one of the bytes in accept, or NULL if no such byte is found.
 */
char *
strpbrk(const char *string, const char *accept);

/**
 * @brief Finds the first occurrence of the substring needle in the string haystack.
 *
 * If needle is an empty string, the result will be the same pointer as the one passed to haystack
 * This function uses KMP algorithm, and thereby uses more memory space
 *
 * @param haystack the target string where we look for a pattern
 * @param needle pattern to look for
 * @return a pointer to the beginning of the located substring, or NULL if the substring is not found
 */
char *
strstr(const char *haystack, const char *needle);

/**
 * @brief Extracts a token from a string
 *
 * The strtok_r() function breaks a string into a sequence of zero or more nonempty tokens.
 *
 * On the first call to strtok_r(), str should point to the string to be parsed, and the value of saveptr is ignored (modified
 * internally). In each subsequent call that should parse the same string, str must be NULL and saveptr should be unchanged since
 * the previous call.
 *
 * The caller may specify different strings in delim in successive calls that parse the same string. For instance, if string is
 * "a,b,c d e f,g", by calling strtok_r() only with "," delimiter will create 4 tokens, but if the user calls it 2 times with ","
 * delimiter and then the rest with " " delimiter, 6 tokens wil be created : "a","b","c","d","e","f,g" ("f,g" is indeed one token
 * as strtok_r() was called with " " delimiter).
 *
 * Each call to strtok_r() returns a pointer to a null-terminated string containing the next token. This string does not include
 * the delimiting byte. If no more tokens are found, strtok_r() returns NULL.
 *
 * A sequence of calls to strtok_r() that operate on the same string maintains a pointer that determines the point from which to
 * start searching for the next token. The first call to strtok_r() sets this pointer to point to the first byte of the string.
 * The start of the next token is determined by scanning forward for the next nondelimiter byte in str. If such a byte is found,
 * it is taken as the start of the next token. If no such byte is found, then there are no more tokens, and strtok_r() returns
 * NULL. (A string that is empty or that contains only delimiters will thus cause strtok() to return NULL on the first call.)
 *
 * The end of each token is found by scanning forward until either the next delimiter byte is found or until the terminating null
 * byte ('\0') is encountered. If a delimiter byte is found, it is OVERWRITTEN with a null byte to terminate the current token,
 * and strtok_r() saves a pointer to the following byte; that pointer will be used as the starting point when searching for the
 * next token. In this case, strtok_r() returns a pointer to the start of the found token.
 *
 * From the above description, it follows that a sequence of two or more contiguous delimiter bytes in the parsed string is
 * considered to be a single delimiter, and that delimiter bytes at the start or end of the string are ignored. Put another way:
 * the tokens returned by strtok() are always nonempty strings.
 *
 * Different strings may be parsed concurrently using sequences of calls to strtok_r() that specify different saveptr arguments.
 *
 * Warning :
 *      strtok_r modifies str
 *      identity of delimiter bytes is lost (i.e. most of them will be replaced by '\0' in str)
 *      if NULL is passed as the first argument for the FIRST call, *saveptr (not saveptr) should also be NULL
 *
 * @param str string to extract tokens from
 * @param delim string that contains bytes that would serve as delimiters
 * @param saveptr pointer used internally by strtok_r in order to maintain context between successive calls that parse the same
 * string
 * @return a pointer to the beginning of the next token terminated by '\0' or NULL if there are no more tokens
 */
char *
strtok_r(char *str, const char *delim, char **saveptr);

/**
 * @brief Extracts a token from a string
 *
 * If *stringp is NULL, the strsep() function returns NULL and does nothing else. Otherwise, this function finds the first token
 * in the string *stringp, that is delimited by one of the bytes in the string delim. This token is terminated by overwriting the
 * delimiter with a null byte ('\0'), and *stringp is updated to point past the token. In case no delimiter was found, the token
 * is taken to be the entire string *stringp, and *stringp is made NULL.
 *
 * Warning:
 *      strsep() modifies its first parameter
 *      identity of delimiter bytes is lost (they will be replaced by '\0' in *stringp)
 *
 * @param stringp pointer to a string (because string will be modified) to extract tokens from
 * @param delim string that contains bytes that would serve as delimiters
 * @return a pointer to the found null-terminated token, that is, it returns the original value of *stringp
 */
char *
strsep(char **stringp, const char *delim);

#endif /* DPUSYSCORE_STRING_H */
