
#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdint.h>

#include <omp.h>

#include "../../support/common.h"
#include "../../support/graph.h"
#include "../../support/params.h"
#include "../../support/timer.h"
#include "../../support/utils.h"

int main(int argc, char** argv) {

    // Process parameters
    struct Params p = input_params(argc, argv);

    // Initialize BFS data structures
    PRINT_INFO(p.verbosity >= 1, "Reading graph %s", p.fileName);
    struct COOGraph cooGraph = readCOOGraph(p.fileName);
    PRINT_INFO(p.verbosity >= 1, "    Graph has %d nodes and %d edges", cooGraph.numNodes, cooGraph.numEdges);
    struct CSRGraph csrGraph = coo2csr(cooGraph);
    uint32_t* nodeLevel = (uint32_t*) malloc(csrGraph.numNodes*sizeof(uint32_t));
    uint32_t* nodeLevelRef = (uint32_t*) malloc(csrGraph.numNodes*sizeof(uint32_t));
    for(uint32_t i = 0; i < csrGraph.numNodes; ++i) {
        nodeLevel[i] = UINT32_MAX; // Unreachable
        nodeLevelRef[i] = UINT32_MAX; // Unreachable
    }
    uint32_t srcNode = 0;

    // Initialize frontier double buffers
    uint32_t* buffer1 = (uint32_t*) malloc(csrGraph.numNodes*sizeof(uint32_t));
    uint32_t* buffer2 = (uint32_t*) malloc(csrGraph.numNodes*sizeof(uint32_t));
    uint32_t* prevFrontier = buffer1;
    uint32_t* currFrontier = buffer2;

    // Calculating result on CPU
    PRINT_INFO(p.verbosity >= 1, "Calculating result on CPU (OpenMP)");
    omp_set_num_threads(4);
    Timer timer;
    startTimer(&timer);
    nodeLevel[srcNode] = 0;
    prevFrontier[0] = srcNode;
    uint32_t numPrevFrontier = 1;
    for(uint32_t level = 1; numPrevFrontier > 0; ++level) {

        uint32_t numCurrFrontier = 0;

        // Visit nodes in the previous frontier
        #pragma omp parallel for
        for(uint32_t i = 0; i < numPrevFrontier; ++i) {
            uint32_t node = prevFrontier[i];
            for(uint32_t edge = csrGraph.nodePtrs[node]; edge < csrGraph.nodePtrs[node + 1]; ++edge) {
                uint32_t neighbor = csrGraph.neighborIdxs[edge];
                uint32_t justVisited = 0;
                #pragma omp critical
                {
                    if(nodeLevel[neighbor] == UINT32_MAX) { // Node not previously visited
                        nodeLevel[neighbor] = level;
                        justVisited = 1;
                    }
                }
                if(justVisited) {
                    uint32_t currFrontierIdx;
                    #pragma omp critical
                    {
                        currFrontierIdx = numCurrFrontier++;
                    }
                    currFrontier[currFrontierIdx] = neighbor;
                }
            }
        }

        // Swap buffers
        uint32_t* tmp = prevFrontier;
        prevFrontier = currFrontier;
        currFrontier = tmp;
        numPrevFrontier = numCurrFrontier;

    }
    stopTimer(&timer);
    if(p.verbosity == 0) PRINT("%f", getElapsedTime(timer)*1e3);
    PRINT_INFO(p.verbosity >= 1, "Elapsed time: %f ms", getElapsedTime(timer)*1e3);

    // Calculating result on CPU sequentially
    PRINT_INFO(p.verbosity >= 1, "Calculating result on CPU (sequential)");
    startTimer(&timer);
    nodeLevelRef[srcNode] = 0;
    prevFrontier[0] = srcNode;
    numPrevFrontier = 1;
    for(uint32_t level = 1; numPrevFrontier > 0; ++level) {

        uint32_t numCurrFrontier = 0;

        // Visit nodes in the previous frontier
        for(uint32_t i = 0; i < numPrevFrontier; ++i) {
            uint32_t node = prevFrontier[i];
            for(uint32_t edge = csrGraph.nodePtrs[node]; edge < csrGraph.nodePtrs[node + 1]; ++edge) {
                uint32_t neighbor = csrGraph.neighborIdxs[edge];
                uint32_t justVisited = 0;
                if(nodeLevelRef[neighbor] == UINT32_MAX) { // Node not previously visited
                    nodeLevelRef[neighbor] = level;
                    justVisited = 1;
                }
                if(justVisited) {
                    uint32_t currFrontierIdx;
                    currFrontierIdx = numCurrFrontier++;
                    currFrontier[currFrontierIdx] = neighbor;
                }
            }
        }

        // Swap buffers
        uint32_t* tmp = prevFrontier;
        prevFrontier = currFrontier;
        currFrontier = tmp;
        numPrevFrontier = numCurrFrontier;

    }
    stopTimer(&timer);
    if(p.verbosity == 0) PRINT("%f", getElapsedTime(timer)*1e3);
    PRINT_INFO(p.verbosity >= 1, "Elapsed time: %f ms", getElapsedTime(timer)*1e3);

    // Verifying result
    PRINT_INFO(p.verbosity >= 1, "Verifying the result");
    for(uint32_t nodeIdx = 0; nodeIdx < csrGraph.numNodes; ++nodeIdx) {
        if(nodeLevel[nodeIdx] != nodeLevelRef[nodeIdx]) {
            PRINT_ERROR("Mismatch at node %u (CPU sequential result = level %u, CPU parallel result = level %u)", nodeIdx, nodeLevelRef[nodeIdx], nodeLevel[nodeIdx]);
        }
    }


    // Deallocate data structures
    freeCOOGraph(cooGraph);
    freeCSRGraph(csrGraph);
    free(nodeLevel);
    free(buffer1);
    free(buffer2);

    return 0;

}

