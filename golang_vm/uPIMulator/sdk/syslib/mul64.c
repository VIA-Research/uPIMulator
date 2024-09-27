/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/*
 * 64x64 multiplication emulation.
 *
 * A relatively fast emulation of 64x64 multiplication using byte multipliers.
 * Basically, the two operands X and Y are seen as byte polynomials:
 *  - X = X0.2^0 + X1.2^8 + X2.2^16 + X3.2^24 + X4.2^32 + X5.2^40 + X6.2^48 + X7.2^56
 *  - Y = Y0.2^0 + Y1.2^8 + Y2.2^16 + Y3.2^24 + Y4.2^32 + Y5.2^40 + Y6.2^48 + Y7.2^56
 *
 * The product Z is expressed as a similar polynomial. Since the result is 64 bits,
 * the function drops any coefficient for a power greater than 56, hence the following
 * formula:
 *  Z = (X0.Y0).2^0
 *      + (X0.Y1 + X1.Y0).2^8
 *      + (X0.Y2 + X2.Y0 + X1.Y1).2^16
 *      + (X0.Y3 + X1.Y2 + X2.Y1 + X3.Y0).2^24
 *      + (X0.Y4 + X1.Y3 + X2.Y2 + X3.Y1 + X4.Y0).2^32
 *      etc.
 *
 * Each individual produce is computed with the native built-in 8x8 instructions.
 * Resulting processing time is in the magnitude of 150 instructions.
 *
 * The two operands are found in __D0 and the first kernel nano-stack entry.
 * The result goes into __R0 (lsbits) and __R1 (msbits).
 * Also, __R2 contains the return address register, instead of __RET__.
 */
#include <stdint.h>

static uint16_t
_mul00(uint32_t a, uint32_t b)
{
#ifndef DPU
    return (a & 0xff) * (b & 0xff);
#else
    uint32_t r;
    __asm__ volatile("mul_ul_ul %[rc_wr32], %[ra_r32], %[rb_wr32]" : [rc_wr32] "=r"(r) : [ra_r32] "r"(a), [rb_wr32] "r"(b) :);
    return r;
#endif
}

static uint16_t
_mul01(uint32_t a, uint32_t b)
{
#ifndef DPU
    return (a & 0xff) * ((b >> 8) & 0xff);
#else
    uint32_t r;
    __asm__ volatile("mul_ul_uh %[rc_wr32], %[ra_r32], %[rb_wr32]" : [rc_wr32] "=r"(r) : [ra_r32] "r"(a), [rb_wr32] "r"(b) :);
    return r;
#endif
}

#define _mul02(a, b) _mul00(a, (b >> 16))
#define _mul03(a, b) _mul01(a, (b >> 16))

static uint16_t
_mul11(uint32_t a, uint32_t b)
{
#ifndef DPU
    return ((a >> 8) & 0xff) * ((b >> 8) & 0xff);
#else
    uint32_t r;
    __asm__ volatile("mul_uh_uh %[rc_wr32], %[ra_r32], %[rb_wr32]" : [rc_wr32] "=r"(r) : [ra_r32] "r"(a), [rb_wr32] "r"(b) :);
    return r;
#endif
}

static uint16_t
_mul12(uint32_t a, uint32_t b)
{
#ifndef DPU
    return ((a >> 8) & 0xff) * ((b >> 16) & 0xff);
#else
    uint32_t r = (b >> 16);
    __asm__ volatile("mul_uh_ul %[rc_wr32], %[ra_r32], %[rb_wr32]" : [rc_wr32] "=r"(r) : [ra_r32] "r"(a), [rb_wr32] "r"(r) :);
    return r;
#endif
}

#define _mul13(a, b) _mul11(a, (b >> 16))
#define _mul22(a, b) _mul00((a >> 16), (b >> 16))
#define _mul23(a, b) _mul01((a >> 16), (b >> 16))
#define _mul33(a, b) _mul11((a >> 16), (b >> 16))

#define mulx0y0(xl, yl) _mul00(xl, yl)
#define mulx0y1(xl, yl) _mul01(xl, yl)
#define mulx0y2(xl, yl) _mul02(xl, yl)
#define mulx0y3(xl, yl) _mul03(xl, yl)
#define mulx0y4(xl, yh) _mul00(xl, yh)
#define mulx0y5(xl, yh) _mul01(xl, yh)
#define mulx0y6(xl, yh) _mul02(xl, yh)
#define mulx0y7(xl, yh) _mul03(xl, yh)

#define mulx1y1(xl, yl) _mul11(xl, yl)
#define mulx1y2(xl, yl) _mul12(xl, yl)
#define mulx1y3(xl, yl) _mul13(xl, yl)
#define mulx1y4(xl, yh) _mul01(yh, xl)
#define mulx1y5(xl, yh) _mul11(xl, yh)
#define mulx1y6(xl, yh) _mul12(xl, yh)

#define mulx2y2(xl, yl) _mul22(xl, yl)
#define mulx2y3(xl, yl) _mul23(xl, yl)
#define mulx2y4(xl, yh) _mul02(yh, xl)
#define mulx2y5(xl, yh) _mul12(yh, xl)

#define mulx3y3(xl, yl) _mul33(xl, yl)
#define mulx3y4(xl, yh) _mul03(yh, xl)

// Symmetry...
#define mulx1y0(xl, yl) mulx0y1(yl, xl)
#define mulx2y0(xl, yl) mulx0y2(yl, xl)
#define mulx2y1(xl, yl) mulx1y2(yl, xl)
#define mulx3y0(xl, yl) mulx0y3(yl, xl)
#define mulx3y1(xl, yl) mulx1y3(yl, xl)
#define mulx3y2(xl, yl) mulx2y3(yl, xl)
#define mulx4y0(xh, yl) mulx0y4(yl, xh)
#define mulx4y1(xh, yl) mulx1y4(yl, xh)
#define mulx4y2(xh, yl) mulx2y4(yl, xh)
#define mulx4y3(xh, yl) mulx3y4(yl, xh)
#define mulx5y0(xh, yl) mulx0y5(yl, xh)
#define mulx5y1(xh, yl) mulx1y5(yl, xh)
#define mulx5y2(xh, yl) mulx2y5(yl, xh)
#define mulx6y0(xh, yl) mulx0y6(yl, xh)
#define mulx6y1(xh, yl) mulx1y6(yl, xh)
#define mulx7y0(xh, yl) mulx0y7(yl, xh)

uint64_t
__muldi3(uint64_t x, uint64_t y)
{
    uint32_t xl = x;
    uint32_t xh = ((uint64_t)x >> 32);
    uint32_t yl = y;
    uint32_t yh = ((uint64_t)y >> 32);

    // Each fragment of the product.
    uint32_t p0, p1, p2, p3, p4, p5, p6, p7, rh;
    uint64_t rl;

    p0 = mulx0y0(xl, yl);
    rl = (uint64_t)p0;
    p1 = mulx0y1(xl, yl) + mulx1y0(xl, yl);
    rl += ((uint64_t)p1 << 8);
    p2 = mulx0y2(xl, yl) + mulx2y0(xl, yl) + mulx1y1(xl, yl);
    rl += ((uint64_t)p2 << 16);
    p3 = mulx0y3(xl, yl) + mulx3y0(xl, yl) + mulx1y2(xl, yl) + mulx2y1(xl, yl);
    rl += ((uint64_t)p3 << 24);
    p4 = mulx0y4(xl, yh) + mulx4y0(xh, yl) + mulx1y3(xl, yl) + mulx3y1(xl, yl) + mulx2y2(xl, yl);
    rh = p4;
    p5 = mulx0y5(xl, yh) + mulx5y0(xh, yl) + mulx1y4(xl, yh) + mulx4y1(xh, yl) + mulx2y3(xl, yl) + mulx3y2(xl, yl);
    rh += p5 << 8;
    p6 = mulx0y6(xl, yh) + mulx6y0(xh, yl) + mulx1y5(xl, yh) + mulx5y1(xh, yl) + mulx2y4(xl, yh) + mulx4y2(xh, yl)
        + mulx3y3(xl, yl);
    rh += p6 << 16;
    p7 = mulx0y7(xl, yh) + mulx7y0(xh, yl) + mulx1y6(xl, yh) + mulx6y1(xh, yl) + mulx2y5(xl, yh) + mulx5y2(xh, yl)
        + mulx3y4(xl, yh) + mulx4y3(xh, yl);
    rh += p7 << 24;

    return rl + (((uint64_t)rh) << 32);
}
