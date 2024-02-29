from benchmark._base_benchmark import _BaseBenchmark


class SCAN_SSA(_BaseBenchmark):
    def __init__(
        self,
        num_dpus: int,
        data_prep_param: int,
        hw_kernel: float,
        hw_cpu_to_dpu: float,
        hw_dpu_to_cpu: float,
        hw_dpu_to_dpu: float,
        sim_kernel: float,
    ):
        super().__init__(
            "SCAN-SSA", num_dpus, data_prep_param, hw_kernel, hw_cpu_to_dpu, hw_dpu_to_cpu, hw_dpu_to_dpu, sim_kernel
        )

    def calculate_cpu_to_dpu_bytes(self) -> int:
        input_argument_size_1 = 16
        input_argument_size_2 = 8

        vector_length = self._data_prep_param
        elem_size = 8

        mram_transfer_size = (vector_length // self._num_dpus) * elem_size

        return input_argument_size_1 + mram_transfer_size + input_argument_size_2

    def calculate_dpu_to_cpu_bytes(self) -> int:
        vector_length = self._data_prep_param
        elem_size = 8

        mram_transfer_size = (vector_length // self._num_dpus) * elem_size

        return mram_transfer_size

    def calculate_dpu_to_dpu_from_cpu_to_dpu_bytes(self) -> int:
        input_argument_size_2 = 8
        return input_argument_size_2

    def calculate_dpu_to_dpu_from_dpu_to_cpu_bytes(self) -> int:
        results_size_1 = 128
        return results_size_1
