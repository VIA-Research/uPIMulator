add_executable(host_mmu_test host.c)
target_include_directories(host_mmu_test PRIVATE ${HOST_INCLUDE_DIRS})
target_link_libraries(host_mmu_test dpu)
