class Datum:
    def __init__(
        self,
        name: str,
        num_dpus: int,
        data_prep_param: int,
        hw_kernel: float,
        hw_cpu_to_dpu: float,
        hw_dpu_to_cpu: float,
        hw_dpu_to_dpu: float,
        sim_kernel: float,
        regression_cpu_to_dpu: float,
        regression_dpu_to_cpu: float,
        regression_dpu_to_dpu: float,
    ):
        self._name: str = name
        self._num_dpus: int = num_dpus
        self._data_prep_param: int = data_prep_param
        self._hw_kernel: float = hw_kernel
        self._hw_cpu_to_dpu: float = hw_cpu_to_dpu
        self._hw_dpu_to_cpu: float = hw_dpu_to_cpu
        self._hw_dpu_to_dpu: float = hw_dpu_to_dpu
        self._sim_kernel: float = sim_kernel
        self._regression_cpu_to_dpu: float = regression_cpu_to_dpu
        self._regression_dpu_to_cpu: float = regression_dpu_to_cpu
        self._regression_dpu_to_dpu: float = regression_dpu_to_dpu

    def name(self) -> str:
        return self._name

    def num_dpus(self) -> int:
        return self._num_dpus

    def data_prep_param(self) -> int:
        return self._data_prep_param

    def hw_kernel(self) -> float:
        return self._hw_kernel

    def hw_cpu_to_dpu(self) -> float:
        return self._hw_cpu_to_dpu

    def hw_dpu_to_cpu(self) -> float:
        return self._hw_dpu_to_cpu

    def hw_dpu_to_dpu(self) -> float:
        return self._hw_dpu_to_dpu

    def hw_communication(self) -> float:
        return self._hw_cpu_to_dpu + self._hw_dpu_to_cpu + self._hw_dpu_to_dpu

    def hw_total(self) -> float:
        return self.hw_kernel() + self.hw_communication()

    def sim_kernel(self) -> float:
        return self._sim_kernel

    def regression_cpu_to_dpu(self) -> float:
        return self._regression_cpu_to_dpu

    def regression_dpu_to_cpu(self) -> float:
        return self._regression_dpu_to_cpu

    def regression_dpu_to_dpu(self) -> float:
        return self._regression_dpu_to_dpu

    def regression_communication(self) -> float:
        return self._regression_cpu_to_dpu + self._regression_dpu_to_cpu + self._regression_dpu_to_dpu

    def sim_reg_total(self) -> float:
        return self.sim_kernel() + self.regression_communication()

    def kernel_err(self) -> float:
        return self._sim_kernel - self._hw_kernel

    def cpu_to_dpu_err(self) -> float:
        return self._regression_cpu_to_dpu - self._hw_cpu_to_dpu

    def dpu_to_cpu_err(self) -> float:
        return self._regression_dpu_to_cpu - self._hw_dpu_to_cpu

    def dpu_to_dpu_err(self) -> float:
        return self._regression_dpu_to_dpu - self._hw_dpu_to_dpu

    def communication_err(self) -> float:
        return self.regression_communication() - self.hw_communication()

    def total_err(self) -> float:
        return self.kernel_err() + self.communication_err()
