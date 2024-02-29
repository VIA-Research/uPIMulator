/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <buddy_alloc.h>
#include <alloc.h> //for mem_alloc
#include <stddef.h> //for size_t
#include <string.h> //for memset
#include <dpuruntime.h>
#include <errno.h>
#include <defs.h>
#include <stdbool.h>
#include <built_ins.h>
#include <atomic_bit.h>
#include <attributes.h>

static unsigned char __buddy_init_done = 0;
static int *__buddy_blocks = 0;
static void *__buddy_heap_start = 0;
static unsigned int __BUDDY_SIZE_OF_HEAP__ = 0;
static unsigned char __BUDDY_MAX_POWER__ = 0;
static unsigned char __BUDDY_NUMBER_OF_LEVELS__ = 0;

ATOMIC_BIT_INIT(__buddy_lock);

#define __BUDDY_DEPTH_LEVELS__ 3
#define __BUDDY_SHIFT_ADDRESS_TO_INDEX__ 4

static inline unsigned int
next_power_of_2(int x)
{
    // in order to find the size of the block to allocate we
    // count the number of leading zeros to get the correct
    // log2(size)
    // if size has only one "1" bit then we keep this power of 2
    // otherwise, we add 1

    unsigned int power_of_2 = 31 - count_leading_zeros(x);

    if (count_population(x) != 1) {
        power_of_2++;
    }

    return power_of_2;
}

/*
 * Note 1:
 *   each bit represents the state of a block. It can equal either 0 or 1 and can represent 4 different states of a block :
 *   target = 1 & buddy = 0                                      -> free                     : target block is [free] and can be
 * allocated without the need to cut the bigger block in half target = 1 & buddy = 1                                      -> not
 * in use               : neither target nor its buddy are in use and bigger blocks must be checked
 *
 *   target = 0 & both successors = 1 [not in use]               -> allocated                : this block is completely allocated
 *   target = 0 & one or both of the successors = 0 [allocated]  -> partially allocated      : block is partially allocated i.e.
 * at least one of the sub-blocks is [allocated]
 *
 * Note 2:
 *  Indexes of __buddy_blocks are indeed quite bizarre, which was a mistake during the conception stage. It doesn't change much.
 *  Indexe_in_level start from 0 from left to right, but not for the 5 least significant bits where they still start from 0 but
 * from right to left
 */
void *
safe_buddy_alloc(size_t size)
{
    if ((size == 0) || (size > __BUDDY_SIZE_OF_HEAP__)) {
        errno = EINVAL;
        return NULL;
    }

    // we replace the size by the smallest 2 to the power of X such as it is greater than or equals the size
    // afterwards, we take X as power_of_2 and thereby find the correct level to search for free blocks:
    // blocks_level = __BUDDY_MAX_POWER__ - power_of_2
    // for example if we try to allocate a block of 64 bytes (64=2^6), blocks_level = __BUDDY_MAX_POWER__ - 6

    unsigned int power_of_2 = next_power_of_2(size);
    int blocks_level = __BUDDY_MAX_POWER__ - power_of_2;

    // if the size is smaller than the minimal size of
    // block divided by 2, then we can't allocate a block
    // so small and we have to allocate the block of the
    // minimal allowed size
    if (blocks_level >= __BUDDY_NUMBER_OF_LEVELS__) {
        blocks_level = __BUDDY_NUMBER_OF_LEVELS__ - 1;
    }

    // we browse all levels until we find the smallest free block that is big enough to contain "size" bytes
    for (int blocks_level_current = blocks_level; blocks_level_current >= 0; --blocks_level_current) {
        // we initialise index_in_level as the biggest index permitted in level
        unsigned int index_in_level = 1 << blocks_level_current;

        // levels 0-4 are in the same bitfield,
        // level 5 consists of a single bitfield
        // all other levels consist of 1<<(#level - 5)
        // bitfields
        unsigned int initial_number_of_current_bitfield = 0;
        unsigned int loaded_case_mask;

        switch (blocks_level_current) {
            default:
                initial_number_of_current_bitfield = ((1 << blocks_level_current) >> 5) - 1;
                loaded_case_mask = 0xFFFFFFFF;
                break;
            case 0:
                loaded_case_mask = 0x40000000;
                break;
            case 1:
                loaded_case_mask = 0x30000000;
                break;
            case 2:
                loaded_case_mask = 0x0F000000;
                break;
            case 3:
                loaded_case_mask = 0x00FF0000;
                break;
            case 4:
                loaded_case_mask = 0x0000FFFF;
                break;
        }

        for (int number_of_current_bitfield = initial_number_of_current_bitfield; number_of_current_bitfield >= 0;
             --number_of_current_bitfield, index_in_level -= 32) {
            // sizes of levels in __buddy_blocks is a geometric
            // series and thereby we can easily calculate the
            // number of bitfields that precede the current level
            unsigned int real_index = (1 << blocks_level_current) >> 5;
            // we load the bitfield
            unsigned int *initial_loaded_case_address
                = (unsigned int *)(__buddy_blocks + number_of_current_bitfield + real_index);
            // if we are in the very first bitfield
            // that contains first 5 levels, then we
            // need to make sure that buddy_free will
            // ignore the bits that represent blocks
            // that don't belong to blocks_level_current
            unsigned int loaded_case = *initial_loaded_case_address & loaded_case_mask;

            // this formula gives the number of zeros
            // before the first pair of free/allocated blocks
            // present in the current bitfield. If it equals
            // 32, then no such pair is present.
            unsigned int lz_before_first_pair_tmp = ((loaded_case << 1) ^ loaded_case) & 0xAAAAAAAA;

            if (lz_before_first_pair_tmp != 0) {
                unsigned int lz_before_first_pair = count_leading_zeros(lz_before_first_pair_tmp);
                index_in_level = index_in_level - lz_before_first_pair - 1;

                // if we are in the very first bitfield,
                // then we need to take into account the special
                // positioning of the 5 levels of blocks
                // inside this bitfield
                if (blocks_level_current < 5) {
                    index_in_level += 1 << blocks_level_current;
                }

                unsigned int highlight_target_bit = 1 << index_in_level;

                // if current level is among first 4
                //(5th is unnecessary to consider)
                // then we need to shift the mask
                if (blocks_level_current < 4) {
                    // we want to find the value of the offset to shift the mask for target and its buddy
                    // offset = 32 - 2^(current_lvl+1)
                    highlight_target_bit = highlight_target_bit << (32 - (2 << blocks_level_current));
                }

                // we have the position of a pair of
                // free/allocated blocks, but we might
                // need to shift this position by 1
                // if the target's buddy is the potential
                // block to allocate
                loaded_case = loaded_case & highlight_target_bit;

                if (loaded_case == 0) {
                    highlight_target_bit = highlight_target_bit >> 1;
                    index_in_level--;
                }

                // loaded_case have been modified if the first bitfield is handled
                loaded_case = *initial_loaded_case_address;

                // we mark the target bit as allocated
                loaded_case -= highlight_target_bit;
                *initial_loaded_case_address = loaded_case;

                // if we have already had a free block of the necessary size,
                // we return its address, otherwise we jump to
                //__buddy_alloc_break_loop and start cutting the blocks
                // in half until we get a block of the required size
                if (blocks_level == blocks_level_current) {
                    return __buddy_heap_start + (index_in_level << (__BUDDY_MAX_POWER__ - blocks_level));
                }

                // we have already set the block to 0 [partially allocated]
                // so we descend and start with the next one
                blocks_level_current++;

                while (blocks_level_current <= blocks_level) {
                    index_in_level = index_in_level << 1;
                    highlight_target_bit = 1 << index_in_level;

                    // first 5 levels are stored in the same
                    // bitfield. We need to handle this
                    // special case and shift the mask
                    int blocks_level_clamped = blocks_level_current - 4;

                    if (blocks_level_clamped < 0) {
                        highlight_target_bit = highlight_target_bit << (32 - (2 << blocks_level_current));
                    }

                    // we load the bitfield

                    // real_index is the offset in 32-bit words from the first bitfield in __buddy_blocks
                    unsigned int real_index = 0;
                    blocks_level_clamped--;
                    if (blocks_level_clamped >= 0) {
                        real_index = 1 << blocks_level_clamped;
                        blocks_level_clamped = index_in_level >> 5;
                        real_index += blocks_level_clamped;
                    }
                    loaded_case = __buddy_blocks[real_index];

                    // we mark the bit as allocated
                    __buddy_blocks[real_index] = loaded_case - highlight_target_bit;

                    // we continue to descend
                    blocks_level_current++;
                }

                return __buddy_heap_start + (index_in_level << (__BUDDY_MAX_POWER__ - blocks_level));
            }
        }
    }

    // if no level contains a big enough block, we return NULL
    errno = ENOMEM;
    return NULL;
}

typedef struct _buddy_search_context_t {
    unsigned int target_level;
    unsigned int real_index;
    unsigned int highlight_target_bit;
    unsigned int highlight_buddy_bit;
} * buddy_search_context_t;

int
buddy_search_for_pointer(void *pointer, buddy_search_context_t context)
{
    // if the pointer is not aligned to 64 bits, then it is corrupted
    // if the pointer is outside of the heap, we can do nothing
    if (((((unsigned int)pointer) & 7) != 0) || (pointer < __buddy_heap_start)
        || (pointer > (__buddy_heap_start + __BUDDY_SIZE_OF_HEAP__ - 1))) {
        errno = EINVAL;
        return -1;
    }

    // we transform the real address into an index for __buddy_blocks
    // index_in_level = (pointer - START_OF_HEAP)>>__BUDDY_SHIFT_ADDRESS_TO_INDEX__;
    unsigned int index_in_level = (pointer - __buddy_heap_start) >> __BUDDY_SHIFT_ADDRESS_TO_INDEX__;
    // we start to search for the pointer from the lowest level
    unsigned int target_level = __BUDDY_NUMBER_OF_LEVELS__ - 1;

    while (true) {
        // knowing the index_in_level we can calculate
        // highlight_target_bit and highlight_buddy_bit
        unsigned int highlight_target_bit = 1 << index_in_level;
        unsigned int highlight_buddy_bit;

        if ((count_leading_zeros(highlight_target_bit) & 1) == 0) {
            highlight_buddy_bit = highlight_target_bit >> 1;
        } else {
            highlight_buddy_bit = highlight_target_bit << 1;
        }

        // first 5 levels are stored in the
        // same bitfield. We need to handle
        // this special case and shift
        // the masks
        int target_level_clamped = target_level - 4;

        if (target_level_clamped < 0) {
            // we want to find the value of the offset to shift the mask for target and its buddy
            // offset = 32 - 2^(current_lvl+1)
            unsigned int offset = 32 - (2 << target_level);

            highlight_buddy_bit = highlight_buddy_bit << offset;
            highlight_target_bit = highlight_target_bit << offset;
        }

        // we load a bitfield corresponding to the index
        // in __buddy_blocks
        // real_index is the offset in 32-bit words from the first bitfield in __buddy_blocks
        unsigned int real_index = 0;
        target_level_clamped--;

        if (target_level_clamped >= 0) {
            real_index = (1 << target_level_clamped);
            target_level_clamped = index_in_level >> 5;
            real_index = real_index + target_level_clamped;
        }

        unsigned int loaded_case = __buddy_blocks[real_index];

        // We search for an allocated block
        // if we are on the highest level, we quit.
        if (((loaded_case & highlight_target_bit) == 0) || (target_level <= 0)) {
            context->target_level = target_level;
            context->real_index = real_index;
            context->highlight_target_bit = highlight_target_bit;
            context->highlight_buddy_bit = highlight_buddy_bit;
            return index_in_level;
        }

        // Condition required to avoid unwanted release
        if ((loaded_case & highlight_buddy_bit) == 0) {
            errno = EINVAL;
            return -1;
        }

        // we rise to the higher level
        index_in_level = index_in_level >> 1;
        target_level--;
    }
}

static void
buddy_free_fusion_of_blocks(unsigned int index, buddy_search_context_t context)
{
    unsigned int real_index = context->real_index;
    unsigned int current_level_freeing = context->target_level;
    unsigned int highlight_target_bit = context->highlight_target_bit;
    unsigned int highlight_buddy_bit = context->highlight_buddy_bit;
    unsigned int loaded_case = __buddy_blocks[real_index];
    // if we are at the highest level, there is nothing to fuse
    while (current_level_freeing > 0) {
        // fusion occurs only when both blocks are free
        if (((loaded_case & highlight_target_bit) == 0) | ((loaded_case & highlight_buddy_bit) == 0)) {
            return;
        }

        // index of predecessor = index of successor >> 1
        index = index >> 1;
        // we rise to the higher level
        current_level_freeing--;

        // knowing the index_in_level we can calculate
        // highlight_target_bit and highlight_buddy_bit
        highlight_target_bit = 1 << index;

        if ((count_leading_zeros(highlight_target_bit) & 1) == 0) {
            highlight_buddy_bit = highlight_target_bit >> 1;
        } else {
            highlight_buddy_bit = highlight_target_bit << 1;
        }

        // first 5 levels are stored in the same bitfield.
        // We need to handle this special case and shift
        // the masks
        int target_level_clamped = current_level_freeing - 4;

        if (target_level_clamped < 0) {
            // we want to find the value of the offset to shift the mask for target and its buddy
            // offset = 32 - 2^(current_lvl+1)
            unsigned int offset = 32 - (2 << current_level_freeing);

            highlight_buddy_bit = highlight_buddy_bit << offset;
            highlight_target_bit = highlight_target_bit << offset;
        }

        // we load a bitfield corresponding
        // to the index in __buddy_blocks
        //
        // As both sub-blocks were freed, their "father"
        // must be marked as free.
        real_index = 0;
        target_level_clamped--;

        if (target_level_clamped >= 0) {
            real_index = (1 << target_level_clamped);
            target_level_clamped = index >> 5;
            real_index = real_index + target_level_clamped;
        }

        loaded_case = __buddy_blocks[real_index] | highlight_target_bit;
        __buddy_blocks[real_index] = loaded_case;
    }
}

void
safe_buddy_free(void *pointer)
{
    struct _buddy_search_context_t context;

    int index = buddy_search_for_pointer(pointer, &context);

    if (index != -1) {
        __buddy_blocks[context.real_index] |= context.highlight_target_bit;
        buddy_free_fusion_of_blocks(index, &context);
    }
}

// noinline, because part of grind tracked functions
void __noinline *
buddy_alloc(size_t size)
{
    ATOMIC_BIT_ACQUIRE(__buddy_lock);
    void *result = safe_buddy_alloc(size);
    ATOMIC_BIT_RELEASE(__buddy_lock);
    return result;
}

// noinline, because part of grind tracked functions
void __noinline
buddy_free(void *pointer)
{
    ATOMIC_BIT_ACQUIRE(__buddy_lock);
    safe_buddy_free(pointer);
    ATOMIC_BIT_RELEASE(__buddy_lock);
}

/*if the size of the smallest block needs to be changed, then only 2 things need to change :
 *  __BUDDY_DEPTH_LEVELS__ here in buddy_init.c
 *  __BUDDY_SHIFT_ADDRESS_TO_INDEX__ in buddy_defs.s
 *
 *  When the minimal size of a block needs to be 32 bytes, then these two constants must equal 4 and 5 correspondingly
 *  If minimal size needs to be 16 bytes, then these two constants must equal 3 and 4 correspondingly
 *  If minimal size needs to be 8 bytes, then these two constants must equal 2 and 3 correspondingly
 *
 *  Also, certain tests (Global, Reset, LevelByLevel) should also be changed.
 */

// TODO: Right now the size must be of power of 2 : 2048, 4096, 8192 and it won't work with other sizes.
// TODO: There is no point in accepting other sizes, as the whole idea of having buddy allocation is based on it.
// TODO: Right now the possibility of failure of mem_alloc is not taken into account of.
// noinline, because part of grind tracked functions
void __noinline
buddy_init(size_t size_of_heap)
{
    if (__buddy_init_done == 0) {
        ATOMIC_BIT_RELEASE(__buddy_lock);

        __BUDDY_SIZE_OF_HEAP__ = size_of_heap;

        unsigned int power_of_2 = count_leading_zeros(size_of_heap);
        __BUDDY_MAX_POWER__ = (31 - power_of_2);

        __BUDDY_NUMBER_OF_LEVELS__ = __BUDDY_MAX_POWER__ - __BUDDY_DEPTH_LEVELS__;

        unsigned int blocks_in_buddy_blocks = (1 << (__BUDDY_NUMBER_OF_LEVELS__ - 5)) << 2;
        __buddy_blocks = mem_alloc(size_of_heap + blocks_in_buddy_blocks);

        __buddy_heap_start = __buddy_blocks + (blocks_in_buddy_blocks >> 2);

        __buddy_blocks[0] = 0x7fffffff; // all bits except for the very first one must be set to 1
        memset(&__buddy_blocks[1], 0xff, blocks_in_buddy_blocks - 4); // in order to initialize the __buddy_blocks structure
        __buddy_init_done = 1;
    }
}

// noinline, because part of grind tracked functions
void __noinline
buddy_reset()
{
    ATOMIC_BIT_ACQUIRE(__buddy_lock);

    __buddy_blocks[0] = 0x7fffffff; // all bits except for the very first one must be set to 1
    memset(&__buddy_blocks[1],
        0xff,
        ((1 << (__BUDDY_NUMBER_OF_LEVELS__ - 5)) << 2) - 4); // in order to initialize the __buddy_blocks structure

    ATOMIC_BIT_RELEASE(__buddy_lock);
}
