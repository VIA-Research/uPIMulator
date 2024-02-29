from typing import List, Optional

import numpy as np

from abi.word.immediate import Immediate
from abi.word.representation import Representation
from assembler.data_prep.bin import Bin
from encoder.byte import Byte
from util.config_loader import ConfigLoader


class MLPDataPrep:
    def __init__(self, num_tasklets: int, data_prep_param: List[int], num_dpus: int):
        assert 0 < num_tasklets < ConfigLoader.max_num_tasklets()

        self._num_tasklets: int = num_tasklets
        self._num_dpus: int = num_dpus

        self._m_size: int = data_prep_param[0]
        self._n_size: int = data_prep_param[0]

        self._num_layers: int = 3

        self._num_executions: int = self._num_layers

        assert self._m_size % self._num_dpus == 0

        self._n_size_pad: List[int] = [
            (self._n_size if self._n_size % 2 == 0 else self._n_size + 1) for _ in range(self._num_dpus)
        ]
        self._nr_rows: List[int] = [self._m_size // self._num_dpus for _ in range(self._num_dpus)]
        self._max_rows: List[int] = [
            (self._nr_rows[i] if self._nr_rows[i] % 2 == 0 else self._nr_rows[i] + 1) for i in range(self._num_dpus)
        ]

        # weights
        self._buffer_a: List[List[List[int]]] = [
            [[0 if i % 100 < 98 else (layer + i) % 2 for i in range(self._n_size)] for _ in range(self._m_size)]
            for layer in range(self._num_layers)
        ]

        # input activations
        self._buffer_b: List[List[int]] = [
            [0 if i % 50 < 48 else i % 2 for i in range(self._n_size)] for _ in range(self._num_layers)
        ]

        # output activations
        self._buffer_c: List[List[List[int]]] = [
            [[0 for _ in range(self._m_size)] for _ in range(self._num_dpus)] for _ in range(self._num_layers)
        ]
        for layer in range(self._num_layers):
            for dpu_id in range(self._num_dpus):
                start_row: int = self._nr_rows[dpu_id] * dpu_id
                end_row: int = self._nr_rows[dpu_id] * (dpu_id + 1)

                self._buffer_c[layer][dpu_id] = np.matmul(
                    self._buffer_a[layer][start_row:end_row], self._buffer_b[layer]
                ).tolist()
                for i in range(len(self._buffer_c[layer][dpu_id])):
                    if self._buffer_c[layer][dpu_id][i] < 0:
                        self._buffer_c[layer][dpu_id][i] = 0

                if layer < self._num_layers - 1:
                    self._buffer_b[layer + 1][start_row:end_row] = self._buffer_c[layer][dpu_id]

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

        for row in self._buffer_a[execution][start_row:end_row]:
            for element in row:
                element_immediate = Immediate(Representation.UNSIGNED, 32, element)
                bytes_ += element_immediate.to_bytes()

        for element in self._buffer_b[execution]:
            element_immediate = Immediate(Representation.UNSIGNED, 32, element)
            bytes_ += element_immediate.to_bytes()

        return Bin(bytes_)

    def output_dpu_mram_heap_pointer_name(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        start_row: int = self._nr_rows[dpu_id] * dpu_id
        end_row: int = self._nr_rows[dpu_id] * (dpu_id + 1)

        for row in self._buffer_a[execution][start_row:end_row]:
            for element in row:
                element_immediate = Immediate(Representation.UNSIGNED, 32, element)
                bytes_ += element_immediate.to_bytes()

        for element in self._buffer_b[execution]:
            element_immediate = Immediate(Representation.UNSIGNED, 32, element)
            bytes_ += element_immediate.to_bytes()

        for element in self._buffer_c[execution][dpu_id]:
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
