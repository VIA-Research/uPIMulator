/*
* BFS with multiple tasklets
*
*/
#include <stdio.h>

#include <alloc.h>
#include <barrier.h>
#include <defs.h>
#include <mram.h>
#include <mutex.h>
#include <perfcounter.h>

#include "dpu-utils.h"
#include "../support/common.h"

BARRIER_INIT(my_barrier, NR_TASKLETS);

BARRIER_INIT(bfsBarrier, NR_TASKLETS);
MUTEX_INIT(nextFrontierMutex);

// main
int main() {

    if(me() == 0) {
        mem_reset(); // Reset the heap
    }
    // Barrier
    barrier_wait(&my_barrier);

    // Load parameters
    uint32_t params_m = (uint32_t) DPU_MRAM_HEAP_POINTER;
    struct DPUParams* params_w = (struct DPUParams*) mem_alloc(ROUND_UP_TO_MULTIPLE_OF_8(sizeof(struct DPUParams)));
    mram_read((__mram_ptr void const*)params_m, params_w, ROUND_UP_TO_MULTIPLE_OF_8(sizeof(struct DPUParams)));

    // Extract parameters
    uint32_t numGlobalNodes = params_w->numNodes;
    uint32_t startNodeIdx = params_w->dpuStartNodeIdx;
    uint32_t numNodes = params_w->dpuNumNodes;
    uint32_t nodePtrsOffset = params_w->dpuNodePtrsOffset;
    uint32_t level = params_w->level;
    uint32_t nodePtrs_m = params_w->dpuNodePtrs_m;
    uint32_t neighborIdxs_m = params_w->dpuNeighborIdxs_m;
    uint32_t nodeLevel_m = params_w->dpuNodeLevel_m;
    uint32_t visited_m = params_w->dpuVisited_m;
    uint32_t currentFrontier_m = params_w->dpuCurrentFrontier_m;
    uint32_t nextFrontier_m = params_w->dpuNextFrontier_m;

    if(numNodes > 0) {

        // Sanity check
        if(me() == 0) {
            if(numGlobalNodes%64 != 0) {
                //PRINT_ERROR("The number of nodes in the graph is not a multiple of 64!");
            }
            if(startNodeIdx%64 != 0 || numNodes%64 != 0) {
                //PRINT_ERROR("The number of nodes assigned to the DPU is not aligned to or a multiple of 64!");
            }
        }

        // Allocate WRAM cache for each tasklet to use throughout
        uint64_t* cache_w = mem_alloc(sizeof(uint64_t));

        // Update current frontier and visited list based on the next frontier from the previous iteration
        for(uint32_t nodeTileIdx = me(); nodeTileIdx < numGlobalNodes/64; nodeTileIdx += NR_TASKLETS) {

            // Get the next frontier tile from MRAM
            uint64_t nextFrontierTile = load8B(nextFrontier_m, nodeTileIdx, cache_w);

            // Process next frontier tile if it is not empty 
            if(nextFrontierTile) {

                // Mark everything that was previously added to the next frontier as visited
                uint64_t visitedTile = load8B(visited_m, nodeTileIdx, cache_w);
                visitedTile |= nextFrontierTile;
                store8B(visitedTile, visited_m, nodeTileIdx, cache_w);

                // Clear the next frontier
                store8B(0, nextFrontier_m, nodeTileIdx, cache_w);

            }

            // Extract the current frontier from the previous next frontier and update node levels
            uint32_t startTileIdx = startNodeIdx/64;
            uint32_t numTiles = numNodes/64;
            if(startTileIdx <= nodeTileIdx && nodeTileIdx < startTileIdx + numTiles) {

                // Update current frontier
                store8B(nextFrontierTile, currentFrontier_m, nodeTileIdx - startTileIdx, cache_w);

                // Update node levels
                if(nextFrontierTile) {
                    for(uint32_t node = nodeTileIdx*64; node < (nodeTileIdx + 1)*64; ++node) {
                        if(isSet(nextFrontierTile, node%64)) {
                            store4B(level, nodeLevel_m, node - startNodeIdx, cache_w); // No false sharing so no need for locks
                        }
                    }
                }
            }

        }

        // Wait until all tasklets have updated the current frontier
        barrier_wait(&bfsBarrier);

        // Identify tasklet's nodes
        uint32_t numNodesPerTasklet = (numNodes + NR_TASKLETS - 1)/NR_TASKLETS;
        uint32_t taskletNodesStart = me()*numNodesPerTasklet;
        uint32_t taskletNumNodes;
        if(taskletNodesStart > numNodes) {
            taskletNumNodes = 0;
        } else if(taskletNodesStart + numNodesPerTasklet > numNodes) {
            taskletNumNodes = numNodes - taskletNodesStart;
        } else {
            taskletNumNodes = numNodesPerTasklet;
        }

        // Visit neighbors of the current frontier
        mutex_id_t mutexID = MUTEX_GET(nextFrontierMutex);
        for(uint32_t node = taskletNodesStart; node < taskletNodesStart + taskletNumNodes; ++node) {
            uint32_t nodeTileIdx = node/64;
            uint64_t currentFrontierTile = load8B(currentFrontier_m, nodeTileIdx, cache_w); // TODO: Optimize: load tile then loop over nodes in the tile
            if(isSet(currentFrontierTile, node%64)) { // If the node is in the current frontier
                // Visit its neighbors
                uint32_t nodePtr = load4B(nodePtrs_m, node, cache_w) - nodePtrsOffset;
                uint32_t nextNodePtr = load4B(nodePtrs_m, node + 1, cache_w) - nodePtrsOffset; // TODO: Optimize: might be in the same 8B as nodePtr
                for(uint32_t i = nodePtr; i < nextNodePtr; ++i) {
                    uint32_t neighbor = load4B(neighborIdxs_m, i, cache_w); // TODO: Optimize: sequential access to neighbors can use sequential reader
                    uint32_t neighborTileIdx = neighbor/64;
                    uint64_t visitedTile = load8B(visited_m, neighborTileIdx, cache_w);
                    if(!isSet(visitedTile, neighbor%64)) { // Neighbor not previously visited
                        // Add neighbor to next frontier
                        mutex_lock(mutexID); // TODO: Optimize: use more locks to reduce contention
                        uint64_t nextFrontierTile = load8B(nextFrontier_m, neighborTileIdx, cache_w);
                        setBit(nextFrontierTile, neighbor%64);
                        store8B(nextFrontierTile, nextFrontier_m, neighborTileIdx, cache_w);
                        mutex_unlock(mutexID);
                    }
                }
            }
        }

    }

    return 0;
}
