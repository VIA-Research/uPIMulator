SET(BL 10)

set(CMAKE_C_COMPILER "/root/upmem/bin/dpu-upmem-dpurte-clang")
set(CMAKE_C_FLAGS "-w -I/root/uPIMulator/benchmark/HST-S/support -O2 -S -DNR_TASKLETS=${NR_TASKLETS} -DBL=${BL}")

file(GLOB_RECURSE SRCS *.c)

add_executable(HST-S_device ${SRCS})

