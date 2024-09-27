SET(BL 10)
SET(TYPE INT64)

set(CMAKE_C_COMPILER "/root/upmem-2023.2.0-Linux-x86_64/bin/dpu-upmem-dpurte-clang")
set(CMAKE_C_FLAGS "-w -I/root/uPIMulator/benchmark/SCAN-RSS/support -O2 -S -DNR_TASKLETS=${NR_TASKLETS} -DBL=${BL} -D${TYPE}")

file(GLOB_RECURSE SRCS *.c)

add_executable(SCAN-RSS_device ${SRCS})

