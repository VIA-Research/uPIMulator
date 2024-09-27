set(CMAKE_C_COMPILER "/root/upmem-2023.2.0-Linux-x86_64/bin/dpu-upmem-dpurte-clang")
set(CMAKE_C_FLAGS "-O3 -S -DNR_TASKLETS=${NR_TASKLETS}")

include_directories("/root/uPIMulator/sdk/misc")
include_directories("/root/uPIMulator/sdk/stdlib")
include_directories("/root/uPIMulator/sdk/syslib")

file(GLOB_RECURSE SRCS *.c)

add_library(syslib ${SRCS})
