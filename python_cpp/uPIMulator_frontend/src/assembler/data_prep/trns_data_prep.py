from typing import List, Optional

import numpy as np

from abi.word.immediate import Immediate
from abi.word.representation import Representation
from assembler.data_prep.bin import Bin
from encoder.byte import Byte
from initializer.int_initializer import IntInitializer
from util.config_loader import ConfigLoader


class TRNSDataPrep:
    def __init__(self, num_tasklets: int, data_prep_param: List[int], num_dpus: int):
        assert 0 < num_tasklets < ConfigLoader.max_num_tasklets()

        self._num_tasklets: int = num_tasklets
        self._num_dpus: int = num_dpus

        if self._num_dpus > 1:
            self._N_: int = 64  # Should be greater than or equal to "self._num_dpus" for strong scaling to fully utilize all DPUs.
            self._n: int = 8
            self._M_: int = data_prep_param[0]
            self._m: int = 4
        else:
            self._N_: int = 1
            self._n: int = 4
            self._M_: int = data_prep_param[0]
            self._m: int = 16
        
        self._is_strong_scaling = True  # True --> strong scaling / False --> weak scaling
        self._N_ = self._N_ if self._is_strong_scaling else self._N_ * self._num_dpus

        self._num_active_dpus_at_begining: int = self._num_dpus if self._N_ > self._num_dpus else self._N_

        self._num_executions: int = (
            2 * (self._N_ // self._num_active_dpus_at_begining) if self._is_strong_scaling else 2 * self._N_
        )

        self._buffer_a: List[List[List[int]]] = [
            [[IntInitializer.value_by_range(0, 100) for _ in range(self._n)] for _ in range(self._M_ * self._m)]
            for _ in range(self._N_)
        ]

        self._done: List[int] = []
        if (self._M_ * self._n) // 8 == 0:
            for _ in range(8):
                self._done.append(0)
        else:
            for _ in range(self._M_ * self._n):
                self._done.append(0)

        self._buffer_c: List[List[List[int]]] = [[] for _ in range(self._N_)]
        for n in range(self._N_):
            self._buffer_c[n] = np.transpose(self._buffer_a[n]).tolist()

        self._kernel: List[int] = [0, 1]

    def num_executions(self) -> int:
        return self._num_executions

    def num_dpus(self) -> int:
        return self._num_dpus

    def input_dpu_mram_heap_pointer_name(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        if (execution % 2) != 0:
            return None
        else:
            bytes_: List[Byte] = []

            for row in self._buffer_a[(self._num_active_dpus_at_begining * (execution // 2)) + dpu_id]:
                for element in row:
                    element_immediate = Immediate(Representation.UNSIGNED, 64, element)
                    bytes_ += element_immediate.to_bytes()
            for element in self._done:
                element_immediate = Immediate(Representation.UNSIGNED, 8, element)
                bytes_ += element_immediate.to_bytes()

        return Bin(bytes_)

    def output_dpu_mram_heap_pointer_name(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        if (execution % 2) != 1:
            return None
        else:
            bytes_: List[Byte] = []

            for row in self._buffer_c[(self._num_active_dpus_at_begining * (execution // 2)) + dpu_id]:
                for element in row:
                    element_immediate = Immediate(Representation.UNSIGNED, 64, element)
                    bytes_ += element_immediate.to_bytes()

            return Bin(bytes_)

    def dpu_input_arguments(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        m_immediate = Immediate(Representation.UNSIGNED, 32, self._m)
        bytes_ += m_immediate.to_bytes()

        n_immediate = Immediate(Representation.UNSIGNED, 32, self._n)
        bytes_ += n_immediate.to_bytes()

        M_immediate = Immediate(Representation.UNSIGNED, 32, self._M_)
        bytes_ += M_immediate.to_bytes()

        kernel_immediate = Immediate(Representation.UNSIGNED, 32, self._kernel[execution % 2])
        bytes_ += kernel_immediate.to_bytes()

        return Bin(bytes_)

    def dpu_results(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        return None
