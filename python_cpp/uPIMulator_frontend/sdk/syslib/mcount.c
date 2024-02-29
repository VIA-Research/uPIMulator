/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

void __attribute__((naked, noinline, no_instrument_function)) mcount(void)
{
    // Please see ret_mcount comment regarding why mcount references ret_mcount.
    __asm__ volatile("jump ret_mcount");
}

void __attribute__((naked, noinline, no_instrument_function)) ret_mcount(void)
{
    // ret_mcount is used in statistics mode, mcount *must* reference ret_mcount
    // so that ret_mcount symbol is not gc (remember that we patch the binary
    // when copying it to iram).
    __asm__ volatile("sh id4, thread_profiling, r23\n"
                     "jump r23");
}
