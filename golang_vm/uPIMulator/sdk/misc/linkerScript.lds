/* Copyright 2020 UPMEM. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/*
 * Simple linker script for embedded programs.
 */
ENTRY(__bootstrap)
MEMORY {
	iram (rx) : ORIGIN = 0x80000000, LENGTH = 32K
}
SECTIONS
{
  /*
   * Merge every text section into .text.
   */
  .text : {
    *(.text)
  } > iram
}

