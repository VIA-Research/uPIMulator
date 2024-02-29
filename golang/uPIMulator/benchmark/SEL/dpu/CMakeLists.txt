SET(BL 10)

set(CMAKE_C_COMPILER "/root/upmem-2023.2.0-Linux-x86_64/bin/dpu-upmem-dpurte-clang")
set(CMAKE_C_FLAGS "-w -I/root/uPIMulator/benchmark/SEL/support -O2 -S -DNR_TASKLETS=${NR_TASKLETS} -DBL=${BL}")

file(GLOB_RECURSE SRCS *.c)

add_executable(SEL_device ${SRCS})

