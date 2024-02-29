from benchmark._base_benchmark import _BaseBenchmark


class UNI(_BaseBenchmark):
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
            "UNI", num_dpus, data_prep_param, hw_kernel, hw_cpu_to_dpu, hw_dpu_to_cpu, hw_dpu_to_dpu, sim_kernel
        )

    def calculate_cpu_to_dpu_bytes(self) -> int:
        input_argument_size = 8

        vector_length = self._data_prep_param
        elem_size = 8

        mram_transfer_size = (vector_length // self._num_dpus) * elem_size

        return input_argument_size + mram_transfer_size

    def calculate_dpu_to_cpu_bytes(self) -> int:
        results_size = 384

        vector_length = self._data_prep_param
        elem_size = 8

        mram_transfer_size = round((vector_length // 2 + 1) / self._num_dpus) * elem_size

        return results_size + mram_transfer_size

    def calculate_dpu_to_dpu_from_cpu_to_dpu_bytes(self) -> int:
        return 0

    def calculate_dpu_to_dpu_from_dpu_to_cpu_bytes(self) -> int:
        return 0
