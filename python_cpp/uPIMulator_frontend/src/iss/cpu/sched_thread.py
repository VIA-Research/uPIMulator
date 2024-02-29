import os
from typing import List

from encoder.byte import Byte
from iss.cpu.init_thread import InitThread
from iss.dpu.dpu import DPU
from iss.dram.mram_command import MRAMCommand
from util.path_collector import PathCollector


class SchedThread:
    def __init__(self, benchmark: str, num_tasklets: int, dpu: DPU):
        self._benchmark: str = benchmark
        self._num_tasklets: int = num_tasklets
        self._dpu: DPU = dpu

    def num_executions(self) -> int:
        return InitThread.load_num_executions(self._benchmark, self._num_tasklets)

    def sched(self, execution: int) -> None:
        self._dma_transfer_input_mram_heap_pointer_name(execution)
        self._dma_transfer_dpu_input_arguments(execution)

    def check(self, execution: int) -> None:
        self._dma_transfer_output_mram_heap_pointer_name(execution)
        self._dma_transfer_dpu_results(execution)

    def cycle(self) -> None:
        pass

    def _dma_transfer_input_mram_heap_pointer_name(self, execution: int) -> None:
        input_dpu_mram_heap_pointer_name_bin_filepath = os.path.join(
            PathCollector.bin_path_in_local(),
            f"{self._benchmark}.{self._num_tasklets}",
            f"input_dpu_mram_heap_pointer_name.{execution}.bin",
        )

        if os.path.exists(input_dpu_mram_heap_pointer_name_bin_filepath):
            sys_used_mram_end_pointer = InitThread.load_sys_used_mram_end_pointer(self._benchmark, self._num_tasklets)
            input_dpu_mram_heap_pointer_name_bytes = InitThread.load_bytes(
                input_dpu_mram_heap_pointer_name_bin_filepath
            )
            self._dpu.dma().host_dma_transfer_to_mram(sys_used_mram_end_pointer, input_dpu_mram_heap_pointer_name_bytes)

    def _dma_transfer_dpu_input_arguments(self, execution: int) -> None:
        dpu_input_arguments_bin_filepath = os.path.join(
            PathCollector.bin_path_in_local(),
            f"{self._benchmark}.{self._num_tasklets}",
            f"dpu_input_arguments.{execution}.bin",
        )

        if os.path.exists(dpu_input_arguments_bin_filepath):
            dpu_input_arguments_pointer = InitThread.load_dpu_input_arguments_pointer(
                self._benchmark, self._num_tasklets
            )

            if dpu_input_arguments_pointer != -1:
                dpu_input_arguments_bytes = InitThread.load_bytes(dpu_input_arguments_bin_filepath)
                self._dpu.dma().host_dma_transfer_to_wram(dpu_input_arguments_pointer, dpu_input_arguments_bytes)

    def _dma_transfer_output_mram_heap_pointer_name(self, execution: int) -> None:
        output_dpu_mram_heap_pointer_name_bin_filepath = os.path.join(
            PathCollector.bin_path_in_local(),
            f"{self._benchmark}.{self._num_tasklets}",
            f"output_dpu_mram_heap_pointer_name.{execution}.bin",
        )

        if os.path.exists(output_dpu_mram_heap_pointer_name_bin_filepath):
            output_dpu_mram_heap_pointer_name_bytes = InitThread.load_bytes(
                output_dpu_mram_heap_pointer_name_bin_filepath
            )

            sys_used_mram_end_pointer = InitThread.load_sys_used_mram_end_pointer(self._benchmark, self._num_tasklets)

            mram_command = MRAMCommand(
                MRAMCommand.Operation.READ, sys_used_mram_end_pointer, len(output_dpu_mram_heap_pointer_name_bytes),
            )
            assert self._dpu.mram().can_push()
            self._dpu.mram().push(mram_command)
            assert self._dpu.mram().can_pop()
            assert self._dpu.mram().pop() == mram_command

            mram_command_bytes: List[Byte] = []
            for data_word in mram_command.data_words():
                mram_command_bytes += data_word.to_bytes()

            for output_dpu_mram_heap_pointer_name_byte, mram_command_byte in zip(
                output_dpu_mram_heap_pointer_name_bytes, mram_command_bytes
            ):
                assert output_dpu_mram_heap_pointer_name_byte.value() == mram_command_byte.value()

    def _dma_transfer_dpu_results(self, execution: int) -> None:
        dpu_results_bin_filepath = os.path.join(
            PathCollector.bin_path_in_local(),
            f"{self._benchmark}.{self._num_tasklets}",
            f"dpu_results.{execution}.bin",
        )

        if os.path.exists(dpu_results_bin_filepath):
            dpu_results_bytes = InitThread.load_bytes(dpu_results_bin_filepath)

            dpu_results_pointer = InitThread.load_dpu_results_pointer(self._benchmark, self._num_tasklets)

            wram_bytes = self._dpu.dma().host_dma_transfer_from_wram(dpu_results_pointer, len(dpu_results_bytes))

            for dpu_results_byte, wram_byte in zip(dpu_results_bytes, wram_bytes):
                assert dpu_results_byte.value() == wram_byte.value()
