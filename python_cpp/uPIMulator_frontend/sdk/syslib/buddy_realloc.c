/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

// todo integrate in buddy_alloc.c, when we have reduced the buddy_alloc object file size

#include <buddy_alloc.h>
#include <stddef.h>
#include <alloc.h>
#include <string.h> //for memcpy
#include <dpuruntime.h>
#include <atomic_bit.h>
#include <attributes.h>

ATOMIC_BIT_EXTERN(__buddy_lock);

typedef struct _buddy_search_context_t {
    unsigned int target_level;
    unsigned int real_index;
    unsigned int highlight_target_bit;
    unsigned int highlight_buddy_bit;
} * buddy_search_context_t;

extern void *
safe_buddy_alloc(size_t size);
extern void
safe_buddy_free(void *ptr);
extern int
buddy_search_for_pointer(void *ptr, buddy_search_context_t context);

static int
buddy_sizeofblock(void *pointer)
{
    // We get the pointer (address) as the parameter and look for any allocated block that starts at this address.
    // If it is allocated, it will be found and its size will be returned.
    // If it is currently non allocated, buddy_sizeofblock will do nothing.
    struct _buddy_search_context_t dummy;
    int index = buddy_search_for_pointer(pointer, &dummy);

    if (index == -1) {
        return -1;
    }

    return 1 << (12 - index);
}

// noinline, because part of grind tracked functions
void *__noinline
buddy_realloc(void *ptr, size_t size)
{
    ATOMIC_BIT_ACQUIRE(__buddy_lock);
    void *result = ptr;
    if (ptr == NULL) { // if ptr == NULL, then buddy_realloc must behave as buddy_alloc
        result = safe_buddy_alloc(size);
        ATOMIC_BIT_RELEASE(__buddy_lock);
        return result; //
    }

    if (size == 0) { // if size == 0 and ptr != NULL, then buddy_realloc behaves as buddy_free
        safe_buddy_free(ptr); //
        ATOMIC_BIT_RELEASE(__buddy_lock); //
        return ptr; //
    }

    size_t size_block = buddy_sizeofblock(ptr);

    if (size_block == ((size_t)-1)) { // size_block is set to -1 if ptr was not found among
        ATOMIC_BIT_RELEASE(__buddy_lock); // allocated pointers and that there is nothing to do
        return NULL;
    }

    if (size <= (size_block >> 1)) { // if newly allocated block is smaller than the currently allocated block
        size_block = size; // we will only copy "size" bytes
    } else if (size <= size_block) { // if newly allocated block is of the same size as the currently allocated block
        ATOMIC_BIT_RELEASE(__buddy_lock); // then there is no reason to do anything
        return ptr;
    }

    safe_buddy_free(ptr); // newly allocated block is either bigger or smaller than the currently allocated
    result = safe_buddy_alloc(size); // block and we need to call buddy_alloc to be sure that external fragmentation
    memcpy(result, ptr, size_block); // is avoided.

    ATOMIC_BIT_RELEASE(__buddy_lock);
    return result;
}
