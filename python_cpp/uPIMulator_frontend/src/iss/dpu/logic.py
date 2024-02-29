import logging
import math
from typing import Optional

from abi.isa.flag import Flag
from abi.isa.instruction.condition import Condition
from abi.isa.instruction.instruction import Instruction
from abi.isa.instruction.op_code import OpCode
from abi.isa.instruction.suffix import Suffix
from abi.isa.register.gp_register import GPRegister
from abi.word.data_word import DataWord
from abi.word.instruction_word import InstructionWord
from abi.word.representation import Representation
from converter.instruction_converter import InstructionConverter
from iss.dpu.alu import ALU
from iss.dpu.decoder import Decoder
from iss.dpu.dispatcher import Dispatcher
from iss.dpu.dma import DMA
from iss.dpu.scheduler import Scheduler
from iss.dpu.thread import Thread
from iss.dram.mram import MRAM
from iss.sram.atomic import Atomic
from iss.sram.iram import IRAM
from iss.sram.wram import WRAM
from util.config_loader import ConfigLoader


class Logic:
    def __init__(
        self,
        scheduler: Scheduler,
        atomic: Atomic,
        iram: IRAM,
        wram: WRAM,
        mram: MRAM,
        dma: DMA,
        dispatcher: Dispatcher,
    ):
        self._scheduler: Scheduler = scheduler
        self._atomic: Atomic = atomic
        self._iram: IRAM = iram
        self._wram: WRAM = wram
        self._mram: MRAM = mram
        self._dma: DMA = dma
        self._dispatcher: Dispatcher = dispatcher
        self._logger: logging.Logger = logging.getLogger("iss")

        self._cur_thread: Optional[Thread] = None

    def cycle(self) -> None:
        self._cur_thread = self._scheduler.schedule()
        if self._cur_thread is not None:
            instruction_word = self._iram.read(self._cur_thread.register_file().read_pc())
            instruction = Decoder.decode(instruction_word)

            self._print_instruction(instruction)
            self._execute_instruction(instruction)
            self._print_register_file()

    def _print_instruction(self, instruction: Instruction) -> None:
        assert self._cur_thread is not None
        self._logger.info(f"[{self._cur_thread.id_()}] {InstructionConverter.convert_to_string(instruction)}")

    def _print_register_file(self) -> None:
        assert self._cur_thread is not None
        for i in range(ConfigLoader.num_gp_registers()):
            self._logger.info(f"r{i}: {self._cur_thread.register_file().read(GPRegister(i), Representation.SIGNED)}")

    def _execute_instruction(self, instruction: Instruction) -> None:
        if instruction.suffix() == Suffix.RICI:
            self._execute_rici(instruction)
        elif instruction.suffix() == Suffix.RRI:
            self._execute_rri(instruction)
        elif instruction.suffix() == Suffix.RRIC:
            self._execute_rric(instruction)
        elif instruction.suffix() == Suffix.RRICI:
            self._execute_rrici(instruction)
        elif instruction.suffix() == Suffix.RRIF:
            self._execute_rrif(instruction)
        elif instruction.suffix() == Suffix.RRR:
            self._execute_rrr(instruction)
        elif instruction.suffix() == Suffix.RRRC:
            self._execute_rrrc(instruction)
        elif instruction.suffix() == Suffix.RRRCI:
            self._execute_rrrci(instruction)
        elif instruction.suffix() == Suffix.ZRI:
            self._execute_zri(instruction)
        elif instruction.suffix() == Suffix.ZRIC:
            self._execute_zric(instruction)
        elif instruction.suffix() == Suffix.ZRICI:
            self._execute_zrici(instruction)
        elif instruction.suffix() == Suffix.ZRIF:
            self._execute_zrif(instruction)
        elif instruction.suffix() == Suffix.ZRR:
            self._execute_zrr(instruction)
        elif instruction.suffix() == Suffix.ZRRC:
            self._execute_zrrc(instruction)
        elif instruction.suffix() == Suffix.ZRRCI:
            self._execute_zrrci(instruction)
        elif instruction.suffix() == Suffix.S_RRI:
            self._execute_s_rri(instruction)
        elif instruction.suffix() == Suffix.S_RRIC:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.S_RRICI:
            self._execute_s_rrici(instruction)
        elif instruction.suffix() == Suffix.S_RRIF:
            self._execute_s_rrif(instruction)
        elif instruction.suffix() == Suffix.S_RRR:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.S_RRRC:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.S_RRRCI:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.U_RRI:
            self._execute_u_rri(instruction)
        elif instruction.suffix() == Suffix.U_RRIC:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.U_RRICI:
            self._execute_u_rrici(instruction)
        elif instruction.suffix() == Suffix.U_RRIF:
            self._execute_u_rrif(instruction)
        elif instruction.suffix() == Suffix.U_RRR:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.U_RRRC:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.U_RRRCI:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.RR:
            self._execute_rr(instruction)
        elif instruction.suffix() == Suffix.RRC:
            self._execute_rrc(instruction)
        elif instruction.suffix() == Suffix.RRCI:
            self._execute_rrci(instruction)
        elif instruction.suffix() == Suffix.ZR:
            self._execute_zr(instruction)
        elif instruction.suffix() == Suffix.ZRC:
            self._execute_zrc(instruction)
        elif instruction.suffix() == Suffix.ZRCI:
            self._execute_zrci(instruction)
        elif instruction.suffix() == Suffix.S_RR:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.S_RRC:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.S_RRCI:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.U_RR:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.U_RRC:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.U_RRCI:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.DRDICI:
            self._execute_drdici(instruction)
        elif instruction.suffix() == Suffix.RRRI:
            self._execute_rrri(instruction)
        elif instruction.suffix() == Suffix.RRRICI:
            self._execute_rrrici(instruction)
        elif instruction.suffix() == Suffix.ZRRI:
            self._execute_zrri(instruction)
        elif instruction.suffix() == Suffix.ZRRICI:
            self._execute_zrrici(instruction)
        elif instruction.suffix() == Suffix.S_RRRI:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.S_RRRICI:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.U_RRRI:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.U_RRRICI:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.RIR:
            self._execute_rir(instruction)
        elif instruction.suffix() == Suffix.RIRC:
            self._execute_rirc(instruction)
        elif instruction.suffix() == Suffix.RIRCI:
            self._execute_rirci(instruction)
        elif instruction.suffix() == Suffix.ZIR:
            self._execute_zir(instruction)
        elif instruction.suffix() == Suffix.ZIRC:
            self._execute_zirc(instruction)
        elif instruction.suffix() == Suffix.ZIRCI:
            self._execute_zirci(instruction)
        elif instruction.suffix() == Suffix.S_RIRC:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.S_RIRCI:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.U_RIRC:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.U_RIRCI:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.R:
            self._execute_r(instruction)
        elif instruction.suffix() == Suffix.RCI:
            self._execute_rci(instruction)
        elif instruction.suffix() == Suffix.Z:
            self._execute_z(instruction)
        elif instruction.suffix() == Suffix.ZCI:
            self._execute_zci(instruction)
        elif instruction.suffix() == Suffix.S_R:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.S_RCI:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.U_R:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.U_RCI:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.CI:
            self._execute_ci(instruction)
        elif instruction.suffix() == Suffix.I:
            self._execute_i(instruction)
        elif instruction.suffix() == Suffix.DDCI:
            self._execute_ddci(instruction)
        elif instruction.suffix() == Suffix.ERRI:
            self._execute_erri(instruction)
        elif instruction.suffix() == Suffix.S_ERRI:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.U_ERRI:
            raise NotImplementedError
        elif instruction.suffix() == Suffix.EDRI:
            self._execute_edri(instruction)
        elif instruction.suffix() == Suffix.ERII:
            self._execute_erii(instruction)
        elif instruction.suffix() == Suffix.ERIR:
            self._execute_erir(instruction)
        elif instruction.suffix() == Suffix.ERID:
            self._execute_erid(instruction)
        elif instruction.suffix() == Suffix.DMA_RRI:
            self._execute_dma_rri(instruction)
        else:
            raise ValueError

    def _execute_rici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RICIOpCodes

        if instruction.op_code() in Instruction.AcquireRICIOpCodes:
            self._execute_acquire_rici(instruction)
        elif instruction.op_code() in Instruction.ReleaseRICIOpCodes:
            self._execute_release_rici(instruction)
        elif instruction.op_code() in Instruction.BootRICIOpCodes:
            self._execute_boot_rici(instruction)
        else:
            raise ValueError

    def _execute_acquire_rici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AcquireRICIOpCodes
        assert instruction.suffix() == Suffix.RICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(
            instruction.ra(),
            Representation.UNSIGNED,
        )
        imm = instruction.imm().value()
        atomic_address = ALU.atomic_address_hash(ra, imm)

        can_acquire = self._atomic.can_acquire(atomic_address)
        if can_acquire:
            self._atomic.acquire(atomic_address, self._cur_thread.id_())

        self._cur_thread.register_file().clear_conditions()
        self._set_acquire_cc(not can_acquire)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(not can_acquire, False)

    def _execute_release_rici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.ReleaseRICIOpCodes
        assert instruction.suffix() == Suffix.RICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(
            instruction.ra(),
            Representation.UNSIGNED,
        )
        imm = instruction.imm().value()
        atomic_address = ALU.atomic_address_hash(ra, imm)

        can_release = self._atomic.can_release(atomic_address, self._cur_thread.id_())
        if can_release:
            self._atomic.release(atomic_address, self._cur_thread.id_())

        self._cur_thread.register_file().clear_conditions()
        self._set_release_cc(not can_release)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(not can_release, False)

    def _execute_boot_rici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.BootRICIOpCodes
        assert instruction.suffix() == Suffix.RICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(
            instruction.ra(),
            Representation.UNSIGNED,
        )
        imm = instruction.imm().value()
        thread_id = ALU.atomic_address_hash(ra, imm)

        self._cur_thread.register_file().clear_conditions()
        if instruction.op_code() == OpCode.BOOT:
            can_boot = self._scheduler.boot(thread_id)
            self._set_boot_cc(ra, not can_boot)
            self._set_flags(not can_boot, False)
        elif instruction.op_code() == OpCode.RESUME:
            can_resume = self._scheduler.awake(thread_id)
            self._set_boot_cc(ra, not can_resume)
            self._set_flags(not can_resume, False)
        else:
            raise ValueError

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

    def _execute_rri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRIOpCodes

        if instruction.op_code() in Instruction.AddRRIOpCodes:
            self._execute_add_rri(instruction)
        elif instruction.op_code() in Instruction.AsrRRIOpCodes:
            self._execute_asr_rri(instruction)
        elif instruction.op_code() in Instruction.CallRRIOpCodes:
            self._execute_call_rri(instruction)
        else:
            raise ValueError

    def _execute_add_rri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AddRRIOpCodes
        assert instruction.suffix() == Suffix.RRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ADD:
            result, carry, _ = ALU.add(ra, imm)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, _ = ALU.addc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.AND:
            result, carry = ALU.and_(ra, imm), False
        elif instruction.op_code() == OpCode.OR:
            result, carry = ALU.or_(ra, imm), False
        elif instruction.op_code() == OpCode.XOR:
            result, carry = ALU.xor(ra, imm), False
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._cur_thread.register_file().write(instruction.rc(), result)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_asr_rri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AsrRRIOpCodes
        assert instruction.suffix() == Suffix.RRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ASR:
            result = ALU.asr(ra, imm)
        elif instruction.op_code() == OpCode.LSL:
            result = ALU.lsl(ra, imm)
        elif instruction.op_code() == OpCode.LSL1:
            result = ALU.lsl1(ra, imm)
        elif instruction.op_code() == OpCode.LSL1X:
            result = ALU.lsl1x(ra, imm)
        elif instruction.op_code() == OpCode.LSLX:
            result = ALU.lslx(ra, imm)
        elif instruction.op_code() == OpCode.LSR:
            result = ALU.lsr(ra, imm)
        elif instruction.op_code() == OpCode.LSR1:
            result = ALU.lsr1(ra, imm)
        elif instruction.op_code() == OpCode.LSR1X:
            result = ALU.lsr1x(ra, imm)
        elif instruction.op_code() == OpCode.LSRX:
            result = ALU.lsrx(ra, imm)
        elif instruction.op_code() == OpCode.ROL:
            result = ALU.rol(ra, imm)
        elif instruction.op_code() == OpCode.ROR:
            result = ALU.ror(ra, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.rc(), result)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_call_rri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.CallRRIOpCodes
        assert instruction.suffix() == Suffix.RRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if imm == 0:
            callee_address, carry, _ = ALU.add(ra, imm)
        else:
            callee_address, carry, _ = ALU.add(ra * InstructionWord().size(), imm)

        self._cur_thread.register_file().clear_conditions()

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write(instruction.rc(), pc + InstructionWord().size())

        self._cur_thread.register_file().write_pc(callee_address)

        self._set_flags(callee_address, carry)

    def _execute_rric(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRICOpCodes

        if instruction.op_code() in Instruction.AddRRICOpCodes:
            self._execute_add_rric(instruction)
        elif instruction.op_code() in Instruction.AsrRRICOpCodes:
            self._execute_asr_rric(instruction)
        elif instruction.op_code() in Instruction.SubRRICOpCodes:
            self._execute_sub_rric(instruction)
        else:
            raise ValueError

    def _execute_add_rric(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AddRRICOpCodes
        assert instruction.suffix() == Suffix.RRIC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ADD:
            result, carry, _ = ALU.add(ra, imm)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, _ = ALU.addc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.AND:
            result, carry = ALU.and_(ra, imm), False
        elif instruction.op_code() == OpCode.ANDN:
            result, carry = ALU.andn(ra, imm), False
        elif instruction.op_code() == OpCode.NAND:
            result, carry = ALU.nand(ra, imm), False
        elif instruction.op_code() == OpCode.NOR:
            result, carry = ALU.nor(ra, imm), False
        elif instruction.op_code() == OpCode.NXOR:
            result, carry = ALU.nxor(ra, imm), False
        elif instruction.op_code() == OpCode.OR:
            result, carry = ALU.or_(ra, imm), False
        elif instruction.op_code() == OpCode.ORN:
            result, carry = ALU.orn(ra, imm), False
        elif instruction.op_code() == OpCode.XOR:
            result, carry = ALU.xor(ra, imm), False
        elif instruction.op_code() == OpCode.HASH:
            result, carry = ALU.hash(ra, imm), False
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_set_cc(result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            self._cur_thread.register_file().write(instruction.rc(), 1)
        else:
            self._cur_thread.register_file().write(instruction.rc(), 0)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_asr_rric(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AsrRRICOpCodes
        assert instruction.suffix() == Suffix.RRIC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ASR:
            result = ALU.asr(ra, imm)
        elif instruction.op_code() == OpCode.LSL:
            result = ALU.lsl(ra, imm)
        elif instruction.op_code() == OpCode.LSL1:
            result = ALU.lsl1(ra, imm)
        elif instruction.op_code() == OpCode.LSL1X:
            result = ALU.lsl1x(ra, imm)
        elif instruction.op_code() == OpCode.LSLX:
            result = ALU.lslx(ra, imm)
        elif instruction.op_code() == OpCode.LSR:
            result = ALU.lsr(ra, imm)
        elif instruction.op_code() == OpCode.LSR1:
            result = ALU.lsr1(ra, imm)
        elif instruction.op_code() == OpCode.LSR1X:
            result = ALU.lsr1x(ra, imm)
        elif instruction.op_code() == OpCode.LSRX:
            result = ALU.lsrx(ra, imm)
        elif instruction.op_code() == OpCode.ROL:
            result = ALU.rol(ra, imm)
        elif instruction.op_code() == OpCode.ROR:
            result = ALU.ror(ra, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_set_cc(result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            self._cur_thread.register_file().write(instruction.rc(), 1)
        else:
            self._cur_thread.register_file().write(instruction.rc(), 0)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_sub_rric(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.SubRRICOpCodes
        assert instruction.suffix() == Suffix.RRIC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.SUB:
            result, carry, overflow = ALU.sub(ra, imm)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, overflow = ALU.subc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_ext_sub_set_cc(ra, imm, result, carry, overflow)

        if self._cur_thread.register_file().condition(instruction.condition()):
            self._cur_thread.register_file().write(instruction.rc(), 1)
        else:
            self._cur_thread.register_file().write(instruction.rc(), 0)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_rrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRICIOpCodes

        if instruction.op_code() in Instruction.AddRRICIOpCodes:
            self._execute_add_rrici(instruction)
        elif instruction.op_code() in Instruction.AndRRICIOpCodes:
            self._execute_and_rrici(instruction)
        elif instruction.op_code() in Instruction.AsrRRICIOpCodes:
            self._execute_asr_rrici(instruction)
        elif instruction.op_code() in Instruction.SubRRICIOpCodes:
            self._execute_sub_rrici(instruction)
        else:
            raise ValueError

    def _execute_add_rrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AddRRICIOpCodes
        assert instruction.suffix() == Suffix.RRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ADD:
            result, carry, overflow = ALU.add(ra, imm)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, overflow = ALU.addc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_add_nz_cc(ra, result, carry, overflow)

        self._cur_thread.register_file().write(instruction.rc(), result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_and_rrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AndRRICIOpCodes
        assert instruction.suffix() == Suffix.RRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.AND:
            result = ALU.and_(ra, imm)
        elif instruction.op_code() == OpCode.ANDN:
            result = ALU.andn(ra, imm)
        elif instruction.op_code() == OpCode.NAND:
            result = ALU.nand(ra, imm)
        elif instruction.op_code() == OpCode.NOR:
            result = ALU.nor(ra, imm)
        elif instruction.op_code() == OpCode.NXOR:
            result = ALU.nxor(ra, imm)
        elif instruction.op_code() == OpCode.OR:
            result = ALU.or_(ra, imm)
        elif instruction.op_code() == OpCode.ORN:
            result = ALU.orn(ra, imm)
        elif instruction.op_code() == OpCode.XOR:
            result = ALU.xor(ra, imm)
        elif instruction.op_code() == OpCode.HASH:
            result = ALU.hash(ra, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_nz_cc(ra, result)

        self._cur_thread.register_file().write(instruction.rc(), result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_asr_rrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AsrRRICIOpCodes
        assert instruction.suffix() == Suffix.RRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ASR:
            result = ALU.asr(ra, imm)
        elif instruction.op_code() == OpCode.LSL:
            result = ALU.lsl(ra, imm)
        elif instruction.op_code() == OpCode.LSL1:
            result = ALU.lsl1(ra, imm)
        elif instruction.op_code() == OpCode.LSL1X:
            result = ALU.lsl1x(ra, imm)
        elif instruction.op_code() == OpCode.LSLX:
            result = ALU.lslx(ra, imm)
        elif instruction.op_code() == OpCode.LSR:
            result = ALU.lsr(ra, imm)
        elif instruction.op_code() == OpCode.LSR1:
            result = ALU.lsr1(ra, imm)
        elif instruction.op_code() == OpCode.LSR1X:
            result = ALU.lsr1x(ra, imm)
        elif instruction.op_code() == OpCode.LSRX:
            result = ALU.lsrx(ra, imm)
        elif instruction.op_code() == OpCode.ROL:
            result = ALU.rol(ra, imm)
        elif instruction.op_code() == OpCode.ROR:
            result = ALU.ror(ra, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_imm_shift_nz_cc(ra, result)

        self._cur_thread.register_file().write(instruction.rc(), result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_sub_rrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.SubRRICIOpCodes
        assert instruction.suffix() == Suffix.RRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.SUB:
            result, carry, overflow = ALU.sub(ra, imm)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, overflow = ALU.subc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_sub_nz_cc(ra, imm, result, carry, overflow)

        self._cur_thread.register_file().write(instruction.rc(), result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_rrif(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRIFOpCodes
        assert instruction.suffix() == Suffix.RRIF
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ADD:
            result, carry, _ = ALU.add(ra, imm)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, _ = ALU.addc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.AND:
            result, carry = ALU.and_(ra, imm), False
        elif instruction.op_code() == OpCode.ANDN:
            result, carry = ALU.andn(ra, imm), False
        elif instruction.op_code() == OpCode.NAND:
            result, carry = ALU.nand(ra, imm), False
        elif instruction.op_code() == OpCode.NOR:
            result, carry = ALU.nor(ra, imm), False
        elif instruction.op_code() == OpCode.NXOR:
            result, carry = ALU.nxor(ra, imm), False
        elif instruction.op_code() == OpCode.OR:
            result, carry = ALU.or_(ra, imm), False
        elif instruction.op_code() == OpCode.ORN:
            result, carry = ALU.orn(ra, imm), False
        elif instruction.op_code() == OpCode.SUB:
            result, carry, _ = ALU.sub(ra, imm)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, _ = ALU.subc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.XOR:
            result, carry = ALU.xor(ra, imm), False
        elif instruction.op_code() == OpCode.HASH:
            result, carry = ALU.hash(ra, imm), False
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.rc(), result)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_rrr(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRROpCodes
        assert instruction.suffix() == Suffix.RRR
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.ADD:
            result, carry, _ = ALU.add(ra, rb)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, _ = ALU.addc(ra, rb, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.AND:
            result, carry = ALU.and_(ra, rb), False
        elif instruction.op_code() == OpCode.ANDN:
            result, carry = ALU.andn(ra, rb), False
        elif instruction.op_code() == OpCode.ASR:
            result, carry = ALU.asr(ra, rb), False
        elif instruction.op_code() == OpCode.CMPB4:
            result, carry = ALU.cmpb4(ra, rb), False
        elif instruction.op_code() == OpCode.LSL:
            result, carry = ALU.lsl(ra, rb), False
        elif instruction.op_code() == OpCode.LSL1:
            result, carry = ALU.lsl1(ra, rb), False
        elif instruction.op_code() == OpCode.LSL1X:
            result, carry = ALU.lsl1x(ra, rb), False
        elif instruction.op_code() == OpCode.LSLX:
            result, carry = ALU.lslx(ra, rb), False
        elif instruction.op_code() == OpCode.LSR:
            result, carry = ALU.lsr(ra, rb), False
        elif instruction.op_code() == OpCode.LSR1:
            result, carry = ALU.lsr1(ra, rb), False
        elif instruction.op_code() == OpCode.LSR1X:
            result, carry = ALU.lsr1x(ra, rb), False
        elif instruction.op_code() == OpCode.LSRX:
            result, carry = ALU.lsrx(ra, rb), False
        elif instruction.op_code() == OpCode.ROL:
            result, carry = ALU.rol(ra, rb), False
        elif instruction.op_code() == OpCode.ROR:
            result, carry = ALU.ror(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_SH:
            result, carry = ALU.mul_sh_sh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_SL:
            result, carry = ALU.mul_sh_sl(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_UH:
            result, carry = ALU.mul_sh_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_UL:
            result, carry = ALU.mul_sh_ul(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_SH:
            result, carry = ALU.mul_sl_sh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_SL:
            result, carry = ALU.mul_sl_sl(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_UH:
            result, carry = ALU.mul_sl_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_UL:
            result, carry = ALU.mul_sl_ul(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UH_UH:
            result, carry = ALU.mul_uh_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UH_UL:
            result, carry = ALU.mul_uh_ul(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UL_UH:
            result, carry = ALU.mul_ul_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UL_UL:
            result, carry = ALU.mul_ul_ul(ra, rb), False
        elif instruction.op_code() == OpCode.NAND:
            result, carry = ALU.nand(ra, rb), False
        elif instruction.op_code() == OpCode.NOR:
            result, carry = ALU.nor(ra, rb), False
        elif instruction.op_code() == OpCode.NXOR:
            result, carry = ALU.nxor(ra, rb), False
        elif instruction.op_code() == OpCode.OR:
            result, carry = ALU.or_(ra, rb), False
        elif instruction.op_code() == OpCode.ORN:
            result, carry = ALU.orn(ra, rb), False
        elif instruction.op_code() == OpCode.ROL:
            result, carry = ALU.rol(ra, rb), False
        elif instruction.op_code() == OpCode.ROR:
            result, carry = ALU.ror(ra, rb), False
        elif instruction.op_code() == OpCode.RSUB:
            result, carry, _ = ALU.sub(rb, ra)
        elif instruction.op_code() == OpCode.RSUBC:
            result, carry, _ = ALU.subc(rb, ra, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.SUB:
            result, carry, _ = ALU.sub(ra, rb)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, _ = ALU.subc(ra, rb, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.XOR:
            result, carry = ALU.xor(ra, rb), False
        elif instruction.op_code() == OpCode.HASH:
            result, carry = ALU.hash(ra, rb), False
        elif instruction.op_code() == OpCode.CALL:
            result, carry, _ = ALU.add(ra, rb)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        if instruction.op_code() == OpCode.CALL:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write(instruction.rc(), pc + InstructionWord().size())

            self._cur_thread.register_file().write_pc(result)
        else:
            self._cur_thread.register_file().write(instruction.rc(), result)

            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_rrrc(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRRCOpCodes

        if instruction.op_code() in Instruction.AddRRRCOpCodes:
            self._execute_add_rrrc(instruction)
        elif instruction.op_code() in Instruction.RsubRRRCOpCodes:
            self._execute_rsub_rrrc(instruction)
        elif instruction.op_code() in Instruction.SubRRRCOpCodes:
            self._execute_sub_rrrc(instruction)
        else:
            raise ValueError

    def _execute_add_rrrc(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AddRRRCOpCodes
        assert instruction.suffix() == Suffix.RRRC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.ADD:
            result, carry, _ = ALU.add(ra, rb)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, _ = ALU.addc(ra, rb, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.AND:
            result, carry = ALU.and_(ra, rb), False
        elif instruction.op_code() == OpCode.ANDN:
            result, carry = ALU.andn(ra, rb), False
        elif instruction.op_code() == OpCode.ASR:
            result, carry = ALU.asr(ra, rb), False
        elif instruction.op_code() == OpCode.CMPB4:
            result, carry = ALU.cmpb4(ra, rb), False
        elif instruction.op_code() == OpCode.LSL:
            result, carry = ALU.lsl(ra, rb), False
        elif instruction.op_code() == OpCode.LSL1:
            result, carry = ALU.lsl1(ra, rb), False
        elif instruction.op_code() == OpCode.LSL1X:
            result, carry = ALU.lsl1x(ra, rb), False
        elif instruction.op_code() == OpCode.LSLX:
            result, carry = ALU.lslx(ra, rb), False
        elif instruction.op_code() == OpCode.LSR:
            result, carry = ALU.lsr(ra, rb), False
        elif instruction.op_code() == OpCode.LSR1:
            result, carry = ALU.lsr1(ra, rb), False
        elif instruction.op_code() == OpCode.LSR1X:
            result, carry = ALU.lsr1x(ra, rb), False
        elif instruction.op_code() == OpCode.LSRX:
            result, carry = ALU.lsrx(ra, rb), False
        elif instruction.op_code() == OpCode.ROL:
            result, carry = ALU.rol(ra, rb), False
        elif instruction.op_code() == OpCode.ROR:
            result, carry = ALU.ror(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_SH:
            result, carry = ALU.mul_sh_sh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_SL:
            result, carry = ALU.mul_sh_sl(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_UH:
            result, carry = ALU.mul_sh_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_UL:
            result, carry = ALU.mul_sh_ul(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_SH:
            result, carry = ALU.mul_sl_sh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_SL:
            result, carry = ALU.mul_sl_sl(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_UH:
            result, carry = ALU.mul_sl_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_UL:
            result, carry = ALU.mul_sl_ul(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UH_UH:
            result, carry = ALU.mul_uh_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UH_UL:
            result, carry = ALU.mul_uh_ul(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UL_UH:
            result, carry = ALU.mul_ul_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UL_UL:
            result, carry = ALU.mul_ul_ul(ra, rb), False
        elif instruction.op_code() == OpCode.NAND:
            result, carry = ALU.nand(ra, rb), False
        elif instruction.op_code() == OpCode.NOR:
            result, carry = ALU.nor(ra, rb), False
        elif instruction.op_code() == OpCode.NXOR:
            result, carry = ALU.nxor(ra, rb), False
        elif instruction.op_code() == OpCode.OR:
            result, carry = ALU.or_(ra, rb), False
        elif instruction.op_code() == OpCode.ORN:
            result, carry = ALU.orn(ra, rb), False
        elif instruction.op_code() == OpCode.ROL:
            result, carry = ALU.rol(ra, rb), False
        elif instruction.op_code() == OpCode.ROR:
            result, carry = ALU.ror(ra, rb), False
        elif instruction.op_code() == OpCode.XOR:
            result, carry = ALU.xor(ra, rb), False
        elif instruction.op_code() == OpCode.HASH:
            result, carry = ALU.hash(ra, rb), False
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_set_cc(result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            self._cur_thread.register_file().write(instruction.rc(), 1)
        else:
            self._cur_thread.register_file().write(instruction.rc(), 0)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_rsub_rrrc(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RsubRRRCOpCodes
        assert instruction.suffix() == Suffix.RRRC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.RSUB:
            result, carry, _ = ALU.sub(rb, ra)
        elif instruction.op_code() == OpCode.RSUBC:
            result, carry, _ = ALU.subc(rb, ra, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_sub_set_cc(ra, rb, result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            self._cur_thread.register_file().write(instruction.rc(), 1)
        else:
            self._cur_thread.register_file().write(instruction.rc(), 0)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_sub_rrrc(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.SubRRRCOpCodes
        assert instruction.suffix() == Suffix.RRRC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.SUB:
            result, carry, overflow = ALU.sub(ra, rb)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, overflow = ALU.subc(ra, rb, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_ext_sub_set_cc(ra, rb, result, carry, overflow)

        if self._cur_thread.register_file().condition(instruction.condition()):
            self._cur_thread.register_file().write(instruction.rc(), 1)
        else:
            self._cur_thread.register_file().write(instruction.rc(), 0)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_rrrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRRCIOpCodes

        if instruction.op_code() in Instruction.AddRRRCIOpCodes:
            self._execute_add_rrrci(instruction)
        elif instruction.op_code() in Instruction.AndRRRCIOpCodes:
            self._execute_and_rrrci(instruction)
        elif instruction.op_code() in Instruction.AsrRRRCIOpCodes:
            self._execute_asr_rrrci(instruction)
        elif instruction.op_code() in Instruction.MulRRRCIOpCodes:
            self._execute_mul_rrrci(instruction)
        elif instruction.op_code() in Instruction.RsubRRRCIOpCodes:
            self._execute_rsub_rrrci(instruction)
        else:
            raise ValueError

    def _execute_add_rrrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AddRRRCIOpCodes
        assert instruction.suffix() == Suffix.RRRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.ADD:
            result, carry, overflow = ALU.add(ra, rb)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, overflow = ALU.addc(ra, rb, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_add_nz_cc(ra, result, carry, overflow)

        self._cur_thread.register_file().write(instruction.rc(), result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_and_rrrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AndRRRCIOpCodes
        assert instruction.suffix() == Suffix.RRRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.AND:
            result = ALU.and_(ra, rb)
        elif instruction.op_code() == OpCode.ANDN:
            result = ALU.andn(ra, rb)
        elif instruction.op_code() == OpCode.NAND:
            result = ALU.nand(ra, rb)
        elif instruction.op_code() == OpCode.NOR:
            result = ALU.nor(ra, rb)
        elif instruction.op_code() == OpCode.NXOR:
            result = ALU.nxor(ra, rb)
        elif instruction.op_code() == OpCode.OR:
            result = ALU.or_(ra, rb)
        elif instruction.op_code() == OpCode.ORN:
            result = ALU.orn(ra, rb)
        elif instruction.op_code() == OpCode.XOR:
            result = ALU.xor(ra, rb)
        elif instruction.op_code() == OpCode.HASH:
            result = ALU.hash(ra, rb)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_nz_cc(ra, result)

        self._cur_thread.register_file().write(instruction.rc(), result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_asr_rrrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AsrRRRCIOpCodes
        assert instruction.suffix() == Suffix.RRRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.ASR:
            result = ALU.asr(ra, rb)
        elif instruction.op_code() == OpCode.CMPB4:
            result = ALU.cmpb4(ra, rb)
        elif instruction.op_code() == OpCode.LSL:
            result = ALU.lsl(ra, rb)
        elif instruction.op_code() == OpCode.LSL1:
            result = ALU.lsl1(ra, rb)
        elif instruction.op_code() == OpCode.LSL1X:
            result = ALU.lsl1x(ra, rb)
        elif instruction.op_code() == OpCode.LSLX:
            result = ALU.lslx(ra, rb)
        elif instruction.op_code() == OpCode.LSR:
            result = ALU.lsr(ra, rb)
        elif instruction.op_code() == OpCode.LSR1:
            result = ALU.lsr1(ra, rb)
        elif instruction.op_code() == OpCode.LSR1X:
            result = ALU.lsr1x(ra, rb)
        elif instruction.op_code() == OpCode.LSRX:
            result = ALU.lsrx(ra, rb)
        elif instruction.op_code() == OpCode.ROL:
            result = ALU.rol(ra, rb)
        elif instruction.op_code() == OpCode.ROR:
            result = ALU.ror(ra, rb)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_shift_nz_cc(ra, result)

        self._cur_thread.register_file().write(instruction.rc(), result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_mul_rrrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.MulRRRCIOpCodes
        assert instruction.suffix() == Suffix.RRRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.MUL_SH_SH:
            result = ALU.mul_sh_sh(ra, rb)
        elif instruction.op_code() == OpCode.MUL_SH_SL:
            result = ALU.mul_sh_sl(ra, rb)
        elif instruction.op_code() == OpCode.MUL_SH_UH:
            result = ALU.mul_sh_uh(ra, rb)
        elif instruction.op_code() == OpCode.MUL_SH_UL:
            result = ALU.mul_sh_ul(ra, rb)
        elif instruction.op_code() == OpCode.MUL_SL_SH:
            result = ALU.mul_sl_sh(ra, rb)
        elif instruction.op_code() == OpCode.MUL_SL_SL:
            result = ALU.mul_sl_sl(ra, rb)
        elif instruction.op_code() == OpCode.MUL_SL_UH:
            result = ALU.mul_sl_uh(ra, rb)
        elif instruction.op_code() == OpCode.MUL_SL_UL:
            result = ALU.mul_sl_ul(ra, rb)
        elif instruction.op_code() == OpCode.MUL_UH_UH:
            result = ALU.mul_uh_uh(ra, rb)
        elif instruction.op_code() == OpCode.MUL_UH_UL:
            result = ALU.mul_uh_ul(ra, rb)
        elif instruction.op_code() == OpCode.MUL_UL_UH:
            result = ALU.mul_ul_uh(ra, rb)
        elif instruction.op_code() == OpCode.MUL_UL_UL:
            result = ALU.mul_ul_ul(ra, rb)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_mul_nz_cc(ra, result)

        self._cur_thread.register_file().write(instruction.rc(), result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_rsub_rrrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RsubRRRCIOpCodes
        assert instruction.suffix() == Suffix.RRRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.RSUB:
            result, carry, overflow = ALU.sub(rb, ra)
        elif instruction.op_code() == OpCode.RSUBC:
            result, carry, overflow = ALU.subc(rb, ra, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.SUB:
            result, carry, overflow = ALU.sub(ra, rb)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, overflow = ALU.subc(ra, rb, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_sub_nz_cc(ra, rb, result, carry, overflow)

        self._cur_thread.register_file().write(instruction.rc(), result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_zri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRIOpCodes

        if instruction.op_code() in Instruction.AddRRICIOpCodes:
            self._execute_add_zri(instruction)
        elif instruction.op_code() in Instruction.AsrRRIOpCodes:
            self._execute_asr_zri(instruction)
        elif instruction.op_code() in Instruction.CallRRIOpCodes:
            self._execute_call_zri(instruction)
        else:
            raise ValueError

    def _execute_add_zri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AddRRIOpCodes
        assert instruction.suffix() == Suffix.ZRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ADD:
            result, carry, _ = ALU.add(ra, imm)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, _ = ALU.addc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.AND:
            result, carry = ALU.and_(ra, imm), False
        elif instruction.op_code() == OpCode.OR:
            result, carry = ALU.or_(ra, imm), False
        elif instruction.op_code() == OpCode.XOR:
            result, carry = ALU.xor(ra, imm), False
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_asr_zri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AsrRRIOpCodes
        assert instruction.suffix() == Suffix.ZRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ASR:
            result = ALU.asr(ra, imm)
        elif instruction.op_code() == OpCode.LSL:
            result = ALU.lsl(ra, imm)
        elif instruction.op_code() == OpCode.LSL1:
            result = ALU.lsl1(ra, imm)
        elif instruction.op_code() == OpCode.LSL1X:
            result = ALU.lsl1x(ra, imm)
        elif instruction.op_code() == OpCode.LSLX:
            result = ALU.lslx(ra, imm)
        elif instruction.op_code() == OpCode.LSR:
            result = ALU.lsr(ra, imm)
        elif instruction.op_code() == OpCode.LSR1:
            result = ALU.lsr1(ra, imm)
        elif instruction.op_code() == OpCode.LSR1X:
            result = ALU.lsr1x(ra, imm)
        elif instruction.op_code() == OpCode.LSRX:
            result = ALU.lsrx(ra, imm)
        elif instruction.op_code() == OpCode.ROL:
            result = ALU.rol(ra, imm)
        elif instruction.op_code() == OpCode.ROR:
            result = ALU.ror(ra, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_call_zri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.CallRRIOpCodes
        assert instruction.suffix() == Suffix.ZRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if imm == 0:
            callee_address, carry, _ = ALU.add(ra, imm)
        else:
            callee_address, carry, _ = ALU.add(ra * InstructionWord().size(), imm)

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write_pc(callee_address)

        self._set_flags(callee_address, carry)

    def _execute_zric(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRICOpCodes

        if instruction.op_code() in Instruction.AddRRICOpCodes:
            self._execute_add_zric(instruction)
        elif instruction.op_code() in Instruction.AsrRRICOpCodes:
            self._execute_asr_zric(instruction)
        elif instruction.op_code() in Instruction.SubRRICOpCodes:
            self._execute_sub_zric(instruction)
        else:
            raise ValueError

    def _execute_add_zric(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AddRRICOpCodes
        assert instruction.suffix() == Suffix.ZRIC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ADD:
            result, carry, _ = ALU.add(ra, imm)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, _ = ALU.addc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.AND:
            result, carry = ALU.and_(ra, imm), False
        elif instruction.op_code() == OpCode.ANDN:
            result, carry = ALU.andn(ra, imm), False
        elif instruction.op_code() == OpCode.NAND:
            result, carry = ALU.nand(ra, imm), False
        elif instruction.op_code() == OpCode.NOR:
            result, carry = ALU.nor(ra, imm), False
        elif instruction.op_code() == OpCode.NXOR:
            result, carry = ALU.nxor(ra, imm), False
        elif instruction.op_code() == OpCode.OR:
            result, carry = ALU.or_(ra, imm), False
        elif instruction.op_code() == OpCode.ORN:
            result, carry = ALU.orn(ra, imm), False
        elif instruction.op_code() == OpCode.XOR:
            result, carry = ALU.xor(ra, imm), False
        elif instruction.op_code() == OpCode.HASH:
            result, carry = ALU.hash(ra, imm), False
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_set_cc(result)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_asr_zric(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AsrRRICOpCodes
        assert instruction.suffix() == Suffix.ZRIC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ASR:
            result = ALU.asr(ra, imm)
        elif instruction.op_code() == OpCode.LSL:
            result = ALU.lsl(ra, imm)
        elif instruction.op_code() == OpCode.LSL1:
            result = ALU.lsl1(ra, imm)
        elif instruction.op_code() == OpCode.LSL1X:
            result = ALU.lsl1x(ra, imm)
        elif instruction.op_code() == OpCode.LSLX:
            result = ALU.lslx(ra, imm)
        elif instruction.op_code() == OpCode.LSR:
            result = ALU.lsr(ra, imm)
        elif instruction.op_code() == OpCode.LSR1:
            result = ALU.lsr1(ra, imm)
        elif instruction.op_code() == OpCode.LSR1X:
            result = ALU.lsr1x(ra, imm)
        elif instruction.op_code() == OpCode.LSRX:
            result = ALU.lsrx(ra, imm)
        elif instruction.op_code() == OpCode.ROL:
            result = ALU.rol(ra, imm)
        elif instruction.op_code() == OpCode.ROR:
            result = ALU.ror(ra, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_set_cc(result)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_sub_zric(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.SubRRICOpCodes
        assert instruction.suffix() == Suffix.RRIC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.SUB:
            result, carry, overflow = ALU.sub(ra, imm)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, overflow = ALU.subc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_ext_sub_set_cc(ra, imm, result, carry, overflow)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_zrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRICIOpCodes

        if instruction.op_code() in Instruction.AddRRICIOpCodes:
            self._execute_add_zrici(instruction)
        elif instruction.op_code() in Instruction.AndRRICIOpCodes:
            self._execute_and_zrici(instruction)
        elif instruction.op_code() in Instruction.AsrRRICIOpCodes:
            self._execute_asr_zrici(instruction)
        elif instruction.op_code() in Instruction.SubRRICIOpCodes:
            self._execute_sub_zrici(instruction)
        else:
            raise ValueError

    def _execute_add_zrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AddRRICIOpCodes
        assert instruction.suffix() == Suffix.ZRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ADD:
            result, carry, overflow = ALU.add(ra, imm)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, overflow = ALU.addc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_add_nz_cc(ra, result, carry, overflow)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_and_zrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AndRRICIOpCodes
        assert instruction.suffix() == Suffix.ZRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.AND:
            result = ALU.and_(ra, imm)
        elif instruction.op_code() == OpCode.ANDN:
            result = ALU.andn(ra, imm)
        elif instruction.op_code() == OpCode.NAND:
            result = ALU.nand(ra, imm)
        elif instruction.op_code() == OpCode.NOR:
            result = ALU.nor(ra, imm)
        elif instruction.op_code() == OpCode.NXOR:
            result = ALU.nxor(ra, imm)
        elif instruction.op_code() == OpCode.OR:
            result = ALU.or_(ra, imm)
        elif instruction.op_code() == OpCode.ORN:
            result = ALU.orn(ra, imm)
        elif instruction.op_code() == OpCode.XOR:
            result = ALU.xor(ra, imm)
        elif instruction.op_code() == OpCode.HASH:
            result = ALU.hash(ra, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_nz_cc(ra, result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_asr_zrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AsrRRICIOpCodes
        assert instruction.suffix() == Suffix.ZRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ASR:
            result = ALU.asr(ra, imm)
        elif instruction.op_code() == OpCode.LSL:
            result = ALU.lsl(ra, imm)
        elif instruction.op_code() == OpCode.LSL1:
            result = ALU.lsl1(ra, imm)
        elif instruction.op_code() == OpCode.LSL1X:
            result = ALU.lsl1x(ra, imm)
        elif instruction.op_code() == OpCode.LSLX:
            result = ALU.lslx(ra, imm)
        elif instruction.op_code() == OpCode.LSR:
            result = ALU.lsr(ra, imm)
        elif instruction.op_code() == OpCode.LSR1:
            result = ALU.lsr1(ra, imm)
        elif instruction.op_code() == OpCode.LSR1X:
            result = ALU.lsr1x(ra, imm)
        elif instruction.op_code() == OpCode.LSRX:
            result = ALU.lsrx(ra, imm)
        elif instruction.op_code() == OpCode.ROL:
            result = ALU.rol(ra, imm)
        elif instruction.op_code() == OpCode.ROR:
            result = ALU.ror(ra, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_imm_shift_nz_cc(ra, result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_sub_zrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.SubRRICIOpCodes
        assert instruction.suffix() == Suffix.ZRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.SUB:
            result, carry, overflow = ALU.sub(ra, imm)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, overflow = ALU.subc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_sub_nz_cc(ra, imm, result, carry, overflow)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_zrif(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRIFOpCodes
        assert instruction.suffix() == Suffix.ZRIF
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ADD:
            result, carry, _ = ALU.add(ra, imm)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, _ = ALU.addc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.AND:
            result, carry = ALU.and_(ra, imm), False
        elif instruction.op_code() == OpCode.ANDN:
            result, carry = ALU.andn(ra, imm), False
        elif instruction.op_code() == OpCode.NAND:
            result, carry = ALU.nand(ra, imm), False
        elif instruction.op_code() == OpCode.NOR:
            result, carry = ALU.nor(ra, imm), False
        elif instruction.op_code() == OpCode.NXOR:
            result, carry = ALU.nxor(ra, imm), False
        elif instruction.op_code() == OpCode.OR:
            result, carry = ALU.or_(ra, imm), False
        elif instruction.op_code() == OpCode.ORN:
            result, carry = ALU.orn(ra, imm), False
        elif instruction.op_code() == OpCode.SUB:
            result, carry, _ = ALU.sub(ra, imm)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, _ = ALU.subc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.XOR:
            result, carry = ALU.xor(ra, imm), False
        elif instruction.op_code() == OpCode.HASH:
            result, carry = ALU.hash(ra, imm), False
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_zrr(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRROpCodes
        assert instruction.suffix() == Suffix.ZRR
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.ADD:
            result, carry, _ = ALU.add(ra, rb)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, _ = ALU.addc(ra, rb, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.AND:
            result, carry = ALU.and_(ra, rb), False
        elif instruction.op_code() == OpCode.ANDN:
            result, carry = ALU.andn(ra, rb), False
        elif instruction.op_code() == OpCode.ASR:
            result, carry = ALU.asr(ra, rb), False
        elif instruction.op_code() == OpCode.CMPB4:
            result, carry = ALU.cmpb4(ra, rb), False
        elif instruction.op_code() == OpCode.LSL:
            result, carry = ALU.lsl(ra, rb), False
        elif instruction.op_code() == OpCode.LSL1:
            result, carry = ALU.lsl1(ra, rb), False
        elif instruction.op_code() == OpCode.LSL1X:
            result, carry = ALU.lsl1x(ra, rb), False
        elif instruction.op_code() == OpCode.LSLX:
            result, carry = ALU.lslx(ra, rb), False
        elif instruction.op_code() == OpCode.LSR:
            result, carry = ALU.lsr(ra, rb), False
        elif instruction.op_code() == OpCode.LSR1:
            result, carry = ALU.lsr1(ra, rb), False
        elif instruction.op_code() == OpCode.LSR1X:
            result, carry = ALU.lsr1x(ra, rb), False
        elif instruction.op_code() == OpCode.LSRX:
            result, carry = ALU.lsrx(ra, rb), False
        elif instruction.op_code() == OpCode.ROL:
            result, carry = ALU.rol(ra, rb), False
        elif instruction.op_code() == OpCode.ROR:
            result, carry = ALU.ror(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_SH:
            result, carry = ALU.mul_sh_sh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_SL:
            result, carry = ALU.mul_sh_sl(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_UH:
            result, carry = ALU.mul_sh_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_UL:
            result, carry = ALU.mul_sh_ul(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_SH:
            result, carry = ALU.mul_sl_sh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_SL:
            result, carry = ALU.mul_sl_sl(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_UH:
            result, carry = ALU.mul_sl_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_UL:
            result, carry = ALU.mul_sl_ul(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UH_UH:
            result, carry = ALU.mul_uh_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UH_UL:
            result, carry = ALU.mul_uh_ul(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UL_UH:
            result, carry = ALU.mul_ul_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UL_UL:
            result, carry = ALU.mul_ul_ul(ra, rb), False
        elif instruction.op_code() == OpCode.NAND:
            result, carry = ALU.nand(ra, rb), False
        elif instruction.op_code() == OpCode.NOR:
            result, carry = ALU.nor(ra, rb), False
        elif instruction.op_code() == OpCode.NXOR:
            result, carry = ALU.nxor(ra, rb), False
        elif instruction.op_code() == OpCode.OR:
            result, carry = ALU.or_(ra, rb), False
        elif instruction.op_code() == OpCode.ORN:
            result, carry = ALU.orn(ra, rb), False
        elif instruction.op_code() == OpCode.ROL:
            result, carry = ALU.rol(ra, rb), False
        elif instruction.op_code() == OpCode.ROR:
            result, carry = ALU.ror(ra, rb), False
        elif instruction.op_code() == OpCode.RSUB:
            result, carry, _ = ALU.sub(rb, ra)
        elif instruction.op_code() == OpCode.RSUBC:
            result, carry, _ = ALU.subc(rb, ra, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.SUB:
            result, carry, _ = ALU.sub(ra, rb)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, _ = ALU.subc(ra, rb, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.XOR:
            result, carry = ALU.xor(ra, rb), False
        elif instruction.op_code() == OpCode.HASH:
            result, carry = ALU.hash(ra, rb), False
        elif instruction.op_code() == OpCode.CALL:
            result, carry, _ = ALU.add(ra, rb)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        if instruction.op_code() == OpCode.CALL:
            self._cur_thread.register_file().write_pc(result)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_zrrc(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRRCOpCodes

        if instruction.op_code() in Instruction.AddRRRCOpCodes:
            self._execute_add_zrrc(instruction)
        elif instruction.op_code() in Instruction.RsubRRRCOpCodes:
            self._execute_rsub_zrrc(instruction)
        elif instruction.op_code() in Instruction.SubRRRCOpCodes:
            self._execute_sub_zrrc(instruction)
        else:
            raise ValueError

    def _execute_add_zrrc(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AddRRRCOpCodes
        assert instruction.suffix() == Suffix.ZRRC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.ADD:
            result, carry, _ = ALU.add(ra, rb)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, _ = ALU.addc(ra, rb, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.AND:
            result, carry = ALU.and_(ra, rb), False
        elif instruction.op_code() == OpCode.ANDN:
            result, carry = ALU.andn(ra, rb), False
        elif instruction.op_code() == OpCode.ASR:
            result, carry = ALU.asr(ra, rb), False
        elif instruction.op_code() == OpCode.CMPB4:
            result, carry = ALU.cmpb4(ra, rb), False
        elif instruction.op_code() == OpCode.LSL:
            result, carry = ALU.lsl(ra, rb), False
        elif instruction.op_code() == OpCode.LSL1:
            result, carry = ALU.lsl1(ra, rb), False
        elif instruction.op_code() == OpCode.LSL1X:
            result, carry = ALU.lsl1x(ra, rb), False
        elif instruction.op_code() == OpCode.LSLX:
            result, carry = ALU.lslx(ra, rb), False
        elif instruction.op_code() == OpCode.LSR:
            result, carry = ALU.lsr(ra, rb), False
        elif instruction.op_code() == OpCode.LSR1:
            result, carry = ALU.lsr1(ra, rb), False
        elif instruction.op_code() == OpCode.LSR1X:
            result, carry = ALU.lsr1x(ra, rb), False
        elif instruction.op_code() == OpCode.LSRX:
            result, carry = ALU.lsrx(ra, rb), False
        elif instruction.op_code() == OpCode.ROL:
            result, carry = ALU.rol(ra, rb), False
        elif instruction.op_code() == OpCode.ROR:
            result, carry = ALU.ror(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_SH:
            result, carry = ALU.mul_sh_sh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_SL:
            result, carry = ALU.mul_sh_sl(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_UH:
            result, carry = ALU.mul_sh_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SH_UL:
            result, carry = ALU.mul_sh_ul(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_SH:
            result, carry = ALU.mul_sl_sh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_SL:
            result, carry = ALU.mul_sl_sl(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_UH:
            result, carry = ALU.mul_sl_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_SL_UL:
            result, carry = ALU.mul_sl_ul(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UH_UH:
            result, carry = ALU.mul_uh_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UH_UL:
            result, carry = ALU.mul_uh_ul(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UL_UH:
            result, carry = ALU.mul_ul_uh(ra, rb), False
        elif instruction.op_code() == OpCode.MUL_UL_UL:
            result, carry = ALU.mul_ul_ul(ra, rb), False
        elif instruction.op_code() == OpCode.NAND:
            result, carry = ALU.nand(ra, rb), False
        elif instruction.op_code() == OpCode.NOR:
            result, carry = ALU.nor(ra, rb), False
        elif instruction.op_code() == OpCode.NXOR:
            result, carry = ALU.nxor(ra, rb), False
        elif instruction.op_code() == OpCode.OR:
            result, carry = ALU.or_(ra, rb), False
        elif instruction.op_code() == OpCode.ORN:
            result, carry = ALU.orn(ra, rb), False
        elif instruction.op_code() == OpCode.ROL:
            result, carry = ALU.rol(ra, rb), False
        elif instruction.op_code() == OpCode.ROR:
            result, carry = ALU.ror(ra, rb), False
        elif instruction.op_code() == OpCode.XOR:
            result, carry = ALU.xor(ra, rb), False
        elif instruction.op_code() == OpCode.HASH:
            result, carry = ALU.hash(ra, rb), False
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_set_cc(result)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_rsub_zrrc(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RsubRRRCOpCodes
        assert instruction.suffix() == Suffix.ZRRC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.RSUB:
            result, carry, _ = ALU.sub(rb, ra)
        elif instruction.op_code() == OpCode.RSUBC:
            result, carry, _ = ALU.subc(rb, ra, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_sub_set_cc(ra, rb, result)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_sub_zrrc(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.SubRRRCOpCodes
        assert instruction.suffix() == Suffix.ZRRC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.SUB:
            result, carry, overflow = ALU.sub(ra, rb)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, overflow = ALU.subc(ra, rb, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_ext_sub_set_cc(ra, rb, result, carry, overflow)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_zrrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRRCIOpCodes

        if instruction.op_code() in Instruction.AddRRRCIOpCodes:
            self._execute_add_zrrci(instruction)
        elif instruction.op_code() in Instruction.AndRRRCIOpCodes:
            self._execute_and_zrrci(instruction)
        elif instruction.op_code() in Instruction.AsrRRRCIOpCodes:
            self._execute_asr_zrrci(instruction)
        elif instruction.op_code() in Instruction.MulRRRCIOpCodes:
            self._execute_mul_zrrci(instruction)
        elif instruction.op_code() in Instruction.RsubRRRCIOpCodes:
            self._execute_rsub_zrrci(instruction)
        else:
            raise ValueError

    def _execute_add_zrrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AddRRRCIOpCodes
        assert instruction.suffix() == Suffix.ZRRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.ADD:
            result, carry, overflow = ALU.add(ra, rb)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, overflow = ALU.addc(ra, rb, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_add_nz_cc(ra, result, carry, overflow)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_and_zrrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AndRRRCIOpCodes
        assert instruction.suffix() == Suffix.ZRRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.AND:
            result = ALU.and_(ra, rb)
        elif instruction.op_code() == OpCode.ANDN:
            result = ALU.andn(ra, rb)
        elif instruction.op_code() == OpCode.NAND:
            result = ALU.nand(ra, rb)
        elif instruction.op_code() == OpCode.NOR:
            result = ALU.nor(ra, rb)
        elif instruction.op_code() == OpCode.NXOR:
            result = ALU.nxor(ra, rb)
        elif instruction.op_code() == OpCode.OR:
            result = ALU.or_(ra, rb)
        elif instruction.op_code() == OpCode.ORN:
            result = ALU.orn(ra, rb)
        elif instruction.op_code() == OpCode.XOR:
            result = ALU.xor(ra, rb)
        elif instruction.op_code() == OpCode.HASH:
            result = ALU.hash(ra, rb)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_nz_cc(ra, result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_asr_zrrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AsrRRRCIOpCodes
        assert instruction.suffix() == Suffix.ZRRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.ASR:
            result = ALU.asr(ra, rb)
        elif instruction.op_code() == OpCode.CMPB4:
            result = ALU.cmpb4(ra, rb)
        elif instruction.op_code() == OpCode.LSL:
            result = ALU.lsl(ra, rb)
        elif instruction.op_code() == OpCode.LSL1:
            result = ALU.lsl1(ra, rb)
        elif instruction.op_code() == OpCode.LSL1X:
            result = ALU.lsl1x(ra, rb)
        elif instruction.op_code() == OpCode.LSLX:
            result = ALU.lslx(ra, rb)
        elif instruction.op_code() == OpCode.LSR:
            result = ALU.lsr(ra, rb)
        elif instruction.op_code() == OpCode.LSR1:
            result = ALU.lsr1(ra, rb)
        elif instruction.op_code() == OpCode.LSR1X:
            result = ALU.lsr1x(ra, rb)
        elif instruction.op_code() == OpCode.LSRX:
            result = ALU.lsrx(ra, rb)
        elif instruction.op_code() == OpCode.ROL:
            result = ALU.rol(ra, rb)
        elif instruction.op_code() == OpCode.ROR:
            result = ALU.ror(ra, rb)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_shift_nz_cc(ra, result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_mul_zrrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.MulRRRCIOpCodes
        assert instruction.suffix() == Suffix.ZRRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.MUL_SH_SH:
            result = ALU.mul_sh_sh(ra, rb)
        elif instruction.op_code() == OpCode.MUL_SH_SL:
            result = ALU.mul_sh_sl(ra, rb)
        elif instruction.op_code() == OpCode.MUL_SH_UH:
            result = ALU.mul_sh_uh(ra, rb)
        elif instruction.op_code() == OpCode.MUL_SH_UL:
            result = ALU.mul_sh_ul(ra, rb)
        elif instruction.op_code() == OpCode.MUL_SL_SH:
            result = ALU.mul_sl_sh(ra, rb)
        elif instruction.op_code() == OpCode.MUL_SL_SL:
            result = ALU.mul_sl_sl(ra, rb)
        elif instruction.op_code() == OpCode.MUL_SL_UH:
            result = ALU.mul_sl_uh(ra, rb)
        elif instruction.op_code() == OpCode.MUL_SL_UL:
            result = ALU.mul_sl_ul(ra, rb)
        elif instruction.op_code() == OpCode.MUL_UH_UH:
            result = ALU.mul_uh_uh(ra, rb)
        elif instruction.op_code() == OpCode.MUL_UH_UL:
            result = ALU.mul_uh_ul(ra, rb)
        elif instruction.op_code() == OpCode.MUL_UL_UH:
            result = ALU.mul_ul_uh(ra, rb)
        elif instruction.op_code() == OpCode.MUL_UL_UL:
            result = ALU.mul_ul_ul(ra, rb)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_mul_nz_cc(ra, result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_rsub_zrrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RsubRRRCIOpCodes
        assert instruction.suffix() == Suffix.ZRRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.RSUB:
            result, carry, overflow = ALU.sub(rb, ra)
        elif instruction.op_code() == OpCode.RSUBC:
            result, carry, _ = ALU.subc(rb, ra, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.SUB:
            result, carry, overflow = ALU.sub(ra, rb)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, overflow = ALU.subc(ra, rb, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_sub_nz_cc(ra, rb, result, carry, overflow)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_s_rri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRIOpCodes

        if instruction.op_code() in Instruction.AddRRIOpCodes:
            self._execute_add_s_rri(instruction)
        elif instruction.op_code() in Instruction.AsrRRIOpCodes:
            self._execute_asr_s_rri(instruction)
        else:
            raise ValueError

    def _execute_add_s_rri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AddRRIOpCodes
        assert instruction.suffix() == Suffix.S_RRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ADD:
            result, carry, _ = ALU.add(ra, imm)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, _ = ALU.addc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.AND:
            result, carry = ALU.and_(ra, imm), False
        elif instruction.op_code() == OpCode.OR:
            result, carry = ALU.or_(ra, imm), False
        elif instruction.op_code() == OpCode.XOR:
            result, carry = ALU.xor(ra, imm), False
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.dc(), ALU.signed_extension(result))

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_asr_s_rri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AsrRRIOpCodes
        assert instruction.suffix() == Suffix.S_RRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ASR:
            result = ALU.asr(ra, imm)
        elif instruction.op_code() == OpCode.LSL:
            result = ALU.lsl(ra, imm)
        elif instruction.op_code() == OpCode.LSL1:
            result = ALU.lsl1(ra, imm)
        elif instruction.op_code() == OpCode.LSL1X:
            result = ALU.lsl1x(ra, imm)
        elif instruction.op_code() == OpCode.LSLX:
            result = ALU.lslx(ra, imm)
        elif instruction.op_code() == OpCode.LSR:
            result = ALU.lsr(ra, imm)
        elif instruction.op_code() == OpCode.LSR1:
            result = ALU.lsr1(ra, imm)
        elif instruction.op_code() == OpCode.LSR1X:
            result = ALU.lsr1x(ra, imm)
        elif instruction.op_code() == OpCode.LSRX:
            result = ALU.lsrx(ra, imm)
        elif instruction.op_code() == OpCode.ROL:
            result = ALU.rol(ra, imm)
        elif instruction.op_code() == OpCode.ROR:
            result = ALU.ror(ra, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.dc(), ALU.signed_extension(result))

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_s_rrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRICIOpCodes

        if instruction.op_code() in Instruction.AddRRICIOpCodes:
            self._execute_add_s_rrici(instruction)
        elif instruction.op_code() in Instruction.AndRRICIOpCodes:
            self._execute_and_s_rrici(instruction)
        elif instruction.op_code() in Instruction.AsrRRICIOpCodes:
            self._execute_asr_s_rrici(instruction)
        elif instruction.op_code() in Instruction.SubRRICIOpCodes:
            self._execute_sub_s_rrici(instruction)
        else:
            raise ValueError

    def _execute_add_s_rrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AddRRICIOpCodes
        assert instruction.suffix() == Suffix.S_RRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ADD:
            result, carry, overflow = ALU.add(ra, imm)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, overflow = ALU.addc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_add_nz_cc(ra, result, carry, overflow)

        self._cur_thread.register_file().write(instruction.dc(), ALU.signed_extension(result))

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_and_s_rrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AndRRICIOpCodes
        assert instruction.suffix() == Suffix.S_RRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.AND:
            result = ALU.and_(ra, imm)
        elif instruction.op_code() == OpCode.OR:
            result = ALU.or_(ra, imm)
        elif instruction.op_code() == OpCode.XOR:
            result = ALU.xor(ra, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_nz_cc(ra, result)

        self._cur_thread.register_file().write(instruction.dc(), ALU.signed_extension(result))

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_asr_s_rrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AsrRRICIOpCodes
        assert instruction.suffix() == Suffix.S_RRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ASR:
            result = ALU.asr(ra, imm)
        elif instruction.op_code() == OpCode.LSL:
            result = ALU.lsl(ra, imm)
        elif instruction.op_code() == OpCode.LSL1:
            result = ALU.lsl1(ra, imm)
        elif instruction.op_code() == OpCode.LSL1X:
            result = ALU.lsl1x(ra, imm)
        elif instruction.op_code() == OpCode.LSLX:
            result = ALU.lslx(ra, imm)
        elif instruction.op_code() == OpCode.LSR:
            result = ALU.lsr(ra, imm)
        elif instruction.op_code() == OpCode.LSR1:
            result = ALU.lsr1(ra, imm)
        elif instruction.op_code() == OpCode.LSR1X:
            result = ALU.lsr1x(ra, imm)
        elif instruction.op_code() == OpCode.LSRX:
            result = ALU.lsrx(ra, imm)
        elif instruction.op_code() == OpCode.ROL:
            result = ALU.rol(ra, imm)
        elif instruction.op_code() == OpCode.ROR:
            result = ALU.ror(ra, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_imm_shift_nz_cc(ra, result)

        self._cur_thread.register_file().write(instruction.dc(), ALU.signed_extension(result))

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_sub_s_rrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.SubRRICIOpCodes
        assert instruction.suffix() == Suffix.S_RRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.SUB:
            result, carry, overflow = ALU.sub(ra, imm)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, overflow = ALU.subc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_sub_nz_cc(ra, imm, result, carry, overflow)

        self._cur_thread.register_file().write(instruction.dc(), ALU.signed_extension(result))

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_s_rrif(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRIFOpCodes
        assert instruction.suffix() == Suffix.S_RRIF
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ADD:
            result, carry, _ = ALU.add(ra, imm)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, _ = ALU.addc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.AND:
            result, carry = ALU.and_(ra, imm), False
        elif instruction.op_code() == OpCode.ANDN:
            result, carry = ALU.andn(ra, imm), False
        elif instruction.op_code() == OpCode.NAND:
            result, carry = ALU.nand(ra, imm), False
        elif instruction.op_code() == OpCode.NOR:
            result, carry = ALU.nor(ra, imm), False
        elif instruction.op_code() == OpCode.NXOR:
            result, carry = ALU.nxor(ra, imm), False
        elif instruction.op_code() == OpCode.OR:
            result, carry = ALU.or_(ra, imm), False
        elif instruction.op_code() == OpCode.ORN:
            result, carry = ALU.orn(ra, imm), False
        elif instruction.op_code() == OpCode.SUB:
            result, carry, _ = ALU.sub(ra, imm)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, _ = ALU.subc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.XOR:
            result, carry = ALU.xor(ra, imm), False
        elif instruction.op_code() == OpCode.HASH:
            result, carry = ALU.hash(ra, imm), False
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.dc(), ALU.signed_extension(result))

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_u_rri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRIOpCodes

        if instruction.op_code() in Instruction.AddRRIOpCodes:
            self._execute_add_u_rri(instruction)
        elif instruction.op_code() in Instruction.AsrRRIOpCodes:
            self._execute_asr_u_rri(instruction)
        else:
            raise ValueError

    def _execute_add_u_rri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AddRRIOpCodes
        assert instruction.suffix() == Suffix.U_RRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ADD:
            result, carry, _ = ALU.add(ra, imm)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, _ = ALU.addc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.AND:
            result, carry = ALU.and_(ra, imm), False
        elif instruction.op_code() == OpCode.OR:
            result, carry = ALU.or_(ra, imm), False
        elif instruction.op_code() == OpCode.XOR:
            result, carry = ALU.xor(ra, imm), False
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.dc(), ALU.unsigned_extension(result))

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_asr_u_rri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AsrRRIOpCodes
        assert instruction.suffix() == Suffix.U_RRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ASR:
            result = ALU.asr(ra, imm)
        elif instruction.op_code() == OpCode.LSL:
            result = ALU.lsl(ra, imm)
        elif instruction.op_code() == OpCode.LSL1:
            result = ALU.lsl1(ra, imm)
        elif instruction.op_code() == OpCode.LSL1X:
            result = ALU.lsl1x(ra, imm)
        elif instruction.op_code() == OpCode.LSLX:
            result = ALU.lslx(ra, imm)
        elif instruction.op_code() == OpCode.LSR:
            result = ALU.lsr(ra, imm)
        elif instruction.op_code() == OpCode.LSR1:
            result = ALU.lsr1(ra, imm)
        elif instruction.op_code() == OpCode.LSR1X:
            result = ALU.lsr1x(ra, imm)
        elif instruction.op_code() == OpCode.LSRX:
            result = ALU.lsrx(ra, imm)
        elif instruction.op_code() == OpCode.ROL:
            result = ALU.rol(ra, imm)
        elif instruction.op_code() == OpCode.ROR:
            result = ALU.ror(ra, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.dc(), ALU.unsigned_extension(result))

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_u_rrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRICIOpCodes

        if instruction.op_code() in Instruction.AddRRICIOpCodes:
            self._execute_add_u_rrici(instruction)
        elif instruction.op_code() in Instruction.AndRRICIOpCodes:
            self._execute_and_u_rrici(instruction)
        elif instruction.op_code() in Instruction.AsrRRICIOpCodes:
            self._execute_asr_u_rrici(instruction)
        elif instruction.op_code() in Instruction.SubRRICIOpCodes:
            self._execute_sub_u_rrici(instruction)
        else:
            raise ValueError

    def _execute_add_u_rrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AddRRICIOpCodes
        assert instruction.suffix() == Suffix.U_RRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ADD:
            result, carry, overflow = ALU.add(ra, imm)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, overflow = ALU.addc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_add_nz_cc(ra, result, carry, overflow)

        self._cur_thread.register_file().write(instruction.dc(), ALU.unsigned_extension(result))

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_and_u_rrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AndRRICIOpCodes
        assert instruction.suffix() == Suffix.U_RRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.AND:
            result = ALU.and_(ra, imm)
        elif instruction.op_code() == OpCode.OR:
            result = ALU.or_(ra, imm)
        elif instruction.op_code() == OpCode.XOR:
            result = ALU.xor(ra, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_nz_cc(ra, result)

        self._cur_thread.register_file().write(instruction.dc(), ALU.unsigned_extension(result))

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_asr_u_rrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.AsrRRICIOpCodes
        assert instruction.suffix() == Suffix.U_RRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ASR:
            result = ALU.asr(ra, imm)
        elif instruction.op_code() == OpCode.LSL:
            result = ALU.lsl(ra, imm)
        elif instruction.op_code() == OpCode.LSL1:
            result = ALU.lsl1(ra, imm)
        elif instruction.op_code() == OpCode.LSL1X:
            result = ALU.lsl1x(ra, imm)
        elif instruction.op_code() == OpCode.LSLX:
            result = ALU.lslx(ra, imm)
        elif instruction.op_code() == OpCode.LSR:
            result = ALU.lsr(ra, imm)
        elif instruction.op_code() == OpCode.LSR1:
            result = ALU.lsr1(ra, imm)
        elif instruction.op_code() == OpCode.LSR1X:
            result = ALU.lsr1x(ra, imm)
        elif instruction.op_code() == OpCode.LSRX:
            result = ALU.lsrx(ra, imm)
        elif instruction.op_code() == OpCode.ROL:
            result = ALU.rol(ra, imm)
        elif instruction.op_code() == OpCode.ROR:
            result = ALU.ror(ra, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_imm_shift_nz_cc(ra, result)

        self._cur_thread.register_file().write(instruction.dc(), ALU.unsigned_extension(result))

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_sub_u_rrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.SubRRICIOpCodes
        assert instruction.suffix() == Suffix.U_RRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.SUB:
            result, carry, overflow = ALU.sub(ra, imm)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, overflow = ALU.subc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_sub_nz_cc(ra, imm, result, carry, overflow)

        self._cur_thread.register_file().write(instruction.dc(), ALU.unsigned_extension(result))

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_u_rrif(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRIFOpCodes
        assert instruction.suffix() == Suffix.U_RRIF
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.ADD:
            result, carry, _ = ALU.add(ra, imm)
        elif instruction.op_code() == OpCode.ADDC:
            result, carry, _ = ALU.addc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.AND:
            result, carry = ALU.and_(ra, imm), False
        elif instruction.op_code() == OpCode.ANDN:
            result, carry = ALU.andn(ra, imm), False
        elif instruction.op_code() == OpCode.NAND:
            result, carry = ALU.nand(ra, imm), False
        elif instruction.op_code() == OpCode.NOR:
            result, carry = ALU.nor(ra, imm), False
        elif instruction.op_code() == OpCode.NXOR:
            result, carry = ALU.nxor(ra, imm), False
        elif instruction.op_code() == OpCode.OR:
            result, carry = ALU.or_(ra, imm), False
        elif instruction.op_code() == OpCode.ORN:
            result, carry = ALU.orn(ra, imm), False
        elif instruction.op_code() == OpCode.SUB:
            result, carry, _ = ALU.sub(ra, imm)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, _ = ALU.subc(ra, imm, self._cur_thread.register_file().flag(Flag.CARRY))
        elif instruction.op_code() == OpCode.XOR:
            result, carry = ALU.xor(ra, imm), False
        elif instruction.op_code() == OpCode.HASH:
            result, carry = ALU.hash(ra, imm), False
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.dc(), ALU.unsigned_extension(result))

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_u_rrr(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRROpCodes
        assert instruction.suffix() == Suffix.U_RRR
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)

        if instruction.op_code() == OpCode.LSR:
            result = ALU.lsr(ra, rb)
        elif instruction.op_code() == OpCode.LSL:
            result = ALU.lsl(ra, rb)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.dc(), ALU.unsigned_extension(result))

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_rr(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RROpCodes
        assert instruction.suffix() == Suffix.RR
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)

        if instruction.op_code() == OpCode.CAO:
            result = ALU.cao(ra)
        elif instruction.op_code() == OpCode.CLO:
            result = ALU.clo(ra)
        elif instruction.op_code() == OpCode.CLS:
            result = ALU.cls(ra)
        elif instruction.op_code() == OpCode.CLZ:
            result = ALU.clz(ra)
        elif instruction.op_code() == OpCode.EXTSB:
            result = ALU.extsb(ra)
        elif instruction.op_code() == OpCode.EXTSH:
            result = ALU.extsh(ra)
        elif instruction.op_code() == OpCode.EXTUH:
            result = ALU.extuh(ra)
        elif instruction.op_code() == OpCode.SATS:
            result = ALU.sats(ra)
        elif instruction.op_code() == OpCode.TIME_CFG:
            raise NotImplementedError
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.rc(), result)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_rrc(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRCOpCodes
        assert instruction.suffix() == Suffix.RRC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)

        if instruction.op_code() == OpCode.CAO:
            result = ALU.cao(ra)
        elif instruction.op_code() == OpCode.CLO:
            result = ALU.clo(ra)
        elif instruction.op_code() == OpCode.CLS:
            result = ALU.cls(ra)
        elif instruction.op_code() == OpCode.CLZ:
            result = ALU.clz(ra)
        elif instruction.op_code() == OpCode.EXTSB:
            result = ALU.extsb(ra)
        elif instruction.op_code() == OpCode.EXTSH:
            result = ALU.extsh(ra)
        elif instruction.op_code() == OpCode.EXTUH:
            result = ALU.extuh(ra)
        elif instruction.op_code() == OpCode.SATS:
            result = ALU.sats(ra)
        elif instruction.op_code() == OpCode.TIME_CFG:
            raise NotImplementedError
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_set_cc(result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            self._cur_thread.register_file().write(instruction.rc(), 1)
        else:
            self._cur_thread.register_file().write(instruction.rc(), 0)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_rrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRCIOpCodes

        if instruction.op_code() in Instruction.CaoRRCIOpCodes:
            self._execute_cao_rrci(instruction)
        elif instruction.op_code() in Instruction.ExtsbRRCIOpCodes:
            self._execute_extsb_rrci(instruction)
        elif instruction.op_code() in Instruction.TimeCfgRRCIOpCodes:
            self._execute_time_cfg_rrci(instruction)
        else:
            raise ValueError

    def _execute_cao_rrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.CaoRRCIOpCodes
        assert instruction.suffix() == Suffix.RRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)

        if instruction.op_code() == OpCode.CAO:
            result = ALU.cao(ra)
        elif instruction.op_code() == OpCode.CLO:
            result = ALU.clo(ra)
        elif instruction.op_code() == OpCode.CLS:
            result = ALU.cls(ra)
        elif instruction.op_code() == OpCode.CLZ:
            result = ALU.clz(ra)
        elif instruction.op_code() == OpCode.EXTSB:
            result = ALU.extsb(ra)
        elif instruction.op_code() == OpCode.EXTSH:
            result = ALU.extsh(ra)
        elif instruction.op_code() == OpCode.EXTUH:
            result = ALU.extuh(ra)
        elif instruction.op_code() == OpCode.SATS:
            result = ALU.sats(ra)
        elif instruction.op_code() == OpCode.TIME_CFG:
            raise NotImplementedError
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_count_nz_cc(ra, result)

        self._cur_thread.register_file().write(instruction.rc(), result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_extsb_rrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.ExtsbRRCIOpCodes
        assert instruction.suffix() == Suffix.RRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)

        if instruction.op_code() == OpCode.EXTSB:
            result = ALU.extsb(ra)
        elif instruction.op_code() == OpCode.EXTSH:
            result = ALU.extsh(ra)
        elif instruction.op_code() == OpCode.EXTUH:
            result = ALU.extuh(ra)
        elif instruction.op_code() == OpCode.SATS:
            result = ALU.sats(ra)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_nz_cc(ra, result)

        self._cur_thread.register_file().write(instruction.rc(), result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_time_cfg_rrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.TimeCfgRRCIOpCodes
        assert instruction.suffix() == Suffix.RRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)

        if instruction.op_code() == OpCode.TIME_CFG:
            result = ALU.sats(ra)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_nz_cc(ra, result)

        self._cur_thread.register_file().write(instruction.rc(), result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_zr(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RROpCodes
        assert instruction.suffix() == Suffix.ZR
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)

        if instruction.op_code() == OpCode.CAO:
            result = ALU.cao(ra)
        elif instruction.op_code() == OpCode.CLO:
            result = ALU.clo(ra)
        elif instruction.op_code() == OpCode.CLS:
            result = ALU.cls(ra)
        elif instruction.op_code() == OpCode.CLZ:
            result = ALU.clz(ra)
        elif instruction.op_code() == OpCode.EXTSB:
            result = ALU.extsb(ra)
        elif instruction.op_code() == OpCode.EXTSH:
            result = ALU.extsh(ra)
        elif instruction.op_code() == OpCode.EXTUH:
            result = ALU.extuh(ra)
        elif instruction.op_code() == OpCode.SATS:
            result = ALU.sats(ra)
        elif instruction.op_code() == OpCode.TIME_CFG:
            raise NotImplementedError
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_zrc(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRCOpCodes
        assert instruction.suffix() == Suffix.ZRC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)

        if instruction.op_code() == OpCode.CAO:
            result = ALU.cao(ra)
        elif instruction.op_code() == OpCode.CLO:
            result = ALU.clo(ra)
        elif instruction.op_code() == OpCode.CLS:
            result = ALU.cls(ra)
        elif instruction.op_code() == OpCode.CLZ:
            result = ALU.clz(ra)
        elif instruction.op_code() == OpCode.EXTSB:
            result = ALU.extsb(ra)
        elif instruction.op_code() == OpCode.EXTSH:
            result = ALU.extsh(ra)
        elif instruction.op_code() == OpCode.EXTUH:
            result = ALU.extuh(ra)
        elif instruction.op_code() == OpCode.SATS:
            result = ALU.sats(ra)
        elif instruction.op_code() == OpCode.TIME_CFG:
            raise NotImplementedError
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_set_cc(result)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_zrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRCIOpCodes

        if instruction.op_code() in Instruction.CaoRRCIOpCodes:
            self._execute_cao_rrci(instruction)
        elif instruction.op_code() in Instruction.ExtsbRRCIOpCodes:
            self._execute_extsb_rrci(instruction)
        elif instruction.op_code() in Instruction.TimeCfgRRCIOpCodes:
            self._execute_time_cfg_rrci(instruction)
        else:
            raise ValueError

    def _execute_cao_zrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.CaoRRCIOpCodes
        assert instruction.suffix() == Suffix.ZRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)

        if instruction.op_code() == OpCode.CAO:
            result = ALU.cao(ra)
        elif instruction.op_code() == OpCode.CLO:
            result = ALU.clo(ra)
        elif instruction.op_code() == OpCode.CLS:
            result = ALU.cls(ra)
        elif instruction.op_code() == OpCode.CLZ:
            result = ALU.clz(ra)
        elif instruction.op_code() == OpCode.EXTSB:
            result = ALU.extsb(ra)
        elif instruction.op_code() == OpCode.EXTSH:
            result = ALU.extsh(ra)
        elif instruction.op_code() == OpCode.EXTUH:
            result = ALU.extuh(ra)
        elif instruction.op_code() == OpCode.SATS:
            result = ALU.sats(ra)
        elif instruction.op_code() == OpCode.TIME_CFG:
            raise NotImplementedError
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_count_nz_cc(ra, result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_extsb_zrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.ExtsbRRCIOpCodes
        assert instruction.suffix() == Suffix.ZRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)

        if instruction.op_code() == OpCode.EXTSB:
            result = ALU.extsb(ra)
        elif instruction.op_code() == OpCode.EXTSH:
            result = ALU.extsh(ra)
        elif instruction.op_code() == OpCode.EXTUH:
            result = ALU.extuh(ra)
        elif instruction.op_code() == OpCode.SATS:
            result = ALU.sats(ra)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_log_nz_cc(ra, result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_time_cfg_zrci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.TimeCfgRRCIOpCodes
        assert instruction.suffix() == Suffix.ZRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)

        if instruction.op_code() == OpCode.TIME_CFG:
            result = ALU.sats(ra)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_true_cc()

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_drdici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.DRDICIOpCodes
        assert instruction.suffix() == Suffix.DRDICI
        assert self._cur_thread is not None

        if instruction.op_code() == OpCode.DIV_STEP:
            self._execute_div_step_drdici(instruction)
        elif instruction.op_code() == OpCode.MUL_STEP:
            self._execute_mul_step_drdici(instruction)
        else:
            raise ValueError

    def _execute_div_step_drdici(self, instruction: Instruction) -> None:
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        dbe = self._cur_thread.register_file().read(instruction.db().even_register(), Representation.SIGNED)
        dbo = self._cur_thread.register_file().read(instruction.db().odd_register(), Representation.SIGNED)
        imm = instruction.imm().value()

        dbo_data_word = DataWord()
        dbo_data_word.set_value(dbo)

        ra_shift_data_word = DataWord()
        ra_shift_data_word.set_value(ALU.lsl(ra, imm))

        result, _, _ = ALU.sub(dbo, ALU.lsl(ra, imm))

        if dbo_data_word.value(Representation.UNSIGNED) >= ra_shift_data_word.value(Representation.UNSIGNED):
            dce = ALU.lsl1(dbe, 1)
            dco = result
        else:
            dce = ALU.lsl(dbe, 1)
            dco = self._cur_thread.register_file().read(instruction.dc().odd_register(), Representation.SIGNED)

        self._cur_thread.register_file().clear_conditions()
        self._set_div_cc(ra)

        self._cur_thread.register_file().write(instruction.dc().even_register(), dce)
        self._cur_thread.register_file().write(instruction.dc().odd_register(), dco)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, False)

    def _execute_mul_step_drdici(self, instruction: Instruction) -> None:
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        dbe = self._cur_thread.register_file().read(instruction.db().even_register(), Representation.SIGNED)
        dbo = self._cur_thread.register_file().read(instruction.db().odd_register(), Representation.SIGNED)
        imm = instruction.imm().value()

        result1 = ALU.lsr(dbe, 1)
        result2, _, _ = ALU.sub(ALU.and_(dbe, 1), 1)

        if result2 == 0:
            dco, _, _ = ALU.add(dbo, ALU.lsl(ra, imm))
        else:
            dco = self._cur_thread.register_file().read(instruction.dc().odd_register(), Representation.SIGNED)
        dce = ALU.lsr(dbe, 1)

        self._cur_thread.register_file().clear_conditions()
        self._set_boot_cc(ra, result1)

        self._cur_thread.register_file().write(instruction.dc().even_register(), dce)
        self._cur_thread.register_file().write(instruction.dc().odd_register(), dco)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result1, False)

    def _execute_rrri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRRIOpCodes
        assert instruction.suffix() == Suffix.RRRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.LSL_ADD:
            result, carry, _ = ALU.lsl_add(ra, rb, imm)
        elif instruction.op_code() == OpCode.LSL_SUB:
            result, carry, _ = ALU.lsl_sub(ra, rb, imm)
        elif instruction.op_code() == OpCode.LSR_ADD:
            result, carry, _ = ALU.lsr_add(ra, rb, imm)
        elif instruction.op_code() == OpCode.ROL_ADD:
            result, carry, _ = ALU.rol_add(ra, rb, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.rc(), result)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_rrrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRRICIOpCodes
        assert instruction.suffix() == Suffix.RRRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.LSL_ADD:
            result, carry, _ = ALU.lsl_add(ra, rb, imm)
        elif instruction.op_code() == OpCode.LSL_SUB:
            result, carry, _ = ALU.lsl_sub(ra, rb, imm)
        elif instruction.op_code() == OpCode.LSR_ADD:
            result, carry, _ = ALU.lsr_add(ra, rb, imm)
        elif instruction.op_code() == OpCode.ROL_ADD:
            result, carry, _ = ALU.rol_add(ra, rb, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_div_nz_cc(ra)

        self._cur_thread.register_file().write(instruction.rc(), result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_zrri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRRIOpCodes
        assert instruction.suffix() == Suffix.ZRRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.LSL_ADD:
            result, carry, _ = ALU.lsl_add(ra, rb, imm)
        elif instruction.op_code() == OpCode.LSL_SUB:
            result, carry, _ = ALU.lsl_sub(ra, rb, imm)
        elif instruction.op_code() == OpCode.LSR_ADD:
            result, carry, _ = ALU.lsr_add(ra, rb, imm)
        elif instruction.op_code() == OpCode.ROL_ADD:
            result, carry, _ = ALU.rol_add(ra, rb, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.rc(), result)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_zrrici(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RRRICIOpCodes
        assert instruction.suffix() == Suffix.ZRRICI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.LSL_ADD:
            result, carry, _ = ALU.lsl_add(ra, rb, imm)
        elif instruction.op_code() == OpCode.LSL_SUB:
            result, carry, _ = ALU.lsl_sub(ra, rb, imm)
        elif instruction.op_code() == OpCode.LSR_ADD:
            result, carry, _ = ALU.lsr_add(ra, rb, imm)
        elif instruction.op_code() == OpCode.ROL_ADD:
            result, carry, _ = ALU.rol_add(ra, rb, imm)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_div_nz_cc(ra)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_rir(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RIROpCodes
        assert instruction.suffix() == Suffix.RIR
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.SUB:
            result, carry, _ = ALU.sub(imm, ra)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, _ = ALU.subc(imm, ra, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.rc(), result)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_rirc(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RIRCOpCodes
        assert instruction.suffix() == Suffix.RIRC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.SUB:
            result, carry, _ = ALU.sub(imm, ra)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, _ = ALU.subc(imm, ra, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_sub_set_cc(ra, imm, result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            self._cur_thread.register_file().write(instruction.rc(), 1)
        else:
            self._cur_thread.register_file().write(instruction.rc(), 0)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_rirci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RIRCIOpCodes
        assert instruction.suffix() == Suffix.RIRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.SUB:
            result, carry, overflow = ALU.sub(imm, ra)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, overflow = ALU.subc(imm, ra, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_sub_nz_cc(ra, imm, result, carry, overflow)

        self._cur_thread.register_file().write(instruction.rc(), result)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_zir(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RIROpCodes
        assert instruction.suffix() == Suffix.ZIR
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.SUB:
            result, carry, _ = ALU.sub(imm, ra)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, _ = ALU.subc(imm, ra, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_zirc(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RIRCOpCodes
        assert instruction.suffix() == Suffix.ZIRC
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.SUB:
            result, carry, _ = ALU.sub(imm, ra)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, _ = ALU.subc(imm, ra, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_sub_set_cc(ra, imm, result)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_zirci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.RIRCIOpCodes
        assert instruction.suffix() == Suffix.ZIRCI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        imm = instruction.imm().value()

        if instruction.op_code() == OpCode.SUB:
            result, carry, overflow = ALU.sub(imm, ra)
        elif instruction.op_code() == OpCode.SUBC:
            result, carry, overflow = ALU.subc(imm, ra, self._cur_thread.register_file().flag(Flag.CARRY))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()
        self._set_sub_nz_cc(ra, imm, result, carry, overflow)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

        self._set_flags(result, carry)

    def _execute_r(self, instruction: Instruction) -> None:
        raise NotImplementedError

    def _execute_rci(self, instruction: Instruction) -> None:
        raise NotImplementedError

    def _execute_z(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.ROpCodes or instruction.op_code() == OpCode.NOP
        assert instruction.suffix() == Suffix.Z
        assert self._cur_thread is not None

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

    def _execute_zci(self, instruction: Instruction) -> None:
        raise NotImplementedError

    def _execute_ci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.CIOpCodes
        assert instruction.suffix() == Suffix.CI
        assert self._cur_thread is not None

        assert instruction.op_code() == OpCode.STOP
        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)

            self._cur_thread.set_thread_state(Thread.State.SLEEP)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

            self._cur_thread.set_thread_state(Thread.State.SLEEP)

    def _execute_i(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.IOpCodes
        assert instruction.suffix() == Suffix.I
        assert self._cur_thread is not None

        assert instruction.op_code() == OpCode.FAULT

        # TODO(bongjoon.hyun@gmail.com): this behavior must be simulated by a host CPU thread
        raise ValueError

    def _execute_ddci(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.DDCIOpCodes
        assert instruction.suffix() == Suffix.DDCI
        assert self._cur_thread is not None

        if instruction.op_code() == OpCode.MOVD:
            self._execute_movd_ddci(instruction)
        elif instruction.op_code() == OpCode.SWAPD:
            self._execute_swapd_ddci(instruction)
        else:
            raise ValueError

    def _execute_movd_ddci(self, instruction: Instruction) -> None:
        assert instruction.op_code() == OpCode.MOVD
        assert instruction.suffix() == Suffix.DDCI
        assert self._cur_thread is not None

        dbe = self._cur_thread.register_file().read(instruction.db().even_register(), Representation.SIGNED)
        dbo = self._cur_thread.register_file().read(instruction.db().odd_register(), Representation.SIGNED)

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.dc().even_register(), dbe)
        self._cur_thread.register_file().write(instruction.dc().odd_register(), dbo)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

    def _execute_swapd_ddci(self, instruction: Instruction) -> None:
        assert instruction.op_code() == OpCode.SWAPD
        assert instruction.suffix() == Suffix.DDCI
        assert self._cur_thread is not None

        dbe = self._cur_thread.register_file().read(instruction.db().even_register(), Representation.SIGNED)
        dbo = self._cur_thread.register_file().read(instruction.db().odd_register(), Representation.SIGNED)

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.dc().even_register(), dbo)
        self._cur_thread.register_file().write(instruction.dc().odd_register(), dbe)

        if self._cur_thread.register_file().condition(instruction.condition()):
            pc = instruction.pc().value()
            self._cur_thread.register_file().write_pc(pc)
        else:
            pc = self._cur_thread.register_file().read_pc()
            self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

    def _execute_erri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.ERRIOpCodes
        assert instruction.suffix() == Suffix.ERRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        off = instruction.off().value()

        address, _, _ = ALU.add(ra, off)

        if instruction.op_code() == OpCode.LBS:
            result = self._dispatcher.lbs(address)
        elif instruction.op_code() == OpCode.LBU:
            result = self._dispatcher.lbu(address)
        elif instruction.op_code() == OpCode.LHS:
            result = self._dispatcher.lhs(address)
        elif instruction.op_code() == OpCode.LHU:
            result = self._dispatcher.lhu(address)
        elif instruction.op_code() == OpCode.LW:
            result = self._dispatcher.lw(address)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.rc(), result)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

    def _execute_edri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.EDRIOpCodes
        assert instruction.suffix() == Suffix.EDRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        off = instruction.off().value()

        address, _, _ = ALU.add(ra, off)

        if instruction.op_code() == OpCode.LD:
            result = self._dispatcher.ld(address)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        self._cur_thread.register_file().write(instruction.dc(), result)

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

    def _execute_erii(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.ERIIOpCodes
        assert instruction.suffix() == Suffix.ERII
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        off = instruction.off().value()
        imm = instruction.imm().value()

        address, _, _ = ALU.add(ra, off)

        if instruction.op_code() == OpCode.SB:
            self._dispatcher.sb(address, imm)
        elif instruction.op_code() == OpCode.SB_ID:
            id_ = self._cur_thread.id_()
            self._dispatcher.sb(address, ALU.or_(id_, imm))
        elif instruction.op_code() == OpCode.SH:
            self._dispatcher.sh(address, imm)
        elif instruction.op_code() == OpCode.SH_ID:
            id_ = self._cur_thread.id_()
            self._dispatcher.sh(address, ALU.or_(id_, imm))
        elif instruction.op_code() == OpCode.SW:
            self._dispatcher.sw(address, imm)
        elif instruction.op_code() == OpCode.SW_ID:
            id_ = self._cur_thread.id_()
            self._dispatcher.sw(address, ALU.or_(id_, imm))
        elif instruction.op_code() == OpCode.SD:
            self._dispatcher.sd(address, imm)
        elif instruction.op_code() == OpCode.SD_ID:
            id_ = self._cur_thread.id_()
            self._dispatcher.sd(address, ALU.or_(id_, imm))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

    def _execute_erir(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.ERIROpCodes
        assert instruction.suffix() == Suffix.ERIR
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        off = instruction.off().value()
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)
        rb_word = DataWord()
        rb_word.set_value(rb)

        address, _, _ = ALU.add(ra, off)

        if instruction.op_code() == OpCode.SB:
            self._dispatcher.sb(address, rb_word.bit_slice(Representation.UNSIGNED, 0, 8))
        elif instruction.op_code() == OpCode.SH:
            self._dispatcher.sh(address, rb_word.bit_slice(Representation.UNSIGNED, 0, 16))
        elif instruction.op_code() == OpCode.SW:
            self._dispatcher.sw(address, rb_word.bit_slice(Representation.UNSIGNED, 0, 32))
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

    def _execute_erid(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.ERIDOpCodes
        assert instruction.suffix() == Suffix.ERID
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        off = instruction.off().value()
        db = self._cur_thread.register_file().read(instruction.db(), Representation.SIGNED)

        address, _, _ = ALU.add(ra, off)

        if instruction.op_code() == OpCode.SD:
            self._dispatcher.sd(address, db)
        else:
            raise ValueError

        self._cur_thread.register_file().clear_conditions()

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

    def _execute_dma_rri(self, instruction: Instruction) -> None:
        assert instruction.op_code() in Instruction.DMARRIOpCodes

        if instruction.op_code() == OpCode.LDMA:
            self._execute_ldma(instruction)
        elif instruction.op_code() == OpCode.LDMAI:
            self._execute_ldmai(instruction)
        elif instruction.op_code() == OpCode.SDMA:
            self._execute_sdma(instruction)
        else:
            raise ValueError

    def _execute_ldma(self, instruction: Instruction) -> None:
        assert instruction.op_code() == OpCode.LDMA
        assert instruction.suffix() == Suffix.DMA_RRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)
        imm = instruction.imm().value()

        wram_end_address = ConfigLoader.wram_offset() + ConfigLoader.wram_size()
        wram_end_address_width = math.floor(math.log2(wram_end_address)) + 1
        wram_mask = 2**wram_end_address_width - 1
        wram_address = ALU.and_(ra, wram_mask)

        mram_end_address = ConfigLoader.mram_offset() + ConfigLoader.mram_size()
        mram_end_address_width = math.floor(math.log2(mram_end_address)) + 1
        mram_mask = 2**mram_end_address_width - 1
        mram_address = ALU.and_(rb, mram_mask)

        min_access_granularity = ConfigLoader.min_access_granularity()

        num_bytes = (1 + ALU.and_(imm + ALU.and_(ALU.lsr(ra, 24), 255), 255)) * min_access_granularity

        self._dma.dpu_dma_transfer_from_mram_to_wram(mram_address, wram_address, num_bytes)

        self._cur_thread.register_file().clear_conditions()

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

    def _execute_ldmai(self, instruction: Instruction) -> None:
        raise NotImplementedError

    def _execute_sdma(self, instruction: Instruction) -> None:
        assert instruction.op_code() == OpCode.SDMA
        assert instruction.suffix() == Suffix.DMA_RRI
        assert self._cur_thread is not None

        ra = self._cur_thread.register_file().read(instruction.ra(), Representation.SIGNED)
        rb = self._cur_thread.register_file().read(instruction.rb(), Representation.SIGNED)
        imm = instruction.imm().value()

        wram_end_address = ConfigLoader.wram_offset() + ConfigLoader.wram_size()
        wram_end_address_width = math.floor(math.log2(wram_end_address)) + 1
        wram_mask = 2**wram_end_address_width - 1
        wram_address = ALU.and_(ra, wram_mask)

        mram_end_address = ConfigLoader.mram_offset() + ConfigLoader.mram_size()
        mram_end_address_width = math.floor(math.log2(mram_end_address)) + 1
        mram_mask = 2**mram_end_address_width - 1
        mram_address = ALU.and_(rb, mram_mask)

        min_access_granularity = ConfigLoader.min_access_granularity()

        num_bytes = (1 + ALU.and_(imm + ALU.and_(ALU.lsr(ra, 24), 255), 255)) * min_access_granularity

        self._dma.dpu_dma_transfer_from_wram_to_mram(wram_address, mram_address, num_bytes)

        self._cur_thread.register_file().clear_conditions()

        pc = self._cur_thread.register_file().read_pc()
        self._cur_thread.register_file().write_pc(pc + InstructionWord().size())

    def _set_flags(self, result: int, carry: bool) -> None:
        assert self._cur_thread is not None

        if result == 0:
            self._cur_thread.register_file().set_flag(Flag.ZERO)
        else:
            self._cur_thread.register_file().clear_flag(Flag.ZERO)

        if carry:
            self._cur_thread.register_file().set_flag(Flag.CARRY)
        else:
            self._cur_thread.register_file().clear_flag(Flag.CARRY)

    def _set_acquire_cc(self, result: int) -> None:
        assert self._cur_thread is not None

        if result == 0:
            self._cur_thread.register_file().set_condition(Condition.Z)
        else:
            self._cur_thread.register_file().set_condition(Condition.NZ)

    def _set_add_nz_cc(self, operand1: int, result: int, carry: bool, overflow: bool) -> None:
        assert self._cur_thread is not None

        if result == 0:
            self._cur_thread.register_file().set_condition(Condition.Z)
        else:
            self._cur_thread.register_file().set_condition(Condition.NZ)

        if carry:
            self._cur_thread.register_file().set_condition(Condition.C)
        else:
            self._cur_thread.register_file().set_condition(Condition.NC)

        if result == 0 and self._cur_thread.register_file().flag(Flag.ZERO):
            self._cur_thread.register_file().set_condition(Condition.XZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.XNZ)

        if overflow:
            self._cur_thread.register_file().set_condition(Condition.OV)
        else:
            self._cur_thread.register_file().set_condition(Condition.NOV)

        if result >= 0:
            self._cur_thread.register_file().set_condition(Condition.PL)
        else:
            self._cur_thread.register_file().set_condition(Condition.MI)

        if operand1 == 0:
            self._cur_thread.register_file().set_condition(Condition.SZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.NOV)

        if operand1 >= 0:
            self._cur_thread.register_file().set_condition(Condition.SPL)
        else:
            self._cur_thread.register_file().set_condition(Condition.SMI)

        result_data_word = DataWord()
        result_data_word.set_value(result)

        if result_data_word.bit(6):
            self._cur_thread.register_file().set_condition(Condition.NC5)
        if result_data_word.bit(7):
            self._cur_thread.register_file().set_condition(Condition.NC6)
        if result_data_word.bit(8):
            self._cur_thread.register_file().set_condition(Condition.NC7)
        if result_data_word.bit(9):
            self._cur_thread.register_file().set_condition(Condition.NC8)
        if result_data_word.bit(10):
            self._cur_thread.register_file().set_condition(Condition.NC9)
        if result_data_word.bit(11):
            self._cur_thread.register_file().set_condition(Condition.NC10)
        if result_data_word.bit(12):
            self._cur_thread.register_file().set_condition(Condition.NC11)
        if result_data_word.bit(13):
            self._cur_thread.register_file().set_condition(Condition.NC12)
        if result_data_word.bit(14):
            self._cur_thread.register_file().set_condition(Condition.NC13)
        if result_data_word.bit(15):
            self._cur_thread.register_file().set_condition(Condition.NC14)

    def _set_boot_cc(self, operand1: int, result: int) -> None:
        assert self._cur_thread is not None

        if result == 0:
            self._cur_thread.register_file().set_condition(Condition.Z)
        else:
            self._cur_thread.register_file().set_condition(Condition.NZ)

        if result == 0 and self._cur_thread.register_file().flag(Flag.ZERO):
            self._cur_thread.register_file().set_condition(Condition.XZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.XNZ)

        if operand1 == 0:
            self._cur_thread.register_file().set_condition(Condition.SZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.SNZ)

        if operand1 >= 0:
            self._cur_thread.register_file().set_condition(Condition.SPL)
        else:
            self._cur_thread.register_file().set_condition(Condition.SMI)

    def _set_count_nz_cc(self, operand1: int, result: int) -> None:
        assert self._cur_thread is not None

        if result == 0:
            self._cur_thread.register_file().set_condition(Condition.Z)
        else:
            self._cur_thread.register_file().set_condition(Condition.NZ)

        if result == 0 and self._cur_thread.register_file().flag(Flag.ZERO):
            self._cur_thread.register_file().set_condition(Condition.XZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.XNZ)

        if operand1 == 0:
            self._cur_thread.register_file().set_condition(Condition.SZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.SNZ)

        if operand1 >= 0:
            self._cur_thread.register_file().set_condition(Condition.SPL)
        else:
            self._cur_thread.register_file().set_condition(Condition.SMI)

        if result == DataWord().width():
            self._cur_thread.register_file().set_condition(Condition.MAX)
        else:
            self._cur_thread.register_file().set_condition(Condition.NMAX)

    def _set_div_cc(self, operand1: int) -> None:
        assert self._cur_thread is not None

        if operand1 == 0:
            self._cur_thread.register_file().set_condition(Condition.SZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.SNZ)

        if operand1 >= 0:
            self._cur_thread.register_file().set_condition(Condition.SPL)
        else:
            self._cur_thread.register_file().set_condition(Condition.SMI)

    def _set_div_nz_cc(self, operand1: int) -> None:
        assert self._cur_thread is not None

        if operand1 == 0:
            self._cur_thread.register_file().set_condition(Condition.SZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.SNZ)

        if operand1 >= 0:
            self._cur_thread.register_file().set_condition(Condition.SPL)
        else:
            self._cur_thread.register_file().set_condition(Condition.SMI)

    def _set_ext_sub_set_cc(self, operand1: int, operand2: int, result: int, carry: bool, overflow: bool) -> None:
        assert self._cur_thread is not None

        if result == 0:
            self._cur_thread.register_file().set_condition(Condition.Z)
        else:
            self._cur_thread.register_file().set_condition(Condition.NZ)

        if carry:
            self._cur_thread.register_file().set_condition(Condition.C)
        else:
            self._cur_thread.register_file().set_condition(Condition.NC)

        if result == 0 and self._cur_thread.register_file().flag(Flag.ZERO):
            self._cur_thread.register_file().set_condition(Condition.XZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.XNZ)

        if overflow:
            self._cur_thread.register_file().set_condition(Condition.OV)
        else:
            self._cur_thread.register_file().set_condition(Condition.NOV)

        if result >= 0:
            self._cur_thread.register_file().set_condition(Condition.PL)
        else:
            self._cur_thread.register_file().set_condition(Condition.MI)

        if operand1 == operand2:
            self._cur_thread.register_file().set_condition(Condition.EQ)
        else:
            self._cur_thread.register_file().set_condition(Condition.NEQ)

        if operand1 == 0:
            self._cur_thread.register_file().set_condition(Condition.SZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.SNZ)

        if operand1 >= 0:
            self._cur_thread.register_file().set_condition(Condition.SPL)
        else:
            self._cur_thread.register_file().set_condition(Condition.SMI)

        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        if data_word1.value(Representation.UNSIGNED) < data_word2.value(Representation.UNSIGNED):
            self._cur_thread.register_file().set_condition(Condition.LTU)

        if data_word1.value(Representation.UNSIGNED) <= data_word2.value(Representation.UNSIGNED):
            self._cur_thread.register_file().set_condition(Condition.LEU)

        if data_word1.value(Representation.UNSIGNED) > data_word2.value(Representation.UNSIGNED):
            self._cur_thread.register_file().set_condition(Condition.GTU)

        if data_word1.value(Representation.UNSIGNED) >= data_word2.value(Representation.UNSIGNED):
            self._cur_thread.register_file().set_condition(Condition.GEU)

        if data_word1.value(Representation.SIGNED) < data_word2.value(Representation.SIGNED):
            self._cur_thread.register_file().set_condition(Condition.LTS)

        if data_word1.value(Representation.SIGNED) <= data_word2.value(Representation.SIGNED):
            self._cur_thread.register_file().set_condition(Condition.LES)

        if data_word1.value(Representation.SIGNED) > data_word2.value(Representation.SIGNED):
            self._cur_thread.register_file().set_condition(Condition.GTS)

        if data_word1.value(Representation.SIGNED) >= data_word2.value(Representation.SIGNED):
            self._cur_thread.register_file().set_condition(Condition.GES)

        if carry or self._cur_thread.register_file().flag(Flag.ZERO):
            self._cur_thread.register_file().set_condition(Condition.XLEU)

        if carry and not self._cur_thread.register_file().flag(Flag.ZERO):
            self._cur_thread.register_file().set_condition(Condition.XGTU)

        if self._cur_thread.register_file().flag(Flag.ZERO) and (result < 0 or overflow):
            self._cur_thread.register_file().set_condition(Condition.XLES)

        if not self._cur_thread.register_file().flag(Flag.ZERO) and (result >= 0 or overflow):
            self._cur_thread.register_file().set_condition(Condition.XGTS)

    def _set_imm_shift_nz_cc(self, operand1: int, result: int) -> None:
        assert self._cur_thread is not None

        if result == 0:
            self._cur_thread.register_file().set_condition(Condition.Z)
        else:
            self._cur_thread.register_file().set_condition(Condition.NZ)

        if result == 0 and self._cur_thread.register_file().flag(Flag.ZERO):
            self._cur_thread.register_file().set_condition(Condition.XZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.XNZ)

        if result % 2 == 0:
            self._cur_thread.register_file().set_condition(Condition.E)
        else:
            self._cur_thread.register_file().set_condition(Condition.O)

        if result >= 0:
            self._cur_thread.register_file().set_condition(Condition.PL)
        else:
            self._cur_thread.register_file().set_condition(Condition.MI)

        if operand1 == 0:
            self._cur_thread.register_file().set_condition(Condition.SZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.SNZ)

        if operand1 % 2 == 0:
            self._cur_thread.register_file().set_condition(Condition.SE)
        else:
            self._cur_thread.register_file().set_condition(Condition.SO)

        if operand1 >= 0:
            self._cur_thread.register_file().set_condition(Condition.SPL)
        else:
            self._cur_thread.register_file().set_condition(Condition.SMI)

    def _set_log_nz_cc(self, operand1: int, result: int) -> None:
        assert self._cur_thread is not None

        if result == 0:
            self._cur_thread.register_file().set_condition(Condition.Z)
        else:
            self._cur_thread.register_file().set_condition(Condition.NZ)

        if result == 0 and self._cur_thread.register_file().flag(Flag.ZERO):
            self._cur_thread.register_file().set_condition(Condition.XZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.XNZ)

        if result >= 0:
            self._cur_thread.register_file().set_condition(Condition.PL)
        else:
            self._cur_thread.register_file().set_condition(Condition.MI)

        if operand1 == 0:
            self._cur_thread.register_file().set_condition(Condition.SZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.SNZ)

        if operand1 >= 0:
            self._cur_thread.register_file().set_condition(Condition.SPL)
        else:
            self._cur_thread.register_file().set_condition(Condition.SMI)

    def _set_log_set_cc(self, result: int) -> None:
        assert self._cur_thread is not None

        if result == 0:
            self._cur_thread.register_file().set_condition(Condition.Z)
        else:
            self._cur_thread.register_file().set_condition(Condition.NZ)

        if result == 0 and self._cur_thread.register_file().flag(Flag.ZERO):
            self._cur_thread.register_file().set_condition(Condition.XZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.XNZ)

    def _set_mul_nz_cc(self, operand1: int, result: int) -> None:
        assert self._cur_thread is not None

        if result == 0:
            self._cur_thread.register_file().set_condition(Condition.Z)
        else:
            self._cur_thread.register_file().set_condition(Condition.NZ)

        if result == 0 and self._cur_thread.register_file().flag(Flag.ZERO):
            self._cur_thread.register_file().set_condition(Condition.XZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.XNZ)

        if operand1 == 0:
            self._cur_thread.register_file().set_condition(Condition.SZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.SNZ)

        if operand1 >= 0:
            self._cur_thread.register_file().set_condition(Condition.SPL)
        else:
            self._cur_thread.register_file().set_condition(Condition.SMI)

        if result < 256:
            self._cur_thread.register_file().set_condition(Condition.SMALL)
        else:
            self._cur_thread.register_file().set_condition(Condition.LARGE)

    def _set_release_cc(self, result: int) -> None:
        assert self._cur_thread is not None

        if result != 0:
            self._cur_thread.register_file().set_condition(Condition.NZ)

    def _set_shift_nz_cc(self, operand1: int, result: int) -> None:
        assert self._cur_thread is not None

        if result == 0:
            self._cur_thread.register_file().set_condition(Condition.Z)
        else:
            self._cur_thread.register_file().set_condition(Condition.NZ)

        if result == 0 and self._cur_thread.register_file().flag(Flag.ZERO):
            self._cur_thread.register_file().set_condition(Condition.XZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.XNZ)

        if result % 2 == 0:
            self._cur_thread.register_file().set_condition(Condition.E)
        else:
            self._cur_thread.register_file().set_condition(Condition.O)

        if result >= 0:
            self._cur_thread.register_file().set_condition(Condition.PL)
        else:
            self._cur_thread.register_file().set_condition(Condition.MI)

        if operand1 == 0:
            self._cur_thread.register_file().set_condition(Condition.SZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.SNZ)

        if operand1 % 2 == 0:
            self._cur_thread.register_file().set_condition(Condition.SE)
        else:
            self._cur_thread.register_file().set_condition(Condition.SO)

        if operand1 >= 0:
            self._cur_thread.register_file().set_condition(Condition.SPL)
        else:
            self._cur_thread.register_file().set_condition(Condition.SMI)

        data_word1 = DataWord()
        data_word1.set_value(operand1)

        if data_word1.bit(5):
            self._cur_thread.register_file().set_condition(Condition.SH32)
        else:
            self._cur_thread.register_file().set_condition(Condition.NSH32)

    def _set_sub_nz_cc(self, operand1: int, operand2: int, result: int, carry: bool, overflow: bool) -> None:
        assert self._cur_thread is not None

        if result == 0:
            self._cur_thread.register_file().set_condition(Condition.Z)
        else:
            self._cur_thread.register_file().set_condition(Condition.NZ)

        if carry:
            self._cur_thread.register_file().set_condition(Condition.C)
        else:
            self._cur_thread.register_file().set_condition(Condition.NC)

        if result == 0 and self._cur_thread.register_file().flag(Flag.ZERO):
            self._cur_thread.register_file().set_condition(Condition.XZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.XNZ)

        if overflow:
            self._cur_thread.register_file().set_condition(Condition.OV)
        else:
            self._cur_thread.register_file().set_condition(Condition.NOV)

        if result >= 0:
            self._cur_thread.register_file().set_condition(Condition.PL)
        else:
            self._cur_thread.register_file().set_condition(Condition.MI)

        if operand1 == operand2:
            self._cur_thread.register_file().set_condition(Condition.EQ)
        else:
            self._cur_thread.register_file().set_condition(Condition.NEQ)

        if operand1 >= 0:
            self._cur_thread.register_file().set_condition(Condition.SPL)
        else:
            self._cur_thread.register_file().set_condition(Condition.SMI)

        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        if data_word1.value(Representation.UNSIGNED) < data_word2.value(Representation.UNSIGNED):
            self._cur_thread.register_file().set_condition(Condition.LTU)

        if data_word1.value(Representation.UNSIGNED) <= data_word2.value(Representation.UNSIGNED):
            self._cur_thread.register_file().set_condition(Condition.LEU)

        if data_word1.value(Representation.UNSIGNED) > data_word2.value(Representation.UNSIGNED):
            self._cur_thread.register_file().set_condition(Condition.GTU)

        if data_word1.value(Representation.UNSIGNED) >= data_word2.value(Representation.UNSIGNED):
            self._cur_thread.register_file().set_condition(Condition.GEU)

        if data_word1.value(Representation.SIGNED) < data_word2.value(Representation.SIGNED):
            self._cur_thread.register_file().set_condition(Condition.LTS)

        if data_word1.value(Representation.SIGNED) <= data_word2.value(Representation.SIGNED):
            self._cur_thread.register_file().set_condition(Condition.LES)

        if data_word1.value(Representation.SIGNED) > data_word2.value(Representation.SIGNED):
            self._cur_thread.register_file().set_condition(Condition.GTS)

        if data_word1.value(Representation.SIGNED) >= data_word2.value(Representation.SIGNED):
            self._cur_thread.register_file().set_condition(Condition.GES)

        if carry or self._cur_thread.register_file().flag(Flag.ZERO):
            self._cur_thread.register_file().set_condition(Condition.XLEU)

        if carry and not self._cur_thread.register_file().flag(Flag.ZERO):
            self._cur_thread.register_file().set_condition(Condition.XGTU)

        if self._cur_thread.register_file().flag(Flag.ZERO) and (result < 0 or overflow):
            self._cur_thread.register_file().set_condition(Condition.XLES)

        if not self._cur_thread.register_file().flag(Flag.ZERO) and (result >= 0 or overflow):
            self._cur_thread.register_file().set_condition(Condition.XGTS)

    def _set_sub_set_cc(self, operand1: int, operand2: int, result: int) -> None:
        assert self._cur_thread is not None

        if result == 0:
            self._cur_thread.register_file().set_condition(Condition.Z)
        else:
            self._cur_thread.register_file().set_condition(Condition.NZ)

        if result == 0 and self._cur_thread.register_file().flag(Flag.ZERO):
            self._cur_thread.register_file().set_condition(Condition.XZ)
        else:
            self._cur_thread.register_file().set_condition(Condition.XNZ)

        if operand1 == operand2:
            self._cur_thread.register_file().set_condition(Condition.EQ)
        else:
            self._cur_thread.register_file().set_condition(Condition.NEQ)

    def _set_true_cc(self):
        pass
