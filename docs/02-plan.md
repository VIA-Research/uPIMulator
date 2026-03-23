# Plan (16 weeks / 4 months)

## Week 1: Exploration & code map
- Locate: IRAM load, fetch, JUMP execute, cycle counters, program/task instantiation.

## Weeks 2-3: IRAM protection MVP
- trampoline install/reinstall at 0x0000
- guard load for user .text
- tests: oversized .text triggers FAULT_IRAM_OVERFLOW

## Weeks 4-7: MRAM paging MVP
- per-process page table
- sys_mram_read/write
- page faults + allocation
- microbench cycles

## Weeks 8-11: WRAM safe alloc + validate
- WRAM boundary enforcement
- sys_sbrk + malloc/free
- sys_validate + cache + epoch invalidation

## Weeks 12-13: Evaluation
- 6-8 PrIM benchmarks subset
- adversarial tests (OOB, overflow IRAM)
- overhead <20% target discussion

## Weeks 14-16: Report & slides
- write-up + plots + cleanup
