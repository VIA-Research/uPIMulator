# Design: Software MMU prototype (IRAM/WRAM/MRAM) in uPIMulator

## 1) Entities
### Thread (existing)
uPIMulator has a `Thread` (tasklet) abstraction with regfile + PC.

### Process (to add)
Minimal fProc-like structure:
- pid, state
- saved thread contexts (regs/pc per tasklet or a subset)
- MRAM page table (vpn->ppn)
- WRAM heap metadata (heap_end, epoch)
- fault/exit info

### Kernel/Runtime (to add)
- process table
- current_pid
- scheduler (round robin)
- context switch (SWAP)
- syscall dispatcher (Option A trampoline)
- fault handling + statistics

## 2) IRAM (Priority #1)
### 2.1 Trampoline invariant (Option A)
- IRAM[0x0000] contains `JUMP dispatcher`.
- Installed at boot + re-installed after every context switch.

### 2.2 Loader guard
- When loading a process `.text` into user IRAM:
  - check `text_size <= IRAM_USER_SIZE`
  - if violation: mark process FAILED with FAULT_IRAM_OVERFLOW and refuse to run it.

## 3) MRAM virtualization (paging)
- PAGE_SIZE = 64KB
- single-level page table per process
- translation: vpn=vaddr/PAGE_SIZE, off=vaddr%PAGE_SIZE
- if unmapped -> allocate a physical MRAM page (ppn) and map.

Provide syscalls:
- sys_mram_read(vaddr, size, wram_dst)
- sys_mram_write(vaddr, size, wram_src)

## 4) WRAM protection + safe alloc
### 4.1 Boundary
- define WRAM_USER_BASE and WRAM_BOUNDARY
- user accesses must satisfy: base <= addr < boundary

### 4.2 Safe sbrk/malloc/free
- sys_sbrk(delta) updates heap_end only if:
  - heap_end + delta < WRAM_BOUNDARY
  - and does not collide with reserved stacks

### 4.3 Batch validation (amortized)
- sys_validate(ptr, len, perm)
- optional per-thread cache: (base,limit,perm,epoch)
- epoch++ on malloc/free/sbrk and on context switch.

## 5) Fault model
Fault codes (initial):
- FAULT_IRAM_OVERFLOW
- FAULT_WRAM_OOB
- FAULT_WRAM_HEAP_COLLISION
- FAULT_VALIDATE_FAIL
- FAULT_MRAM_INVALID

Policy:
- soft fault -> kill process (FAILED) + scheduler continues.
