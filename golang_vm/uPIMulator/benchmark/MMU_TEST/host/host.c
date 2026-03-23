#include <stdio.h>
#include <stdint.h>
#include "dpu.h"

// Minimal host program for MMU_TEST benchmark
// Just loads and runs the DPU task

int main(int argc, char *argv[]) {
  struct dpu_set_t set, dpu;
  
  if (dpu_alloc(1, NULL, &set) != DPU_OK) {
    printf("dpu_alloc failed\n");
    return 1;
  }
  
  if (dpu_load(set, "task_mmu_test", NULL) != DPU_OK) {
    printf("dpu_load failed\n");
    dpu_free(set);
    return 1;
  }
  
  if (dpu_launch(set, DPU_SYNCHRONOUS) != DPU_OK) {
    printf("dpu_launch failed\n");
    dpu_free(set);
    return 1;
  }
  
  printf("MMU_TEST completed\n");
  
  dpu_free(set);
  return 0;
}
