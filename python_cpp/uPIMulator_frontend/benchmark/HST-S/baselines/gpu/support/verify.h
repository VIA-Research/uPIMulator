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

#include "common.h"
#include <math.h>
#include <string.h>

inline int compare_output(unsigned int *outp, unsigned int *outpCPU, int bins) {
    for(int i = 0; i < bins; i++) {
        if(outp[i] != outpCPU[i]) {
            printf("Test failed\n");
            exit(EXIT_FAILURE);
        }
    }
    return 0;
}

// Sequential implementation for comparison purposes
inline void HistogramCPU(unsigned int *histo, unsigned int *data, int size, int bins) {
    for(int i = 0; i < size; i++) {
        // Read pixel
        unsigned int d = ((data[i] * bins) >> 12);
        // Vote in histogram
        histo[d]++;
    }
}

inline void verify(unsigned int *histo, unsigned int *input, int size, int bins) {
    unsigned int *gold = (unsigned int *)malloc(bins * sizeof(unsigned int));
    memset(gold, 0, bins * sizeof(unsigned int));
    HistogramCPU(gold, input, size, bins);
    compare_output(histo, gold, bins);
    free(gold);
}
