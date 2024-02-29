from typing import List, Optional

import numpy as np

from abi.word.immediate import Immediate
from abi.word.representation import Representation
from assembler.data_prep.bin import Bin
from encoder.byte import Byte
from initializer.int_initializer import IntInitializer
from util.config_loader import ConfigLoader


class VADataPrep:
    def __init__(self, num_tasklets: int, data_prep_param: List[int], num_dpus: int):
        assert 0 < num_tasklets < ConfigLoader.max_num_tasklets()

        self._num_tasklets: int = num_tasklets
        self._num_dpus: int = num_dpus

        # NOTE(bongjoon.hyun@gmail.com): PrIM default buffer size is 2621440
        self._buffer_size: int = data_prep_param[0]  # 1048576

        self._num_executions: int = 1

        elem_size = 4

        self._is_strong_scaling = True  # True --> strong scaling / False --> weak scaling
        input_size = self._buffer_size if self._is_strong_scaling else self._buffer_size * self._num_dpus

        if (input_size * elem_size) % 8 != 0:
            input_size_8bytes = (input_size // 8) * 8 + 8
        else:
            input_size_8bytes = input_size

        input_size_dpu = ((input_size) - 1) // (self._num_dpus) + 1

        if (input_size_dpu * elem_size) % 8 != 0:
            self._input_size_dpu_8bytes = (input_size_dpu // 8) * 8 + 8
        else:
            self._input_size_dpu_8bytes = input_size_dpu

        self._buffer_a: List[int] = [0 for _ in range(self._input_size_dpu_8bytes * num_dpus)]
        self._buffer_b: List[int] = [0 for _ in range(self._input_size_dpu_8bytes * num_dpus)]

        for i in range(input_size):
            self._buffer_a[i] = IntInitializer.value_by_range(0, 2**31)
            self._buffer_b[i] = IntInitializer.value_by_range(0, 2**31)

        self._buffer_c: List[int] = list(np.add(self._buffer_a, self._buffer_b))

        self._size: List[int] = [self._input_size_dpu_8bytes * elem_size for _ in range(self._num_dpus - 1)]
        self._size.append((input_size_8bytes - (self._input_size_dpu_8bytes * (self._num_dpus - 1))) * elem_size)

        self._transfer_size: List[int] = [self._input_size_dpu_8bytes * elem_size for _ in range(self._num_dpus)]

        self._kernel: List[int] = [0 for _ in range(self._num_dpus)]

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
            element_immediate = Immediate(Representation.UNSIGNED, 32, self._buffer_a[start_elem + i])
            bytes_ += element_immediate.to_bytes()

        for i in range(self._input_size_dpu_8bytes):
            element_immediate = Immediate(Representation.UNSIGNED, 32, self._buffer_b[start_elem + i])
            bytes_ += element_immediate.to_bytes()

        return Bin(bytes_)

    def output_dpu_mram_heap_pointer_name(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        start_elem = self._input_size_dpu_8bytes * dpu_id

        for i in range(self._input_size_dpu_8bytes):
            element_immediate = Immediate(Representation.UNSIGNED, 32, self._buffer_a[start_elem + i])
            bytes_ += element_immediate.to_bytes()

        for i in range(self._input_size_dpu_8bytes):
            element_immediate = Immediate(Representation.UNSIGNED, 32, self._buffer_c[start_elem + i])
            bytes_ += element_immediate.to_bytes()

        return Bin(bytes_)

    def dpu_input_arguments(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        size_immediate = Immediate(Representation.UNSIGNED, 32, self._size[dpu_id])
        bytes_ += size_immediate.to_bytes()

        transfer_size_immediate = Immediate(Representation.UNSIGNED, 32, self._transfer_size[dpu_id])
        bytes_ += transfer_size_immediate.to_bytes()

        kernel_immediate = Immediate(Representation.UNSIGNED, 32, self._kernel[dpu_id])
        bytes_ += kernel_immediate.to_bytes()

        return Bin(bytes_)

    def dpu_results(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        return None
