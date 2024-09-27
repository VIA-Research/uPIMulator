/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef _DPUSYSCORE_STDALIGN_H_
#define _DPUSYSCORE_STDALIGN_H_

/**
 * @file stdalign.h
 * @brief Defines align macros.
 */

/**
 * @def alignas
 * @brief _Alignas specifier.
 */
#define alignas _Alignas

/**
 * @def alignof
 * @brief _Alignof operator.
 */
#define alignof _Alignof

/**
 * @def __alignas_is_defined
 * @brief Whether the alignas macro is defined.
 */
#define __alignas_is_defined 1

/**
 * @def __alignof_is_defined
 * @brief Whether the alignof macro is defined.
 */
#define __alignof_is_defined 1

#endif /* _DPUSYSCORE_STDALIGN_H_ */
