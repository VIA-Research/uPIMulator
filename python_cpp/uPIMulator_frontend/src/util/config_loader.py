class ConfigLoader:
    def __init__(self):
        pass

    @staticmethod
    def atomic_address_width() -> int:
        return 32

    @staticmethod
    def atomic_data_width() -> int:
        return 32

    @staticmethod
    def atomic_offset() -> int:
        return 0

    @staticmethod
    def atomic_size() -> int:
        return 256

    @staticmethod
    def iram_address_width() -> int:
        return 32

    @staticmethod
    def iram_data_width() -> int:
        return 96

    @staticmethod
    def iram_offset() -> int:
        return 384 * 1024

    @staticmethod
    def iram_size() -> int:
        return 48 * 1024

    @staticmethod
    def wram_address_width() -> int:
        return 32

    @staticmethod
    def wram_data_width() -> int:
        return 32

    @staticmethod
    def wram_offset() -> int:
        return 512

    @staticmethod
    def wram_size() -> int:
        return 128 * 1024

    @staticmethod
    def stack_size() -> int:
        return 2 * 1024

    @staticmethod
    def heap_size() -> int:
        return 4 * 1024

    @staticmethod
    def mram_address_width() -> int:
        return 32

    @staticmethod
    def mram_data_width() -> int:
        return 32

    @staticmethod
    def mram_offset() -> int:
        return 512 * 1024

    @staticmethod
    def mram_size() -> int:
        return 64 * 1024 * 1024

    @staticmethod
    def num_gp_registers() -> int:
        return 24

    @staticmethod
    def max_num_tasklets() -> int:
        return 24

    @staticmethod
    def min_access_granularity() -> int:
        return 8
