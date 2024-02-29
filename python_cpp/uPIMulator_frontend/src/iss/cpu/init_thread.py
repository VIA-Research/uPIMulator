import os
from typing import List, Tuple

from encoder.byte import Byte
from iss.dpu.dpu import DPU
from util.config_loader import ConfigLoader
from util.path_collector import PathCollector


class InitThread:
    def __init__(self, benchmark: str, num_tasklets: int, dpu: DPU):
        self._benchmark: str = benchmark
        self._num_tasklets: int = num_tasklets
        self._dpu: DPU = dpu

    def init(self) -> None:
        self._dma_transfer_atomic()
        self._dma_transfer_iram()
        self._dma_transfer_wram()
        self._dma_transfer_mram()

    def launch(self) -> None:
        for thread in self._dpu.threads():
            bootstrap = ConfigLoader.iram_offset()
            thread.register_file().write_pc(bootstrap)

        self._dpu.boot()

    def cycle(self) -> None:
        pass

    def _dma_transfer_atomic(self) -> None:
        atomic_bin_filepath = os.path.join(
            PathCollector.bin_path_in_local(), f"{self._benchmark}.{self._num_tasklets}", "atomic.bin",
        )
        atomic_address = ConfigLoader.atomic_offset()
        atomic_bytes = InitThread.load_bytes(atomic_bin_filepath)
        self._dpu.dma().host_dma_transfer_to_atomic(atomic_address, atomic_bytes)

    def _dma_transfer_iram(self) -> None:
        iram_bin_filepath = os.path.join(
            PathCollector.bin_path_in_local(), f"{self._benchmark}.{self._num_tasklets}", "iram.bin",
        )
        iram_address = ConfigLoader.iram_offset()
        iram_bytes = InitThread.load_bytes(iram_bin_filepath)
        self._dpu.dma().host_dma_transfer_to_iram(iram_address, iram_bytes)

    def _dma_transfer_wram(self) -> None:
        wram_bin_filepath = os.path.join(
            PathCollector.bin_path_in_local(), f"{self._benchmark}.{self._num_tasklets}", "wram.bin",
        )
        wram_address = ConfigLoader.wram_offset()
        wram_bytes = InitThread.load_bytes(wram_bin_filepath)
        self._dpu.dma().host_dma_transfer_to_wram(wram_address, wram_bytes)

    def _dma_transfer_mram(self) -> None:
        mram_bin_filepath = os.path.join(
            PathCollector.bin_path_in_local(), f"{self._benchmark}.{self._num_tasklets}", "mram.bin",
        )
        mram_address = ConfigLoader.mram_offset()
        mram_bytes = InitThread.load_bytes(mram_bin_filepath)
        self._dpu.dma().host_dma_transfer_to_mram(mram_address, mram_bytes)

    @staticmethod
    def load_sys_used_mram_end_pointer(benchmark: str, num_tasklets: int) -> int:
        return InitThread.load_dpu_transfer_pointer(benchmark, num_tasklets)[0]

    @staticmethod
    def load_dpu_input_arguments_pointer(benchmark: str, num_tasklets: int) -> int:
        return InitThread.load_dpu_transfer_pointer(benchmark, num_tasklets)[1]

    @staticmethod
    def load_dpu_results_pointer(benchmark: str, num_tasklets: int) -> int:
        return InitThread.load_dpu_transfer_pointer(benchmark, num_tasklets)[2]

    @staticmethod
    def load_sys_end_pointer(benchmark: str, num_tasklets: int) -> int:
        return InitThread.load_dpu_transfer_pointer(benchmark, num_tasklets)[3]

    @staticmethod
    def load_bytes(bin_filepath: str) -> List[Byte]:
        bytes_: List[Byte] = []
        with open(bin_filepath) as file:
            for line in file.readlines():
                bytes_.append(Byte(int(line)))
        return bytes_

    @staticmethod
    def load_dpu_transfer_pointer(benchmark: str, num_tasklets: int) -> Tuple[int, int, int, int]:
        dpu_transfer_pointer_filepath = os.path.join(
            PathCollector.bin_path_in_local(), f"{benchmark}.{num_tasklets}", "dpu_transfer_pointer.bin",
        )

        with open(dpu_transfer_pointer_filepath) as file:
            lines = file.readlines()

            sys_used_mram_end, dpu_input_arguments, dpu_results, sys_end = (
                int(lines[0]),
                int(lines[1]),
                int(lines[2]),
                int(lines[3]),
            )
            return sys_used_mram_end, dpu_input_arguments, dpu_results, sys_end

    @staticmethod
    def load_num_executions(benchmark: str, num_tasklets: int) -> int:
        num_executions_filepath = os.path.join(
            PathCollector.bin_path_in_local(),
            f"{benchmark}.{num_tasklets}",
            "num_executions.bin",
        )

        with open(num_executions_filepath) as file:
            lines = file.readlines()
            num_executions = int(lines[0])
            return num_executions