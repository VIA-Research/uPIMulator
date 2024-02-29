
#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdint.h>

#include "../../support/common.h"
#include "../../support/graph.h"
#include "../../support/params.h"
#include "../../support/timer.h"
#include "../../support/utils.h"

__global__ void bfs_kernel(CSRGraph csrGraph, uint32_t* nodeLevel, uint32_t* prevFrontier, uint32_t* currFrontier, uint32_t numPrevFrontier, uint32_t* numCurrFrontier,  uint32_t level) {
    uint32_t i = blockIdx.x*blockDim.x + threadIdx.x;
    if(i < numPrevFrontier) {
        uint32_t node = prevFrontier[i];
        for(uint32_t edge = csrGraph.nodePtrs[node]; edge < csrGraph.nodePtrs[node + 1]; ++edge) {
            uint32_t neighbor = csrGraph.neighborIdxs[edge];
            if(atomicCAS(&nodeLevel[neighbor], UINT32_MAX, level) == UINT32_MAX) { // Node not previously visited
                uint32_t currFrontierIdx = atomicAdd(numCurrFrontier, 1);
                currFrontier[currFrontierIdx] = neighbor;
            }
        }
    }
}

int main(int argc, char** argv) {

    // Process parameters
    struct Params p = input_params(argc, argv);

    // Initialize BFS data structures
    PRINT_INFO(p.verbosity >= 1, "Reading graph %s", p.fileName);
    struct COOGraph cooGraph = readCOOGraph(p.fileName);
    PRINT_INFO(p.verbosity >= 1, "    Graph has %d nodes and %d edges", cooGraph.numNodes, cooGraph.numEdges);
    struct CSRGraph csrGraph = coo2csr(cooGraph);
    uint32_t* nodeLevel_cpu = (uint32_t*) malloc(csrGraph.numNodes*sizeof(uint32_t));
    uint32_t* nodeLevel_gpu = (uint32_t*) malloc(csrGraph.numNodes*sizeof(uint32_t));
    for(uint32_t i = 0; i < csrGraph.numNodes; ++i) {
        nodeLevel_cpu[i] = UINT32_MAX; // Unreachable
        nodeLevel_gpu[i] = UINT32_MAX; // Unreachable
    }
    uint32_t srcNode = 0;

    // Allocate GPU memory
    CSRGraph csrGraph_d;
    csrGraph_d.numNodes = csrGraph.numNodes;
    csrGraph_d.numEdges = csrGraph.numEdges;
    cudaMalloc((void**) &csrGraph_d.nodePtrs, (csrGraph_d.numNodes + 1)*sizeof(uint32_t));
    cudaMalloc((void**) &csrGraph_d.neighborIdxs, csrGraph_d.numEdges*sizeof(uint32_t));
    uint32_t* nodeLevel_d;
    cudaMalloc((void**) &nodeLevel_d, csrGraph_d.numNodes*sizeof(uint32_t));
    uint32_t* buffer1_d;
    cudaMalloc((void**) &buffer1_d, csrGraph_d.numNodes*sizeof(uint32_t));
    uint32_t* buffer2_d;
    cudaMalloc((void**) &buffer2_d, csrGraph_d.numNodes*sizeof(uint32_t));
    uint32_t* numCurrFrontier_d;
    cudaMalloc((void**) &numCurrFrontier_d, sizeof(uint32_t));
    uint32_t* prevFrontier_d = buffer1_d;
    uint32_t* currFrontier_d = buffer2_d;

    // Copy data to GPU
    cudaMemcpy(csrGraph_d.nodePtrs, csrGraph.nodePtrs, (csrGraph_d.numNodes + 1)*sizeof(uint32_t), cudaMemcpyHostToDevice);
    cudaMemcpy(csrGraph_d.neighborIdxs, csrGraph.neighborIdxs, csrGraph_d.numEdges*sizeof(uint32_t), cudaMemcpyHostToDevice);
    nodeLevel_gpu[srcNode] = 0;
    cudaMemcpy(nodeLevel_d, nodeLevel_gpu, csrGraph_d.numNodes*sizeof(uint32_t), cudaMemcpyHostToDevice);
    cudaMemcpy(prevFrontier_d, &srcNode, sizeof(uint32_t), cudaMemcpyHostToDevice);
    cudaDeviceSynchronize();

    // Calculating result on GPU
    PRINT_INFO(p.verbosity >= 1, "Calculating result on GPU");
    Timer timer;
    startTimer(&timer);
    uint32_t numPrevFrontier = 1;
    uint32_t numThreadsPerBlock = 256;
    for(uint32_t level = 1; numPrevFrontier > 0; ++level) {

        // Visit nodes in previous frontier
        cudaMemset(numCurrFrontier_d, 0, sizeof(uint32_t));
        uint32_t numBlocks = (numPrevFrontier + numThreadsPerBlock - 1)/numThreadsPerBlock;
        bfs_kernel <<< numBlocks, numThreadsPerBlock >>> (csrGraph_d, nodeLevel_d, prevFrontier_d, currFrontier_d, numPrevFrontier, numCurrFrontier_d, level);

        // Swap buffers
        uint32_t* tmp = prevFrontier_d;
        prevFrontier_d = currFrontier_d;
        currFrontier_d = tmp;
        cudaMemcpy(&numPrevFrontier, numCurrFrontier_d, sizeof(uint32_t), cudaMemcpyDeviceToHost);

    }
    cudaDeviceSynchronize();
    stopTimer(&timer);
    if(p.verbosity == 0) PRINT("%f", getElapsedTime(timer)*1e3);
    PRINT_INFO(p.verbosity >= 1, "Elapsed time: %f ms", getElapsedTime(timer)*1e3);

    // Copy data from GPU
    cudaMemcpy(nodeLevel_gpu, nodeLevel_d, csrGraph_d.numNodes*sizeof(uint32_t), cudaMemcpyDeviceToHost);
    cudaDeviceSynchronize();

    // Initialize frontier double buffers for CPU
    uint32_t* buffer1 = (uint32_t*) malloc(csrGraph.numNodes*sizeof(uint32_t));
    uint32_t* buffer2 = (uint32_t*) malloc(csrGraph.numNodes*sizeof(uint32_t));
    uint32_t* prevFrontier = buffer1;
    uint32_t* currFrontier = buffer2;

    // Calculating result on CPU
    PRINT_INFO(p.verbosity >= 1, "Calculating result on CPU");
    nodeLevel_cpu[srcNode] = 0;
    prevFrontier[0] = srcNode;
    numPrevFrontier = 1;
    for(uint32_t level = 1; numPrevFrontier > 0; ++level) { 

        uint32_t numCurrFrontier = 0;

        // Visit nodes in the previous frontier
        for(uint32_t i = 0; i < numPrevFrontier; ++i) {
            uint32_t node = prevFrontier[i];
            for(uint32_t edge = csrGraph.nodePtrs[node]; edge < csrGraph.nodePtrs[node + 1]; ++edge) {
                uint32_t neighbor = csrGraph.neighborIdxs[edge];
                if(nodeLevel_cpu[neighbor] == UINT32_MAX) { // Node not previously visited
                    nodeLevel_cpu[neighbor] = level;
                    currFrontier[numCurrFrontier] = neighbor;
                    ++numCurrFrontier;
                }
            }
        }

        // Swap buffers
        uint32_t* tmp = prevFrontier;
        prevFrontier = currFrontier;
        currFrontier = tmp;
        numPrevFrontier = numCurrFrontier;

    }

    // Verify result
    PRINT_INFO(p.verbosity >= 1, "Verifying the result");
    for(uint32_t i = 0; i < csrGraph.numNodes; ++i) {
        if(nodeLevel_cpu[i] != nodeLevel_gpu[i]) {
            printf("Mismatch detected at node %u (CPU result = %u, GPU result = %u)\n", i, nodeLevel_cpu[i], nodeLevel_gpu[i]);
            exit(0);
        }
    }

    // Deallocate data structures
    freeCOOGraph(cooGraph);
    freeCSRGraph(csrGraph);
    free(nodeLevel_cpu);
    free(nodeLevel_gpu);
    free(buffer1);
    free(buffer2);

    return 0;

}

