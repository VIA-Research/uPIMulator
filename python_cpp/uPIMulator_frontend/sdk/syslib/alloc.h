/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_ALLOC_H
#define DPUSYSCORE_ALLOC_H

/**
 * @file alloc.h
 * @brief Provides a way to manage heap allocation.
 *
 * @internal The heap is situated after the different kernel structures, local and global variables.
 *           It can grow until reaching the end of the WRAM. A reboot of the DPU reset the Heap.
 *           The current heap pointer can be accessed at the address defined by __HEAP_POINTER__.
 */

#include <stddef.h>

#include <fsb_allocator.h>
#include <buddy_alloc.h>

#include <attributes.h>

/**
 * @fn mem_alloc
 * @brief Allocates a buffer of the given size in the heap.
 *
 * The allocated buffer is aligned on 64 bits, in order to ensure compatibility
 * with the maximum buffer alignment constraint. As a consequence, a buffer
 * allocated with this function is also compatible with data transfers to/from MRAM.
 *
 * @param size the allocated buffer's size, in bytes
 * @throws a fault if there is no memory left
 * @return The allocated buffer address.
 */
void *
mem_alloc(size_t size);

/**
 * @fn mem_reset
 * @brief Resets the heap.
 *
 * Every allocated buffer becomes invalid, since subsequent allocations restart from the beginning
 * of the heap.
 *
 * @return The heap initial address.
 */
void *
mem_reset(void);

#endif /* DPUSYSCORE_ALLOC_H */
