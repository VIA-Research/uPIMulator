/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_BUDDY_ALLOC_H
#define DPUSYSCORE_BUDDY_ALLOC_H

/**
 * @file buddy_alloc.h
 * @brief Dynamic memory allocation and freeing.
 *
 * This library allows to create unique memory space in the heap to allocate and free
 * blocks of data.
 *
 * The memory space is initialized with <code>buddy_init</code>, which must be invoked only once during
 * the program's lifecycle.
 *
 * Functions can then dynamically get and free buffers, using <code>buddy_alloc</code> and <code>buddy_free</code>
 * respectively.
 *
 * In this implementation, the allocatable buffer size is chosen during the first call to <code>buddy_init</code>.
 * Tested sizes : 2048, 4096, 8192, 16384 and 32768 bytes
 * The allocated buffers are properly aligned on DMA transfer constraints, so that they can be
 * used as is in MRAM/WRAM transfer operations.
 */

/*
 * @internal The algorithm used in this implementation is Buddy memory allocation.
 *
 * A particularity of this implementation is that no headers are created and thereby
 * the memory consumption of the heap is reduced. The drawback is a slight slow-down
 * in speed of memory freeing.
 *
 * Warning :
 *   Due to the particularities of the implementation (lack of headers), <code>buddy_free</code>
 *   will always try to find a pointer to free. If the pointer given in the parameter
 *   is not currently allocated by <code>buddy_alloc</code> or <code>buddy_realloc</code>, <code>buddy_free</code>
 *   will do nothing.
 *
 */

#include <stddef.h>

/**
 * @fn buddy_init
 * @brief Allocates size_of_heap bytes for a heap that <code>buddy_alloc</code> can access to.
 *
 * Reserves memory space in the heap used to perform dynamic allocation and release of buffers.
 *
 * @param size_of_heap the size of heap in bytes that <code>buddy_alloc</code> can access to
 */
void
buddy_init(size_t size_of_heap);

/**
 * @fn buddy_reset
 * @brief Resets the heap.
 *
 * Quickly frees all pointers allocated by <code>buddy_alloc</code> or <code>buddy_realloc</code>.
 * Warning : currently buddy_reset() doesn't reset the size of the allocated heap.
 */
void
buddy_reset(void);

/**
 * @fn buddy_alloc
 * @brief Allocates a buffer of the given size in the heap, in a runtime-safe way.
 *
 * The allocated buffer is aligned on 64 bits, in order to ensure compatibility
 * with the maximum buffer alignment constraint. As a consequence, a buffer
 * allocated with this function is also compatible with data transfers to/from MRAM.
 *
 * Due to the idea of the buddy algorithm (to decrease external fragmentation),
 * the allocated blocks will be of size equal to a power of 2. In other words,
 * if the user allocates 33 bytes, 64 bytes will be allocated and when 2049 bytes
 * are requested, 4096 will be allocated. The user might want to take this into
 * account if she/he wishes to minimise the memory consumption.
 *
 * The minimal size of the allocated block is 16 bytes, but can easily be changed in
 * future implementations, so <code>buddy_alloc</code> is mostly adapted to allocating medium and
 * big structures, such as arrays containing more than 8 bytes (in order to make sure
 * that not too much memory space is wasted), binary trees or linked lists.
 *
 * If the <code>size</code> passed in parameter is less or equal to 0 or greater than the size of heap,
 * errno will be set to EINVAL and <code>buddy_alloc</code> will do nothing.
 * If <code>buddy_alloc</code> fails to find enough free memory space to allocate, errno will be
 * set to ENOMEM.
 *
 * @param size the allocated buffer's size, in bytes
 * @return A pointer to the allocated buffer if one was available, NULL otherwise.
 */
void *
buddy_alloc(size_t size);

/**
 * @fn buddy_free
 * @brief Frees a specified pointer, in a runtime-safe way.
 *
 *  Warning :
 *   Due to the particularities of the implementation (lack of headers), <code>buddy_free</code>
 *   will always try to find a pointer to free and will see a pointer to the beginning of the
 *   block in the same way as the pointer to anywhere inside the block. For example, if we have
 *   allocated an int array[10], <code>buddy_free</code> will treat &array[0] the same way as &array[1]
 *   or as the address of any other element inside this array and will free the whole block.
 *
 * If the pointer given in the parameter is not currently allocated by <code>buddy_alloc</code> or
 * <code>buddy_realloc</code>, <code>buddy_free</code> will do nothing.
 *
 * This function frees the memory space pointed to by pointer, which
 * must have been returned by a previous call to <code>buddy_alloc</code> or <code>buddy_realloc</code>
 * If it wasn't or if <code>buddy_free</code> has already been called for this pointer before,
 * then <code>buddy_free</code> will do nothing. If pointer is NULL, no operation is performed.
 * If the pointer passed as a parameter is not aligned to 64 bits or if it is outside
 * of the allocated heap errno will be set to EINVAL.
 * If <code>buddy_free</code> detects the attempt to free a non-allocated pointer, it will equally
 * set errno to EINVAL.
 *
 * @param pointer the pointer to the block to free
 */
void
buddy_free(void *pointer);

/**
 * @fn buddy_realloc
 * @brief Changes the size of the memory block pointed to by <code>ptr</code> to <code>size</code> bytes in a runtime-safe way.
 *
 * The contents will be unchanged in the range from the start of the region up to the minimum of the old and new sizes.
 * If the new <code>size</code> is larger than the old size, the added memory will not be initialized.
 * If <code>ptr</code> is NULL, then the call is equivalent to <code>buddy_alloc(size)</code> for all values of <code>size</code>.
 * If <code>size</code> is equal to zero, and <code>ptr</code> is not NULL, then the call is equivalent to
 * <code>buddy_free(ptr)</code> and the return value will be equal to the pointer passed as the parameter. Unless <code>ptr</code>
 * is NULL, it should have been returned by an earlier call to <code>buddy_alloc()</code> or <code>buddy_realloc()</code>. If it
 * wasn't, then <code>buddy_realloc()</code> will try to find this pointer among the allocated ones, but undefined behavior might
 * occur.
 *
 * If new <code>size</code> is smaller than the old size, then the remaining memory will potentially be released (depends
 * on the size of block).
 *
 * <code>buddy_realloc()</code> internally calls <code>buddy_alloc</code> and <code>buddy_free</code> and thereby will set errno
 * to ENOMEM or EINVAL on failure.
 *
 * @param ptr original pointer
 * @param size the new allocated buffer's size, in bytes
 * @return A new (or the same) pointer to the allocated buffer if one was available, NULL otherwise.
 */
void *
buddy_realloc(void *ptr, size_t size);

#endif /* DPUSYSCORE_BUDDY_ALLOC_H */
