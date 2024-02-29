/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/*
 * Prototype of function can be found here: https://llvm.org/docs/Atomics.html
 */

#include <dpuruntime.h>
#include <atomic_bit.h>

#define ATOMIC_BIT llvm_atomic_functions
ATOMIC_BIT_INIT(ATOMIC_BIT);

#define FOR_ALL_TYPES(fct) fct(1, char) fct(2, short) fct(4, int) fct(8, long long)

#define PROLOGUE(ptr, load, n_type)                                                                                              \
    n_type load;                                                                                                                 \
    ATOMIC_BIT_ACQUIRE(ATOMIC_BIT);                                                                                              \
    load = *ptr;

#define EPILOGUE(load)                                                                                                           \
    ATOMIC_BIT_RELEASE(ATOMIC_BIT);                                                                                              \
    return load;

#define __SYNC_VAL_COMPARE_AND_SWAP_N(n_val, n_type)                                                                             \
    n_type __dpu_sync_val_compare_and_swap_##n_val(volatile n_type *ptr, n_type expected, n_type desired)                        \
    {                                                                                                                            \
        PROLOGUE(ptr, load, n_type);                                                                                             \
        if (load == expected)                                                                                                    \
            *ptr = desired;                                                                                                      \
        EPILOGUE(load);                                                                                                          \
    }
FOR_ALL_TYPES(__SYNC_VAL_COMPARE_AND_SWAP_N)

#define __SYNC_LOCK_TEST_AND_SET_N(n_val, n_type)                                                                                \
    n_type __dpu_sync_lock_test_and_set_##n_val(volatile n_type *ptr, n_type val)                                                \
    {                                                                                                                            \
        PROLOGUE(ptr, load, n_type);                                                                                             \
        *ptr = val;                                                                                                              \
        EPILOGUE(load);                                                                                                          \
    }
FOR_ALL_TYPES(__SYNC_LOCK_TEST_AND_SET_N)

#define __SYNC_FETCH_AND_DO_N(fct, fct_name, n_val, n_type)                                                                      \
    n_type __dpu_sync_fetch_and_##fct_name##_##n_val(volatile n_type *ptr, n_type val)                                           \
    {                                                                                                                            \
        PROLOGUE(ptr, load, n_type);                                                                                             \
        *ptr = fct(load, val, n_type);                                                                                           \
        EPILOGUE(load);                                                                                                          \
    }

#define DO_ADD(a, b, n_type) ((a) + (b))
#define __SYNC_FETCH_AND_ADD_N(n_val, n_type) __SYNC_FETCH_AND_DO_N(DO_ADD, add, n_val, n_type)
FOR_ALL_TYPES(__SYNC_FETCH_AND_ADD_N)

#define DO_SUB(a, b, n_type) ((a) - (b))
#define __SYNC_FETCH_AND_SUB_N(n_val, n_type) __SYNC_FETCH_AND_DO_N(DO_SUB, sub, n_val, n_type)
FOR_ALL_TYPES(__SYNC_FETCH_AND_SUB_N)

#define DO_AND(a, b, n_type) ((a) & (b))
#define __SYNC_FETCH_AND_AND_N(n_val, n_type) __SYNC_FETCH_AND_DO_N(DO_AND, and, n_val, n_type)
FOR_ALL_TYPES(__SYNC_FETCH_AND_AND_N)

#define DO_OR(a, b, n_type) ((a) | (b))
#define __SYNC_FETCH_AND_OR_N(n_val, n_type) __SYNC_FETCH_AND_DO_N(DO_OR, or, n_val, n_type)
FOR_ALL_TYPES(__SYNC_FETCH_AND_OR_N)

#define DO_XOR(a, b, n_type) ((a) ^ (b))
#define __SYNC_FETCH_AND_XOR_N(n_val, n_type) __SYNC_FETCH_AND_DO_N(DO_XOR, xor, n_val, n_type)
FOR_ALL_TYPES(__SYNC_FETCH_AND_XOR_N)

#define DO_NAND(a, b, n_type) (~((a) & (b)))
#define __SYNC_FETCH_AND_NAND_N(n_val, n_type) __SYNC_FETCH_AND_DO_N(DO_NAND, nand, n_val, n_type)
FOR_ALL_TYPES(__SYNC_FETCH_AND_NAND_N)

#define DO_MAX(a, b, n_type) ((a) > (b) ? (a) : (b))
#define __SYNC_FETCH_AND_MAX_N(n_val, n_type) __SYNC_FETCH_AND_DO_N(DO_MAX, max, n_val, n_type)
FOR_ALL_TYPES(__SYNC_FETCH_AND_MAX_N)

#define DO_UMAX(a, b, n_type) (((unsigned n_type)(a)) > ((unsigned n_type)(b)) ? (a) : (b))
#define __SYNC_FETCH_AND_UMAX_N(n_val, n_type) __SYNC_FETCH_AND_DO_N(DO_UMAX, umax, n_val, n_type)
FOR_ALL_TYPES(__SYNC_FETCH_AND_UMAX_N)

#define DO_MIN(a, b, n_type) ((a) > (b) ? (a) : (b))
#define __SYNC_FETCH_AND_MIN_N(n_val, n_type) __SYNC_FETCH_AND_DO_N(DO_MIN, min, n_val, n_type)
FOR_ALL_TYPES(__SYNC_FETCH_AND_MIN_N)

#define DO_UMIN(a, b, n_type) (((unsigned n_type)(a)) > ((unsigned n_type)(b)) ? (a) : (b))
#define __SYNC_FETCH_AND_UMIN_N(n_val, n_type) __SYNC_FETCH_AND_DO_N(DO_UMIN, umin, n_val, n_type)
FOR_ALL_TYPES(__SYNC_FETCH_AND_UMIN_N)
