from typing import Set, Tuple

import numpy as np
import scipy

from io_.excel_reader import ExcelReader
from regression.datum import Datum


class Model:
    def __init__(self, benchmarks: Set[ExcelReader.Benchmark], cpu_to_dpu_bw: float, dpu_to_cpu_bw: float):
        self._cpu_to_dpu_bw: float = cpu_to_dpu_bw
        self._dpu_to_cpu_bw: float = dpu_to_cpu_bw

        self._data = set()
        for benchmark in benchmarks:
            regression_cpu_to_dpu = (benchmark.calculate_cpu_to_dpu_bytes() * (10**3)) / (
                cpu_to_dpu_bw * (2**30)
            )
            regression_dpu_to_cpu = (benchmark.calculate_dpu_to_cpu_bytes() * (10**3)) / (
                dpu_to_cpu_bw * (2**30)
            )
            regression_dpu_to_dpu = (benchmark.calculate_dpu_to_dpu_from_cpu_to_dpu_bytes() * (10**3)) / (
                cpu_to_dpu_bw * (2**30)
            ) + (benchmark.calculate_dpu_to_dpu_from_dpu_to_cpu_bytes() * (10**3)) / (
                dpu_to_cpu_bw * (2**30)
            )
            self._data.add(
                Datum(
                    benchmark.name(),
                    benchmark.num_dpus(),
                    benchmark.data_prep_param(),
                    benchmark.hw_kerenl(),
                    benchmark.hw_cpu_to_dpu(),
                    benchmark.hw_dpu_to_cpu(),
                    benchmark.hw_dpu_to_dpu(),
                    benchmark.sim_kernel(),
                    regression_cpu_to_dpu,
                    regression_dpu_to_cpu,
                    regression_dpu_to_dpu,
                )
            )

    def data(self) -> Set[Datum]:
        return self._data

    def regress_kernel(self) -> Tuple[float, float]:
        hw_kernels = [datum.hw_kernel() for datum in self._data]
        sim_kernels = [datum.sim_kernel() for datum in self._data]

        slope, intercept, r_value, p_value, std_err = scipy.stats.linregress(hw_kernels, sim_kernels)
        mae = float(np.mean([abs(datum.kernel_err() / datum.hw_kernel()) for datum in self._data]))

        return r_value**2, mae

    def regress_cpu_to_dpu(self):
        hw_cpu_to_dpus = [datum.hw_cpu_to_dpu() for datum in self._data]
        regression_cpu_to_dpus = [datum.regression_cpu_to_dpu() for datum in self._data]

        slope, intercept, r_value, p_value, std_err = scipy.stats.linregress(hw_cpu_to_dpus, regression_cpu_to_dpus)
        mae = float(np.mean([abs(datum.cpu_to_dpu_err() / datum.hw_cpu_to_dpu()) for datum in self._data]))

        return r_value**2, mae

    def regress_dpu_to_cpu(self):
        hw_dpu_to_cpus = [datum.hw_dpu_to_cpu() for datum in self._data]
        regression_dpu_to_cpus = [datum.regression_dpu_to_cpu() for datum in self._data]

        slope, intercept, r_value, p_value, std_err = scipy.stats.linregress(hw_dpu_to_cpus, regression_dpu_to_cpus)
        mae = float(np.mean([abs(datum.dpu_to_cpu_err() / datum.hw_dpu_to_cpu()) for datum in self._data]))

        return r_value**2, mae

    def regress_dpu_to_dpu(self):
        hw_dpu_to_dpus = [datum.hw_dpu_to_dpu() for datum in self._data]
        regression_dpu_to_dpus = [datum.regression_dpu_to_dpu() for datum in self._data]

        slope, intercept, r_value, p_value, std_err = scipy.stats.linregress(hw_dpu_to_dpus, regression_dpu_to_dpus)
        mae = float(
            np.mean(
                [
                    abs(datum.dpu_to_dpu_err() / datum.hw_dpu_to_dpu())
                    for datum in self._data
                    if datum.hw_dpu_to_dpu() != 0
                ]
            )
        )

        return r_value**2, mae

    def regress_communication(self):
        hw_communications = [datum.hw_communication() for datum in self._data]
        regression_communications = [datum.regression_communication() for datum in self._data]

        slope, intercept, r_value, p_value, std_err = scipy.stats.linregress(
            hw_communications, regression_communications
        )
        mae = float(np.mean([abs(datum.communication_err() / datum.hw_communication()) for datum in self._data]))

        return r_value**2, mae

    def regress_total(self):
        hw_totals = [datum.hw_total() for datum in self._data]
        regression_totals = [datum.sim_reg_total() for datum in self._data]

        slope, intercept, r_value, p_value, std_err = scipy.stats.linregress(hw_totals, regression_totals)
        mae = float(np.mean([abs(datum.total_err() / datum.hw_total()) for datum in self._data]))

        return r_value**2, mae
