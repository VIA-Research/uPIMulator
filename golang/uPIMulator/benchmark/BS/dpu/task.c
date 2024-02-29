/*
* Binary Search with multiple tasklets
*
*/
#include <stdint.h>
#include <stdio.h>
#include <defs.h>
#include <mram.h>
#include <alloc.h>
#include <mram.h>
#include <barrier.h>
#include <perfcounter.h>
#include "common.h"

__host dpu_arguments_t DPU_INPUT_ARGUMENTS;
__host dpu_results_t DPU_RESULTS[NR_TASKLETS];

// Search
DTYPE __attribute__ ((noinline)) search(DTYPE *bufferA, DTYPE searching_for) {
  DTYPE found = -2;
  if(bufferA[0] <= searching_for)
  {
    found = -1;
    for (uint32_t i = 0; i < BLOCK_SIZE / sizeof(DTYPE); i++){
      if(bufferA[i] == searching_for)
      {
        found = i;
        break;
      }
    }
  }
  return found;
}

BARRIER_INIT(my_barrier, NR_TASKLETS);

extern int main_kernel1(void);

int(*kernels[nr_kernels])(void) = {main_kernel1};

int main(void){
  // Kernel
  return kernels[DPU_INPUT_ARGUMENTS.kernel]();
}

// main_kernel1
int main_kernel1() {
  unsigned int tasklet_id = me();
  #if PRINT
  printf("tasklet_id = %u\n", tasklet_id);
  #endif
  if(tasklet_id == 0){
    mem_reset(); // Reset the heap
  }
  // Barrier
  barrier_wait(&my_barrier);

  DTYPE searching_for, found;
  uint64_t input_size = DPU_INPUT_ARGUMENTS.input_size;

  // Address of the current processing block in MRAM
  uint32_t start_mram_block_addr_A       = (uint32_t) DPU_MRAM_HEAP_POINTER;
  uint32_t start_mram_block_addr_aux     = start_mram_block_addr_A;
  uint32_t end_mram_block_addr_A         = start_mram_block_addr_A + sizeof(DTYPE) * input_size;
  uint32_t current_mram_block_addr_query = end_mram_block_addr_A + tasklet_id * (DPU_INPUT_ARGUMENTS.slice_per_dpu / NR_TASKLETS) * sizeof(DTYPE);

  // Initialize a local cache to store the MRAM block
  DTYPE *cache_A     = (DTYPE *) mem_alloc(BLOCK_SIZE);
  DTYPE *cache_aux_A = (DTYPE *) mem_alloc(BLOCK_SIZE);
  DTYPE *cache_aux_B = (DTYPE *) mem_alloc(BLOCK_SIZE);

  dpu_results_t *result = &DPU_RESULTS[tasklet_id];

  // TODO(bongjoon.hyun@gmail.com): original PrIM benchmark uses uint64_t for targets' type
  for(uint32_t targets = 0; targets < (DPU_INPUT_ARGUMENTS.slice_per_dpu / NR_TASKLETS); targets++)
  {
    found = -1;

    mram_read((__mram_ptr void const *) current_mram_block_addr_query, &searching_for, 8);
    current_mram_block_addr_query += 8;

    bool end = false;

    // Initialize input vector boundaries
    start_mram_block_addr_A    = (uint32_t) DPU_MRAM_HEAP_POINTER;
    start_mram_block_addr_aux  = start_mram_block_addr_A;
    end_mram_block_addr_A      = start_mram_block_addr_A + sizeof(DTYPE) * input_size;

    uint32_t current_mram_block_addr_A = start_mram_block_addr_A;

    // Bring first and last values to WRAM
    mram_read((__mram_ptr void const *) current_mram_block_addr_A, cache_aux_A, BLOCK_SIZE);
    mram_read((__mram_ptr void const *) (end_mram_block_addr_A - BLOCK_SIZE * sizeof(DTYPE)),   cache_aux_B, BLOCK_SIZE);

    current_mram_block_addr_A = (start_mram_block_addr_A + end_mram_block_addr_A) / 2;
    while(!end)
    {
      // Load cache with current MRAM block
      mram_read((__mram_ptr void const *) current_mram_block_addr_A, cache_A, BLOCK_SIZE);

      // Search inside block
      found = search(cache_A, searching_for);

      // If found > -1, we found the searching_for query
      if(found > -1)
      {
        result->found = found + (current_mram_block_addr_A - start_mram_block_addr_aux) / sizeof(DTYPE);
        break;
      }

      // If found == -2, we need to discard right part of the input vector
      if(found == -2)
      {
        end_mram_block_addr_A     = current_mram_block_addr_A;
        current_mram_block_addr_A = (current_mram_block_addr_A + start_mram_block_addr_A) / 2;
      }

      // If found == -1, we need to discard left part of the input vector
      else if (found == -1)
      {
        start_mram_block_addr_A   = current_mram_block_addr_A;
        current_mram_block_addr_A = (current_mram_block_addr_A + end_mram_block_addr_A) / 2;
      }

      // Start boundary check
      if(current_mram_block_addr_A < (start_mram_block_addr_aux + BLOCK_SIZE))
      {
        end = true;
        mram_read((__mram_ptr void const *) current_mram_block_addr_A, cache_A, BLOCK_SIZE);
        found = search(cache_A, searching_for);

        if(found > -1)
        {
          end = true;
          result->found = found + (current_mram_block_addr_A - start_mram_block_addr_aux) / sizeof(DTYPE);
        }
      }

      // End boundary check
      if(current_mram_block_addr_A > (end_mram_block_addr_A - BLOCK_SIZE))
      {
        end = true;
        mram_read((__mram_ptr void const *) end_mram_block_addr_A - BLOCK_SIZE, cache_A, BLOCK_SIZE);
        found = search(cache_A, searching_for);

        if(found > -1)
        {
          result->found = found + (current_mram_block_addr_A - start_mram_block_addr_aux) / sizeof(DTYPE);
        }
      }
    }
  }
  return 0;
}
