#include <alloc.h>
#include <barrier.h>
#include <defs.h>
#include <mram.h>

__host int size_per_dpu;

BARRIER_INIT(my_barrier, NR_TASKLETS);

void vector_addition(int *A, int *B, int *C, int size_per_tasklet) {
  for (int i = 0; i < size_per_tasklet / sizeof(int); i++) {
    C[i] = A[i] + B[i];
  }
}

int main() {
  int tasklet_id = me();
  if (tasklet_id == 0) {
    mem_reset();
  }
  barrier_wait(&my_barrier);

  int size_per_tasklet = size_per_dpu / NR_TASKLETS;

  int *A_mram = (int *)(DPU_MRAM_HEAP_POINTER + tasklet_id * size_per_tasklet);
  int *B_mram = (int *)(DPU_MRAM_HEAP_POINTER + size_per_dpu + tasklet_id * size_per_tasklet);
  int *C_mram = (int *)(DPU_MRAM_HEAP_POINTER + 2 * size_per_dpu + tasklet_id * size_per_tasklet);

  int *A_wram = (int *)mem_alloc(size_per_tasklet);
  int *B_wram = (int *)mem_alloc(size_per_tasklet);
  int *C_wram = (int *)mem_alloc(size_per_tasklet);

  mram_read((__mram_ptr void *)A_mram, A_wram, size_per_tasklet);
  mram_read((__mram_ptr void *)B_mram, B_wram, size_per_tasklet);

  vector_addition(A_wram, B_wram, C_wram, size_per_tasklet);

  mram_write(C_wram, (__mram_ptr void *)C_mram, size_per_tasklet);

  return 0;
}
