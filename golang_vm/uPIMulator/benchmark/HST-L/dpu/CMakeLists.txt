SET(BL 10)
SET(NR_HISTO 1)

set(CMAKE_C_COMPILER "/root/upmem-2023.2.0-Linux-x86_64/bin/dpu-upmem-dpurte-clang")
set(CMAKE_C_FLAGS "-w -I/root/uPIMulator/benchmark/HST-L/support -O2 -S -DNR_TASKLETS=${NR_TASKLETS} -DBL=${BL} -DNR_HISTO=${NR_HISTO}")

file(GLOB_RECURSE SRCS *.c)

add_executable(HST-L_device ${SRCS})

