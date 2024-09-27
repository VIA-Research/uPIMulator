/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

void __attribute__((naked, used, section(".text.__bootstrap"))) __bootstrap()
{
    // Preconditions:
    //  - MRAM offset is a multiple of 8
    //  - Buffer size is a multiple of 8

    __asm__ volatile("  sd zero, 0, d0\n" // Saving context
                     "  sd zero, 8, d2\n"
                     "  or r0, zero, 0, ?xnz, . + 2\n"
                     "  or r0, r0, 0x2\n"
                     "  addc r0, r0, 0\n"
                     "  sw zero, 16, r0\n"
                     "  lw r0, zero, 20\n" // MRAM offset, must be patched by the Host
                     "resume_start:\n"
                     "  lw r1, zero, 24\n" // Buffer size, must be patched by the Host
                     "  move r2, 32\n" // Wram offset
                     "  move r3, 2048\n" // Transfer size
                     "  transfer_loop:\n"
                     "  jltu r1, r3, last_transfer\n"
                     "  ldma r2, r0, 255\n" // Can be patched by the Host into a SDMA to write MRAM
                     "  add r0, r0, r3\n"
                     "  add r2, r2, r3\n"
                     "  sub r1, r1, r3, true, transfer_loop\n"
                     "last_transfer:\n"
                     "  jz r1, end\n"
                     "  lsr r3, r1, 3\n"
                     "  add r3, r3, -1\n"
                     "  lsl_add r2, r2, r3, 24\n"
                     "  ldma r2, r0, 0\n" // Can be patched by the Host into a SDMA to write MRAM
                     "  add r0, r0, r1\n"
                     "end:\n"
                     "  lw r2, zero, 28\n" // Restoring context if needed
                     "  jnz r2, . + 2\n"
                     "  stop true, resume_start\n"
                     "  ld d2, zero, 8\n"
                     "  lw r0, zero, 16\n"
                     "  add r1, r0, r0\n"
                     "  add r0, r0, r1\n"
                     "  call zero, r0, . + 1\n"
                     "  add r0, zero, 0x00000001; ld d0, zero, 0; stop true, 0\n" // ... restore Z = 0, C = 0
                     "  add r0, mneg, 0x80000001; ld d0, zero, 0; stop true, 0\n" // ... restore Z = 0, C = 1
                     "  add r0, zero, 0x00000000; ld d0, zero, 0; stop true, 0\n" // ... restore Z = 1, C = 0
                     "  add r0, mneg, 0x80000000; ld d0, zero, 0; stop true, 0\n" // ... restore Z = 1, C = 1
    );
}
