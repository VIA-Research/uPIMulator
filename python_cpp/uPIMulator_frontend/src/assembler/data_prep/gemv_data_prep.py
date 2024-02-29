from typing import List, Optional

import numpy as np

from abi.word.immediate import Immediate
from abi.word.representation import Representation
from assembler.data_prep.bin import Bin
from encoder.byte import Byte
from initializer.int_initializer import IntInitializer
from util.config_loader import ConfigLoader


class GEMVDataPrep:
    def __init__(self, num_tasklets: int, data_prep_param: List[int], num_dpus: int):
        assert 0 < num_tasklets < ConfigLoader.max_num_tasklets()

        self._num_tasklets: int = num_tasklets
        self._num_dpus: int = num_dpus

        self._m_size: int = data_prep_param[0]
        self._n_size: int = 64

        self._num_executions: int = 1

        assert self._m_size % self._num_dpus == 0

        self._n_size_pad: List[int] = [
            (self._n_size if self._n_size % 2 == 0 else self._n_size + 1) for _ in range(self._num_dpus)
        ]
        self._nr_rows: List[int] = [self._m_size // self._num_dpus for _ in range(self._num_dpus)]
        self._max_rows: List[int] = [
            (self._nr_rows[i] if self._nr_rows[i] % 2 == 0 else self._nr_rows[i] + 1) for i in range(self._num_dpus)
        ]

        self._buffer_a: List[List[int]] = [
            [IntInitializer.value_by_range(0, 50) for _ in range(self._n_size)] for _ in range(self._m_size)
        ]

        self._buffer_b: List[int] = [IntInitializer.value_by_range(0, 50) for _ in range(self._n_size)]

        self._buffer_c: List[List[int]] = [[] for _ in range(self._num_dpus)]

        for dpu_id in range(self._num_dpus):
            start_row: int = self._nr_rows[dpu_id] * dpu_id
            end_row: int = self._nr_rows[dpu_id] * (dpu_id + 1)
            self._buffer_c[dpu_id] = np.matmul(self._buffer_a[start_row:end_row], self._buffer_b).tolist()

    def num_executions(self) -> int:
        return self._num_executions

    def num_dpus(self) -> int:
        return self._num_dpus

    def input_dpu_mram_heap_pointer_name(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        start_row: int = self._nr_rows[dpu_id] * dpu_id
        end_row: int = self._nr_rows[dpu_id] * (dpu_id + 1)

        for row in self._buffer_a[start_row:end_row]:
            for element in row:
                element_immediate = Immediate(Representation.UNSIGNED, 32, element)
                bytes_ += element_immediate.to_bytes()

        for element in self._buffer_b:
            element_immediate = Immediate(Representation.UNSIGNED, 32, element)
            bytes_ += element_immediate.to_bytes()

        return Bin(bytes_)

    def output_dpu_mram_heap_pointer_name(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        start_row: int = self._nr_rows[dpu_id] * dpu_id
        end_row: int = self._nr_rows[dpu_id] * (dpu_id + 1)

        for row in self._buffer_a[start_row:end_row]:
            for element in row:
                element_immediate = Immediate(Representation.UNSIGNED, 32, element)
                bytes_ += element_immediate.to_bytes()

        for element in self._buffer_b:
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

        n_size_immediate = Immediate(Representation.UNSIGNED, 32, self._n_size)
        bytes_ += n_size_immediate.to_bytes()

        n_size_pad_immediate = Immediate(Representation.UNSIGNED, 32, self._n_size_pad[dpu_id])
        bytes_ += n_size_pad_immediate.to_bytes()

        nr_rows_immediate = Immediate(Representation.UNSIGNED, 32, self._nr_rows[dpu_id])
        bytes_ += nr_rows_immediate.to_bytes()

        max_rows_immediate = Immediate(Representation.UNSIGNED, 32, self._max_rows[dpu_id])
        bytes_ += max_rows_immediate.to_bytes()

        return Bin(bytes_)

    def dpu_results(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        return None
