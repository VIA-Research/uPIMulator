import math
from typing import Set, Union

import pandas as pd

from benchmark.bs import BS
from benchmark.gemv import GEMV
from benchmark.hst_l import HST_L
from benchmark.hst_s import HST_S
from benchmark.mlp import MLP
from benchmark.red import RED
from benchmark.scan_rss import SCAN_RSS
from benchmark.scan_ssa import SCAN_SSA
from benchmark.sel import SEL
from benchmark.trns import TRNS
from benchmark.ts import TS
from benchmark.uni import UNI
from benchmark.va import VA


class ExcelReader:
    Benchmark = Union[BS, GEMV, HST_L, HST_S, MLP, RED, SCAN_RSS, SCAN_SSA, SEL, TRNS, TS, UNI, VA]

    def __init__(self, excel_filepath: str):
        self._excel_filepath: str = excel_filepath

        self._data_frame: pd.DataFrame = pd.DataFrame(pd.read_excel(excel_filepath))

        self._benchmarks: Set[ExcelReader.Benchmark] = set()
        for (
            benchmark,
            num_dpus,
            data_prep_param,
            hw_kernel,
            hw_cpu_to_dpu,
            hw_dpu_to_cpu,
            hw_dpu_to_dpu,
            sim_kernel,
        ) in zip(
            self._data_frame["benchmark"],
            self._data_frame["num_dpus"],
            self._data_frame["data_prep_param"],
            self._data_frame["hw_kernel"],
            self._data_frame["hw_cpu_to_dpu"],
            self._data_frame["hw_dpu_to_cpu"],
            self._data_frame["hw_dpu_to_dpu"],
            self._data_frame["sim_kernel"],
        ):
            if (
                not (
                    math.isnan(hw_kernel)
                    or math.isnan(hw_cpu_to_dpu)
                    or math.isnan(hw_dpu_to_cpu)
                    or math.isnan(hw_dpu_to_dpu)
                    or math.isnan(sim_kernel)
                )
                # and hw_kernel + hw_cpu_to_dpu + hw_dpu_to_cpu + hw_dpu_to_dpu <= 500
            ):
                if benchmark == "BS":
                    self._benchmarks.add(
                        BS(
                            num_dpus,
                            data_prep_param,
                            hw_kernel,
                            hw_cpu_to_dpu,
                            hw_dpu_to_cpu,
                            hw_dpu_to_dpu,
                            sim_kernel,
                        )
                    )
                elif benchmark == "GEMV":
                    self._benchmarks.add(
                        GEMV(
                            num_dpus,
                            data_prep_param,
                            hw_kernel,
                            hw_cpu_to_dpu,
                            hw_dpu_to_cpu,
                            hw_dpu_to_dpu,
                            sim_kernel,
                        )
                    )
                elif benchmark == "HST-L":
                    self._benchmarks.add(
                        HST_L(
                            num_dpus,
                            data_prep_param,
                            hw_kernel,
                            hw_cpu_to_dpu,
                            hw_dpu_to_cpu,
                            hw_dpu_to_dpu,
                            sim_kernel,
                        )
                    )
                elif benchmark == "HST-S":
                    self._benchmarks.add(
                        HST_S(
                            num_dpus,
                            data_prep_param,
                            hw_kernel,
                            hw_cpu_to_dpu,
                            hw_dpu_to_cpu,
                            hw_dpu_to_dpu,
                            sim_kernel,
                        )
                    )
                elif benchmark == "MLP":
                    self._benchmarks.add(
                        MLP(
                            num_dpus,
                            data_prep_param,
                            hw_kernel,
                            hw_cpu_to_dpu,
                            hw_dpu_to_cpu,
                            hw_dpu_to_dpu,
                            sim_kernel,
                        )
                    )
                elif benchmark == "RED":
                    self._benchmarks.add(
                        RED(
                            num_dpus,
                            data_prep_param,
                            hw_kernel,
                            hw_cpu_to_dpu,
                            hw_dpu_to_cpu,
                            hw_dpu_to_dpu,
                            sim_kernel,
                        )
                    )
                elif benchmark == "SCAN-RSS":
                    self._benchmarks.add(
                        SCAN_RSS(
                            num_dpus,
                            data_prep_param,
                            hw_kernel,
                            hw_cpu_to_dpu,
                            hw_dpu_to_cpu,
                            hw_dpu_to_dpu,
                            sim_kernel,
                        )
                    )
                elif benchmark == "SCAN-SSA":
                    self._benchmarks.add(
                        SCAN_SSA(
                            num_dpus,
                            data_prep_param,
                            hw_kernel,
                            hw_cpu_to_dpu,
                            hw_dpu_to_cpu,
                            hw_dpu_to_dpu,
                            sim_kernel,
                        )
                    )
                elif benchmark == "SEL":
                    self._benchmarks.add(
                        SEL(
                            num_dpus,
                            data_prep_param,
                            hw_kernel,
                            hw_cpu_to_dpu,
                            hw_dpu_to_cpu,
                            hw_dpu_to_dpu,
                            sim_kernel,
                        )
                    )
                elif benchmark == "TRNS":
                    self._benchmarks.add(
                        TRNS(
                            num_dpus,
                            data_prep_param,
                            hw_kernel,
                            hw_cpu_to_dpu,
                            hw_dpu_to_cpu,
                            hw_dpu_to_dpu,
                            sim_kernel,
                        )
                    )
                elif benchmark == "TS":
                    self._benchmarks.add(
                        TS(
                            num_dpus,
                            data_prep_param,
                            hw_kernel,
                            hw_cpu_to_dpu,
                            hw_dpu_to_cpu,
                            hw_dpu_to_dpu,
                            sim_kernel,
                        )
                    )
                elif benchmark == "UNI":
                    self._benchmarks.add(
                        UNI(
                            num_dpus,
                            data_prep_param,
                            hw_kernel,
                            hw_cpu_to_dpu,
                            hw_dpu_to_cpu,
                            hw_dpu_to_dpu,
                            sim_kernel,
                        )
                    )
                elif benchmark == "VA":
                    self._benchmarks.add(
                        VA(
                            num_dpus,
                            data_prep_param,
                            hw_kernel,
                            hw_cpu_to_dpu,
                            hw_dpu_to_cpu,
                            hw_dpu_to_dpu,
                            sim_kernel,
                        )
                    )
                else:
                    raise ValueError

    def benchmarks(self) -> Set[Union[Benchmark]]:
        return self._benchmarks
