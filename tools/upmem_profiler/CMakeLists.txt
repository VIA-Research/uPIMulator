cmake_minimum_required(VERSION 3.16)

#set(CMAKE_CXX_COMPILER "/opt/rh/devtoolset-8/root/usr/bin/g++")
set(CMAKE_CXX_COMPILER "/usr/bin/g++")
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_FLAGS "-O3")
#set(CMAKE_CXX_FLAGS  "-g")

link_libraries(stdc++fs)

project(upmem_profiler)

add_subdirectory(src)
