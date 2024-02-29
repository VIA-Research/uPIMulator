import math
from typing import List, Optional

from abi.word.immediate import Immediate
from abi.word.representation import Representation
from assembler.data_prep.bin import Bin
from encoder.byte import Byte
from util.config_loader import ConfigLoader


class Result:
    def __init__(self, t_count=0, first=0, last=0):
        self._t_count: int = t_count
        self._first: int = first
        self._last: int = last


class UNIDataPrep:
    def __init__(self, num_tasklets: int, data_prep_param: List[int], num_dpus: int):
        assert 0 < num_tasklets < ConfigLoader.max_num_tasklets()

        self._num_tasklets: int = num_tasklets
        self._num_dpus: int = num_dpus

        self._size: int = data_prep_param[0]

        self._num_executions: int = 1

        elem_size = 8

        regs = 128

        self._is_strong_scaling = True  # True --> strong scaling / False --> weak scaling
        input_size = self._size if self._is_strong_scaling else self._size * self._num_dpus

        input_size_dpu = ((input_size) - 1) // (self._num_dpus) + 1

        if input_size_dpu % (num_tasklets * regs) != 0:
            self._input_size_dpu_round = math.ceil(input_size_dpu / (num_tasklets * regs)) * (num_tasklets * regs)
        else:
            self._input_size_dpu_round = input_size_dpu

        self._buffer_a: List[int] = [0 for _ in range(self._input_size_dpu_round * self._num_dpus)]
        for i in range(self._input_size_dpu_round * self._num_dpus):
            if i < input_size:
                if i % 2 == 0:
                    self._buffer_a[i] = i
                else:
                    self._buffer_a[i] = i + 1
            else:
                self._buffer_a[i] = self._buffer_a[input_size - 1]

        self._buffer_c: List[List[int]] = [
            [0 for _ in range(self._input_size_dpu_round)] for _ in range(self._num_dpus)
        ]

        for dpu_id in range(self._num_dpus):
            start_elem = self._input_size_dpu_round * dpu_id

            self._buffer_c[dpu_id][0] = self._buffer_a[start_elem]

            self._pos = 1

            for i in range(1, self._input_size_dpu_round):
                if self._buffer_a[start_elem + i] != self._buffer_a[start_elem + i - 1]:
                    self._buffer_c[dpu_id][self._pos] = self._buffer_a[start_elem + i]
                    self._pos += 1

        self._input_size_dpu: List[int] = [self._input_size_dpu_round * elem_size for _ in range(self._num_dpus)]
        self._kernel: List[int] = [0 for _ in range(self._num_dpus)]

        self._result: List[List[Result]] = [[] for _ in range(self._num_dpus)]
        for dpu_id in range(self._num_dpus):
            start_elem = self._input_size_dpu_round * dpu_id
            for tasklet_id in range(self._num_tasklets):
                if tasklet_id == 0 and tasklet_id != self._num_tasklets - 1:
                    self._result[dpu_id].append(Result(0, self._buffer_a[start_elem], 0))
                elif tasklet_id == self._num_tasklets - 1:
                    self._result[dpu_id].append(Result(self._pos, 0, self._buffer_c[dpu_id][self._pos - 1]))
                else:
                    self._result[dpu_id].append(Result())

    def num_executions(self) -> int:
        return self._num_executions

    def num_dpus(self) -> int:
        return self._num_dpus

    def input_dpu_mram_heap_pointer_name(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        start_elem = self._input_size_dpu_round * dpu_id

        for i in range(self._input_size_dpu_round):
            element_immediate = Immediate(Representation.UNSIGNED, 64, self._buffer_a[start_elem + i])
            bytes_ += element_immediate.to_bytes()
        return Bin(bytes_)

    def output_dpu_mram_heap_pointer_name(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        start_elem = self._input_size_dpu_round * dpu_id

        for i in range(self._input_size_dpu_round):
            element_immediate = Immediate(Representation.UNSIGNED, 64, self._buffer_a[start_elem + i])
            bytes_ += element_immediate.to_bytes()

        for element in self._buffer_c[dpu_id]:
            element_immediate = Immediate(Representation.UNSIGNED, 64, element)
            bytes_ += element_immediate.to_bytes()

        return Bin(bytes_)

    def dpu_input_arguments(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        input_size_dpu_immediate = Immediate(Representation.UNSIGNED, 32, self._input_size_dpu[dpu_id])
        bytes_ += input_size_dpu_immediate.to_bytes()

        kernel_immediate = Immediate(Representation.UNSIGNED, 32, self._kernel[dpu_id])
        bytes_ += kernel_immediate.to_bytes()

        return Bin(bytes_)

    def dpu_results(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        for element in self._result[dpu_id]:
            t_count_immediate = Immediate(Representation.UNSIGNED, 64, element._t_count)
            bytes_ += t_count_immediate.to_bytes()

            first_immediate = Immediate(Representation.UNSIGNED, 64, element._first)
            bytes_ += first_immediate.to_bytes()

            last_immediate = Immediate(Representation.UNSIGNED, 64, element._last)
            bytes_ += last_immediate.to_bytes()

        return Bin(bytes_)
