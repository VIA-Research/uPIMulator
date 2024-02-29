from abi.isa.instruction.instruction import Instruction
from abi.word.instruction_word import InstructionWord
from encoder.instruction_encoder import InstructionEncoder


class Decoder:
    def __init__(self):
        pass

    @staticmethod
    def decode(instruction_word: InstructionWord) -> Instruction:
        bytes_ = instruction_word.to_bytes()
        return InstructionEncoder.decode(bytes_)
