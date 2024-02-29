from typing import Set

import pandas as pd

from regression.datum import Datum


class ExcelWriter:
    def __init__(self, excel_filepath: str, data: Set[Datum]):
        self._data_frame: pd.DataFrame = pd.DataFrame(
            [
                [
                    datum.name(),
                    datum.num_dpus(),
                    datum.data_prep_param(),
                    datum.hw_total(),
                    datum.sim_reg_total(),
                    datum.sim_kernel(),
                    datum.regression_cpu_to_dpu(),
                    datum.regression_dpu_to_cpu(),
                    datum.regression_dpu_to_dpu(),
                ]
                for datum in data
            ],
            columns=[
                "benchmark",
                "num_dpus",
                "data_prep_param",
                "hw_total",
                "sim_reg_total",
                "sim_kernel",
                "regression_cpu_to_dpu",
                "regression_dpu_to_cpu",
                "regression_dpu_to_dpu",
            ],
        )

        self._data_frame.to_excel(excel_filepath, sheet_name="Sheet1")
