cmake_minimum_required(VERSION 3.16)

include_directories(.)

file(GLOB_RECURSE SRCS *.cc)

add_executable(uPIMulator ${SRCS})
