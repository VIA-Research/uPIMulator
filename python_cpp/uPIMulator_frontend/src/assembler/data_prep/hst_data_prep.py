from typing import List, Optional

from abi.word.immediate import Immediate
from abi.word.representation import Representation
from assembler.data_prep.bin import Bin
from encoder.byte import Byte
from initializer.int_initializer import IntInitializer
from util.config_loader import ConfigLoader


class HSTDataPrep:
    def __init__(self, num_tasklets: int, data_prep_param: List[int], num_dpus: int):
        assert 0 < num_tasklets < ConfigLoader.max_num_tasklets()

        self._num_tasklets: int = num_tasklets
        self._num_dpus: int = num_dpus

        self._size: int = data_prep_param[0]  
        self._num_bins: int = 256

        self._num_executions: int = 1

        elem_size = 4

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

        self._buffer_a: List[int] = [IntInitializer.value_by_range(0, 4096) for _ in range(input_size)]

        depth = 12

        self._buffer_c: List[List[int]] = [[0 for _ in range(self._num_bins)] for _ in range(self._num_dpus)]
        for dpu_id in range(self._num_dpus):
            start_elem = self._input_size_dpu_8bytes * dpu_id
            end_elem = self._input_size_dpu_8bytes * (dpu_id + 1)
            for elem in self._buffer_a[start_elem:end_elem]:
                self._buffer_c[dpu_id][(elem * self._num_bins) >> depth] += 1

        elem_size = 4

        input_size = self._size
        if (input_size * elem_size) % 8 != 0:
            input_size_8bytes = (self._size // 8) * 8 + 8
        else:
            input_size_8bytes = self._size

        self._dpu_arg_size: List[int] = [0 for _ in range(self._num_dpus)]
        for dpu_id in range(self._num_dpus):
            if dpu_id != self._num_dpus - 1:
                self._dpu_arg_size[dpu_id] = self._input_size_dpu_8bytes * elem_size
            else:
                self._dpu_arg_size[dpu_id] = (
                    input_size_8bytes - (self._input_size_dpu_8bytes * (self._num_dpus - 1))
                ) * elem_size
        self._transfer_size: int = self._input_size_dpu_8bytes * elem_size
        self._kernel: int = 0

    def num_executions(self) -> int:
        return self._num_executions

    def num_dpus(self) -> int:
        return self._num_dpus

    def input_dpu_mram_heap_pointer_name(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        start_elem = self._input_size_dpu_8bytes * dpu_id
        end_elem = self._input_size_dpu_8bytes * (dpu_id + 1)

        for element in self._buffer_a[start_elem:end_elem]:
            element_immediate = Immediate(Representation.UNSIGNED, 32, element)
            bytes_ += element_immediate.to_bytes()

        return Bin(bytes_)

    def output_dpu_mram_heap_pointer_name(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        start_elem = self._input_size_dpu_8bytes * dpu_id
        end_elem = self._input_size_dpu_8bytes * (dpu_id + 1)

        for element in self._buffer_a[start_elem:end_elem]:
            element_immediate = Immediate(Representation.UNSIGNED, 32, element)
            bytes_ += element_immediate.to_bytes()

        for element in self._buffer_c[dpu_id]:
            element_immediate = Immediate(Representation.UNSIGNED, 32, element)
            bytes_ += element_immediate.to_bytes()

        return Bin(bytes_)

    def dpu_input_arguments(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        dpu_arg_size_immediate = Immediate(Representation.UNSIGNED, 32, self._dpu_arg_size[dpu_id])
        bytes_ += dpu_arg_size_immediate.to_bytes()

        transfer_immediate = Immediate(Representation.UNSIGNED, 32, self._transfer_size)
        bytes_ += transfer_immediate.to_bytes()

        num_bins_immediate = Immediate(Representation.UNSIGNED, 32, self._num_bins)
        bytes_ += num_bins_immediate.to_bytes()

        kernel_immediate = Immediate(Representation.UNSIGNED, 32, self._kernel)
        bytes_ += kernel_immediate.to_bytes()

        return Bin(bytes_)

    def dpu_results(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        return None
