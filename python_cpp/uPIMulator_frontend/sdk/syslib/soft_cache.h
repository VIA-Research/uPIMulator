/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_SOFT_CACHE_H
#define DPUSYSCORE_SOFT_CACHE_H

/**
 * @file soft_cache.h
 * @brief Software cache
 *
 * The software cache mechanism emulates a hardware cache to transparently load and store data from and to the
 * MRAM.
 *
 * This mechanism is quite slow, thus would only be used during the development process, to simplify the code.
 *
 * This module defines:
 *
 *  - A procedure to start the software cache, by creating a "virtual TLB" in the system, along with an area in WRAM to contain
 *    the cached MRAM pages
 *  - A procedure to flush the cache at the end of an execution, ensuring that the data in MRAM are consistent with the cached
 *    data
 *  - The special C directive "__mram", used to declare a pointer directly representing a buffer in MRAM.
 *
 * An MRAM pointer is mapped by the caching system. As a consequence, any access to data within this buffer is trapped
 * by a cache load or store procedure, transparently performing the required memory transactions to fetch and write back
 * the data.
 *
 */

#include <attributes.h>

#endif /* DPUSYSCORE_SOFT_CACHE_H */
