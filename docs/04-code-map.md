# Code map (to be filled during exploration)

## Known hooks (already identified)
- WRAM accesses: OperandCollector -> wram.Read/Write
- Instruction dispatch: Logic.Execute (LW/SW/JUMP)
- DMA WRAM transfers: Dma.TransferToWram/TransferFromWram

## TODO (search targets)
- Where IRAM is loaded from `iram.bin`
- Where fetch reads IRAM (PC -> IRAM.Read)
- Where JUMP is executed, to detect `jump 0x0000` (syscall trampoline)
- Where cycle counters and stats are updated
- How programs/tasks are instantiated (multi-process representation)

## Notes
- Keep this file updated with exact paths + function names + line numbers.
