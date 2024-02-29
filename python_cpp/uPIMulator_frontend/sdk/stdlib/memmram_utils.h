/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef _DPUSYSCORE_MEMMRAM_UTILS_H_
#define _DPUSYSCORE_MEMMRAM_UTILS_H_

#define ALIGN_MASK(x, mask) (((x) + (mask)) & ~(mask))
#define ALIGN(x, a) ALIGN_MASK((x), (a)-1)
#define DMA_ALIGNMENT 8
#define DMA_OFF_MASK (DMA_ALIGNMENT - 1)
#define DMA_ALIGNED(x) ALIGN(x, DMA_ALIGNMENT)

#define MIN(a, b) ((a) < (b) ? (a) : (b))

#define MRAM_CACHE_SIZE 8

#endif /* _DPUSYSCORE_MEMMRAM_UTILS_H_ */
