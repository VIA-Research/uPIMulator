#include <defs.h>

// Force a jump to address 0x0000 to trigger the MMU syscall hook
void mmu_test_force_syscall(void) {
  // Create a function pointer to address 0x0000
  // This will force an indirect CALL that targets 0x0000
  typedef void (*syscall_t)(void);
  syscall_t syscall_entry = (syscall_t)0x0000;
  
  // Do some work before the syscall
  volatile int x = 42;
  x = x + 1;
  
  // Force the call to 0x0000 - this should hit the MMU hook
  syscall_entry();
}

int main(void) {
  // Run a few iterations to measure syscall overhead
  for (int i = 0; i < 3; i++) {
    mmu_test_force_syscall();
  }
  return 0;
}
