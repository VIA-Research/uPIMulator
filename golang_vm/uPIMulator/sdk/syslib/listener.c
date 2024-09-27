/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include <stdint.h>
#include <stdbool.h>

#include <sysdef.h>
#include <defs.h>

extern bool
fetch_request(uint32_t fifo_info, uint32_t *request, uint32_t request_size);
extern bool
fifo_is_full(uint32_t fifo_info, uint32_t request_size);
extern void
fifo_produce(uint32_t fifo_info, uint32_t *request, uint32_t request_size);

static inline uint32_t
fifo_sys_fetch_info(uint32_t fid)
{
    extern uint32_t __sys_fifo_sys_table;
    return (&__sys_fifo_sys_table)[fid];
}

static inline uint32_t
fifo_fetch_info(sysname_t recipient)
{
    extern uint32_t __sys_fifo_table_ptr;
    return *((uint32_t *)((&__sys_fifo_table_ptr)[recipient] & 0xFFFF));
}

static inline sysname_t
fetch_recipient(uint32_t id)
{
    return id >> 24;
}

static inline sysname_t
fetch_request_id(uint32_t id)
{
    return id & 0x00FFFFFF;
}

void
__sys_internal_listener_loop(uint32_t *request, uint32_t request_size)
{
    sysname_t id = me();
    uint32_t self_fifo_info;
    uint32_t from_host_fifo_info;
    uint32_t to_host_fifo_info;

    self_fifo_info = fifo_fetch_info(id);
    from_host_fifo_info = fifo_sys_fetch_info(0);
    to_host_fifo_info = fifo_sys_fetch_info(1);

    while (true) {
        if (fetch_request(self_fifo_info, request, request_size)) {
            while (fifo_is_full(to_host_fifo_info, request_size)) {
                // Waiting for the recipient to read some of its pending requests...
                // Do we want to add some "sleep" here?
            }

            fifo_produce(to_host_fifo_info, request, request_size);
        }
        if (fetch_request(from_host_fifo_info, request, request_size)) {
            sysname_t recipient = fetch_recipient(request[0]);

            /* If a message is sent to the listener from the host, we interpret it as a shutdown order. */
            if (recipient == id)
                break;

            request[0] = fetch_request_id(request[0]);

            extern void internal_actor_send(uint32_t recipient, uint32_t * request, uint32_t request_size);
            internal_actor_send(recipient, request, request_size);
        }

        // Waiting for some request...
        // Do we want to add some "sleep" here?
    }
}
