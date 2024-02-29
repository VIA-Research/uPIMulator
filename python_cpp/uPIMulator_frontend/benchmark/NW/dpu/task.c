/**
* Needleman-Wunsch with multiple tasklets
*
*/
#include <stdint.h>
#include <stdio.h>
#include <defs.h>
#include <mram.h>
#include <alloc.h>
#include <perfcounter.h>
#include <barrier.h>

#include "../support/common.h"

__host dpu_arguments_t DPU_INPUT_ARGUMENTS;

// Barrier
BARRIER_INIT(my_barrier, NR_TASKLETS);

// main
int main() {
    unsigned int tasklet_id = me();
    if (tasklet_id == 0){ // Initialize once the cycle counter
        mem_reset(); // Reset the heap
    }
    // Barrier
    barrier_wait(&my_barrier);
    uint32_t nblocks = DPU_INPUT_ARGUMENTS.nblocks;
    uint32_t active_blocks = DPU_INPUT_ARGUMENTS.active_blocks;
    uint32_t penalty = DPU_INPUT_ARGUMENTS.penalty;
#if PRINT
    printf("tasklet_id = %d, nblocks = %d \n", tasklet_id, nblocks);
#endif
	
    uint32_t mram_base_addr_input_itemsets = (uint32_t) (DPU_MRAM_HEAP_POINTER);
    uint32_t mram_base_addr_ref = (uint32_t) (DPU_MRAM_HEAP_POINTER + nblocks * (BL+1) * (BL+2) * sizeof(int32_t));
    if (nblocks != active_blocks)
        mram_base_addr_ref = (uint32_t) (DPU_MRAM_HEAP_POINTER + active_blocks * (BL+1) * (BL+2) * sizeof(int32_t));

    int32_t *cache_input = mem_alloc((BL_IN+1) * (BL_IN+2) * sizeof(int32_t));
    int32_t *cache_ref = mem_alloc(BL_IN * BL_IN * sizeof(int32_t));
    uint32_t REP = BL/BL_IN;
    uint32_t chunks;
    uint32_t mod;
    uint32_t start;
    uint32_t addr_input;
    uint32_t addr_ref;
    uint32_t cache_input_offset;

    for (uint32_t bl = 0; bl < nblocks; bl++) {

        // Top-left computation
        for(uint32_t blk = 0; blk <= REP; blk++) {
            
            // Partition chunks/subblocks of the diagonal to tasklets 
            chunks = blk / NR_TASKLETS; 
            mod = blk % NR_TASKLETS;
            if (tasklet_id < mod)
                chunks++;
            if (mod > 0) {
                if(tasklet_id < mod)
                    start = tasklet_id * chunks;
                else
                    start = mod * (chunks + 1) + (tasklet_id - mod) * chunks;
            } else
                start = tasklet_id * chunks;
            
            // Compute all assigned chunks  
            for (uint32_t bl_indx = 0; bl_indx < chunks; bl_indx++) {
                int t_index_x = start + bl_indx;
                int t_index_y = blk - 1 - t_index_x; 
                
                // Move input from MRAM to WRAM
                addr_input =  mram_base_addr_input_itemsets + (t_index_x * (BL+2) * BL_IN * sizeof(int32_t)) + (t_index_y * BL_IN * sizeof(int32_t));
                cache_input_offset = (BL_IN+2);
                mram_read((__mram_ptr void const *) addr_input, (void *) cache_input, (BL_IN+2) * sizeof(int32_t)); 
                addr_input += ((BL+2) * sizeof(int32_t));
                for (int i = 1; i < BL_IN + 1; i++) {
                    mram_read((__mram_ptr void const *) addr_input, (void *) (cache_input + cache_input_offset), (2) * sizeof(int32_t)); 
                    cache_input_offset += (BL_IN+2); 
                    addr_input += ((BL+2) * sizeof(int32_t));
                }

                addr_ref = mram_base_addr_ref + (t_index_x * BL * BL_IN * sizeof(int32_t)) +  (t_index_y * BL_IN * sizeof(int32_t));
                cache_input_offset = 0;
                for (int i = 0; i < BL_IN; i++) {
                    mram_read((__mram_ptr void const *) addr_ref, (void *) (cache_ref + cache_input_offset), (BL_IN) * sizeof(int32_t)); 
                    cache_input_offset += BL_IN; 
                    addr_ref += (BL * sizeof(int32_t));
                }

                // Computation
                for (uint32_t i = 1; i < BL_IN + 1; i++) {
                    for (uint32_t j = 1; j < BL_IN + 1; j++) {
                        cache_input[i*(BL_IN+2) + j] = maximum(cache_input[(i-1)*(BL_IN+2) + j - 1] + cache_ref[(i-1)*BL_IN + j-1],
                                                cache_input[i*(BL_IN+2) + j - 1] - penalty,
                                                cache_input[(i-1)*(BL_IN+2) + j] - penalty);
                    }
                }

                // Move output from WRAM to MRAM
                addr_input =  mram_base_addr_input_itemsets + (t_index_x * (BL+2) * BL_IN * sizeof(int32_t)) + (t_index_y * BL_IN * sizeof(int32_t));
                cache_input_offset = (BL_IN+2);
                addr_input += ((BL+2) * sizeof(int32_t));
                for (int i = 1; i < BL_IN + 1; i++) {
                    mram_write((cache_input + cache_input_offset), (__mram_ptr void *)  addr_input, (BL_IN+2) * sizeof(int32_t)); 
                    cache_input_offset += (BL_IN+2); 
                    addr_input += ((BL+2) * sizeof(int32_t));
                }

            }
            
            barrier_wait(&my_barrier);
        }
       
        // Bottom-right computation
        for(uint32_t blk = 2; blk <= REP; blk++) {
            // Partition chunks/subblocks of the diagonal to tasklets 
            chunks = (REP - blk + 1) / NR_TASKLETS; 
            mod = (REP - blk + 1) % NR_TASKLETS;
            if (tasklet_id < mod)
                chunks++;
            if (mod > 0){
                if(tasklet_id < mod)
                    start = tasklet_id * chunks;
                else
                    start = mod * (chunks + 1) + (tasklet_id - mod) * chunks;
            } else
                start = tasklet_id * chunks;

            // Compute all assigned chunks  
            for (uint32_t bl_indx = 0; bl_indx < chunks; bl_indx++) {
                int t_index_x = blk - 1 + start + bl_indx;
                int t_index_y = REP + blk - 2 - t_index_x; 

                // Move input from MRAM to WRAM
                addr_input =  mram_base_addr_input_itemsets + (t_index_x * (BL+2) * BL_IN * sizeof(int32_t)) + (t_index_y * BL_IN * sizeof(int32_t));
                cache_input_offset = (BL_IN+2);
                mram_read((__mram_ptr void const *) addr_input, (void *) cache_input, (BL_IN+2) * sizeof(int32_t)); 
                addr_input += ((BL+2) * sizeof(int32_t));
                for (int i = 1; i < BL_IN + 1; i++) {
                    mram_read((__mram_ptr void const *) addr_input, (void *) (cache_input + cache_input_offset), (2) * sizeof(int32_t)); 
                    cache_input_offset += (BL_IN+2); 
                    addr_input += ((BL+2) * sizeof(int32_t));
                }

                addr_ref = mram_base_addr_ref + (t_index_x * BL * BL_IN * sizeof(int32_t)) +  (t_index_y * BL_IN * sizeof(int32_t));
                cache_input_offset = 0;
                for (int i = 0; i < BL_IN; i++) {
                    mram_read((__mram_ptr void const *) addr_ref, (void *) (cache_ref + cache_input_offset), (BL_IN) * sizeof(int32_t)); 
                    cache_input_offset += BL_IN; 
                    addr_ref += (BL * sizeof(int32_t));
                }


                // Computation
                for (int i = 1; i < BL_IN + 1; i++) {
                    for (int j = 1; j < BL_IN + 1; j++) {
                        cache_input[i*(BL_IN+2) + j] = maximum(cache_input[(i-1)*(BL_IN+2) + j - 1] + cache_ref[(i-1)*BL_IN + j-1],
                                                cache_input[i*(BL_IN+2) + j - 1] - penalty,
                                                cache_input[(i-1)*(BL_IN+2) + j] - penalty);
                    }
                }

                // Move output from WRAM to MRAM
                addr_input =  mram_base_addr_input_itemsets + (t_index_x * (BL+2) * BL_IN * sizeof(int32_t)) + (t_index_y * BL_IN * sizeof(int32_t));
                cache_input_offset = (BL_IN+2);
                addr_input += ((BL+2) * sizeof(int32_t));
                for (int i = 1; i < BL_IN + 1; i++) {
                    mram_write(cache_input + cache_input_offset, (__mram_ptr void *)  addr_input, (BL_IN+2) * sizeof(int32_t)); 
                    cache_input_offset += (BL_IN+2); 
                    addr_input += ((BL+2) * sizeof(int32_t));
                }

            }
            
            barrier_wait(&my_barrier);

        }
		
        mram_base_addr_input_itemsets += ((BL+1) * (BL+2) * sizeof(int32_t));
        mram_base_addr_ref += (BL * BL * sizeof(int32_t)); 
    }
    return 0;
}
