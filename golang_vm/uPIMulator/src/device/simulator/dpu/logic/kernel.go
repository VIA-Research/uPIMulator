package logic

import (
	"uPIMulator/src/misc"
)

// Kernel handles syscalls and process management (MVP stub)
type Kernel struct {
	current_pid          int
	stat_factory         *misc.StatFactory
	test_mode            bool // Enable syscall injection for testing
	test_instruction_cnt int64
}

// Init initializes the kernel with a stat factory
func (this *Kernel) Init(stat_factory *misc.StatFactory) {
	this.current_pid = 0
	this.stat_factory = stat_factory
	this.test_mode = false
	this.test_instruction_cnt = 0
}

// SetTestMode enables syscall injection for measurement tests
func (this *Kernel) SetTestMode(enabled bool) {
	this.test_mode = enabled
	this.test_instruction_cnt = 0
}

// InjectTestSyscallIfNeeded checks if we should inject a test syscall
// (every 1000 instructions in test mode)
func (this *Kernel) InjectTestSyscallIfNeeded(thread *Thread) {
	if !this.test_mode {
		return
	}

	this.test_instruction_cnt++

	// Inject syscall every 1000 instructions
	if this.test_instruction_cnt%1000 == 0 {
		this.SyscallDispatcher(thread)
	}
}

// SyscallDispatcher handles syscalls when user code jumps to 0x0000
// For MVP: dispatch based on syscall number (in r0 convention)
// Returns true if syscall was handled (caller should not WritePcReg)
func (this *Kernel) SyscallDispatcher(thread *Thread) bool {
	// Increment syscall counter
	this.stat_factory.Increment("mmu_syscall_hits", 1)

	// TODO (Week 3–4): extract syscall number from registers
	// For now, just a NOP syscall that returns
	syscall_nop := func() {
		this.stat_factory.Increment("mmu_syscall_nop", 1)
	}

	syscall_nop()

	// Use the register-file helper to keep PC stepping semantics consistent.
	thread.RegFile().IncrementPcReg()

	return true // Handled: don't execute normal WritePcReg in ExecuteCallZri
}

// Placeholder for non-0x0000 jumps (future: support multiple dispatchers)
func (this *Kernel) IsDispatcherAddress(addr int64) bool {
	return addr == 0 // For MVP, only 0x0000 is the syscall entry
}
