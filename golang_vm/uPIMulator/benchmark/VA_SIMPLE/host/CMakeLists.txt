set(CMAKE_C_COMPILER "/usr/bin/gcc")
set(CMAKE_C_FLAGS "-std=c99 -O3")

find_package(PkgConfig REQUIRED)
pkg_check_modules(DPU REQUIRED dpu)

file(GLOB_RECURSE SRCS *.c)

add_executable(VA_SIMPLE_host ${SRCS})

target_link_libraries(VA_SIMPLE_host ${DPU_LIBRARIES})
target_include_directories(VA_SIMPLE_host PUBLIC ${DPU_INCLUDE_DIRS})
target_compile_options(VA_SIMPLE_host PUBLIC ${DPU_CFLAGS_OTHER})
