/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef _DPUSYSCORE_STDBOOL_H_
#define _DPUSYSCORE_STDBOOL_H_

/**
 * @file stdbool.h
 * @brief Defines the boolean type.
 */

/**
 * @def __bool_true_false_are_defined
 * @brief Whether the boolean type and values are defined.
 */
#define __bool_true_false_are_defined 1

/**
 * @def bool
 * @brief The boolean type.
 */
#define bool _Bool

/**
 * @def true
 * @brief The <code>true</code> constant, represented by <code>1</code>
 */
#define true 1
/**
 * @def false
 * @brief The <code>false</code> constant, represented by <code>0</code>
 */
#define false 0

#endif /* _DPUSYSCORE_STDBOOL_H_ */
