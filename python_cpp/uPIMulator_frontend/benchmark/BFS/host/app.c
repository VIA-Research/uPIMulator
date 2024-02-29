/**
* app.c
* BFS Host Application Source File
*
*/
#include <dpu.h>
#include <dpu_log.h>

#include <assert.h>
#include <getopt.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "mram-management.h"
#include "../support/common.h"
#include "../support/graph.h"
#include "../support/params.h"
#include "../support/timer.h"
#include "../support/utils.h"

#ifndef ENERGY
#define ENERGY 0
#endif
#if ENERGY
#include <dpu_probe.h>
#endif

#define DPU_BINARY "./bin/dpu_code"

// Main of the Host Application
int main(int argc, char** argv) {

    // Process parameters
    struct Params p = input_params(argc, argv);

    // Timer and profiling
    Timer timer;
    float loadTime = 0.0f, dpuTime = 0.0f, hostTime = 0.0f, retrieveTime = 0.0f;
    #if ENERGY
    struct dpu_probe_t probe;
    DPU_ASSERT(dpu_probe_init("energy_probe", &probe));
    double tenergy=0;
    #endif

    // Allocate DPUs and load binary
    struct dpu_set_t dpu_set, dpu;
    uint32_t numDPUs;
    DPU_ASSERT(dpu_alloc(NR_DPUS, NULL, &dpu_set));
    DPU_ASSERT(dpu_load(dpu_set, DPU_BINARY, NULL));
    DPU_ASSERT(dpu_get_nr_dpus(dpu_set, &numDPUs));
    PRINT_INFO(p.verbosity >= 1, "Allocated %d DPU(s)", numDPUs);

    // Initialize BFS data structures
    PRINT_INFO(p.verbosity >= 1, "Reading graph %s", p.fileName);
    struct COOGraph cooGraph = readCOOGraph(p.fileName);
    PRINT_INFO(p.verbosity >= 1, "    Graph has %d nodes and %d edges", cooGraph.numNodes, cooGraph.numEdges);
    struct CSRGraph csrGraph = coo2csr(cooGraph);
    uint32_t numNodes = csrGraph.numNodes;
    uint32_t* nodePtrs = csrGraph.nodePtrs;
    uint32_t* neighborIdxs = csrGraph.neighborIdxs;
    uint32_t* nodeLevel = calloc(numNodes, sizeof(uint32_t)); // Node's BFS level (initially all 0 meaning not reachable)
    uint64_t* visited = calloc(numNodes/64, sizeof(uint64_t)); // Bit vector with one bit per node
    uint64_t* currentFrontier = calloc(numNodes/64, sizeof(uint64_t)); // Bit vector with one bit per node
    uint64_t* nextFrontier = calloc(numNodes/64, sizeof(uint64_t)); // Bit vector with one bit per node
    setBit(nextFrontier[0], 0); // Initialize frontier to first node
    uint32_t level = 1;

    // Partition data structure across DPUs
    uint32_t numNodesPerDPU = ROUND_UP_TO_MULTIPLE_OF_64((numNodes - 1)/numDPUs + 1);
    PRINT_INFO(p.verbosity >= 1, "Assigning %u nodes per DPU", numNodesPerDPU);
    struct DPUParams dpuParams[numDPUs];
    uint32_t dpuParams_m[numDPUs];
    unsigned int dpuIdx = 0;
    DPU_FOREACH (dpu_set, dpu) {

        // Allocate parameters
        struct mram_heap_allocator_t allocator;
        init_allocator(&allocator);
        dpuParams_m[dpuIdx] = mram_heap_alloc(&allocator, sizeof(struct DPUParams));

        // Find DPU's nodes
        uint32_t dpuStartNodeIdx = dpuIdx*numNodesPerDPU;
        uint32_t dpuNumNodes;
        if(dpuStartNodeIdx > numNodes) {
            dpuNumNodes = 0;
        } else if(dpuStartNodeIdx + numNodesPerDPU > numNodes) {
            dpuNumNodes = numNodes - dpuStartNodeIdx;
        } else {
            dpuNumNodes = numNodesPerDPU;
        }
        dpuParams[dpuIdx].dpuNumNodes = dpuNumNodes;
        PRINT_INFO(p.verbosity >= 2, "    DPU %u:", dpuIdx);
        PRINT_INFO(p.verbosity >= 2, "        Receives %u nodes", dpuNumNodes);

        // Partition edges and copy data
        if(dpuNumNodes > 0) {

            // Find DPU's CSR graph partition
            uint32_t* dpuNodePtrs_h = &nodePtrs[dpuStartNodeIdx];
            uint32_t dpuNodePtrsOffset = dpuNodePtrs_h[0];
            uint32_t* dpuNeighborIdxs_h = neighborIdxs + dpuNodePtrsOffset;
            uint32_t dpuNumNeighbors = dpuNodePtrs_h[dpuNumNodes] - dpuNodePtrsOffset;
            uint32_t* dpuNodeLevel_h = &nodeLevel[dpuStartNodeIdx];

            // Allocate MRAM
            uint32_t dpuNodePtrs_m = mram_heap_alloc(&allocator, (dpuNumNodes + 1)*sizeof(uint32_t));
            uint32_t dpuNeighborIdxs_m = mram_heap_alloc(&allocator, dpuNumNeighbors*sizeof(uint32_t));
            uint32_t dpuNodeLevel_m = mram_heap_alloc(&allocator, dpuNumNodes*sizeof(uint32_t));
            uint32_t dpuVisited_m = mram_heap_alloc(&allocator, numNodes/64*sizeof(uint64_t));
            uint32_t dpuCurrentFrontier_m = mram_heap_alloc(&allocator, dpuNumNodes/64*sizeof(uint64_t));
            uint32_t dpuNextFrontier_m = mram_heap_alloc(&allocator, numNodes/64*sizeof(uint64_t));
            PRINT_INFO(p.verbosity >= 2, "        Total memory allocated is %d bytes", allocator.totalAllocated);

            // Set up DPU parameters
            dpuParams[dpuIdx].numNodes = numNodes;
            dpuParams[dpuIdx].dpuStartNodeIdx = dpuStartNodeIdx;
            dpuParams[dpuIdx].dpuNodePtrsOffset = dpuNodePtrsOffset;
            dpuParams[dpuIdx].level = level;
            dpuParams[dpuIdx].dpuNodePtrs_m = dpuNodePtrs_m;
            dpuParams[dpuIdx].dpuNeighborIdxs_m = dpuNeighborIdxs_m;
            dpuParams[dpuIdx].dpuNodeLevel_m = dpuNodeLevel_m;
            dpuParams[dpuIdx].dpuVisited_m = dpuVisited_m;
            dpuParams[dpuIdx].dpuCurrentFrontier_m = dpuCurrentFrontier_m;
            dpuParams[dpuIdx].dpuNextFrontier_m = dpuNextFrontier_m;

            // Send data to DPU
            PRINT_INFO(p.verbosity >= 2, "        Copying data to DPU");
            startTimer(&timer);
            copyToDPU(dpu, (uint8_t*)dpuNodePtrs_h, dpuNodePtrs_m, (dpuNumNodes + 1)*sizeof(uint32_t));
            copyToDPU(dpu, (uint8_t*)dpuNeighborIdxs_h, dpuNeighborIdxs_m, dpuNumNeighbors*sizeof(uint32_t));
            copyToDPU(dpu, (uint8_t*)dpuNodeLevel_h, dpuNodeLevel_m, dpuNumNodes*sizeof(uint32_t));
            copyToDPU(dpu, (uint8_t*)visited, dpuVisited_m, numNodes/64*sizeof(uint64_t));
            copyToDPU(dpu, (uint8_t*)nextFrontier, dpuNextFrontier_m, numNodes/64*sizeof(uint64_t));
            // NOTE: No need to copy current frontier because it is written before being read
            stopTimer(&timer);
            loadTime += getElapsedTime(timer);

        }

        // Send parameters to DPU
        PRINT_INFO(p.verbosity >= 2, "        Copying parameters to DPU");
        startTimer(&timer);
        copyToDPU(dpu, (uint8_t*)&dpuParams[dpuIdx], dpuParams_m[dpuIdx], sizeof(struct DPUParams));
        stopTimer(&timer);
        loadTime += getElapsedTime(timer);

        ++dpuIdx;

    }
    PRINT_INFO(p.verbosity >= 1, "    CPU-DPU Time: %f ms", loadTime*1e3);

    // Iterate until next frontier is empty
    uint32_t nextFrontierEmpty = 0;
    while(!nextFrontierEmpty) {

        PRINT_INFO(p.verbosity >= 1, "Processing current frontier for level %u", level);

	#if ENERGY
	DPU_ASSERT(dpu_probe_start(&probe));
	#endif
        // Run all DPUs
        PRINT_INFO(p.verbosity >= 1, "    Booting DPUs");
        startTimer(&timer);
        DPU_ASSERT(dpu_launch(dpu_set, DPU_SYNCHRONOUS));
        stopTimer(&timer);
        dpuTime += getElapsedTime(timer);
        PRINT_INFO(p.verbosity >= 2, "    Level DPU Time: %f ms", getElapsedTime(timer)*1e3);
	#if ENERGY
    	DPU_ASSERT(dpu_probe_stop(&probe));
    	double energy;
    	DPU_ASSERT(dpu_probe_get(&probe, DPU_ENERGY, DPU_AVERAGE, &energy));
	tenergy += energy;
	#endif



        // Copy back next frontier from all DPUs and compute their union as the current frontier
        startTimer(&timer);
        dpuIdx = 0;
        DPU_FOREACH (dpu_set, dpu) {
            uint32_t dpuNumNodes = dpuParams[dpuIdx].dpuNumNodes;
            if(dpuNumNodes > 0) {
                if(dpuIdx == 0) {
                    copyFromDPU(dpu, dpuParams[dpuIdx].dpuNextFrontier_m, (uint8_t*)currentFrontier, numNodes/64*sizeof(uint64_t));
                } else {
                    copyFromDPU(dpu, dpuParams[dpuIdx].dpuNextFrontier_m, (uint8_t*)nextFrontier, numNodes/64*sizeof(uint64_t));
                    for(uint32_t i = 0; i < numNodes/64; ++i) {
                        currentFrontier[i] |= nextFrontier[i];
                    }
                }
                ++dpuIdx;
            }
        }

        // Check if the next frontier is empty, and copy data to DPU if not empty
        nextFrontierEmpty = 1;
        for(uint32_t i = 0; i < numNodes/64; ++i) {
            if(currentFrontier[i]) {
                nextFrontierEmpty = 0;
                break;
            }
        }
        if(!nextFrontierEmpty) {
            ++level;
            dpuIdx = 0;
            DPU_FOREACH (dpu_set, dpu) {
                uint32_t dpuNumNodes = dpuParams[dpuIdx].dpuNumNodes;
                if(dpuNumNodes > 0) {
                    // Copy current frontier to all DPUs (place in next frontier and DPU will update visited and copy to current frontier)
                    copyToDPU(dpu, (uint8_t*)currentFrontier, dpuParams[dpuIdx].dpuNextFrontier_m, numNodes/64*sizeof(uint64_t));
                    // Copy new level to DPU
                    dpuParams[dpuIdx].level = level;
                    copyToDPU(dpu, (uint8_t*)&dpuParams[dpuIdx], dpuParams_m[dpuIdx], sizeof(struct DPUParams));
                    ++dpuIdx;
                }
            }
        }
        stopTimer(&timer);
        hostTime += getElapsedTime(timer);
        PRINT_INFO(p.verbosity >= 2, "    Level Inter-DPU Time: %f ms", getElapsedTime(timer)*1e3);

    }
    PRINT_INFO(p.verbosity >= 1, "DPU Kernel Time: %f ms", dpuTime*1e3);
    PRINT_INFO(p.verbosity >= 1, "Inter-DPU Time: %f ms", hostTime*1e3);
    #if ENERGY
    PRINT_INFO(p.verbosity >= 1, "    DPU Energy: %f J", tenergy);
    #endif

    // Copy back node levels
    PRINT_INFO(p.verbosity >= 1, "Copying back the result");
    startTimer(&timer);
    dpuIdx = 0;
    DPU_FOREACH (dpu_set, dpu) {
        uint32_t dpuNumNodes = dpuParams[dpuIdx].dpuNumNodes;
        if(dpuNumNodes > 0) {
            uint32_t dpuStartNodeIdx = dpuIdx*numNodesPerDPU;
            copyFromDPU(dpu, dpuParams[dpuIdx].dpuNodeLevel_m, (uint8_t*)(nodeLevel + dpuStartNodeIdx), dpuNumNodes*sizeof(float));
        }
        ++dpuIdx;
    }
    stopTimer(&timer);
    retrieveTime += getElapsedTime(timer);
    PRINT_INFO(p.verbosity >= 1, "    DPU-CPU Time: %f ms", retrieveTime*1e3);
    if(p.verbosity == 0) PRINT("CPU-DPU Time(ms): %f    DPU Kernel Time (ms): %f    Inter-DPU Time (ms): %f    DPU-CPU Time (ms): %f", loadTime*1e3, dpuTime*1e3, hostTime*1e3, retrieveTime*1e3);

    // Calculating result on CPU
    PRINT_INFO(p.verbosity >= 1, "Calculating result on CPU");
    uint32_t* nodeLevelReference = calloc(numNodes, sizeof(uint32_t)); // Node's BFS level (initially all 0 meaning not reachable)
    memset(nextFrontier, 0, numNodes/64*sizeof(uint64_t));
    setBit(nextFrontier[0], 0); // Initialize frontier to first node
    nextFrontierEmpty = 0;
    level = 1;
    while(!nextFrontierEmpty) {
        // Update current frontier and visited list based on the next frontier from the previous iteration
        for(uint32_t nodeTileIdx = 0; nodeTileIdx < numNodes/64; ++nodeTileIdx) {
            uint64_t nextFrontierTile = nextFrontier[nodeTileIdx];
            currentFrontier[nodeTileIdx] = nextFrontierTile;
            if(nextFrontierTile) {
                visited[nodeTileIdx] |= nextFrontierTile;
                nextFrontier[nodeTileIdx] = 0;
                for(uint32_t node = nodeTileIdx*64; node < (nodeTileIdx + 1)*64; ++node) {
                    if(isSet(nextFrontierTile, node%64)) {
                        nodeLevelReference[node] = level;
                    }
                }
            }
        }
        // Visit neighbors of the current frontier
        nextFrontierEmpty = 1;
        for(uint32_t nodeTileIdx = 0; nodeTileIdx < numNodes/64; ++nodeTileIdx) {
            uint64_t currentFrontierTile = currentFrontier[nodeTileIdx];
            if(currentFrontierTile) {
                for(uint32_t node = nodeTileIdx*64; node < (nodeTileIdx + 1)*64; ++node) {
                    if(isSet(currentFrontierTile, node%64)) { // If the node is in the current frontier
                        // Visit its neighbors
                        uint32_t nodePtr = nodePtrs[node];
                        uint32_t nextNodePtr = nodePtrs[node + 1];
                        for(uint32_t i = nodePtr; i < nextNodePtr; ++i) {
                            uint32_t neighbor = neighborIdxs[i];
                            if(!isSet(visited[neighbor/64], neighbor%64)) { // Neighbor not previously visited
                                // Add neighbor to next frontier
                                setBit(nextFrontier[neighbor/64], neighbor%64);
                                nextFrontierEmpty = 0;
                            }
                        }
                    }
                }
            }
        }
        ++level;
    }

    // Verify the result
    PRINT_INFO(p.verbosity >= 1, "Verifying the result");
    for(uint32_t nodeIdx = 0; nodeIdx < numNodes; ++nodeIdx) {
        if(nodeLevel[nodeIdx] != nodeLevelReference[nodeIdx]) {
            PRINT_ERROR("Mismatch at node %u (CPU result = level %u, DPU result = level %u)", nodeIdx, nodeLevelReference[nodeIdx], nodeLevel[nodeIdx]);
        }
    }

    // Display DPU Logs
    if(p.verbosity >= 2) {
        PRINT_INFO(p.verbosity >= 2, "Displaying DPU Logs:");
        dpuIdx = 0;
        DPU_FOREACH (dpu_set, dpu) {
            PRINT("DPU %u:", dpuIdx);
            DPU_ASSERT(dpu_log_read(dpu, stdout));
            ++dpuIdx;
        }
    }

    // Deallocate data structures
    freeCOOGraph(cooGraph);
    freeCSRGraph(csrGraph);
    free(nodeLevel);
    free(visited);
    free(currentFrontier);
    free(nextFrontier);
    free(nodeLevelReference);

    return 0;

}

