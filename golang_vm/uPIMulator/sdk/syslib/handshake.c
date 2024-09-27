/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <defs.h>
#include <errno.h>
#include <dpuconst.h>
#include <dpuruntime.h>
#include <atomic_bit.h>

unsigned char __handshake_array[NR_THREADS] = { [0 ...(NR_THREADS - 1)] = __EMPTY_WAIT_QUEUE };

ATOMIC_BIT_INIT(__handshake)[NR_THREADS];

#define __acquire_handshake(off) __ATOMIC_BIT_ACQUIRE(off + (ATOMIC_BIT_GET(__handshake) - &__atomic_start_addr), 0)
#define __release_handshake(off) __ATOMIC_BIT_RELEASE(off + (ATOMIC_BIT_GET(__handshake) - &__atomic_start_addr), 0)

void
handshake_notify(void)
{
    thread_id_t tid = me();
    unsigned char info;
    __acquire_handshake(tid);
    info = __handshake_array[tid];

    if (unlikely(info == __EMPTY_WAIT_QUEUE)) {
        __handshake_array[tid] = tid;
        __release_handshake(tid);
        __stop();
    } else {
        __resume(info, "0");
        __handshake_array[tid] = __EMPTY_WAIT_QUEUE;
        __release_handshake(tid);
    }
}

int
handshake_wait_for(unsigned int notifier)
{
    thread_id_t tid = me();

    unsigned char thread = (unsigned char)notifier;

    __acquire_handshake(thread);
    unsigned char info = __handshake_array[thread];

    if (unlikely(info == __EMPTY_WAIT_QUEUE)) {
        __handshake_array[thread] = tid;
        __release_handshake(thread);
        __stop();
    } else {
        if (unlikely(info != thread)) {
            errno = EALREADY;
            __release_handshake(thread);
            return EALREADY;
        } else {
            __resume(info, "0");
            __handshake_array[thread] = __EMPTY_WAIT_QUEUE;
        }

        __release_handshake(thread);
    }

    return 0;
}
