import math
from typing import Dict, Set

from abi.binary.executable import Executable
from abi.label.label import Label
from abi.section.section_name import SectionName
from abi.word.instruction_word import InstructionWord
from util.config_loader import ConfigLoader


class LinkerScript:
    def __init__(self, num_tasklets: int):
        self._num_tasklets: int = num_tasklets

        self._constants: Dict[str, int] = self._init_constants()
        self._symbols: Set[Label] = self._init_symbols()

    def constant(self, constant_name: str) -> int:
        return self._constants[constant_name]

    def constant_names(self) -> Set[str]:
        return set(self._constants.keys())

    def symbol(self, symbol_name: str) -> Label:
        for symbol in self._symbols:
            if symbol.name() == symbol_name:
                return symbol
        raise ValueError

    def symbol_names(self) -> Set[str]:
        return {symbol.name() for symbol in self._symbols}

    def assign_address(self, executable: Executable) -> None:
        self._assign_iram(executable)
        self._assign_atomic(executable)
        self._assign_wram(executable)
        self._assign_mram(executable)

    def _assign_iram(self, executable: Executable) -> None:
        cur_address = ConfigLoader.iram_offset()
        bootstrap = executable.section(SectionName.TEXT, "__bootstrap")

        assert bootstrap is not None
        assert bootstrap.address() is None

        bootstrap.set_address(cur_address)
        cur_address += bootstrap.size()
        assert bootstrap.size() % InstructionWord().size() == 0

        text_section = executable.section(SectionName.TEXT, "")
        if text_section is not None:
            assert text_section.address() is None

            text_section.set_address(cur_address)
            cur_address += text_section.size()

            assert text_section.size() % InstructionWord().size() == 0

        for section in executable.sections(SectionName.TEXT):
            if section.name() != "__bootstrap" and section.name() != "":
                assert section.address() is None

                section.set_address(cur_address)
                cur_address += section.size()

                assert section.size() % InstructionWord().size() == 0

        assert cur_address < (ConfigLoader.iram_offset() + ConfigLoader.iram_size())

    def _assign_atomic(self, executable: Executable) -> None:
        cur_address = ConfigLoader.atomic_offset()
        self.symbol("__atomic_start_addr").set_address(cur_address)

        # TODO(bongjoon.hyun@gmail.com): the original UPMEM linker script adds 200 to the cur_address
        # cur_address += 200

        self.symbol("__atomic_used_addr").set_address(cur_address)
        for section in executable.sections(SectionName.ATOMIC):
            assert section.address() is None

            section.set_address(cur_address)
            cur_address += section.size()

        self.symbol("__atomic_end_addr").set_address(cur_address)

        assert cur_address < (ConfigLoader.atomic_offset() + ConfigLoader.atomic_size())

    def _assign_wram(self, executable: Executable) -> None:
        cur_address = ConfigLoader.wram_offset()

        sys_zero = executable.section(SectionName.DATA, "__sys_zero")
        if sys_zero is not None:
            assert sys_zero.address() is None

            sys_zero.set_address(cur_address)
            cur_address += sys_zero.size()

        immediate_memory = executable.section(SectionName.DATA, "immediate_memory")
        if immediate_memory is not None:
            assert immediate_memory.address() is None

            immediate_memory.set_address(cur_address)
            cur_address += immediate_memory.size()

        for section in executable.sections(SectionName.DATA):
            if "immediate_memory." in section.name():
                assert section.address() is None

                section.set_address(cur_address)
                cur_address += section.size()

        self.symbol("__rodata_start_addr").set_address(cur_address)
        rodata_section = executable.section(SectionName.RODATA, "")
        if rodata_section is not None:
            assert rodata_section.address() is None

            rodata_section.set_address(cur_address)
            cur_address += rodata_section.size()

        for section in executable.sections(SectionName.RODATA):
            if section.name() != "":
                assert section.address() is None

                section.set_address(cur_address)
                cur_address += section.size()
        self.symbol("__rodata_end_addr").set_address(cur_address)

        bss_section = executable.section(SectionName.BSS, "")
        if bss_section is not None:
            assert bss_section.address() is None

            bss_section.set_address(cur_address)
            cur_address += bss_section.size()

        for section in executable.sections(SectionName.BSS):
            if section.name() != "":
                assert section.address() is None

                section.set_address(cur_address)
                cur_address += section.size()

        sys_keep = executable.section(SectionName.DATA, "__sys_keep")
        if sys_keep is not None:
            assert sys_keep.address() is None

            sys_keep.set_address(cur_address)
            cur_address += sys_keep.size()

        data_section = executable.section(SectionName.DATA, "")
        if data_section is not None:
            assert data_section.address() is None

            data_section.set_address(cur_address)
            cur_address += data_section.size()

        for section in executable.sections(SectionName.DATA):
            if (
                section.name() != "__sys_zero"
                and section.name() != "__sys_keep"
                and section.name() != ""
                and "immediate_memory" not in section.name()
            ):
                assert section.address() is None

                section.set_address(cur_address)
                cur_address += section.size()

        dpu_host = executable.section(SectionName.DPU_HOST, "")
        if dpu_host is not None:
            assert dpu_host.address() is None

            dpu_host.set_address(cur_address)
            cur_address += dpu_host.size()

        for i in range(ConfigLoader.max_num_tasklets()):
            self.symbol(f"__sys_stack_thread_{i}").set_address(cur_address)
            cur_address += self.constant(f"STACK_SIZE_TASKLET_{i}")

        self.symbol("__sw_cache_buffer").set_address(cur_address)
        cur_address += 8 * self._num_tasklets

        cur_address = (
            math.ceil(cur_address / ConfigLoader.min_access_granularity()) * ConfigLoader.min_access_granularity()
        )
        self.symbol("__sys_heap_pointer_reset").set_address(cur_address)

        assert cur_address < (ConfigLoader.wram_offset() + ConfigLoader.wram_size())

    def _assign_mram(self, executable: Executable) -> None:
        cur_address = ConfigLoader.mram_offset()

        debug_abbrev_section = executable.section(SectionName.DEBUG_ABBREV, "")
        if debug_abbrev_section is not None:
            assert debug_abbrev_section.address() is None

            debug_abbrev_section.set_address(cur_address)
            cur_address += debug_abbrev_section.size()

        debug_frame_section = executable.section(SectionName.DEBUG_FRAME, "")
        if debug_frame_section is not None:
            assert debug_frame_section.address() is None

            debug_frame_section.set_address(cur_address)
            cur_address += debug_frame_section.size()

        debug_info_section = executable.section(SectionName.DEBUG_INFO, "")
        if debug_info_section is not None:
            assert debug_info_section.address() is None

            debug_info_section.set_address(cur_address)
            cur_address += debug_info_section.size()

        debug_line_section = executable.section(SectionName.DEBUG_LINE, "")
        if debug_line_section is not None:
            assert debug_line_section.address() is None

            debug_line_section.set_address(cur_address)
            cur_address += debug_line_section.size()

        debug_loc_section = executable.section(SectionName.DEBUG_LOC, "")
        if debug_loc_section is not None:
            assert debug_loc_section.address() is None

            debug_loc_section.set_address(cur_address)
            cur_address += debug_loc_section.size()

        debug_ranges_section = executable.section(SectionName.DEBUG_RANGES, "")
        if debug_ranges_section is not None:
            assert debug_ranges_section.address() is None

            debug_ranges_section.set_address(cur_address)
            cur_address += debug_ranges_section.size()

        debug_str_section = executable.section(SectionName.DEBUG_STR, "")
        if debug_str_section is not None:
            assert debug_str_section.address() is None

            debug_str_section.set_address(cur_address)
            cur_address += debug_str_section.size()

        stack_sizes_section = executable.section(SectionName.STACK_SIZES, "")
        if stack_sizes_section is not None:
            assert stack_sizes_section.address() is None

            stack_sizes_section.set_address(cur_address)
            cur_address += stack_sizes_section.size()

        noinit = executable.section(SectionName.MRAM, "noinit")
        if noinit is not None:
            assert noinit.address() is None

            noinit.set_address(cur_address)
            cur_address += noinit.size()

        for section in executable.sections(SectionName.MRAM):
            if "noinit." in section.name():
                assert section.address() is None

                section.set_address(cur_address)
                cur_address += section.size()

        mram_section = executable.section(SectionName.MRAM, "")
        if mram_section is not None:
            assert mram_section.address() is None

            mram_section.set_address(cur_address)
            cur_address += mram_section.size()

        for section in executable.sections(SectionName.MRAM):
            if section.name() != "" and "noinit" not in section.name() and "keep" not in section.name():
                assert section.address() is None

                section.set_address(cur_address)
                cur_address += section.size()

        keep = executable.section(SectionName.MRAM, "keep")
        if keep is not None:
            assert keep.address() is None

            keep.set_address(cur_address)
            cur_address += keep.size()

        for section in executable.sections(SectionName.MRAM):
            if "keep." in section.name():
                assert section.address() is None

                section.set_address(cur_address)
                cur_address += section.size()

        cur_address = (
            math.ceil(cur_address / ConfigLoader.min_access_granularity()) * ConfigLoader.min_access_granularity()
        )
        self.symbol("__sys_used_mram_end").set_address(cur_address)

        assert cur_address < (ConfigLoader.mram_offset() + ConfigLoader.mram_size())

    def _init_constants(self) -> Dict[str, int]:
        return {
            "NR_TASKLETS": self._num_tasklets,
            **self._init_stack_size_constants(),
        }

    def _init_stack_size_constants(self) -> Dict[str, int]:
        return {f"STACK_SIZE_TASKLET_{i}": ConfigLoader.stack_size() for i in range(ConfigLoader.max_num_tasklets())}

    def _init_symbols(self) -> Set[Label]:
        return {
            *self._init_atomic_symbols(),
            *self._init_wram_symbols(),
            *self._init_mram_symbols(),
        }

    def _init_atomic_symbols(self) -> Set[Label]:
        return {
            Label("__atomic_start_addr"),
            Label("__atomic_used_addr"),
            Label("__atomic_end_addr"),
        }

    def _init_wram_symbols(self) -> Set[Label]:
        return {
            Label("__rodata_start_addr"),
            Label("__rodata_end_addr"),
            *self._init_stack_symbols(),
            Label("__sw_cache_buffer"),
            Label("__sys_heap_pointer_reset"),
        }

    def _init_stack_symbols(self) -> Set[Label]:
        return {Label(f"__sys_stack_thread_{i}") for i in range(ConfigLoader.max_num_tasklets())}

    def _init_mram_symbols(self) -> Set[Label]:
        return {Label("__sys_used_mram_end")}
