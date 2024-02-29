from abi.isa.instruction.condition import Condition
from abi.isa.instruction.instruction import Instruction
from abi.isa.instruction.op_code import OpCode
from abi.isa.instruction.suffix import Suffix
from encoder.instruction_encoder import InstructionEncoder
from initializer.instruction_initializer import InstructionInitializer


def test_rici():
    for _ in range(100):
        for op_code in Instruction.RICIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RICI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RICI
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_rri():
    for _ in range(100):
        for op_code in Instruction.RRIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RRI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RRI
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()


def test_rric():
    for _ in range(100):
        for op_code in Instruction.RRICOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RRIC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RRIC
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition()


def test_rrici():
    for _ in range(100):
        for op_code in Instruction.RRICIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RRICI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RRICI
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_rrif():
    for _ in range(100):
        for op_code in Instruction.RRIFOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RRIF)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RRIF
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition() == Condition.FALSE


def test_rrr():
    for _ in range(100):
        for op_code in Instruction.RRROpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RRR)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RRR
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()


def test_rrrc():
    for _ in range(100):
        for op_code in Instruction.RRRCOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RRRC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RRRC
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.condition() == dst_instruction.condition()


def test_rrrci():
    for _ in range(100):
        for op_code in Instruction.RRRCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RRRCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RRRCI
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_zri():
    for _ in range(100):
        for op_code in Instruction.RRIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZRI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZRI
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()


def test_zric():
    for _ in range(100):
        for op_code in Instruction.RRICOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZRIC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZRIC
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition()


def test_zrici():
    for _ in range(100):
        for op_code in Instruction.RRICIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZRICI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZRICI
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_zrif():
    for _ in range(100):
        for op_code in Instruction.RRIFOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZRIF)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZRIF
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition() == Condition.FALSE


def test_zrr():
    for _ in range(100):
        for op_code in Instruction.RRROpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZRR)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZRR
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()


def test_zrrc():
    for _ in range(100):
        for op_code in Instruction.RRRCOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZRRC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZRRC
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.condition() == dst_instruction.condition()


def test_zrrci():
    for _ in range(100):
        for op_code in Instruction.RRRCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZRRCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZRRCI
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_s_rri():
    for _ in range(100):
        for op_code in Instruction.RRIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_RRI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_RRI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()


def test_u_rri():
    for _ in range(100):
        for op_code in Instruction.RRIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_RRI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_RRI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()


def test_s_rric():
    for _ in range(100):
        for op_code in Instruction.RRICOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_RRIC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_RRIC
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition()


def test_u_rric():
    for _ in range(100):
        for op_code in Instruction.RRICOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_RRIC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_RRIC
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition()


def test_s_rrici():
    for _ in range(100):
        for op_code in Instruction.RRICOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_RRICI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_RRICI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_u_rrici():
    for _ in range(100):
        for op_code in Instruction.RRICOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_RRICI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_RRICI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_s_rrif():
    for _ in range(100):
        for op_code in Instruction.RRIFOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_RRIF)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_RRIF
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition() == Condition.FALSE


def test_u_rrif():
    for _ in range(100):
        for op_code in Instruction.RRIFOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_RRIF)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_RRIF
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition() == Condition.FALSE


def test_s_rrr():
    for _ in range(100):
        for op_code in Instruction.RRROpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_RRR)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_RRR
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()


def test_u_rrr():
    for _ in range(100):
        for op_code in Instruction.RRROpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_RRR)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_RRR
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()


def test_s_rrrc():
    for _ in range(100):
        for op_code in Instruction.RRRCOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_RRRC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_RRRC
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.condition() == dst_instruction.condition()


def test_u_rrrc():
    for _ in range(100):
        for op_code in Instruction.RRRCOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_RRRC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_RRRC
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.condition() == dst_instruction.condition()


def test_s_rrrci():
    for _ in range(100):
        for op_code in Instruction.RRRCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_RRRCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_RRRCI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_u_rrrci():
    for _ in range(100):
        for op_code in Instruction.RRRCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_RRRCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_RRRCI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_rr():
    for _ in range(100):
        for op_code in Instruction.RROpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RR)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RR
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.ra() == dst_instruction.ra()


def test_rrc():
    for _ in range(100):
        for op_code in Instruction.RRCOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RRC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RRC
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()


def test_rrci():
    for _ in range(100):
        for op_code in Instruction.RRCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RRCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RRCI
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_zr():
    for _ in range(100):
        for op_code in Instruction.RROpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZR)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZR
            assert src_instruction.ra() == dst_instruction.ra()


def test_zrc():
    for _ in range(100):
        for op_code in Instruction.RRCOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZRC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZRC
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()


def test_zrci():
    for _ in range(100):
        for op_code in Instruction.RRCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZRCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZRCI
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_s_rr():
    for _ in range(100):
        for op_code in Instruction.RROpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_RR)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_RR
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()


def test_u_rr():
    for _ in range(100):
        for op_code in Instruction.RROpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_RR)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_RR
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()


def test_s_rrc():
    for _ in range(100):
        for op_code in Instruction.RRCOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_RRC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_RRC
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()


def test_u_rrc():
    for _ in range(100):
        for op_code in Instruction.RRCOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_RRC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_RRC
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()


def test_s_rrci():
    for _ in range(100):
        for op_code in Instruction.RRCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_RRCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_RRCI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_u_rrci():
    for _ in range(100):
        for op_code in Instruction.RRCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_RRCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_RRCI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_drdici():
    for _ in range(100):
        for op_code in Instruction.DRDICIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.DRDICI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.DRDICI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.db() == dst_instruction.db()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_rrri():
    for _ in range(100):
        for op_code in Instruction.RRRIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RRRI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RRRI
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.imm().value() == dst_instruction.imm().value()


def test_rrrici():
    for _ in range(100):
        for op_code in Instruction.RRRICIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RRRICI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RRRICI
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_zrri():
    for _ in range(100):
        for op_code in Instruction.RRRIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZRRI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZRRI
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.imm().value() == dst_instruction.imm().value()


def test_zrrici():
    for _ in range(100):
        for op_code in Instruction.RRRICIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZRRICI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZRRICI
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_s_rrri():
    for _ in range(100):
        for op_code in Instruction.RRRIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_RRRI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_RRRI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.imm().value() == dst_instruction.imm().value()


def test_u_rrri():
    for _ in range(100):
        for op_code in Instruction.RRRIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_RRRI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_RRRI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.imm().value() == dst_instruction.imm().value()


def test_s_rrrici():
    for _ in range(100):
        for op_code in Instruction.RRRIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_RRRICI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_RRRICI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_u_rrrici():
    for _ in range(100):
        for op_code in Instruction.RRRIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_RRRICI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_RRRICI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_rir():
    for _ in range(100):
        for op_code in Instruction.RIROpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RIR)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RIR
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.ra() == dst_instruction.ra()


def test_rirc():
    for _ in range(100):
        for op_code in Instruction.RIRCOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RIRC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RIRC
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()


def test_rirci():
    for _ in range(100):
        for op_code in Instruction.RIRCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RIRCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RIRCI
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_zir():
    for _ in range(100):
        for op_code in Instruction.RIROpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZIR)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZIR
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.ra() == dst_instruction.ra()


def test_zirc():
    for _ in range(100):
        for op_code in Instruction.RIRCOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZIRC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZIRC
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()


def test_zirci():
    for _ in range(100):
        for op_code in Instruction.RIRCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZIRCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZIRCI
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_s_rirc():
    for _ in range(100):
        for op_code in Instruction.RIRCOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_RIRC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_RIRC
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()


def test_u_rirc():
    for _ in range(100):
        for op_code in Instruction.RIRCOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_RIRC)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_RIRC
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()


def test_s_rirci():
    for _ in range(100):
        for op_code in Instruction.RIRCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_RIRCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_RIRCI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_u_rirci():
    for _ in range(100):
        for op_code in Instruction.RIRCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_RIRCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_RIRCI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_r():
    for _ in range(100):
        for op_code in Instruction.ROpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.R)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.R
            assert src_instruction.rc() == dst_instruction.rc()


def test_rci():
    for _ in range(100):
        for op_code in Instruction.RCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.RCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.RCI
            assert src_instruction.rc() == dst_instruction.rc()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_z():
    for _ in range(100):
        for op_code in {*Instruction.ROpCodes, OpCode.NOP}:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.Z)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.Z


def test_zci():
    for _ in range(100):
        for op_code in Instruction.RCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ZCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ZCI
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_s_r():
    for _ in range(100):
        for op_code in Instruction.ROpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_R)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_R
            assert src_instruction.dc() == dst_instruction.dc()


def test_u_r():
    for _ in range(100):
        for op_code in Instruction.ROpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_R)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_R
            assert src_instruction.dc() == dst_instruction.dc()


def test_s_rci():
    for _ in range(100):
        for op_code in Instruction.RCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_RCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_RCI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_u_rci():
    for _ in range(100):
        for op_code in Instruction.RCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_RCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_RCI
            assert src_instruction.dc() == dst_instruction.dc()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_ci():
    for _ in range(100):
        for op_code in Instruction.CIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.CI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.CI
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_i():
    for _ in range(100):
        for op_code in Instruction.IOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.I)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.I
            assert src_instruction.imm().value() == dst_instruction.imm().value()


def test_ddci():
    for _ in range(100):
        for op_code in Instruction.DDCIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.DDCI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.DDCI
            assert src_instruction.dc() == src_instruction.dc()
            assert src_instruction.db() == src_instruction.db()
            assert src_instruction.condition() == dst_instruction.condition()
            assert src_instruction.pc().value() == dst_instruction.pc().value()


def test_erri():
    for _ in range(100):
        for op_code in Instruction.ERRIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ERRI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ERRI
            assert src_instruction.endian() == src_instruction.endian()
            assert src_instruction.rc() == src_instruction.rc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.off().value() == dst_instruction.off().value()


def test_s_erri():
    for _ in range(100):
        for op_code in Instruction.ERRIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.S_ERRI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.S_ERRI
            assert src_instruction.endian() == src_instruction.endian()
            assert src_instruction.dc() == src_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.off().value() == dst_instruction.off().value()


def test_u_erri():
    for _ in range(100):
        for op_code in Instruction.ERRIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.U_ERRI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.U_ERRI
            assert src_instruction.endian() == src_instruction.endian()
            assert src_instruction.dc() == src_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.off().value() == dst_instruction.off().value()


def test_edri():
    for _ in range(100):
        for op_code in Instruction.EDRIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.EDRI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.EDRI
            assert src_instruction.endian() == src_instruction.endian()
            assert src_instruction.dc() == src_instruction.dc()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.off().value() == dst_instruction.off().value()


def test_erii():
    for _ in range(100):
        for op_code in Instruction.ERIIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ERII)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ERII
            assert src_instruction.endian() == src_instruction.endian()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.off().value() == dst_instruction.off().value()
            assert src_instruction.imm().value() == dst_instruction.imm().value()


def test_erir():
    for _ in range(100):
        for op_code in Instruction.ERIROpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ERIR)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ERIR
            assert src_instruction.endian() == src_instruction.endian()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.off().value() == dst_instruction.off().value()
            assert src_instruction.rb() == dst_instruction.rb()


def test_erid():
    for _ in range(100):
        for op_code in Instruction.ERIDOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.ERID)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.ERID
            assert src_instruction.endian() == src_instruction.endian()
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.off().value() == dst_instruction.off().value()
            assert src_instruction.db() == dst_instruction.db()


def test_dma_rri():
    for _ in range(100):
        for op_code in Instruction.DMARRIOpCodes:
            src_instruction = InstructionInitializer.instruction(op_code, Suffix.DMA_RRI)

            bytes_ = InstructionEncoder.encode(src_instruction)
            dst_instruction = InstructionEncoder.decode(bytes_)

            assert src_instruction.op_code() == dst_instruction.op_code() == op_code
            assert src_instruction.suffix() == dst_instruction.suffix() == Suffix.DMA_RRI
            assert src_instruction.ra() == dst_instruction.ra()
            assert src_instruction.rb() == dst_instruction.rb()
            assert src_instruction.imm().value() == dst_instruction.imm().value()
