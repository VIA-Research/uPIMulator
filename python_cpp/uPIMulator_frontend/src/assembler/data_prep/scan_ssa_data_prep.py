import math
from typing import List, Optional

import numpy as np

from abi.word.immediate import Immediate
from abi.word.representation import Representation
from assembler.data_prep.bin import Bin
from encoder.byte import Byte
from initializer.int_initializer import IntInitializer
from util.config_loader import ConfigLoader


class SCANSSADataPrep:
    def __init__(self, num_tasklets: int, data_prep_param: List[int], num_dpus: int):
        assert 0 < num_tasklets < ConfigLoader.max_num_tasklets()

        self._num_tasklets: int = num_tasklets
        self._num_dpus: int = num_dpus

        self._size: int = data_prep_param[0]

        self._num_executions: int = 2

        elem_size = 8

        regs = 128

        self._is_strong_scaling = True  # True --> strong scaling / False --> weak scaling
        input_size = self._size if self._is_strong_scaling else self._size * self._num_dpus

        input_size_dpu = ((input_size) - 1) // (self._num_dpus) + 1

        if input_size_dpu % (num_tasklets * regs) != 0:
            self._input_size_dpu_round = math.ceil(input_size_dpu / (num_tasklets * regs)) * (num_tasklets * regs)
        else:
            self._input_size_dpu_round = input_size_dpu

        self._buffer_a: List[int] = [
            IntInitializer.value_by_range(0, 100) if i < input_size else 0
            for i in range(self._input_size_dpu_round * self._num_dpus)
        ]

        self._buffer_c: List[int] = []
        for i in range(self._input_size_dpu_round * self._num_dpus):
            if i == 0:
                self._buffer_c.append(self._buffer_a[i])
            else:
                self._buffer_c.append(self._buffer_c[i - 1] + self._buffer_a[i])

        self._last_result_value: List[int] = [0 for _ in range(self._num_dpus)]
        self._result_t_count: List[List[int]] = [[0 for _ in range(self._num_tasklets)] for _ in range(self._num_dpus)]

        for dpu_id in range(self._num_dpus):
            start_elem = self._input_size_dpu_round * dpu_id
            end_elem = self._input_size_dpu_round * (dpu_id + 1)
            self._last_result_value[dpu_id] = np.sum(self._buffer_a[start_elem:end_elem])
            self._result_t_count[dpu_id][self._num_tasklets - 1] = self._last_result_value[dpu_id]

        self._dpu_arg_size: int = self._input_size_dpu_round * elem_size
        self._kernel: List[int] = [0, 1]
        self._tcount: List[List[int]] = [
            [0, 0 if dpu_id == 0 else np.sum(self._result_t_count[0:dpu_id])] for dpu_id in range(self._num_dpus)
        ]

    def num_executions(self) -> int:
        return self._num_executions

    def num_dpus(self) -> int:
        return self._num_dpus

    def input_dpu_mram_heap_pointer_name(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        if execution != 0:
            return None
        else:
            bytes_: List[Byte] = []
            start_elem = self._input_size_dpu_round * dpu_id
            end_elem = self._input_size_dpu_round * (dpu_id + 1)
            for element in self._buffer_a[start_elem:end_elem]:
                element_immediate = Immediate(Representation.UNSIGNED, 64, element)
                bytes_ += element_immediate.to_bytes()

        return Bin(bytes_)

    def output_dpu_mram_heap_pointer_name(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        if execution != self._num_executions - 1:
            return None
        else:
            bytes_: List[Byte] = []
            start_elem = self._input_size_dpu_round * dpu_id
            end_elem = self._input_size_dpu_round * (dpu_id + 1)

            for element in self._buffer_a[start_elem:end_elem]:
                element_immediate = Immediate(Representation.UNSIGNED, 64, element)
                bytes_ += element_immediate.to_bytes()

            for element in self._buffer_c[start_elem:end_elem]:
                element_immediate = Immediate(Representation.UNSIGNED, 64, element)
                bytes_ += element_immediate.to_bytes()

            return Bin(bytes_)

    def dpu_input_arguments(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        dpu_arg_size_immediate = Immediate(Representation.UNSIGNED, 32, self._dpu_arg_size)
        bytes_ += dpu_arg_size_immediate.to_bytes()

        kernel_immediate = Immediate(Representation.UNSIGNED, 32, self._kernel[execution])
        bytes_ += kernel_immediate.to_bytes()

        tcount_immediate = Immediate(Representation.UNSIGNED, 64, self._tcount[dpu_id][execution])
        bytes_ += tcount_immediate.to_bytes()

        return Bin(bytes_)

    def dpu_results(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        if execution != 0:
            return None
        else:
            bytes_: List[Byte] = []

            for element in self._result_t_count[dpu_id]:
                element_immediate = Immediate(Representation.UNSIGNED, 64, element)
                bytes_ += element_immediate.to_bytes()

        return Bin(bytes_)
