from typing import List, Optional

import numpy as np

from abi.word.immediate import Immediate
from abi.word.representation import Representation
from assembler.data_prep.bin import Bin
from encoder.byte import Byte
from util.config_loader import ConfigLoader


class BSDataPrep:
    def __init__(self, num_tasklets: int, data_prep_param: List[int], num_dpus: int):
        assert 0 < num_tasklets < ConfigLoader.max_num_tasklets()

        self._num_tasklets: int = num_tasklets
        self._num_dpus: int = num_dpus

        self._size: int = data_prep_param[0]
        self._num_querys = int(data_prep_param[0] // 8)

        self._num_executions: int = 1

        if self._num_querys % (self._num_dpus * self._num_tasklets) != 0:
            self._num_querys = self._num_querys + (
                self._num_dpus * self._num_tasklets - self._num_querys % (self._num_dpus * self._num_tasklets)
            )

        assert self._num_querys % (self._num_dpus * self._num_tasklets) == 0

        self._input_buffer: List[int] = [i + 1 for i in range(self._size)]
        self._query_buffer: List[int] = [i for i in range(self._num_querys)]
        self._result: List[List[int]] = [[0 for _ in range(self._num_tasklets)] for _ in range(self._num_dpus)]

        self._slice_per_dpu = self._num_querys // self._num_dpus
        self._query_per_tasklet = self._slice_per_dpu // self._num_tasklets

        for dpu_id in range(self._num_dpus):
            for tasklet in range(self._num_tasklets):
                for query in range(self._query_per_tasklet):
                    is_found = False
                    l = 0
                    r = self._size - 1
                    while l <= r:
                        m = l + (r - l) // 2
                        if (
                            self._input_buffer[m]
                            == self._query_buffer[
                                query + tasklet * self._query_per_tasklet + dpu_id * self._slice_per_dpu
                            ]
                        ):
                            self._result[dpu_id][tasklet] = m
                            is_found = True
                            break
                        if (
                            self._input_buffer[m]
                            < self._query_buffer[
                                query + tasklet * self._query_per_tasklet + dpu_id * self._slice_per_dpu
                            ]
                        ):
                            l = m + 1
                        else:
                            r = m - 1
                    if is_found == False:
                        self._result[dpu_id][tasklet] = -1

        self._kernel: int = 0

    def num_executions(self) -> int:
        return self._num_executions

    def num_dpus(self) -> int:
        return self._num_dpus

    def input_dpu_mram_heap_pointer_name(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        start_elem = self._slice_per_dpu * dpu_id
        end_elem = self._slice_per_dpu * (dpu_id + 1)

        for element in self._input_buffer:
            element_immediate = Immediate(Representation.UNSIGNED, 64, element)
            bytes_ += element_immediate.to_bytes()

        for element in self._query_buffer[start_elem:end_elem]:
            element_immediate = Immediate(Representation.UNSIGNED, 64, element)
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

        size_immediate = Immediate(Representation.UNSIGNED, 64, self._size)
        bytes_ += size_immediate.to_bytes()

        slice_per_dpu_immediate = Immediate(Representation.UNSIGNED, 64, self._slice_per_dpu)
        bytes_ += slice_per_dpu_immediate.to_bytes()

        kernel_immediate = Immediate(Representation.UNSIGNED, 32, self._kernel)
        bytes_ += kernel_immediate.to_bytes()

        return Bin(bytes_)

    def dpu_results(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []
        for element in self._result[dpu_id]:
            element_immediate = Immediate(Representation.UNSIGNED, 64, element)
            bytes_ += element_immediate.to_bytes()
        return Bin(bytes_)
