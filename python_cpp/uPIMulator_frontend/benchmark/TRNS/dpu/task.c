/*
* 3-step matrix transposition with multiple tasklets
* Acks: Stefano Ballarin (P&S PIM Fall 2020)
*
*/
#include <stdint.h>
#include <stdio.h>
#include <defs.h>
#include <mram.h>
#include <alloc.h>
#include <perfcounter.h>
#include <mutex.h>
#include <barrier.h>

#include "../support/common.h"

__host dpu_arguments_t DPU_INPUT_ARGUMENTS;

uint32_t curr_tile = 0; // protected by MUTEX
uint32_t get_tile();
void read_tile_step2(uint32_t A, uint32_t offset, T* variable, uint32_t m, uint32_t n);
void write_tile_step2(uint32_t A, uint32_t offset, T* variable, uint32_t m, uint32_t n);
void read_tile_step3(uint32_t A, uint32_t offset, T* variable, uint32_t m);
void write_tile_step3(uint32_t A, uint32_t offset, T* variable, uint32_t m);
_Bool get_done(uint32_t done_array_step3, uint32_t address, T* read_done);
_Bool get_and_set_done(uint32_t done_array_step3, uint32_t address, T* read_done);

// Barrier
BARRIER_INIT(my_barrier, NR_TASKLETS);

// Mutexes
MUTEX_INIT(tile_mutex);
MUTEX_INIT(done_mutex);

extern int main_kernel1(void);
extern int main_kernel2(void);

int (*kernels[nr_kernels])(void) = {main_kernel1, main_kernel2};

int main(void) { 
    // Kernel
    return kernels[DPU_INPUT_ARGUMENTS.kernel](); 
}

// Step 2: 0010
int main_kernel1() {
    unsigned int tasklet_id = me();
#if PRINT
    printf("tasklet_id = %u\n", tasklet_id);
#endif
    if (tasklet_id == 0){ // Initialize once the cycle counter
        mem_reset(); // Reset the heap
    }
    // Barrier
    barrier_wait(&my_barrier);

    uint32_t A = (uint32_t)DPU_MRAM_HEAP_POINTER; // A in MRAM
    uint32_t M_ = DPU_INPUT_ARGUMENTS.M_;
    uint32_t m = DPU_INPUT_ARGUMENTS.m;
    uint32_t n = DPU_INPUT_ARGUMENTS.n;

    T* data = (T*) mem_alloc(m * n * sizeof(T));
    T* backup = (T*) mem_alloc(m * n * sizeof(T));

    for(unsigned int tile = tasklet_id; tile < M_; tile += NR_TASKLETS){
        read_tile_step2(A, tile * m * n, data, m, n);
        for (unsigned int i = 0; i < m * n; i++){
            backup[(i * m) - (m * n - 1) * (i / n)] = data[i];
        }
        write_tile_step2(A, tile * m * n, backup, m, n);
    }

    return 0;
}

// Step 3: 0100
int main_kernel2() {
    unsigned int tasklet_id = me();
#if PRINT
    printf("tasklet_id = %u\n", tasklet_id);
#endif
    if (tasklet_id == 0){ // Initialize once the cycle counter
        mem_reset(); // Reset the heap
    }
    // Barrier
    barrier_wait(&my_barrier);

    uint32_t A = (uint32_t)DPU_MRAM_HEAP_POINTER;
    uint32_t m = DPU_INPUT_ARGUMENTS.m;
    uint32_t n = DPU_INPUT_ARGUMENTS.n;
    uint32_t M_ = DPU_INPUT_ARGUMENTS.M_;
    uint32_t done_array = (uint32_t)(DPU_MRAM_HEAP_POINTER + M_ * m * n * sizeof(T));

    const uint32_t tile_max = M_ * n - 1; // Tile id upper bound

    T* data = (T*)mem_alloc(sizeof(T) * m);
    T* backup = (T*)mem_alloc(sizeof(T) * m);
    T* read_done = (T*)mem_alloc(sizeof(T));

    uint32_t tile;
    _Bool done;

    tile = get_tile();

    while (tile < tile_max){
        uint32_t next_in_cycle = ((tile * M_) - tile_max * (tile / n));
        if (next_in_cycle == tile){
            tile = get_tile();
            continue;
        }
        read_tile_step3(A, tile * m, data, m);

        done = get_done(done_array, tile, read_done);
        for(; done == 0; next_in_cycle = ((next_in_cycle * M_) - tile_max * (next_in_cycle / n))){
            read_tile_step3(A, next_in_cycle * m, backup, m);

            done = get_and_set_done(done_array, next_in_cycle, read_done);

            if(!done) {
                write_tile_step3(A, next_in_cycle * m, data, m);
            }
            for(uint32_t i = 0; i < m; i++){
                data[i] = backup[i];
            }
        }
        tile = get_tile();
    }
		
    return 0;
}

// Auxiliary functions
uint32_t __attribute__ ((noinline)) get_tile(){
    mutex_lock(tile_mutex);
    uint32_t value = curr_tile;
    curr_tile++;
    mutex_unlock(tile_mutex);
    return value;
}

void __attribute__ ((noinline)) read_tile_step2(uint32_t A, uint32_t offset, T* variable, uint32_t m, uint32_t n){
    int rest = m * n;
    int transfer;
    while(rest > 0){
        if(rest * sizeof(T) > 2048){
            transfer = 2048 / sizeof(T);
      } else {
            transfer = rest;
      }
      mram_read((__mram_ptr void*)(A + (offset + m * n - rest) * sizeof(T)), variable + (m * n - rest) * sizeof(T), sizeof(T) * transfer);
      rest -= transfer;
    }
}

void __attribute__ ((noinline)) write_tile_step2(uint32_t A, uint32_t offset, T* variable, uint32_t m, uint32_t n){
    int rest = m * n;
    int transfer;
    while(rest > 0){
        if(rest * sizeof(T) > 2048){
            transfer = 2048 / sizeof(T);
      } else {
            transfer = rest;
      }
      mram_write(variable + (m * n - rest) * sizeof(T), (__mram_ptr void*)(A + (offset + m * n - rest) * sizeof(T)), sizeof(T) * transfer);
      rest -= transfer;
    }
}

void __attribute__ ((noinline)) read_tile_step3(uint32_t A, uint32_t offset, T* variable, uint32_t m){
    mram_read((__mram_ptr void*)(A + offset * sizeof(T)), variable, sizeof(T) * m);
}

void __attribute__ ((noinline)) write_tile_step3(uint32_t A, uint32_t offset, T* variable, uint32_t m){
    mram_write(variable, (__mram_ptr void*)(A + offset * sizeof(T)), sizeof(T) * m);
}

_Bool __attribute__ ((noinline)) get_done(uint32_t done_array_step3, uint32_t address, T* read_done){
    uint32_t result;

    mutex_lock(done_mutex);
    mram_read((__mram_ptr void*)(done_array_step3 + address), read_done, sizeof(T));
    result = ((*read_done & (0x01 << (address % sizeof(T)))) != 0);
    mutex_unlock(done_mutex);

    return (_Bool)result;
}

_Bool __attribute__ ((noinline))get_and_set_done(uint32_t done_array_step3, uint32_t address, T* read_done){
    uint32_t result;

    mutex_lock(done_mutex);
    mram_read((__mram_ptr void*)(done_array_step3 + address), read_done, sizeof(T));
    result = ((*read_done & (0x01 << (address % sizeof(T)))) != 0);
    *read_done |= (0x01 << (address % sizeof(T)));
    mram_write(read_done, (__mram_ptr void*)(done_array_step3 + address), sizeof(T));
    mutex_unlock(done_mutex);

    return (_Bool)result;
}
