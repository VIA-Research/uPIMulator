class _BaseBenchmark:
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
    ):
        self._name: str = name
        self._num_dpus: int = num_dpus
        self._data_prep_param: int = data_prep_param
        self._hw_kernel: float = hw_kernel
        self._hw_cpu_to_dpu: float = hw_cpu_to_dpu
        self._hw_dpu_to_cpu: float = hw_dpu_to_cpu
        self._hw_dpu_to_dpu: float = hw_dpu_to_dpu
        self._sim_kernel: float = sim_kernel

    def name(self) -> str:
        return self._name

    def num_dpus(self) -> int:
        return self._num_dpus

    def data_prep_param(self) -> int:
        return self._data_prep_param

    def hw_kerenl(self) -> float:
        return self._hw_kernel

    def hw_cpu_to_dpu(self) -> float:
        return self._hw_cpu_to_dpu

    def hw_dpu_to_cpu(self) -> float:
        return self._hw_dpu_to_cpu

    def hw_dpu_to_dpu(self) -> float:
        return self._hw_dpu_to_dpu

    def sim_kernel(self) -> float:
        return self._sim_kernel

    def calculate_cpu_to_dpu_bytes(self) -> int:
        raise AttributeError

    def calculate_dpu_to_cpu_bytes(self) -> int:
        raise AttributeError

    def calculate_dpu_to_dpu_from_cpu_to_dpu_bytes(self) -> int:
        raise AttributeError

    def calculate_dpu_to_dpu_from_dpu_to_cpu_bytes(self) -> int:
        raise AttributeError
