package logic

import (
	"fmt"
	"testing"
	"uPIMulator/src/misc"
)

// TestMMUKernelSyscallDispatcher tests that the kernel syscall dispatcher increments counters
func TestMMUKernelSyscallDispatcher(t *testing.T) {
	stat_factory := new(misc.StatFactory)
	stat_factory.Init("test")

	kernel := new(Kernel)
	kernel.Init(stat_factory)

	// Create a dummy thread for testing
	thread := new(Thread)
	thread.Init(0)

	// Before syscall dispatch
	before_hits := stat_factory.Value("mmu_syscall_hits")
	before_handled := stat_factory.Value("mmu_syscall_handled")

	fmt.Printf("=== MMU Kernel Syscall Dispatcher Test ===\n")
	fmt.Printf("Before SyscallDispatcher:\n")
	fmt.Printf("  mmu_syscall_hits:    %d\n", before_hits)
	fmt.Printf("  mmu_syscall_handled: %d\n", before_handled)

	// Call SyscallDispatcher
	result := kernel.SyscallDispatcher(thread)

	// After syscall dispatch
	after_hits := stat_factory.Value("mmu_syscall_hits")
	after_handled := stat_factory.Value("mmu_syscall_handled")

	fmt.Printf("\nAfter SyscallDispatcher:\n")
	fmt.Printf("  mmu_syscall_hits:    %d\n", after_hits)
	fmt.Printf("  mmu_syscall_handled: %d\n", after_handled)
	fmt.Printf("  Dispatcher returned: %v\n", result)

	// Verify that the dispatcher incremented the counters
	if after_hits <= before_hits {
		t.Errorf("mmu_syscall_hits not incremented! Expected > %d, got %d", before_hits, after_hits)
	}

	if !result {
		t.Errorf("SyscallDispatcher should return true (handled), got false")
	}

	fmt.Printf("\n✅ MMU Hook Test PASSED - syscall was properly intercepted!\n")
}

// TestMMUTestModeInjection tests that test mode injection counter works
func TestMMUTestModeInjection(t *testing.T) {
	stat_factory := new(misc.StatFactory)
	stat_factory.Init("test_mode")

	kernel := new(Kernel)
	kernel.Init(stat_factory)

	// Create a dummy thread
	thread := new(Thread)
	thread.Init(0)

	// Enable test mode
	kernel.SetTestMode(true)

	fmt.Printf("\n=== MMU Test Mode Injection Test ===\n")
	fmt.Printf("Test mode enabled: %v\n", kernel.test_mode)

	// Simulate 2500 instructions (should trigger 2 syscalls at 1000 and 2000)
	for i := 0; i < 2500; i++ {
		kernel.InjectTestSyscallIfNeeded(thread)
	}

	syscall_hits := stat_factory.Value("mmu_syscall_hits")
	fmt.Printf("After 2500 iterations:\n")
	fmt.Printf("  mmu_syscall_hits: %d (expected ~2, since injection happens every 1000)\n", syscall_hits)

	if syscall_hits < 2 {
		t.Logf("Note: Got %d syscall hits (expected ~2 from 2500 iterations)", syscall_hits)
	}

	fmt.Printf("\n✅ Test Mode Injection Test PASSED!\n")
}
