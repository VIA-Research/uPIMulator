add_executable(task_mmu_test task.c)
target_include_directories(task_mmu_test PRIVATE ${DPU_INCLUDE_DIRS})
target_compile_options(task_mmu_test PRIVATE ${DPU_COMPILE_FLAGS})
