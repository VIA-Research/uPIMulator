/**
* app.c
* NW Host Application Source File 
*
*/
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <dpu.h>
#include <dpu_log.h>
#include <unistd.h>
#include <getopt.h>
#include <assert.h>

#include "../support/common.h"
#include "../support/timer.h"
#include "../support/params.h"

#if ENERGY
#include <dpu_probe.h>
#endif

// Define the DPU Binary path as DPU_BINARY here
#ifndef DPU_BINARY
#define DPU_BINARY "./bin/nw_dpu"
#endif

// Traceback in the host
#if PRINT_FILE
static void traceback(int* traceback_output, char *file, int32_t *input_itemsets, int32_t *reference, unsigned int max_rows, unsigned int max_cols, unsigned int penalty) {
    FILE *fpo = fopen(file, "w"); // Use to print to an output file
#else
static void traceback(int* traceback_output, int32_t *input_itemsets, int32_t *reference, unsigned int max_rows, unsigned int max_cols, unsigned int penalty) {
#endif

    int k = 0;
    for (int i = max_rows - 2,  j = max_rows - 2; i>=0 && j>=0;) {
        int nw = 0, n = 0, w = 0, traceback = 0;
#if PRINT_FILE
        if ( i == (int)max_rows - 2 && j == (int)max_rows - 2 )
            fprintf(fpo, "%d ", input_itemsets[ i * max_cols + j]); //print the first element
#endif

        if (i == 0 && j == 0)
            break;
        if (i > 0 && j > 0) {
            nw = input_itemsets[(i - 1) * max_cols + j - 1];
            w  = input_itemsets[i * max_cols + j - 1];
            n  = input_itemsets[(i - 1) * max_cols + j];
        } else if (i == 0) {
            nw = n = LIMIT;
            w  = input_itemsets[ i * max_cols + j - 1 ];
        } else if (j == 0) {
            nw = w = LIMIT;
            n  = input_itemsets[(i - 1) * max_cols + j];
        } else {
            ;
        }

        int new_nw, new_w, new_n;
        new_nw = nw + reference[i * max_cols + j];
        new_w = w - penalty;
        new_n = n - penalty;

        traceback = maximum(new_nw, new_w, new_n);
        if (traceback == new_nw)
            traceback = nw;
        if (traceback == new_w)
            traceback = w;
        if (traceback == new_n)
            traceback = n;

#if PRINT_FILE
	fprintf(fpo, "%d ", traceback);
#endif 
 	traceback_output[k++] = traceback;

	if (traceback == nw) { 
            i--; 
            j--; 
            continue;
        } else if (traceback == w) {
            j--; 
            continue;
        } else if (traceback == n) {
            i--; 
            continue;
        } else {
            ;
        }
    }

    return;
}

// Compute output in the host
static void nw_host(int32_t *input_itemsets, int32_t *reference, uint64_t max_cols, unsigned int penalty) {

    int32_t *input_itemsets_l = (int32_t *) malloc((BL + 1) * (BL + 1) * sizeof(int32_t));
    int32_t *reference_l = (int32_t *) malloc((BL * BL) * sizeof(int32_t));


    // top-left
    for (uint64_t blk = 1; blk <= (max_cols-1)/BL; blk++) {
        for (uint64_t b_index_x = 0; b_index_x < blk; b_index_x++) {
            uint64_t b_index_y = blk - 1 - b_index_x;

            for (uint64_t i = 0; i < BL; i++){
                for (uint64_t j = 0; j < BL; j++) {
                    reference_l[i*BL + j] = reference[(max_cols-1) * (b_index_y*BL + i) + b_index_x*BL + j];
                }
            }

            for (uint64_t i = 0; i < BL + 1; i++){
                for (uint64_t j = 0; j < BL + 1; j++) {
                    input_itemsets_l[i*(BL + 1) + j] = input_itemsets[max_cols*(b_index_y*BL + i) + b_index_x*BL + j];
                }
            }

            // Computation
            for (uint64_t i = 1; i < BL + 1; i++) {
                for (uint64_t j = 1; j < BL + 1; j++) {
                    input_itemsets_l[i*(BL + 1) + j] = maximum(input_itemsets_l[(i-1)*(BL+1) + j - 1] + reference_l[(i-1)*BL + j - 1],
                            input_itemsets_l[i*(BL+1) + j - 1] - penalty,
                            input_itemsets_l[(i-1)*(BL+1) + j] - penalty);
                }
            }

            for (uint64_t i = 0; i < BL; i++) {
                for (uint64_t j = 0; j < BL; j++) {
                    input_itemsets[max_cols*(b_index_y*BL + i + 1) + b_index_x*BL + j + 1] = input_itemsets_l[(i+1)*(BL+1) + j + 1];
                }
            }

        }

    }

    // bottom-right 
    for (uint64_t blk = 2; blk <= (max_cols-1)/BL; blk++) {
        for (uint64_t b_index_x = blk - 1; b_index_x < (max_cols-1)/BL; b_index_x++) {
            uint64_t b_index_y = (max_cols-1)/BL + blk - 2 - b_index_x;

            for (uint64_t i = 0; i < BL; i++){
                for (uint64_t j = 0; j < BL; j++) {
                    reference_l[i*BL + j] = reference[(max_cols-1)*(b_index_y*BL + i) + b_index_x*BL + j];
                }
            }

            for (uint64_t i = 0; i < BL + 1; i++){
                for (uint64_t j = 0; j < BL + 1; j++) {
                    input_itemsets_l[i*(BL + 1) + j] = input_itemsets[max_cols*(b_index_y*BL + i) + b_index_x*BL + j];
                }
            }

            // Computation
            for (uint64_t i = 1; i < BL + 1; i++) {
                for (uint64_t j = 1; j < BL + 1; j++) {
                    input_itemsets_l[i*(BL + 1) + j] = maximum(input_itemsets_l[(i-1)*(BL+1) + j - 1] + reference_l[(i-1)*BL + j - 1],
                            input_itemsets_l[i*(BL+1) + j - 1] - penalty,
                            input_itemsets_l[(i-1)*(BL+1) + j] - penalty);
                }
            }

            for (uint64_t i = 0; i < BL; i++) {
                for (uint64_t j = 0; j < BL; j++) {
                    input_itemsets[max_cols*(b_index_y*BL + i + 1) + b_index_x*BL + j + 1] = input_itemsets_l[(i+1)*(BL+1) + j + 1];
                }
            }

        }

    }


    free(input_itemsets_l);
    free(reference_l);
    return;
}

// Main of the Host Application
int main(int argc, char **argv) {

    struct Params p = input_params(argc, argv);
    struct dpu_set_t dpu_set, dpu;
    uint32_t nr_of_dpus, max_dpus;

#if ENERGY
    struct dpu_probe_t probe;
    DPU_ASSERT(dpu_probe_init("energy probe", &probe));
#endif

    // Allocate DPUs and load binary
    DPU_ASSERT(dpu_alloc(NR_DPUS, NULL, &dpu_set));
    DPU_ASSERT(dpu_load(dpu_set, DPU_BINARY, NULL));
    DPU_ASSERT(dpu_get_nr_dpus(dpu_set, &nr_of_dpus));
    printf("Allocated %d DPU(s)\n", nr_of_dpus);
    printf("Allocated %d TASKLET(s) per DPU\n", NR_TASKLETS);
#if DYNAMIC
    max_dpus = nr_of_dpus;
#endif

    uint64_t max_rows = p.max_rows + 1;
    uint64_t max_cols = p.max_rows + 1;
    unsigned int penalty = p.penalty;
    int32_t *reference = (int32_t *) malloc(max_rows * max_cols * sizeof(int32_t));
    int32_t *input_itemsets_host = (int32_t *) malloc(max_rows * max_cols * sizeof(int32_t));
    int32_t *input_itemsets = (int32_t *) malloc((max_rows+1) * (max_cols+1) * sizeof(int32_t));
    dpu_arguments_t *input_args = (dpu_arguments_t *) malloc(nr_of_dpus * sizeof(dpu_arguments_t));
    printf("Max size %d\n", p.max_rows);

    // Traceback output
    int32_t* traceback_output = (int32_t *) malloc((max_rows + max_cols) * sizeof(int32_t));
    int32_t* traceback_output_host = (int32_t *) malloc((max_rows + max_cols) * sizeof(int32_t));
    memset(traceback_output, 0, (max_rows + max_cols) * sizeof(int32_t));
    memset(traceback_output_host, 0, (max_rows + max_cols) * sizeof(int32_t));

    // This array is used for dummy/stale CPU-DPU transfers
    int32_t *dummy = (int32_t *) malloc(nr_of_dpus * (BL+2) * sizeof(int32_t));
    unsigned int blocks_per_dpu;
    unsigned int mram_offset = 0;

    // Timer
    Timer timer; 
    Timer long_diagonal_timer; 
#if ENERGY
    double tacc_energy, tacc_time, tavg_time;
    double tavg_energy=0;
#endif

    for (unsigned int rep = 0; rep < p.n_warmup + p.n_reps; rep++) {

        // Initializing inputs are needed at each iteration
        // Initialize input itemsets
        for(unsigned int i = 0; i < max_rows; i++) {
            for (unsigned int j = 0; j < max_cols; j++) {
                input_itemsets_host[i * max_cols + j] = 0; 
            }
        }

        for(unsigned int i = 0; i <= max_rows; i++) {
            for (unsigned int j = 0; j <= max_cols; j++) {
                input_itemsets[i * (max_cols+1) + j] = 0; 
            }
        }

        // Define random sequences
        srand(7);
        for (unsigned int i = 1; i < max_rows; i++) {
            input_itemsets_host[i * max_cols] = rand() % 10 + 1;
        }

        for (unsigned int j = 1; j < max_cols; j++) {
            input_itemsets_host[j] = rand() % 10 + 1;
        }   

        for (unsigned int i = 0; i < max_rows-1; i++) {
            for (unsigned int j = 0; j < max_cols-1; j++) {
                reference[i * (max_cols-1) + j] = blosum62[input_itemsets[(i+1) * max_cols]][input_itemsets[j+1]];
            }
        }

        for (unsigned int i = 1; i < max_rows; i++) {
            input_itemsets_host[i * max_cols] = -i * penalty;
            input_itemsets[i * (max_cols+1)] = -i * penalty;
        }

        for (unsigned int j = 1; j < max_cols; j++) {
            input_itemsets_host[j] = -j * penalty;
            input_itemsets[j] = -j * penalty;
        }

        if (rep >= p.n_warmup)
            start(&timer, 0, rep - p.n_warmup);
        // Computation on host CPU
        nw_host(input_itemsets_host, reference, max_cols, penalty);

        // Print host output
#if PRINT_FILE
        if (rep >= p.n_warmup) {
            char *host_file = "./bin/host_output.txt";
            traceback(traceback_output_host, host_file, input_itemsets_host, reference, max_rows, max_cols, penalty);
        }
#endif
        if (rep >= p.n_warmup)
            stop(&timer, 0);

        // Top-left computation on DPUs
        for (unsigned int blk = 1; blk <= (max_cols-1)/BL; blk++) {
#if DYNAMIC 
            // If nr_of_blocks are lower than max_dpus,
            // set nr_of_dpus to be equal with nr_of_blocks
            unsigned nr_of_blocks = blk;
            if (nr_of_blocks < max_dpus) {
                DPU_ASSERT(dpu_free(dpu_set));
                DPU_ASSERT(dpu_alloc(nr_of_blocks, NULL, &dpu_set));
                DPU_ASSERT(dpu_load(dpu_set, DPU_BINARY, NULL));
                DPU_ASSERT(dpu_get_nr_dpus(dpu_set, &nr_of_dpus));
            } else if (nr_of_dpus == max_dpus) {
                ;
            } else {
                DPU_ASSERT(dpu_free(dpu_set));
                DPU_ASSERT(dpu_alloc(max_dpus, NULL, &dpu_set));
                DPU_ASSERT(dpu_load(dpu_set, DPU_BINARY, NULL));
                DPU_ASSERT(dpu_get_nr_dpus(dpu_set, &nr_of_dpus));
            }
#if PRINT
            printf("Allocated %d DPU(s) for %d (%d) blocks\n", nr_of_dpus, nr_of_blocks, blk);
#endif
#endif

            // Copy data to DPUs
            unsigned int i=0;
            DPU_FOREACH(dpu_set, dpu, i) {
                unsigned int blocks_per_dpu = blk / nr_of_dpus;
                unsigned int active_blocks_per_dpu = blk / nr_of_dpus;
                unsigned int rest_blocks = blk % nr_of_dpus;
                if(i < rest_blocks)
                    blocks_per_dpu++;

                if(rest_blocks != 0)
                    active_blocks_per_dpu++;

                // Copy input arguments to dpu
                input_args[i].nblocks = blocks_per_dpu;
                input_args[i].active_blocks = active_blocks_per_dpu;
                input_args[i].penalty = penalty;
                DPU_ASSERT(dpu_prepare_xfer(dpu, input_args + i));
            } 
            DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, "DPU_INPUT_ARGUMENTS", 0, sizeof(dpu_arguments_t), DPU_XFER_DEFAULT));

            // Copy itemsets to DPUs
            blocks_per_dpu = blk / nr_of_dpus;
            if (blk % nr_of_dpus != 0)
                blocks_per_dpu++;
            mram_offset = 0;


            if (rep >= p.n_warmup) {
                if ((max_cols-1)/BL == 1) 
                    start(&timer, 2, rep - p.n_warmup + blk - 1);
                else 
                    start(&timer, 1, rep - p.n_warmup + blk - 1);
                
                // Timer for longest diagonal
                if (blk == ((max_cols-1)/BL)) {
                    if ((max_cols-1)/BL == 1) 
                        start(&long_diagonal_timer, 2, rep - p.n_warmup);
                    else 
                        start(&long_diagonal_timer, 1, rep - p.n_warmup);
                }
            }

#if PRINT
            uint64_t total_dpu_memory = 0;
            total_dpu_memory = (uint64_t) blocks_per_dpu * (BL+1) * (BL+2) * sizeof(int32_t) + (uint64_t) blocks_per_dpu * BL * BL * sizeof(int32_t);
            printf("Total memory allocated in each DPU %u bytes\n", total_dpu_memory);
#endif
            for (unsigned int bl_indx = 0; bl_indx < blocks_per_dpu; bl_indx++) {
                for (unsigned int bl = 0; bl < BL + 1; bl++) {

                    i = 0;
                    DPU_FOREACH(dpu_set, dpu, i) {
                        unsigned int chunks = blk / nr_of_dpus;
                        unsigned int prev_block_index = 0;
                        unsigned int rest_blocks = blk % nr_of_dpus;
                        if (rest_blocks > 0) {
                            if (i >= rest_blocks) {
                                prev_block_index = rest_blocks * (chunks + 1) + (i - rest_blocks) * chunks;
                            } else {
                                prev_block_index = i * (chunks + 1);
                            }
                        } else {
                            prev_block_index = i * blocks_per_dpu; 
                        }

                        uint64_t input_itemsets_offset = 0;  
                        int32_t *dpu_pointer;  
                        if (i + bl_indx * nr_of_dpus >= blk) {
                            dpu_pointer = dummy;
                            input_itemsets_offset = 0;  
                        } else {
                            uint64_t b_index_x =  prev_block_index + bl_indx;
                            uint64_t b_index_y = blk - 1 - b_index_x;
                            dpu_pointer = input_itemsets;
                            input_itemsets_offset = b_index_y * (max_cols+1) * BL + b_index_x * BL + bl * (max_cols + 1);  
                        }

                        DPU_ASSERT(dpu_prepare_xfer(dpu, dpu_pointer + input_itemsets_offset));
                    }

                    if (bl == 0) 
                        DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, mram_offset, (BL+2) * sizeof(int32_t), DPU_XFER_DEFAULT));
                    else
                        DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, mram_offset, 2 * sizeof(int32_t), DPU_XFER_DEFAULT));
                    mram_offset += ((BL+2) * sizeof(int32_t));

                }
            }
            if (rep >= p.n_warmup) {
                if ((max_cols-1)/BL == 1) 
                    stop(&timer, 2);
                else
                    stop(&timer, 1);
                // Timer for longest diagonal
                if (blk == ((max_cols-1)/BL)) {
                    if ((max_cols-1)/BL == 1) 
                        stop(&long_diagonal_timer, 2);
                    else 
                        stop(&long_diagonal_timer, 1);
                }
            }


            if (rep >= p.n_warmup) {
                start(&timer, 2, rep - p.n_warmup + blk - 1);
                // Timer for longest diagonal
                if (blk == ((max_cols-1)/BL)) {
                    start(&long_diagonal_timer, 2, rep - p.n_warmup);
                }
            }
            // Copy reference to DPUs
            mram_offset = blocks_per_dpu * (BL+1) * (BL+2) * sizeof(int32_t); 
            for (unsigned int bl_indx = 0; bl_indx < blocks_per_dpu; bl_indx++) {
                for (unsigned int bl = 0; bl < BL; bl++) {

                    i = 0;
                    DPU_FOREACH(dpu_set, dpu, i) {
                        unsigned int chunks = blk / nr_of_dpus;
                        unsigned int prev_block_index = 0;
                        unsigned int rest_blocks = blk % nr_of_dpus;
                        if (rest_blocks > 0) {
                            if (i >= rest_blocks) {
                                prev_block_index = rest_blocks * (chunks + 1) + (i - rest_blocks) * chunks;
                            } else {
                                prev_block_index = i * (chunks + 1);
                            }
                        } else {
                            prev_block_index = i * blocks_per_dpu; 
                        }

                        uint64_t reference_offset = 0;  
                        int32_t *dpu_pointer;  
                        if (i + bl_indx * nr_of_dpus >= blk) {
                            dpu_pointer = dummy;
                            reference_offset = 0;  
                        } else {
                            uint64_t b_index_x =  prev_block_index + bl_indx;
                            uint64_t b_index_y = blk - 1 - b_index_x;
                            dpu_pointer = reference;
                            reference_offset = b_index_y * (max_cols - 1) * BL + b_index_x * BL + bl * (max_cols - 1);  
                        }

                        DPU_ASSERT(dpu_prepare_xfer(dpu, dpu_pointer + reference_offset));
                    }
                    DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, mram_offset, BL * sizeof(int32_t), DPU_XFER_DEFAULT));
                    mram_offset += BL * sizeof(int32_t);

                }
            }
            if (rep >= p.n_warmup) {
                stop(&timer, 2);
                if (blk == ((max_cols-1)/BL)) {
                    stop(&long_diagonal_timer, 2);
                }
            }

#if ENERGY
            if (rep >= p.n_warmup) {
                DPU_ASSERT(dpu_probe_start(&probe));
            }
#endif
            if (rep >= p.n_warmup) {
                start(&timer, 3, rep - p.n_warmup + blk - 1);
                // Timer for longest diagonal
                if (blk == ((max_cols-1)/BL)) {
                    start(&long_diagonal_timer, 3, rep - p.n_warmup);
                }
            }
            // Launch kernel on DPUs
            DPU_ASSERT(dpu_launch(dpu_set, DPU_SYNCHRONOUS));
            if (rep >= p.n_warmup) {
                stop(&timer, 3);
                // Timer for longest diagonal
                if (blk == ((max_cols-1)/BL)) {
                    stop(&long_diagonal_timer, 3);
                }
            }
#if ENERGY
            if (rep >= p.n_warmup) {
                DPU_ASSERT(dpu_probe_stop(&probe));
            }
#endif

#if ENERGY
    	    double acc_energy, avg_energy, acc_time, avg_time;
	        DPU_ASSERT(dpu_probe_get(&probe, DPU_ENERGY, DPU_ACCUMULATE, &acc_energy));
    	    DPU_ASSERT(dpu_probe_get(&probe, DPU_ENERGY, DPU_AVERAGE, &avg_energy));
    	    DPU_ASSERT(dpu_probe_get(&probe, DPU_TIME, DPU_ACCUMULATE, &acc_time));
    	    DPU_ASSERT(dpu_probe_get(&probe, DPU_TIME, DPU_AVERAGE, &avg_time));
    	    tavg_energy += avg_energy;
#endif

#if PRINT
            // Display DPU Logs
            DPU_FOREACH(dpu_set, dpu) {
                DPU_ASSERT(dpulog_read_for_dpu(dpu.dpu, stdout));
            }
#endif

            if (rep >= p.n_warmup) {
                start(&timer, 4, rep - p.n_warmup + blk - 1);
                // Timer for longest diagonal
                if (blk == ((max_cols-1)/BL)) {
                    start(&long_diagonal_timer, 4, rep - p.n_warmup);
                }
            }
            // Retrieve results
            // Copy output result to Host CPU
            mram_offset = 0;
            for (unsigned int bl_indx = 0; bl_indx < blocks_per_dpu; bl_indx++) {
                for (unsigned int bl = 0; bl < BL + 1; bl++) {

                    i = 0;
                    DPU_FOREACH(dpu_set, dpu, i) {
                        unsigned int chunks = blk / nr_of_dpus;
                        unsigned int prev_block_index = 0;
                        unsigned int rest_blocks = blk % nr_of_dpus;
                        if (rest_blocks > 0) {
                            if (i >= rest_blocks) {
                                prev_block_index = rest_blocks * (chunks + 1) + (i - rest_blocks) * chunks;
                            } else {
                                prev_block_index = i * (chunks + 1);
                            }
                        } else {
                            prev_block_index = i * blocks_per_dpu; 
                        }

                        uint64_t input_itemsets_offset = 0;  
                        int32_t *dpu_pointer;  
                        if (i + bl_indx * nr_of_dpus >= blk) {
                            dpu_pointer = dummy;
                            input_itemsets_offset = 0;  
                        } else {
                            uint64_t b_index_x =  prev_block_index + bl_indx;
                            uint64_t b_index_y = blk - 1 - b_index_x;
                            dpu_pointer = input_itemsets;
                            input_itemsets_offset = b_index_y * (max_cols+1) * BL + b_index_x * BL + bl * (max_cols + 1);  
                        }

                        if (bl == 0) // Skip the first row of the block
                            continue;
                        DPU_ASSERT(dpu_prepare_xfer(dpu, dpu_pointer + input_itemsets_offset));

                    }
                    if (bl == 0) {
                        mram_offset += (BL+2) * sizeof(int32_t);
                        continue;
                    }
                    DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_FROM_DPU, DPU_MRAM_HEAP_POINTER_NAME, mram_offset, (BL+2) * sizeof(int32_t), DPU_XFER_DEFAULT));
                    mram_offset += (BL+2) * sizeof(int32_t);

                }
            }
            if (rep >= p.n_warmup) {
                stop(&timer, 4);
                // Timer for longest diagonal
                if (blk == ((max_cols-1)/BL)) {
                    stop(&long_diagonal_timer, 4);
                }
            }
        }


        // Bottom-right computation on DPUs
        for (unsigned int blk = 2; blk <= (max_cols-1)/BL; blk++) {
#if DYNAMIC
            // If nr_of_blocks are lower than max_dpus,
            // set nr_of_dpus to be equal with nr_of_blocks
            unsigned nr_of_blocks = (((max_cols-1)/BL) - blk + 1);
            if (nr_of_blocks < max_dpus) {
                DPU_ASSERT(dpu_free(dpu_set));
                DPU_ASSERT(dpu_alloc(nr_of_blocks, NULL, &dpu_set));
                DPU_ASSERT(dpu_load(dpu_set, DPU_BINARY, NULL));
                DPU_ASSERT(dpu_get_nr_dpus(dpu_set, &nr_of_dpus));
            } else if (nr_of_dpus == max_dpus) {
                ;
            } else {
                DPU_ASSERT(dpu_free(dpu_set));
                DPU_ASSERT(dpu_alloc(max_dpus, NULL, &dpu_set));
                DPU_ASSERT(dpu_load(dpu_set, DPU_BINARY, NULL));
                DPU_ASSERT(dpu_get_nr_dpus(dpu_set, &nr_of_dpus));
            }
#if PRINT
            printf("Allocated %d DPU(s) for %d (%d) blocks\n", nr_of_dpus, nr_of_blocks, (((max_cols-1)/BL) - blk + 1));
#endif
#endif

            // Copy data to DPUs
            unsigned int i=0;
            DPU_FOREACH(dpu_set, dpu, i) {
                unsigned int blocks_per_dpu = (((max_cols-1)/BL) - blk + 1) / nr_of_dpus;
                unsigned int active_blocks_per_dpu = (((max_cols-1)/BL) - blk + 1) / nr_of_dpus;
                unsigned int rest_blocks = (((max_cols-1)/BL) - blk + 1) % nr_of_dpus;
                if(i < rest_blocks)
                    blocks_per_dpu++;

                if(rest_blocks != 0)
                    active_blocks_per_dpu++;

                // Copy input arguments to dpu
                input_args[i].nblocks = blocks_per_dpu;
                input_args[i].active_blocks = active_blocks_per_dpu;
                input_args[i].penalty = penalty;
                DPU_ASSERT(dpu_prepare_xfer(dpu, input_args + i));
            } 
            DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, "DPU_INPUT_ARGUMENTS", 0, sizeof(dpu_arguments_t), DPU_XFER_DEFAULT));

            if (rep >= p.n_warmup)
                start(&timer, 1, rep - p.n_warmup + blk - 1);
            // Copy itemsets to DPUs
            unsigned int blocks_per_dpu = (((max_cols-1)/BL) - blk + 1) / nr_of_dpus;
            if ((((max_cols-1)/BL) - blk + 1) % nr_of_dpus != 0)
                blocks_per_dpu++;
#if PRINT
            uint64_t total_dpu_memory = 0;
            total_dpu_memory = (uint64_t) blocks_per_dpu * (BL+1) * (BL+2) * sizeof(int32_t) + (uint64_t) blocks_per_dpu * BL * BL * sizeof(int32_t);
            printf("Total memory allocated in each DPU %u bytes\n", total_dpu_memory);
#endif
            unsigned int mram_offset = 0;
            for (unsigned int bl_indx = 0; bl_indx < blocks_per_dpu; bl_indx++) {
                for (unsigned int bl = 0; bl < BL + 1; bl++) {

                    i = 0;
                    DPU_FOREACH(dpu_set, dpu, i) {
                        unsigned int chunks = (((max_cols-1)/BL) - blk + 1) / nr_of_dpus;
                        unsigned int prev_block_index = 0;
                        unsigned int rest_blocks = (((max_cols-1)/BL) - blk + 1) % nr_of_dpus;
                        if (rest_blocks > 0) {
                            if (i >= rest_blocks) {
                                prev_block_index = rest_blocks * (chunks + 1) + (i - rest_blocks) * chunks;
                            } else {
                                prev_block_index = i * (chunks + 1);
                            }
                        } else {
                            prev_block_index = i * blocks_per_dpu; 
                        }

                        uint64_t input_itemsets_offset = 0;  
                        int32_t *dpu_pointer;  
                        if (i + bl_indx * nr_of_dpus >= (((max_cols-1)/BL) - blk + 1)) {
                            dpu_pointer = dummy;
                            input_itemsets_offset = 0;  
                        } else {
                            uint64_t b_index_x = blk - 1 + prev_block_index + bl_indx;
                            uint64_t b_index_y = (max_cols-1)/BL + blk - 2 - b_index_x;
                            dpu_pointer = input_itemsets;
                            input_itemsets_offset = b_index_y * (max_cols+1) * BL + b_index_x * BL + bl * (max_cols + 1);  
                        }

                        DPU_ASSERT(dpu_prepare_xfer(dpu, dpu_pointer + input_itemsets_offset));
                    }

                    if (bl == 0) 
                        DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, mram_offset, (BL+2) * sizeof(int32_t), DPU_XFER_DEFAULT));
                    else
                        DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, mram_offset, 2 * sizeof(int32_t), DPU_XFER_DEFAULT));
                    mram_offset += (BL+2) * sizeof(int32_t);

                }
            }
            if (rep >= p.n_warmup)
                stop(&timer, 1);


            if (rep >= p.n_warmup)
                start(&timer, 2, rep - p.n_warmup + blk - 1);
            // Copy reference to DPUs
            mram_offset = blocks_per_dpu * (BL+1) * (BL+2) * sizeof(int32_t); 
            for (unsigned int bl_indx = 0; bl_indx < blocks_per_dpu; bl_indx++) {
                for (unsigned int bl = 0; bl < BL; bl++) {

                    i = 0;
                    DPU_FOREACH(dpu_set, dpu, i) {
                        unsigned int chunks = (((max_cols-1)/BL) - blk + 1) / nr_of_dpus;
                        unsigned int prev_block_index = 0;
                        unsigned int rest_blocks = (((max_cols-1)/BL) - blk + 1) % nr_of_dpus;
                        if (rest_blocks > 0) {
                            if (i >= rest_blocks) {
                                prev_block_index = rest_blocks * (chunks + 1) + (i - rest_blocks) * chunks;
                            } else {
                                prev_block_index = i * (chunks + 1);
                            }
                        } else {
                            prev_block_index = i * blocks_per_dpu; 
                        }

                        uint64_t reference_offset = 0;  
                        int32_t *dpu_pointer;  
                        if (i + bl_indx * nr_of_dpus >= (((max_cols-1)/BL) - blk + 1)) {
                            dpu_pointer = dummy;
                            reference_offset = 0;  
                        } else {
                            uint64_t b_index_x = blk - 1 + prev_block_index + bl_indx;
                            uint64_t b_index_y = (max_cols-1)/BL + blk - 2 - b_index_x;
                            dpu_pointer = reference;
                            reference_offset = b_index_y * (max_cols - 1) * BL + b_index_x * BL + bl * (max_cols - 1);  
                        }

                        DPU_ASSERT(dpu_prepare_xfer(dpu, dpu_pointer + reference_offset));
                    }

                    DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_TO_DPU, DPU_MRAM_HEAP_POINTER_NAME, mram_offset, BL * sizeof(int32_t), DPU_XFER_DEFAULT));
                    mram_offset += BL * sizeof(int32_t);

                }
            }
            if (rep >= p.n_warmup)
                stop(&timer, 2);

#if ENERGY
            if (rep >= p.n_warmup) {
                DPU_ASSERT(dpu_probe_start(&probe));
            }
#endif
            if (rep >= p.n_warmup)
                start(&timer, 3, rep - p.n_warmup + blk - 1); // Do not re-initialize the counter
            // Launch kernel on DPUs
            DPU_ASSERT(dpu_launch(dpu_set, DPU_SYNCHRONOUS));
            if (rep >= p.n_warmup)
                stop(&timer, 3);
#if ENERGY
            if (rep >= p.n_warmup) {
                DPU_ASSERT(dpu_probe_stop(&probe));
            }
#endif

#if ENERGY
    	    double acc_energy, avg_energy, acc_time, avg_time;
    	    DPU_ASSERT(dpu_probe_get(&probe, DPU_ENERGY, DPU_ACCUMULATE, &acc_energy));
    	    DPU_ASSERT(dpu_probe_get(&probe, DPU_ENERGY, DPU_AVERAGE, &avg_energy));
    	    DPU_ASSERT(dpu_probe_get(&probe, DPU_TIME, DPU_ACCUMULATE, &acc_time));
    	    DPU_ASSERT(dpu_probe_get(&probe, DPU_TIME, DPU_AVERAGE, &avg_time));
    	    tavg_energy += avg_energy;
#endif

#if PRINT
            // Display DPU Logs
            DPU_FOREACH(dpu_set, dpu) {
                DPU_ASSERT(dpulog_read_for_dpu(dpu.dpu, stdout));
            }
#endif


            if (rep >= p.n_warmup)
                start(&timer, 4, rep - p.n_warmup + blk - 1);
            // Retrieve results
            // Copy output result to Host CPU
            mram_offset = 0;
            for (unsigned int bl_indx = 0; bl_indx < blocks_per_dpu; bl_indx++) {
                for (unsigned int bl = 0; bl < BL + 1; bl++) {

                    i = 0;
                    DPU_FOREACH(dpu_set, dpu, i) {
                        unsigned int chunks = (((max_cols-1)/BL) - blk + 1) / nr_of_dpus;
                        unsigned int prev_block_index = 0;
                        unsigned int rest_blocks = (((max_cols-1)/BL) - blk + 1) % nr_of_dpus;
                        if (rest_blocks > 0) {
                            if (i >= rest_blocks) {
                                prev_block_index = rest_blocks * (chunks + 1) + (i - rest_blocks) * chunks;
                            } else {
                                prev_block_index = i * (chunks + 1);
                            }
                        } else {
                            prev_block_index = i * blocks_per_dpu; 
                        }

                        uint64_t input_itemsets_offset = 0;  
                        int32_t *dpu_pointer;  
                        if (i + bl_indx * nr_of_dpus >= (((max_cols-1)/BL) - blk + 1)) {
                            dpu_pointer = dummy;
                            input_itemsets_offset = 0;  
                        } else {
                            uint64_t b_index_x = blk - 1 + prev_block_index + bl_indx;
                            uint64_t b_index_y = (max_cols-1)/BL + blk - 2 - b_index_x;
                            dpu_pointer = input_itemsets;
                            input_itemsets_offset = b_index_y * (max_cols+1) * BL + b_index_x * BL + bl * (max_cols + 1);  
                        }

                        if (bl == 0) // Skip the first row of the block
                            continue;
                        DPU_ASSERT(dpu_prepare_xfer(dpu, dpu_pointer + input_itemsets_offset));

                    }

                    if (bl == 0) {
                        mram_offset += (BL+2) * sizeof(int32_t);
                        continue;
                    }
                    DPU_ASSERT(dpu_push_xfer(dpu_set, DPU_XFER_FROM_DPU, DPU_MRAM_HEAP_POINTER_NAME, mram_offset, (BL+2) * sizeof(int32_t), DPU_XFER_DEFAULT));
                    mram_offset += (BL+2) * sizeof(int32_t);

                }
            }
            if (rep >= p.n_warmup)
                stop(&timer, 4);


        }

        // Traceback step
        if (rep >= p.n_warmup)
            start(&timer, 1, 1);
#if PRINT_FILE
        char *dpu_file = "./bin/dpu_output.txt";
        traceback(traceback_output, dpu_file, input_itemsets, reference, max_rows+1, max_cols+1, penalty);
#else
        traceback(traceback_output, input_itemsets, reference, max_rows+1, max_cols+1, penalty);
#endif
        if (rep >= p.n_warmup)
            stop(&timer, 1);

    }

    // Print timing results
    printf("CPU version ");
    print(&timer, 0, p.n_reps);
    printf("CPU-DPU ");
    print(&timer, 2, p.n_reps);
    printf("DPU Kernel ");
    print(&timer, 3, p.n_reps);
    printf("Inter-DPU ");
    print(&timer, 1, p.n_reps);
    printf("DPU-CPU ");
    print(&timer, 4, p.n_reps);
    printf("\n");
    printf("Longest Diagonal CPU-DPU ");
    print(&long_diagonal_timer, 2, p.n_reps);
    printf("Longest Diagonal DPU Kernel ");
    print(&long_diagonal_timer, 3, p.n_reps);
    printf("Longest Diagonal Inter-DPU ");
    print(&long_diagonal_timer, 1, p.n_reps);
    printf("Longest Diagonal DPU-CPU ");
    print(&long_diagonal_timer, 4, p.n_reps);
    printf("\n");
    
#if ENERGY
    printf("DPU Energy (J): %f \t ", tavg_energy / p.n_reps);
#endif

    // Check output
    bool status = true;
    for (uint64_t i = 1; i < max_rows; i++) {
        for (uint64_t j = 1; j < max_cols; j++) {
            if (input_itemsets_host[i*max_cols + j] != input_itemsets[i*(max_cols+1) + j]) {
                status = false;
#if PRINT
                printf("%ld (%ld, %ld): %d %d\n", i*max_cols + j, i, j, input_itemsets_host[i*max_cols + j], input_itemsets[i*(max_cols+1) + j]); 
#endif
            } 
        }
    }
    
    if (status) {
        printf("[" ANSI_COLOR_GREEN "OK" ANSI_COLOR_RESET "] Outputs are equal\n");
    } else {
        printf("[" ANSI_COLOR_RED "ERROR" ANSI_COLOR_RESET "] Outputs differ!\n");
    }

    free(input_itemsets_host);
    free(input_itemsets);
    free(reference);
    free(traceback_output);
    free(traceback_output_host);
    DPU_ASSERT(dpu_free(dpu_set));
    return status ? 0 : -1;
    return 0;
}
