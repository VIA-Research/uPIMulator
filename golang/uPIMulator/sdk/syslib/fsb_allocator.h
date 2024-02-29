/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_FBS_ALLOC_H
#define DPUSYSCORE_FBS_ALLOC_H

/**
 * @file fsb_allocator.h
 * @brief Provides a fixed-size block memory allocator.
 *
 * @internal When defining an allocator, the total memory needed will be allocated, using mem_alloc.
 *           The allocator structure is a pointer to the next available block. In each free block, the first four bytes
 *           store a pointer to the following free block, creating a linked list. To allocate a block, we just check
 *           the free pointer, check the next free pointer and update the free pointer accordingly. To free a block, we
 *           just check the free pointer, update it with the newly free block, and update the next pointer of this block
 *           with the previous free pointer.
 *           There is no protection to prevent invalid block to be added to the list. Moreover, the list being in the free
 *           blocks, if there is some memory overflow from a block, the list might be corrupted and totally invalid.
 */

/**
 * @fn fsb_allocator_t
 * @brief A fixed-size block allocator.
 */
typedef void **fsb_allocator_t;

/**
 * @fn fsb_alloc
 * @brief Allocate and initialize a fixed-size block allocator.
 *
 * @param block_size the size of the blocks allocated (will be realigned on 8 bytes, with a minimum of 8 bytes)
 * @param nb_of_blocks the number of blocks allocated
 * @throws a fault if there is no memory left
 * @return The newly allocated and ready-to-use fixed-size block allocator.
 */
fsb_allocator_t
fsb_alloc(unsigned int block_size, unsigned int nb_of_blocks);

/**
 * @fn fsb_get
 * @brief Own a block of the specified fixed-size block allocator, in a runtime-safe way.
 *
 * @param allocator the allocator from which we take the block
 * @return A pointer to the owned block if one was available, NULL otherwise.
 */
void *
fsb_get(fsb_allocator_t allocator);

/**
 * @fn fsb_free
 * @brief Free a block of the specified fixed-size block allocator, in a runtime-safe way.
 *
 * @param allocator the allocator in which we put the block back in
 * @param ptr the pointer to the block to free
 */
void
fsb_free(fsb_allocator_t allocator, void *ptr);

#endif /* DPUSYSCORE_FBS_ALLOC_H */