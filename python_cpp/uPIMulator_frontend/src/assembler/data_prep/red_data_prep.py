from typing import List, Optional

from abi.word.immediate import Immediate
from abi.word.representation import Representation
from assembler.data_prep.bin import Bin
from encoder.byte import Byte
from initializer.int_initializer import IntInitializer
from util.config_loader import ConfigLoader


class Result:
    def __init__(self, cycle=0, t_count=0):
        self._cycle: int = cycle
        self._t_count: int = t_count


class REDDataPrep:
    def __init__(self, num_tasklets: int, data_prep_param: List[int], num_dpus: int):
        assert 0 < num_tasklets < ConfigLoader.max_num_tasklets()
        self._num_tasklets: int = num_tasklets
        self._num_dpus: int = num_dpus

        self._size: int = data_prep_param[0]

        self._num_executions: int = 1

        elem_size = 8

        self._is_strong_scaling = True  # True --> strong scaling / False --> weak scaling
        input_size = self._size if self._is_strong_scaling else self._size * self._num_dpus

        if (input_size * elem_size) % 8 != 0:
            input_size_8bytes = (input_size // 8) * 8 + 8
        else:
            input_size_8bytes = input_size

        input_size_dpu = ((input_size) - 1) // (self._num_dpus) + 1

        if (input_size_dpu * elem_size) % 8 != 0:
            self._input_size_dpu_8bytes = (input_size_dpu // 8) * 8 + 8
        else:
            self._input_size_dpu_8bytes = input_size_dpu

        self._buffer_a: List[int] = [
            IntInitializer.value_by_range(0, 100) for _ in range(self._input_size_dpu_8bytes * self._num_dpus)
        ]

        self._count: List[int] = [0 for _ in range(self._num_dpus)]
        self._dpu_arg_size: List[int] = [0 for _ in range(self._num_dpus)]

        for dpu_id in range(self._num_dpus):
            start_elem = self._input_size_dpu_8bytes * dpu_id
            if dpu_id != self._num_dpus - 1:
                end_elem = self._input_size_dpu_8bytes * (dpu_id + 1)
            else:
                end_elem = input_size

            self._count[dpu_id] = sum(self._buffer_a[start_elem:end_elem])

            if dpu_id != self._num_dpus - 1:
                self._dpu_arg_size[dpu_id] = self._input_size_dpu_8bytes * elem_size
            else:
                self._dpu_arg_size[dpu_id] = (
                    input_size_8bytes - (self._input_size_dpu_8bytes * (self._num_dpus - 1))
                ) * elem_size

        self._kernel: List[int] = [0 for _ in range(self._num_dpus)]
        self._input_t_count: List[int] = [0 for _ in range(self._num_dpus)]

        self._result: List[List[Result]] = [
            [Result(0, self._count[j]) if i == 0 else Result() for i in range(self._num_tasklets)]
            for j in range(self._num_dpus)
        ]

    def num_executions(self) -> int:
        return self._num_executions

    def num_dpus(self) -> int:
        return self._num_dpus

    def input_dpu_mram_heap_pointer_name(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        start_elem = self._input_size_dpu_8bytes * dpu_id

        for i in range(self._input_size_dpu_8bytes):
            element_immediate = Immediate(Representation.UNSIGNED, 64, self._buffer_a[start_elem + i])
            bytes_ += element_immediate.to_bytes()
        return Bin(bytes_)

    def output_dpu_mram_heap_pointer_name(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus
        return None

    def dpu_input_arguments(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []
        dpu_arg_size_immediate = Immediate(Representation.UNSIGNED, 32, self._dpu_arg_size[dpu_id])
        bytes_ += dpu_arg_size_immediate.to_bytes()

        kernel_immediate = Immediate(Representation.UNSIGNED, 32, self._kernel[dpu_id])
        bytes_ += kernel_immediate.to_bytes()

        input_t_count_immediate = Immediate(Representation.UNSIGNED, 64, self._input_t_count[dpu_id])
        bytes_ += input_t_count_immediate.to_bytes()

        return Bin(bytes_)

    def dpu_results(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []
        for element in self._result[dpu_id]:
            cycle_immediate = Immediate(Representation.UNSIGNED, 64, element._cycle)
            bytes_ += cycle_immediate.to_bytes()

            t_count_immediate = Immediate(Representation.UNSIGNED, 64, element._t_count)
            bytes_ += t_count_immediate.to_bytes()
        return Bin(bytes_)
