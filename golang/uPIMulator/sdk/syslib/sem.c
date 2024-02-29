/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <sem.h>

#include <defs.h>
#include <dpuruntime.h>

void
sem_take(struct sem_t *sem)
{
    unsigned char lock = sem->lock;
    __acquire(lock, "0");
    char count = sem->count - 1;
    thread_id_t tid = me();

    if (count < 0) {
        unsigned char last = sem->wait_queue;

        if (last != __EMPTY_WAIT_QUEUE) {
            unsigned char first = __WAIT_QUEUE_TABLE[last];
            __WAIT_QUEUE_TABLE[tid] = first;
            __WAIT_QUEUE_TABLE[last] = tid;
        } else {
            __WAIT_QUEUE_TABLE[tid] = tid;
        }

        sem->wait_queue = tid;
        sem->count = count;
        __release(lock, "0", __AT_NEXT_INSTRUCTION);
        __stop();
    } else {
        sem->count = count;
        __release(lock, "0", __AT_NEXT_INSTRUCTION);
    }
}

void
sem_give(struct sem_t *sem)
{
    unsigned char lock = sem->lock;
    __acquire(lock, "0");
    unsigned char count = sem->count + 1;
    unsigned char last = sem->wait_queue;

    if (last != __EMPTY_WAIT_QUEUE) {
        unsigned char first = __WAIT_QUEUE_TABLE[last];

        if (first == last) {
            sem->wait_queue = __EMPTY_WAIT_QUEUE;
        } else {
            __WAIT_QUEUE_TABLE[last] = __WAIT_QUEUE_TABLE[first];
        }

        __resume(first, "0");
    }

    sem->count = count;
    __release(lock, "0", __AT_NEXT_INSTRUCTION);
}
