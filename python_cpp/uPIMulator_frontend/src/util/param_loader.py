from abi.word.data_word import DataWord
from util.config_loader import ConfigLoader


class ParamLoader:
    def __init__(self):
        pass

    @staticmethod
    def logic_frequency() -> int:
        return 450

    @staticmethod
    def memory_frequency() -> int:
        return 2666

    @staticmethod
    def num_pipeline_stages() -> int:
        return 14

    @staticmethod
    def instruction_scheduling_policy() -> str:
        return "revolver"

    @staticmethod
    def num_revolver_scheduling_cycles() -> int:
        return 11

    @staticmethod
    def memory_scheduling_policy() -> str:
        return "FIFO"

    @staticmethod
    def num_wordlines() -> int:
        return 512

    @staticmethod
    def wordline_size() -> int:
        assert ConfigLoader.mram_size() % ParamLoader.num_wordlines() == 0
        wordline_size = ConfigLoader.mram_size() // ParamLoader.num_wordlines()
        assert wordline_size % DataWord().size() == 0
        return wordline_size

    @staticmethod
    def t_rcd() -> int:
        raise NotImplementedError

    @staticmethod
    def t_ras() -> int:
        raise NotImplementedError

    @staticmethod
    def t_cl() -> int:
        raise NotImplementedError

    @staticmethod
    def t_cwl() -> int:
        raise NotImplementedError

    @staticmethod
    def t_bl() -> int:
        raise NotImplementedError

    @staticmethod
    def t_rp() -> int:
        raise NotImplementedError
