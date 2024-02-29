from typing import List, Optional

import numpy as np

from abi.word.immediate import Immediate
from abi.word.representation import Representation
from assembler.data_prep.bin import Bin
from encoder.byte import Byte
from util.config_loader import ConfigLoader


class TSDataPrep:
    def __init__(self, num_tasklets: int, data_prep_param: List[int], num_dpus: int):
        assert 0 < num_tasklets < ConfigLoader.max_num_tasklets()

        self._num_tasklets: int = num_tasklets
        self._num_dpus: int = num_dpus

        self._ts_size: int = data_prep_param[0]
        self._query_length: int = 64  

        if self._ts_size % (self._num_dpus * self._num_tasklets * self._query_length):
            self._ts_size = self._ts_size + (
                (self._num_dpus * self._num_tasklets * self._query_length)
                - (self._ts_size % (self._num_dpus * self._num_tasklets * self._query_length))
            )

        self._num_executions: int = 1
        self._t_series_buffer: List[int] = [i % 127 for i in range(self._ts_size)]
        self._query_buffer: List[int] = [i % 127 for i in range(self._query_length)]

        self._asigma_buffer: List[int] = [0 for _ in range(self._ts_size)]
        self._amean_buffer: List[int] = [0 for _ in range(self._ts_size)]

        acum_sum_buffer: List[int] = []
        for i in range(self._ts_size):
            if i == 0:
                acum_sum_buffer.append(self._t_series_buffer[0])
            else:
                acum_sum_buffer.append(self._t_series_buffer[i] + acum_sum_buffer[i - 1])

        asqcum_sum_buffer: List[int] = []
        for i in range(self._ts_size):
            if i == 0:
                asqcum_sum_buffer.append(self._t_series_buffer[0] * self._t_series_buffer[0])
            else:
                asqcum_sum_buffer.append(self._t_series_buffer[i] * self._t_series_buffer[i] + asqcum_sum_buffer[i - 1])

        asum_buffer: List[int] = []
        for i in range(self._ts_size - self._query_length + 1):
            if i == 0:
                asum_buffer.append(acum_sum_buffer[self._query_length - 1])
            else:
                asum_buffer.append(acum_sum_buffer[self._query_length + i - 1] - acum_sum_buffer[i - 1])

        asum_sq_buffer: List[int] = []
        for i in range(self._ts_size - self._query_length + 1):
            if i == 0:
                asum_sq_buffer.append(asqcum_sum_buffer[self._query_length - 1])
            else:
                asum_sq_buffer.append(asqcum_sum_buffer[self._query_length + i - 1] - asqcum_sum_buffer[i - 1])

        amean_temp_buffer: List[int] = [
            asum_buffer[i] // self._query_length for i in range(self._ts_size - self._query_length)
        ]
        asigma_sq_buffer: List[int] = [
            asum_sq_buffer[i] // self._query_length - self._amean_buffer[i] * self._amean_buffer[i]
            for i in range(self._ts_size - self._query_length)
        ]

        self._asigma_buffer = [int(np.sqrt(asigma_sq_buffer[i])) for i in range(self._ts_size - self._query_length)]
        self._amean_buffer = [amean_temp_buffer[i] for i in range(self._ts_size - self._query_length)]

        for i in range(self._query_length):
            self._t_series_buffer.append(0)
        for i in range(self._query_length * 2):
            self._asigma_buffer.append(0)
            self._amean_buffer.append(0)

        assert len(self._t_series_buffer) == len(self._amean_buffer)

        queryMean = 0
        for i in range(self._query_length):
            queryMean += self._query_buffer[i]

        queryMean = queryMean / self._query_length
        self._query_mean = int(queryMean)

        queryVariance = 0
        for i in range(self._query_length):
            queryVariance += (self._query_buffer[i] - queryMean) * (self._query_buffer[i] - queryMean)

        queryVariance = queryVariance / self._query_length
        queryStdDev = np.sqrt(queryVariance)
        self._query_std = int(queryStdDev)

        self._slice_per_dpu = self._ts_size // self._num_dpus

        self._min_val: List[List[int]] = [
            [0x7FFFFFFF for _ in range(self._num_tasklets)] for _ in range(self._num_dpus)
        ]
        self._min_idx: List[List[int]] = [[0 for _ in range(self._num_tasklets)] for _ in range(self._num_dpus)]
        self._max_val: List[List[int]] = [[0 for _ in range(self._num_tasklets)] for _ in range(self._num_dpus)]
        self._max_idx: List[List[int]] = [[0 for _ in range(self._num_tasklets)] for _ in range(self._num_dpus)]

        my_start_elem: List[List[int]] = [
            [
                self._slice_per_dpu * dpu_id + i * (self._slice_per_dpu // self._num_tasklets)
                for i in range(self._num_tasklets)
            ]
            for dpu_id in range(self._num_dpus)
        ]
        my_end_elem: List[List[int]] = [
            [
                my_start_elem[dpu_id][i] + (self._slice_per_dpu // self._num_tasklets) - 1
                for i in range(self._num_tasklets)
            ]
            for dpu_id in range(self._num_dpus)
        ]

        for dpu_id in range(self._num_dpus):
            for i in range(self._num_tasklets):
                if my_end_elem[dpu_id][i] > self._slice_per_dpu * (dpu_id + 1) - self._query_length:
                    my_end_elem[dpu_id][i] = self._slice_per_dpu * (dpu_id + 1) - self._query_length

        self._block_size: int = 256
        self._elem_size: int = 4
        increment: int = self._block_size // self._elem_size
        self._dotpip: int = self._block_size // self._elem_size
        iter: int = 0
        for dpu_id in range(self._num_dpus):
            iter = 0
            for tasklet in range(self._num_tasklets):
                for i in range(my_start_elem[dpu_id][tasklet], my_end_elem[dpu_id][tasklet], increment):
                    self._cache_dotprods: List[int] = [0 for d in range(self._dotpip)]

                    for j in range(self._query_length // (increment)):
                        self.dot_product(
                            self._t_series_buffer[i : i + increment],
                            self._t_series_buffer[i + increment : (i + 2 * increment)],
                            self._query_buffer[j * increment : (j + 1) * increment],
                            self._cache_dotprods,
                        )

                    for k in range(increment):
                        distance = 2 * (
                            self._query_length
                            - (
                                self._cache_dotprods[k]
                                - self._query_length
                                * self._amean_buffer[k + iter * increment + dpu_id * self._slice_per_dpu]
                                * self._query_mean
                            )
                            // (
                                self._asigma_buffer[k + iter * increment + dpu_id * self._slice_per_dpu]
                                * self._query_std
                            )
                        )

                        if distance < self._min_val[dpu_id][tasklet]:
                            self._min_val[dpu_id][tasklet] = distance
                            self._min_idx[dpu_id][tasklet] = i + k - (dpu_id * self._slice_per_dpu)
                    iter += 1

        self._exclusion_zone: int = 0
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
        end_elem = self._slice_per_dpu * (dpu_id + 1) + self._query_length

        for element in self._query_buffer:
            element_immediate = Immediate(Representation.UNSIGNED, 32, element)
            bytes_ += element_immediate.to_bytes()

        for element in self._t_series_buffer[start_elem:end_elem]:
            element_immediate = Immediate(Representation.UNSIGNED, 32, element)
            bytes_ += element_immediate.to_bytes()

        for element in self._amean_buffer[start_elem:end_elem]:
            element_immediate = Immediate(Representation.UNSIGNED, 32, element)
            bytes_ += element_immediate.to_bytes()

        for element in self._asigma_buffer[start_elem:end_elem]:
            element_immediate = Immediate(Representation.UNSIGNED, 32, element)
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

        ts_size_immediate = Immediate(Representation.UNSIGNED, 32, self._ts_size)
        bytes_ += ts_size_immediate.to_bytes()

        query_length_immediate = Immediate(Representation.UNSIGNED, 32, self._query_length)
        bytes_ += query_length_immediate.to_bytes()

        query_mean_immediate = Immediate(Representation.UNSIGNED, 32, self._query_mean)
        bytes_ += query_mean_immediate.to_bytes()

        query_std_immediate = Immediate(Representation.UNSIGNED, 32, self._query_std)
        bytes_ += query_std_immediate.to_bytes()

        slice_per_dpu_immediate = Immediate(Representation.UNSIGNED, 32, self._slice_per_dpu)
        bytes_ += slice_per_dpu_immediate.to_bytes()

        exclusion_zone_immediate = Immediate(Representation.UNSIGNED, 32, self._exclusion_zone)
        bytes_ += exclusion_zone_immediate.to_bytes()

        kernel_immediate = Immediate(Representation.UNSIGNED, 32, self._kernel)
        bytes_ += kernel_immediate.to_bytes()

        return Bin(bytes_)

    def dpu_results(self, execution: int, dpu_id: int) -> Optional[Bin]:
        assert 0 <= execution < self._num_executions
        assert 0 <= dpu_id < self._num_dpus

        bytes_: List[Byte] = []

        for tasklet in range(self._num_tasklets):
            min_val_immediate = Immediate(Representation.UNSIGNED, 32, self._min_val[dpu_id][tasklet])
            bytes_ += min_val_immediate.to_bytes()

            min_idx_immediate = Immediate(Representation.UNSIGNED, 32, self._min_idx[dpu_id][tasklet])
            bytes_ += min_idx_immediate.to_bytes()

            max_val_immediate = Immediate(Representation.UNSIGNED, 32, self._max_val[dpu_id][tasklet])
            bytes_ += max_val_immediate.to_bytes()

            max_idx_immediate = Immediate(Representation.UNSIGNED, 32, self._max_idx[dpu_id][tasklet])
            bytes_ += max_idx_immediate.to_bytes()

        return Bin(bytes_)

    def dot_product(
        self,
        vector_a: List[int],
        vector_a_aux: List[int],
        vector_query: List[int],
        vector_result: List[int],
    ):
        for i in range(self._block_size // self._elem_size):
            for j in range(self._dotpip):
                if (j + i) > (self._block_size // self._elem_size) - 1:
                    vector_result[j] += vector_a_aux[(j + i) - self._block_size // self._elem_size] * vector_query[i]
                else:
                    vector_result[j] += vector_a[j + i] * vector_query[i]
