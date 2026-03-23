# Experiments & metrics

## Metrics
- Total cycles per benchmark
- Overhead vs baseline (%)
- Breakdown cycles:
  - syscall cost
  - MRAM translation cost
  - context switch save/load WRAM/IRAM
  - validate calls
- Fault detection:
  - IRAM overflow
  - WRAM OOB

## Microbenchmarks
- empty syscall / dispatcher overhead
- validate(ptr,len) cost
- MRAM read/write 8B: baseline vs software VM

## Macrobenchmarks (6-8)
Choose representative subset covering:
- sync-heavy (e.g., SEL)
- DMA congestion-heavy (e.g., BS)
- compute/memory mix (GEMV, RED, BFS, NW, MLP, VA depending on availability)

## Concurrency scenario
- Run 2 processes with forced time-slicing to stress SWAP context switches.
