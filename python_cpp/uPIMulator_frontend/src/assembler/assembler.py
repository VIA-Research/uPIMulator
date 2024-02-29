import os
from typing import List, Set, Union

from abi.binary.executable import Executable
from abi.label.symbol import Symbol
from abi.section.section import Section
from abi.section.section_name import SectionName
from assembler.data_prep.bin import Bin
from assembler.data_prep.bs_data_prep import BSDataPrep
from assembler.data_prep.gemv_data_prep import GEMVDataPrep
from assembler.data_prep.hst_data_prep import HSTDataPrep
from assembler.data_prep.mlp_data_prep import MLPDataPrep
from assembler.data_prep.red_data_prep import REDDataPrep
from assembler.data_prep.scan_rss_data_prep import SCANRSSDataPrep
from assembler.data_prep.scan_ssa_data_prep import SCANSSADataPrep
from assembler.data_prep.sel_data_prep import SELDataPrep
from assembler.data_prep.trns_data_prep import TRNSDataPrep
from assembler.data_prep.ts_data_prep import TSDataPrep
from assembler.data_prep.uni_data_prep import UNIDataPrep
from assembler.data_prep.va_data_prep import VADataPrep
from encoder.byte import Byte
from linker_.linker_script import LinkerScript
from util.path_collector import PathCollector


class Assembler:
    DataPrep = Union[
        BSDataPrep,
        GEMVDataPrep,
        HSTDataPrep,
        MLPDataPrep,
        REDDataPrep,
        SCANRSSDataPrep,
        SCANSSADataPrep,
        SELDataPrep,
        TRNSDataPrep,
        TSDataPrep,
        UNIDataPrep,
        VADataPrep,
    ]

    def __init__(self):
        pass

    @staticmethod
    def assemble(
        executable: Executable, linker_script: LinkerScript, data_prep_param: List[int], num_dpus: int
    ) -> None:
        Assembler._assemble_atomic(executable, num_dpus)
        Assembler._assemble_iram(executable, num_dpus)
        Assembler._assemble_wram(executable, num_dpus)
        Assembler._assemble_mram(executable, num_dpus)

        Assembler._assemble_global_object_symbols(executable, num_dpus)
        Assembler._assemble_labels(executable, num_dpus)

        Assembler._assemble_dpu_transfer_pointer(executable, linker_script, num_dpus)

        benchmark = Assembler._benchmark(executable)
        num_tasklets = Assembler._num_tasklets(executable)
        data_prep = Assembler.data_prep(benchmark, num_tasklets, data_prep_param, num_dpus)

        Assembler._assemble_input_dpu_mram_heap_pointer_name(executable, data_prep)
        Assembler._assemble_output_dpu_mram_heap_pointer_name(executable, data_prep)
        Assembler._assemble_dpu_input_arguments(executable, data_prep)
        Assembler._assemble_dpu_results(executable, data_prep)
        Assembler._assemble_num_executions(executable, data_prep)

    @staticmethod
    def _assemble_atomic(executable: Executable, num_dpus: int) -> None:
        atomic_sections = Assembler._sort(executable.sections(SectionName.ATOMIC))

        benchmark = Assembler._benchmark(executable)
        num_tasklets = Assembler._num_tasklets(executable)
        atomic_filepath = os.path.join(
            PathCollector.bin_path_in_local(),
            f"{num_dpus}_dpus",
            f"{benchmark}.{num_tasklets}",
            "atomic.bin",
        )

        bytes_: List[Byte] = []
        for atomic_section in atomic_sections:
            bytes_ += atomic_section.to_bytes()
        atomic_bin = Bin(bytes_)
        atomic_bin.dump(atomic_filepath)

    @staticmethod
    def _assemble_iram(executable: Executable, num_dpus: int) -> None:
        text_sections = Assembler._sort(executable.sections(SectionName.TEXT))

        benchmark = Assembler._benchmark(executable)
        num_tasklets = Assembler._num_tasklets(executable)
        iram_filepath = os.path.join(
            PathCollector.bin_path_in_local(), f"{num_dpus}_dpus", f"{benchmark}.{num_tasklets}", "iram.bin"
        )

        bytes_: List[Byte] = []
        for text_section in text_sections:
            bytes_ += text_section.to_bytes()
        iram_bin = Bin(bytes_)
        iram_bin.dump(iram_filepath)

    @staticmethod
    def _assemble_wram(executable: Executable, num_dpus: int) -> None:
        sections = Assembler._sort(
            {
                *executable.sections(SectionName.DATA),
                *executable.sections(SectionName.RODATA),
                *executable.sections(SectionName.BSS),
                *executable.sections(SectionName.DPU_HOST),
            }
        )

        benchmark = Assembler._benchmark(executable)
        num_tasklets = Assembler._num_tasklets(executable)
        wram_filepath = os.path.join(
            PathCollector.bin_path_in_local(), f"{num_dpus}_dpus", f"{benchmark}.{num_tasklets}", "wram.bin"
        )

        bytes_: List[Byte] = []
        for section in sections:
            bytes_ += section.to_bytes()
        wram_bin = Bin(bytes_)
        wram_bin.dump(wram_filepath)

    @staticmethod
    def _assemble_mram(executable: Executable, num_dpus: int) -> None:
        sections = Assembler._sort(
            {
                *executable.sections(SectionName.DEBUG_ABBREV),
                *executable.sections(SectionName.DEBUG_FRAME),
                *executable.sections(SectionName.DEBUG_INFO),
                *executable.sections(SectionName.DEBUG_LINE),
                *executable.sections(SectionName.DEBUG_LOC),
                *executable.sections(SectionName.DEBUG_RANGES),
                *executable.sections(SectionName.DEBUG_STR),
                *executable.sections(SectionName.STACK_SIZES),
                *executable.sections(SectionName.MRAM),
            }
        )

        benchmark = Assembler._benchmark(executable)
        num_tasklets = Assembler._num_tasklets(executable)
        mram_filepath = os.path.join(
            PathCollector.bin_path_in_local(), f"{num_dpus}_dpus", f"{benchmark}.{num_tasklets}", "mram.bin"
        )

        bytes_: List[Byte] = []
        for section in sections:
            bytes_ += section.to_bytes()
        mram_bin = Bin(bytes_)
        mram_bin.dump(mram_filepath)

    @staticmethod
    def _assemble_global_object_symbols(executable: Executable, num_dpus: int) -> None:
        benchmark = Assembler._benchmark(executable)
        num_tasklets = Assembler._num_tasklets(executable)
        global_object_symbols_filepath = os.path.join(
            PathCollector.bin_path_in_local(),
            f"{num_dpus}_dpus",
            f"{benchmark}.{num_tasklets}",
            "global_object_symbols.bin",
        )
        with open(global_object_symbols_filepath, "w") as file:
            lines = ""
            for global_symbol in executable.liveness().global_symbols():
                if executable.liveness().symbol(global_symbol) == Symbol.OBJECT:
                    label = executable.label(global_symbol)
                    assert label is not None
                    symbol_address = label.address()
                    symbol_size = label.size()
                    lines += f"{symbol_address}: {symbol_size}\n"
            file.writelines(lines)

    @staticmethod
    def _assemble_dpu_transfer_pointer(executable: Executable, linker_script: LinkerScript, num_dpus: int) -> None:
        benchmark = Assembler._benchmark(executable)
        num_tasklets = Assembler._num_tasklets(executable)

        dpu_transfer_pointer_filepath = os.path.join(
            PathCollector.bin_path_in_local(),
            f"{num_dpus}_dpus",
            f"{benchmark}.{num_tasklets}",
            "dpu_transfer_pointer.bin",
        )
        with open(dpu_transfer_pointer_filepath, "w") as file:
            lines = ""

            sys_used_mram_end = linker_script.symbol("__sys_used_mram_end")
            assert sys_used_mram_end.address() is not None
            lines += f"{sys_used_mram_end.address()}\n"

            dpu_input_arguments = executable.label("DPU_INPUT_ARGUMENTS")
            if dpu_input_arguments is not None:
                assert dpu_input_arguments.address() is not None
                lines += f"{dpu_input_arguments.address()}\n"
            else:
                lines += "-1\n"

            dpu_results = executable.label("DPU_RESULTS")
            if dpu_results is not None:
                assert dpu_results is not None
                lines += f"{dpu_results.address()}\n"
            else:
                lines += "-1\n"

            sys_end = executable.label("__sys_end")
            assert sys_end is not None
            assert sys_end.address() is not None
            lines += f"{sys_end.address()}\n"

            file.writelines(lines)

    @staticmethod
    def _assemble_input_dpu_mram_heap_pointer_name(executable: Executable, data_prep: DataPrep) -> None:
        benchmark = Assembler._benchmark(executable)
        num_tasklets = Assembler._num_tasklets(executable)

        for execution in range(data_prep.num_executions()):
            for dpu_id in range(data_prep.num_dpus()):
                input_dpu_mram_heap_pointer_name = data_prep.input_dpu_mram_heap_pointer_name(execution, dpu_id)
                if input_dpu_mram_heap_pointer_name is not None:
                    input_dpu_mram_heap_pointer_name_filepath = os.path.join(
                        PathCollector.bin_path_in_local(),
                        f"{data_prep.num_dpus()}_dpus",
                        f"{benchmark}.{num_tasklets}",
                        f"input_dpu_mram_heap_pointer_name.dpu_id{dpu_id}.{execution}.bin",
                    )
                    input_dpu_mram_heap_pointer_name.dump(input_dpu_mram_heap_pointer_name_filepath)

    @staticmethod
    def _assemble_output_dpu_mram_heap_pointer_name(executable: Executable, data_prep: DataPrep) -> None:
        benchmark = Assembler._benchmark(executable)
        num_tasklets = Assembler._num_tasklets(executable)

        for execution in range(data_prep.num_executions()):
            for dpu_id in range(data_prep.num_dpus()):
                output_dpu_mram_heap_pointer_name = data_prep.output_dpu_mram_heap_pointer_name(execution, dpu_id)
                if output_dpu_mram_heap_pointer_name is not None:
                    output_dpu_mram_heap_pointer_name_filepath = os.path.join(
                        PathCollector.bin_path_in_local(),
                        f"{data_prep.num_dpus()}_dpus",
                        f"{benchmark}.{num_tasklets}",
                        f"output_dpu_mram_heap_pointer_name.dpu_id{dpu_id}.{execution}.bin",
                    )
                    output_dpu_mram_heap_pointer_name.dump(output_dpu_mram_heap_pointer_name_filepath)

    @staticmethod
    def _assemble_dpu_input_arguments(executable: Executable, data_prep: DataPrep) -> None:
        benchmark = Assembler._benchmark(executable)
        num_tasklets = Assembler._num_tasklets(executable)

        for execution in range(data_prep.num_executions()):
            for dpu_id in range(data_prep.num_dpus()):
                dpu_input_arguments = data_prep.dpu_input_arguments(execution, dpu_id)
                if dpu_input_arguments is not None:
                    dpu_input_arguments_filepath = os.path.join(
                        PathCollector.bin_path_in_local(),
                        f"{data_prep.num_dpus()}_dpus",
                        f"{benchmark}.{num_tasklets}",
                        f"dpu_input_arguments.dpu_id{dpu_id}.{execution}.bin",
                    )

                    dpu_input_arguments.dump(dpu_input_arguments_filepath)

    @staticmethod
    def _assemble_dpu_results(executable: Executable, data_prep: DataPrep) -> None:
        benchmark = Assembler._benchmark(executable)
        num_tasklets = Assembler._num_tasklets(executable)

        for execution in range(data_prep.num_executions()):
            for dpu_id in range(data_prep.num_dpus()):
                dpu_results = data_prep.dpu_results(execution, dpu_id)
                if dpu_results is not None:
                    dpu_results_filepath = os.path.join(
                        PathCollector.bin_path_in_local(),
                        f"{data_prep.num_dpus()}_dpus",
                        f"{benchmark}.{num_tasklets}",
                        f"dpu_results.dpu_id{dpu_id}.{execution}.bin",
                    )

                    dpu_results.dump(dpu_results_filepath)

    @staticmethod
    def _assemble_labels(executable: Executable, num_dpus: int) -> None:
        benchmark = Assembler._benchmark(executable)
        num_tasklets = Assembler._num_tasklets(executable)
        labels_filepath = os.path.join(
            PathCollector.bin_path_in_local(),
            f"{num_dpus}_dpus",
            f"{benchmark}.{num_tasklets}",
            "labels.bin",
        )
        with open(labels_filepath, "w") as file:
            lines = ""
            for label in executable.labels():
                lines += f"{label.name()}: {label.address()}\n"
            file.writelines(lines)

    @staticmethod
    def _assemble_num_executions(executable: Executable, data_prep: DataPrep) -> None:
        benchmark = Assembler._benchmark(executable)
        num_tasklets = Assembler._num_tasklets(executable)
        num_executions_filepath = os.path.join(
            PathCollector.bin_path_in_local(),
            f"{data_prep.num_dpus()}_dpus",
            f"{benchmark}.{num_tasklets}",
            "num_executions.bin",
        )
        with open(num_executions_filepath, "w") as file:
            lines = f"{data_prep.num_executions()}\n"
            file.writelines(lines)

    @staticmethod
    def data_prep(benchmark: str, num_tasklets: int, data_prep_param: List[int], num_dpus: int) -> DataPrep:
        if benchmark == "BS":
            return BSDataPrep(num_tasklets, data_prep_param, num_dpus)
        elif benchmark == "GEMV":
            return GEMVDataPrep(num_tasklets, data_prep_param, num_dpus)
        elif benchmark == "HST-L" or benchmark == "HST-S":
            return HSTDataPrep(num_tasklets, data_prep_param, num_dpus)
        elif benchmark == "MLP":
            return MLPDataPrep(num_tasklets, data_prep_param, num_dpus)
        elif benchmark == "RED":
            return REDDataPrep(num_tasklets, data_prep_param, num_dpus)
        elif benchmark == "SCAN-RSS":
            return SCANRSSDataPrep(num_tasklets, data_prep_param, num_dpus)
        elif benchmark == "SCAN-SSA":
            return SCANSSADataPrep(num_tasklets, data_prep_param, num_dpus)
        elif benchmark == "SEL":
            return SELDataPrep(num_tasklets, data_prep_param, num_dpus)
        elif benchmark == "TRNS":
            return TRNSDataPrep(num_tasklets, data_prep_param, num_dpus)
        elif benchmark == "TS":
            return TSDataPrep(num_tasklets, data_prep_param, num_dpus)
        elif benchmark == "UNI":
            return UNIDataPrep(num_tasklets, data_prep_param, num_dpus)
        elif benchmark == "VA":
            return VADataPrep(num_tasklets, data_prep_param, num_dpus)
        else:
            raise ValueError

    @staticmethod
    def _benchmark(executable: Executable) -> str:
        return executable.filepath().split(os.path.sep)[-2].split(".")[0]

    @staticmethod
    def _num_tasklets(executable: Executable) -> int:
        return int(executable.filepath().split(os.path.sep)[-2].split(".")[-1])

    @staticmethod
    def _sort(sections: Set[Section]) -> List[Section]:
        def _address(section: Section) -> int:
            address = section.address()
            assert address is not None
            return address

        return sorted(sections, key=_address)
