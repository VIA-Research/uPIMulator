DPU_DIR := dpu
HOST_DIR := host
CPU_BASE_DIR := baselines/cpu
GPU_BASE_DIR := baselines/gpu
BUILDDIR ?= bin
NR_TASKLETS ?= 16
NR_DPUS ?= 1

define conf_filename
	${BUILDDIR}/.NR_DPUS_$(1)_NR_TASKLETS_$(2).conf
endef
CONF := $(call conf_filename,${NR_DPUS},${NR_TASKLETS})

HOST_TARGET := ${BUILDDIR}/host_code
DPU_TARGET := ${BUILDDIR}/dpu_code
CPU_BASE_TARGET := ${BUILDDIR}/cpu_baseline
GPU_BASE_TARGET := ${BUILDDIR}/gpu_baseline

COMMON_INCLUDES := support
HOST_SOURCES := $(wildcard ${HOST_DIR}/*.c)
DPU_SOURCES := $(wildcard ${DPU_DIR}/*.c)
CPU_BASE_SOURCES := $(wildcard ${CPU_BASE_DIR}/*.c)
GPU_BASE_SOURCES := $(wildcard ${GPU_BASE_DIR}/*.cu)

.PHONY: all clean test

__dirs := $(shell mkdir -p ${BUILDDIR})

COMMON_FLAGS := -w -I${COMMON_INCLUDES}
HOST_FLAGS := ${COMMON_FLAGS} -std=c11 -O3 `dpu-pkg-config --cflags --libs dpu` -DNR_TASKLETS=${NR_TASKLETS} -DNR_DPUS=${NR_DPUS} 
DPU_FLAGS := ${COMMON_FLAGS} -O2 -DNR_TASKLETS=${NR_TASKLETS} 
CPU_BASE_FLAGS := -O3 -fopenmp
GPU_BASE_FLAGS := -O3

all: ${HOST_TARGET} ${DPU_TARGET} ${CPU_BASE_TARGET}

gpu: ${GPU_BASE_TARGET}

${CONF}:
	$(RM) $(call conf_filename,*,*)
	touch ${CONF}

${HOST_TARGET}: ${HOST_SOURCES} ${COMMON_INCLUDES} ${CONF}
	$(CC) -o $@ ${HOST_SOURCES} ${HOST_FLAGS}
	$(CC) -S -o ${HOST_TARGET}.S ${HOST_SOURCES} ${HOST_FLAGS}

${DPU_TARGET}: ${DPU_SOURCES} ${COMMON_INCLUDES} ${CONF}
	dpu-upmem-dpurte-clang ${DPU_FLAGS} -o $@ ${DPU_SOURCES}
	dpu-upmem-dpurte-clang -S ${DPU_FLAGS} -o ${DPU_TARGET}.S ${DPU_SOURCES}

${CPU_BASE_TARGET}: ${CPU_BASE_SOURCES}
	$(CC) -o $@ ${CPU_BASE_SOURCES} ${CPU_BASE_FLAGS}

${GPU_BASE_TARGET}: ${GPU_BASE_SOURCES}
	nvcc -o $@ ${GPU_BASE_SOURCES} ${GPU_BASE_FLAGS}

clean:
	$(RM) -r $(BUILDDIR)

test: all
	./${HOST_TARGET}

