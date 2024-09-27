/*
 * Copyright (c) 2016 University of Cordoba and University of Illinois
 * All rights reserved.
 *
 * Developed by:    IMPACT Research Group
 *                  University of Cordoba and University of Illinois
 *                  http://impact.crhc.illinois.edu/
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * with the Software without restriction, including without limitation the 
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
 * sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 *      > Redistributions of source code must retain the above copyright notice,
 *        this list of conditions and the following disclaimers.
 *      > Redistributions in binary form must reproduce the above copyright
 *        notice, this list of conditions and the following disclaimers in the
 *        documentation and/or other materials provided with the distribution.
 *      > Neither the names of IMPACT Research Group, University of Cordoba, 
 *        University of Illinois nor the names of its contributors may be used 
 *        to endorse or promote products derived from this Software without 
 *        specific prior written permission.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE 
 * CONTRIBUTORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS WITH
 * THE SOFTWARE.
 *
 */

#include "kernel.h"
#include <math.h>
#include <thread>
#include <vector>
#include <algorithm>

// CPU threads-----------------------------------------------------------------
void run_cpu_threads_100(T *input, std::atomic_int *finished, std::atomic_int *head, int A, int B, int b, int threads) {
///////////////// Run CPU worker threads /////////////////////////////////
#if PRINT
    printf("Starting %d CPU threads\n", threads);
#endif

    std::vector<std::thread> cpu_threads;
    for(int i = 0; i < threads; i++) {

        cpu_threads.push_back(std::thread([=]() {

            T   data[b];
            T   backup[b];
            int done;
            int m = A * B - 1;
            // Dynamic fetch
            int gid = (head)->fetch_add(1);

            while(gid < m) {
                int next_in_cycle = (gid * A) - m * (gid / B);
                if(next_in_cycle == gid) {
                    // Dynamic fetch
                    gid = (head)->fetch_add(1);
                    continue;
                }
                for(int i = 0; i < b; i++) {
                    data[i] = input[gid * b + i];
                }
                //make sure the read is not cached
                done = (finished + gid)->load();
                for(; done == 0; next_in_cycle = (next_in_cycle * A) - m * (next_in_cycle / B)) {
                    for(int i = 0; i < b; i++) {
                        backup[i] = input[next_in_cycle * b + i];
                    }
                    done = (finished + next_in_cycle)->exchange(1);
                    if(!done) {
                        for(int i = 0; i < b; i++) {
                            input[next_in_cycle * b + i] = data[i];
                        }
                    }
                    for(int i = 0; i < b; i++) {
                        data[i] = backup[i];
                    }
                }
                // Dynamic fetch
                gid = (head)->fetch_add(1);
            }
        }));
    }

    std::for_each(cpu_threads.begin(), cpu_threads.end(), [](std::thread &t) { t.join(); });
}


// CPU threads-----------------------------------------------------------------
void run_cpu_threads_010(T *input, std::atomic_int* head, int a, int b, int tiles, int threads) {
///////////////// Run CPU worker threads /////////////////////////////////
#if PRINT
    printf("Starting %d CPU threads\n", threads);
#endif

    std::vector<std::thread> cpu_threads;
    for(int i = 0; i < threads; i++) {

        cpu_threads.push_back(std::thread([=]() {

            T   tile[a * b];
            int m = a * b - 1;

            // Dynamic fetch
            int gid = (head)->fetch_add(1);

            while(gid < tiles) {
                T* input_array = input + a * b * gid;
                for (int j = 0; j < a * b; j++) {
                    int next = (j * a)-m*(j/b);
                    tile[next] = input_array[j];
                }
                for (int j = 0; j < a * b; j++) {
                    input_array[j] = tile[j];
                }
                // Dynamic fetch
                gid = (head)->fetch_add(1);
		    }
        }));
    }

    std::for_each(cpu_threads.begin(), cpu_threads.end(), [](std::thread &t) { t.join(); });
}
