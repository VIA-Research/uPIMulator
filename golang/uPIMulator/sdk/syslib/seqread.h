/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef DPUSYSCORE_SEQREAD_H
#define DPUSYSCORE_SEQREAD_H

/**
 * @file seqread.h
 * @brief Sequential reading of items in MRAM.
 *
 * A sequential reader allows to parse a contiguous area in MRAM in sequence.
 * For example, if the MRAM contains an array of N structures, a sequential
 * reader on this array will automatically fetch the data into WRAM, thus
 * simplify the iterative loop on the elements.
 *
 * The size of cached area is defined by default but can be overriden by
 * defining this value in SEQREAD_CACHE_SIZE.
 *
 * The use of a sequential reader implies:
 *
 *  - first, to allocate some storage in WRAM to cache the items, using seqread_alloc.
 *  - then to initialize a reader on the MRAM area, via seqread_init
 *  - finally to iterate on the elements, invoking seqread_get whenever a new item is accessed.
 *
 */

#include <stdint.h>
#include <mram.h>
#include <macro_utils.h>

#ifndef SEQREAD_CACHE_SIZE
/**
 * @def SEQREAD_CACHE_SIZE
 * @hideinitializer
 * @brief Size of caches used by seqread.
 */
#define SEQREAD_CACHE_SIZE 256
#endif

_Static_assert(SEQREAD_CACHE_SIZE == 32 || SEQREAD_CACHE_SIZE == 64 || SEQREAD_CACHE_SIZE == 128 || SEQREAD_CACHE_SIZE == 256
        || SEQREAD_CACHE_SIZE == 512 || SEQREAD_CACHE_SIZE == 1024,
    "seqread error: invalid cache size defined");

#define __SEQREAD_FCT(suffix) __CONCAT(__CONCAT(seqread, SEQREAD_CACHE_SIZE), suffix)
#define __SEQREAD_ALLOC __SEQREAD_FCT(_alloc)
#define __SEQREAD_INIT __SEQREAD_FCT(_init)
#define __SEQREAD_GET __SEQREAD_FCT(_get)
#define __SEQREAD_TELL __SEQREAD_FCT(_tell)
#define __SEQREAD_SEEK __SEQREAD_FCT(_seek)

/**
 * @typedef seqreader_buffer_t
 * @brief An buffer to use to initial a sequential reader.
 */
typedef uintptr_t seqreader_buffer_t;

/**
 * @typedef seqreader_t
 * @brief An object used to perform sequential reading of MRAM.
 */
typedef struct {
    seqreader_buffer_t wram_cache;
    uintptr_t mram_addr;
} seqreader_t;

seqreader_buffer_t
__SEQREAD_ALLOC();

/**
 * @fn seqread_alloc
 * @brief Initializes an area in WRAM to cache the read buffers.
 *
 * Notice that this buffer can be re-used for different sequential reads,
 * as long as it is initialized each time to a new buffer in MRAM.
 *
 * @return A pointer to the allocated cache base address.
 */
#define seqread_alloc __SEQREAD_ALLOC

void *
__SEQREAD_INIT(seqreader_buffer_t cache, __mram_ptr void *mram_addr, seqreader_t *reader);

/**
 * @fn seqread_init
 * @brief Creates a sequential reader.
 *
 * The reader is associated to an existing cache in WRAM, created with
 * seqread_alloc and a contiguous area of data in MRAM. The function
 * loads the first pages of data into the cache and provides a pointer
 * to the first byte in cache actually mapping the expected data.
 *
 * Notice that the provided MRAM address does not need to be aligned on
 * any constraint: the routine does the alignment automatically.
 *
 * @param cache the reader's cache in WRAM
 * @param mram_addr the buffer address in MRAM
 * @param reader the sequential reader to init to the supplied MRAM address
 * @return A ptr to the first byte in cache corresponding to the MRAM address
 */
#define seqread_init __SEQREAD_INIT

void *
__SEQREAD_GET(void *ptr, uint32_t inc, seqreader_t *reader);

/**
 * @fn seqread_get
 * @brief Fetches the next item in a sequence.
 *
 * This operation basically consists in incrementing the pointer that goes
 * through the mapped area of memory. The function automatically reloads
 * data from cache if necessary.
 *
 * As a result, the provided pointer to the cache area is set to its new value.
 *
 * The provided increment must be less than SEQREAD_CACHE_SIZE. The reader's
 * behavior is undefined if the increment exceeds this value.
 *
 * @param ptr the incremented pointer
 * @param inc the number of bytes added to this pointer
 * @param reader a pointer to the sequential reader
 * @return The updated pointer value.
 */
#define seqread_get __SEQREAD_GET

void *
__SEQREAD_SEEK(__mram_ptr void *mram_addr, seqreader_t *reader);

/**
 * @fn seqread_seek
 * @brief Set the position of the cache to the supplied MRAM address
 *
 * Update automatically the cache if necessary.
 *
 * @param mram_addr the new buffer address in MRAM
 * @param reader a pointer to the sequential reader
 * @return A ptr to the first byte in cache corresponding to the MRAM address
 */
#define seqread_seek __SEQREAD_SEEK

__mram_ptr void *
__SEQREAD_TELL(void *ptr, seqreader_t *reader);

/**
 * @fn seqread_tell
 * @brief Get the MRAM address corresponding to the supplied ptr in the cache
 *
 * @param ptr a pointer in the cache
 * @param reader a pointer to the sequential reader
 * @return A ptr to the MRAM address corresponding to the supplied pointer in the cache
 */
#define seqread_tell __SEQREAD_TELL

#endif /* DPUSYSCORE_SEQREAD_H */
