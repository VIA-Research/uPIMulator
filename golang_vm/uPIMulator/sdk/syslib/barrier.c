/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <barrier.h>
#include <defs.h>
#include <sysdef.h>
#include <dpuruntime.h>

void
barrier_wait(struct barrier_t *barrier)
{
    unsigned char lock = barrier->lock;
    __acquire(lock, "0");
    unsigned char count = barrier->count;
    unsigned char last = barrier->wait_queue;
    unsigned char first;
    thread_id_t tid = me();

    /* Count = 1 means that I am the last to enter the barrier.
     * Need to wake up everybody.*/
    if (unlikely(count == 1)) {
        if (likely(last != __EMPTY_WAIT_QUEUE)) {
            first = __WAIT_QUEUE_TABLE[last];
            while (first != last) {
                __resume(first, "0");
                first = __WAIT_QUEUE_TABLE[first];
            }
            __resume(first, "0");
            barrier->wait_queue = __EMPTY_WAIT_QUEUE;
            barrier->count = barrier->initial_count;
        }
        __release(lock, "0", __AT_NEXT_INSTRUCTION);
    } else {
        if (unlikely(last == __EMPTY_WAIT_QUEUE)) {
            __WAIT_QUEUE_TABLE[tid] = tid;
        } else {
            first = __WAIT_QUEUE_TABLE[last];
            __WAIT_QUEUE_TABLE[tid] = first;
            __WAIT_QUEUE_TABLE[last] = tid;
        }

        barrier->wait_queue = tid;
        barrier->count = --count;
        __release(lock, "0", __AT_NEXT_INSTRUCTION);
        __stop();
    }
}
