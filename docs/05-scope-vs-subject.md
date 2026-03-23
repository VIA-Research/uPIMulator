# Scope vs Subject (traceability matrix)

This page maps the internship subject requirements to what we implement in the uPIMulator prototype, what is measured, and what is left as future work.

## Context / platform choice
The subject states: *"integrated in Faho or an equivalent prototype"*. We implement an equivalent prototype inside **uPIMulator (golang_vm)** to obtain cycle-level, reproducible measurements without requiring UPMEM hardware and without relying on an unavailable Faho-uPIMulator integration branch.

---

## Requirement 1 — Memory virtualization (WRAM/MRAM): private address space illusion

### What we implement (MVP)
- **MRAM**: true virtual memory by software paging:
  - per-process page table (single-level)
  - VA -> PA translation
  - 64KB pages
  - page fault allocation
- **WRAM**: practical virtualization consistent with the SWAP model:
  - at any time, only one process is resident in WRAM user region
  - context switch saves/restores process WRAM image via MRAM
  - user/kernel boundary enforced

### What we measure
- MRAM access overhead: baseline vs software paging
- context switch overhead (save/load WRAM image)
- overall benchmark overhead vs baseline

### Limitations / future work
- Fine-grained WRAM virtual addressing per LOAD/STORE is not realistic on current UPMEM hardware without compiler instrumentation or hardware support.

---

## Requirement 2 — Memory protection: forbid a process from accessing another process's regions

### What we implement (MVP)
- **MRAM**: strong inter-process protection via per-process translation (same VA maps to different PAs).
- **WRAM/IRAM**: protection via SWAP + boundaries:
  - kernel regions protected from user overwrites
  - process cannot access other processes' WRAM/IRAM concurrently because they are not resident

### What we measure
- fault detection rate for invalid MRAM accesses (if enabled)
- enforcement correctness on boundary tests

### Limitations / future work
- Attack model with pointer forging in C is out of scope for MVP; stronger isolation may require compiler instrumentation (bounds checks) or hardware support.

---

## Requirement 3 — Error detection: trigger a software exception on violations

### What we implement (MVP)
- Soft faults:
  - IRAM overflow at load time
  - WRAM out-of-bounds / heap collision
  - validate() failure
  - MRAM invalid access (policy dependent)
- Fault handling policy:
  - mark the process FAILED
  - stop its execution and schedule the next process
  - record fault statistics

### What we measure
- number of detected faults in adversarial tests
- overhead of checks (validate, sbrk, MRAM translation)

---

## Requirement 4 — Dynamic allocation: safe malloc/free in multi-process context

### What we implement (MVP)
- WRAM allocator with bounds:
  - sys_sbrk
  - minimal malloc/free on top of heap_end
  - epoch invalidation for validation cache
- MRAM allocator:
  - allocation of MRAM pages through the paging subsystem

### What we measure
- allocation overhead (cycles)
- fragmentation / failure cases (qualitative)

---

## Requirement — "Intercept memory accesses of fProcs"
### Clarification (important)
- On real UPMEM DPUs, trapping every WRAM LOAD/STORE is not possible (no hardware trap/MMU).
- Our prototype intercepts **mediated points** (syscalls, allocators, MRAM VM, context switch) and uses **batch validation**.
- In uPIMulator, we may optionally instrument LW/SW as a *measurement-only* baseline to quantify the cost of naive per-access checks.
