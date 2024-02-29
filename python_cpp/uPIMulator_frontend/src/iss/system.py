import logging
import os

from iss.cpu.cpu import CPU
from iss.dpu.dpu import DPU
from util.path_collector import PathCollector


class System:
    def __init__(self, benchmark: str, num_tasklets: int):
        self._benchmark: str = benchmark
        self._num_tasklets: int = num_tasklets
        self._logger: logging.Logger = self._init_logger()

        self._dpu = DPU(num_tasklets)
        self._cpu = CPU(benchmark, num_tasklets, self._dpu)

        self._execution: int = 0

    def init(self) -> None:
        self._cpu.init()
        self._cpu.sched(self._execution)
        self._cpu.launch()

    def fini(self) -> None:
        self._cpu.fini()

    def is_finished(self) -> bool:
        return self._execution == self._cpu.num_executions()

    def is_zombie(self) -> bool:
        return self._dpu.is_zombie()

    def cycle(self):
        self._dpu.cycle()
        self._cpu.cycle()

        if self.is_zombie():
            self._cpu.check(self._execution)
            self._execution += 1

            if not self.is_finished():
                self._cpu.sched(self._execution)
                self._cpu.launch()

    def _init_logger(self) -> logging.Logger:
        logger = logging.getLogger("iss")
        logger.setLevel(logging.DEBUG)

        formatter = logging.Formatter("%(message)s")

        trace_dirpath = PathCollector.trace_path_in_local()
        if not os.path.exists(trace_dirpath):
            os.makedirs(trace_dirpath)

        log_filepath = os.path.join(trace_dirpath, f"{self._benchmark}.{self._num_tasklets}.trace",)
        if os.path.exists(log_filepath):
            os.remove(log_filepath)

        stream_handler = logging.StreamHandler()
        file_handler = logging.FileHandler(filename=log_filepath)

        stream_handler.setLevel(logging.INFO)
        file_handler.setLevel(logging.DEBUG)

        stream_handler.setFormatter(formatter)
        file_handler.setFormatter(formatter)

        logger.addHandler(stream_handler)
        logger.addHandler(file_handler)

        return logger
