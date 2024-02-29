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

inline int compare_output(T *output, T *ref, int dim) {
    int i;
    for(i = 0; i < dim; i++) {
        T diff = fabs(ref[i] - output[i]);
        if((diff - 0.0f) > 0.00001f && diff > 0.01 * fabs(ref[i])) {
            printf("line: %d ref: %f actual: %f diff: %f\n", i, ref[i], output[i], diff);
            exit(EXIT_FAILURE);
        }
    }
    return 0;
}

// Sequential transposition for comparison purposes
//[w][h/t][t] to [h/t][w][t]
static void trns_host(T* input, unsigned int A, unsigned int B, unsigned int b){
   T* output = (T*) malloc(sizeof(T) * A * B * b);
   unsigned int next;
   for (unsigned int j = 0; j < b; j++){
      for (unsigned int i = 0; i < A * B; i++){
         next = (i * A) - (A * B - 1) * (i / B);
         output[next * b + j] = input[i*b+j];
      }
   }
   for (unsigned int k = 0; k < A * B * b; k++){
      input[k] = output[k];
   }
   free(output);
}

inline void verify(T *input2, T *input, int height, int width, int tile_size) {
    trns_host(input, height, width, tile_size);
    compare_output(input2, input, height * width);
}
