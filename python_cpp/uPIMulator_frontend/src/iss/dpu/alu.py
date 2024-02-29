from typing import Tuple

from abi.word.data_word import DataWord
from abi.word.double_data_word import DoubleDataWord
from abi.word.representation import Representation


class ALU:
    def __init__(self):
        pass

    @staticmethod
    def atomic_address_hash(operand1: int, operand2: int) -> int:
        assert operand1 + operand2 < 2 ** 8
        return ALU.add(operand1, operand2)[0]

    @staticmethod
    def add(operand1: int, operand2: int) -> Tuple[int, bool, bool]:
        return ALU.addc(operand1, operand2, False)

    @staticmethod
    def addc(operand1: int, operand2: int, carry_flag: bool) -> Tuple[int, bool, bool]:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result = data_word1.value(Representation.UNSIGNED) + data_word2.value(Representation.UNSIGNED) + int(carry_flag)

        max_unsigned_value = 2 ** DataWord().width() - 1
        if result > max_unsigned_value:
            result %= 2 ** DataWord().width()
            carry = True
        else:
            carry = False

        result_data_word = DataWord()
        result_data_word.set_value(result)

        if data_word1.sign_bit() and data_word2.sign_bit() and not result_data_word.sign_bit():
            overflow = True
        elif not data_word1.sign_bit() and not data_word2.sign_bit() and result_data_word.sign_bit():
            overflow = True
        else:
            overflow = False

        return result, carry, overflow

    @staticmethod
    def sub(operand1: int, operand2: int) -> Tuple[int, bool, bool]:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        if data_word1.value(Representation.UNSIGNED) >= data_word2.value(Representation.UNSIGNED):
            result = data_word1.value(Representation.UNSIGNED) - data_word2.value(Representation.UNSIGNED)
            carry = False
        else:
            result = (
                2 ** DataWord().width()
                + data_word1.value(Representation.UNSIGNED)
                - data_word2.value(Representation.UNSIGNED)
            )
            carry = True

        result_data_word = DataWord()
        result_data_word.set_value(result)

        if data_word1.sign_bit() and not data_word2.sign_bit() and result_data_word.sign_bit():
            overflow = True
        elif not data_word1.sign_bit() and data_word2.sign_bit() and not result_data_word.sign_bit():
            overflow = True
        else:
            overflow = False
        return result, carry, overflow

    @staticmethod
    def subc(operand1: int, operand2: int, carry_flag: bool) -> Tuple[int, bool, bool]:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        # NOTE(bongjoon.hyun@gmail.com): this doesn't make sense logically, but it works
        if data_word1.value(Representation.UNSIGNED) + int(carry_flag) >= data_word2.value(Representation.UNSIGNED):
            result = (
                data_word1.value(Representation.UNSIGNED) - data_word2.value(Representation.UNSIGNED) - int(carry_flag)
            )
            carry = False
        else:
            result = (
                2 ** DataWord().width()
                + data_word1.value(Representation.UNSIGNED)
                - data_word2.value(Representation.UNSIGNED)
                - int(carry_flag)
            )
            carry = True

        result_data_word = DataWord()
        result_data_word.set_value(result)

        if data_word1.sign_bit() and not data_word2.sign_bit() and result_data_word.sign_bit():
            overflow = True
        elif not data_word1.sign_bit() and data_word2.sign_bit() and not result_data_word.sign_bit():
            overflow = True
        else:
            overflow = False
        return result, carry, overflow

    @staticmethod
    def and_(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        for i in range(result_data_word.width()):
            if data_word1.bit(i) and data_word2.bit(i):
                result_data_word.set_bit(i)
            else:
                result_data_word.clear_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def nand(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        for i in range(result_data_word.width()):
            if data_word1.bit(i) and data_word2.bit(i):
                result_data_word.clear_bit(i)
            else:
                result_data_word.set_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def andn(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        for i in range(result_data_word.width()):
            if (not data_word1.bit(i)) and data_word2.bit(i):
                result_data_word.set_bit(i)
            else:
                result_data_word.clear_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def or_(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        for i in range(result_data_word.width()):
            if data_word1.bit(i) or data_word2.bit(i):
                result_data_word.set_bit(i)
            else:
                result_data_word.clear_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def nor(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        for i in range(result_data_word.width()):
            if data_word1.bit(i) or data_word2.bit(i):
                result_data_word.clear_bit(i)
            else:
                result_data_word.set_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def orn(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        for i in range(result_data_word.width()):
            if (not data_word1.bit(i)) or data_word2.bit(i):
                result_data_word.set_bit(i)
            else:
                result_data_word.clear_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def xor(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        for i in range(result_data_word.width()):
            if ((not data_word1.bit(i)) and data_word2.bit(i)) or (data_word1.bit(i) and (not data_word2.bit(i))):
                result_data_word.set_bit(i)
            else:
                result_data_word.clear_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def nxor(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        for i in range(result_data_word.width()):
            if ((not data_word1.bit(i)) and data_word2.bit(i)) or (data_word1.bit(i) and (not data_word2.bit(i))):
                result_data_word.clear_bit(i)
            else:
                result_data_word.set_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def asr(operand: int, shift: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)
        msb = data_word.bit(data_word.width() - 1)

        shift_data_word = DataWord()
        shift_data_word.set_value(shift)
        shift_value = shift_data_word.bit_slice(Representation.UNSIGNED, 0, 5)

        result_data_word = DataWord()
        for i in range(result_data_word.width()):
            if i + shift_value >= result_data_word.width():
                if msb:
                    result_data_word.set_bit(i)
                else:
                    result_data_word.clear_bit(i)
            else:
                if data_word.bit(i + shift_value):
                    result_data_word.set_bit(i)
                else:
                    result_data_word.clear_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def lsl(operand: int, shift: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)

        shift_data_word = DataWord()
        shift_data_word.set_value(shift)
        shift_value = shift_data_word.bit_slice(Representation.UNSIGNED, 0, 5)

        result_data_word = DataWord()
        for i in range(result_data_word.width()):
            if i < shift_value:
                result_data_word.clear_bit(i)
            else:
                if data_word.bit(i - shift_value):
                    result_data_word.set_bit(i)
                else:
                    result_data_word.clear_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def lsl_add(operand1: int, operand2: int, shift: int) -> Tuple[int, bool, bool]:
        return ALU.add(operand1, ALU.lsl(operand2, shift))

    @staticmethod
    def lsl_sub(operand1: int, operand2: int, shift: int) -> Tuple[int, bool, bool]:
        return ALU.sub(operand1, ALU.lsl(operand2, shift))

    @staticmethod
    def lsl1(operand: int, shift: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)

        shift_data_word = DataWord()
        shift_data_word.set_value(shift)
        shift_value = shift_data_word.bit_slice(Representation.UNSIGNED, 0, 5)

        result_data_word = DataWord()
        for i in range(result_data_word.width()):
            if i < shift_value:
                result_data_word.set_bit(i)
            else:
                if data_word.bit(i - shift_value):
                    result_data_word.set_bit(i)
                else:
                    result_data_word.clear_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def lsl1x(operand: int, shift: int) -> int:
        raise NotImplementedError

    @staticmethod
    def lslx(operand: int, shift: int) -> int:
        shift_data_word = DataWord()
        shift_data_word.set_value(shift)
        shift_value = shift_data_word.bit_slice(Representation.UNSIGNED, 0, 5)

        if shift_value == 0:
            return 0
        else:
            return ALU.lsr(operand, 32 - shift_value)

    @staticmethod
    def lsr(operand: int, shift: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)

        shift_data_word = DataWord()
        shift_data_word.set_value(shift)
        shift_value = shift_data_word.bit_slice(Representation.UNSIGNED, 0, 5)
        
        result_data_word = DataWord()
        for i in range(result_data_word.width()):
            if i + shift_value >= result_data_word.width():
                result_data_word.clear_bit(i)
            else:
                if data_word.bit(i + shift_value):
                    result_data_word.set_bit(i)
                else:
                    result_data_word.clear_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def lsr_add(operand1: int, operand2: int, shift: int) -> Tuple[int, bool, bool]:
        return ALU.add(operand1, ALU.lsr(operand2, shift))

    @staticmethod
    def lsr1(operand: int, shift: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)

        shift_data_word = DataWord()
        shift_data_word.set_value(shift)
        shift_value = shift_data_word.bit_slice(Representation.UNSIGNED, 0, 5)

        result_data_word = DataWord()
        for i in range(result_data_word.width()):
            if i + shift_value >= result_data_word.width():
                result_data_word.set_bit(i)
            else:
                if data_word.bit(i + shift_value):
                    result_data_word.set_bit(i)
                else:
                    result_data_word.clear_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def lsr1x(operand: int, shift: int) -> int:
        raise NotImplementedError

    @staticmethod
    def lsrx(operand: int, shift: int) -> int:
        shift_data_word = DataWord()
        shift_data_word.set_value(shift)
        shift_value = shift_data_word.bit_slice(Representation.UNSIGNED, 0, 5)

        if shift_value == 0:
            return 0
        else:
            return ALU.lsl(operand, 32 - shift_value)

    @staticmethod
    def rol(operand: int, shift: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)

        shift_data_word = DataWord()
        shift_data_word.set_value(shift)
        shift_value = shift_data_word.bit_slice(Representation.UNSIGNED, 0, 5)

        result_data_word = DataWord()
        for i in range(result_data_word.width()):
            if i < shift_value:
                if data_word.bit(i + data_word.width() - shift_value):
                    result_data_word.set_bit(i)
                else:
                    result_data_word.clear_bit(i)
            else:
                if data_word.bit(i - shift_value):
                    result_data_word.set_bit(i)
                else:
                    result_data_word.clear_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def rol_add(operand1: int, operand2: int, shift: int) -> Tuple[int, bool, bool]:
        return ALU.add(operand1, ALU.rol(operand2, shift))

    @staticmethod
    def ror(operand: int, shift: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)

        shift_data_word = DataWord()
        shift_data_word.set_value(shift)
        shift_value = shift_data_word.bit_slice(Representation.UNSIGNED, 0, 5)

        result_data_word = DataWord()
        for i in range(result_data_word.width()):
            if i + shift_value >= result_data_word.width():
                if data_word.bit((i + shift_value) % data_word.width()):
                    result_data_word.set_bit(i)
                else:
                    result_data_word.clear_bit(i)
            else:
                if data_word.bit(i + shift_value):
                    result_data_word.set_bit(i)
                else:
                    result_data_word.clear_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def cao(operand: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)

        ones = 0
        for i in range(data_word.width()):
            if data_word.bit(i):
                ones += 1
        return ones

    @staticmethod
    def clo(operand: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)

        leading_ones = 0
        for i in range(data_word.width()):
            if data_word.bit(data_word.width() - 1 - i):
                leading_ones += 1
            else:
                break
        return leading_ones

    @staticmethod
    def cls(operand: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)
        msb = data_word.bit(data_word.width() - 1)

        leading_sign_bits = 0
        for i in range(data_word.width()):
            if data_word.bit(data_word.width() - 1 - i) == msb:
                leading_sign_bits += 1
            else:
                break
        return leading_sign_bits

    @staticmethod
    def clz(operand: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)

        leading_zeros = 0
        for i in range(data_word.width()):
            if not data_word.bit(data_word.width() - 1 - i):
                leading_zeros += 1
            else:
                break
        return leading_zeros

    @staticmethod
    def cmpb4(operand1: int, operand2: int) -> int:
        assert DataWord().width() == 4 * 8

        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        for i in range(4):
            begin = 8 * i
            end = 8 * (i + 1)

            byte1 = data_word1.bit_slice(Representation.UNSIGNED, begin, end)
            byte2 = data_word2.bit_slice(Representation.UNSIGNED, begin, end)
            if byte1 == byte2:
                result_data_word.set_bit_slice(begin, end, 1)
            else:
                result_data_word.set_bit_slice(begin, end, 0)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def extsb(operand: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)
        return data_word.bit_slice(Representation.SIGNED, 0, 8)

    @staticmethod
    def extsh(operand: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)
        return data_word.bit_slice(Representation.SIGNED, 0, 16)

    @staticmethod
    def extub(operand: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)
        return data_word.bit_slice(Representation.UNSIGNED, 0, 8)

    @staticmethod
    def extuh(operand: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)
        return data_word.bit_slice(Representation.UNSIGNED, 0, 16)

    @staticmethod
    def mul_sh_sh(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        result_data_word.set_value(
            data_word1.bit_slice(Representation.SIGNED, 8, 16) * data_word2.bit_slice(Representation.SIGNED, 8, 16)
        )

        return result_data_word.value(Representation.SIGNED)

    @staticmethod
    def mul_sh_sl(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        result_data_word.set_value(
            data_word1.bit_slice(Representation.SIGNED, 8, 16) * data_word2.bit_slice(Representation.SIGNED, 0, 8)
        )

        return result_data_word.value(Representation.SIGNED)

    @staticmethod
    def mul_sh_uh(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        result_data_word.set_value(
            data_word1.bit_slice(Representation.SIGNED, 8, 16) * data_word2.bit_slice(Representation.UNSIGNED, 8, 16)
        )

        return result_data_word.value(Representation.SIGNED)

    @staticmethod
    def mul_sh_ul(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        result_data_word.set_value(
            data_word1.bit_slice(Representation.SIGNED, 8, 16) * data_word2.bit_slice(Representation.UNSIGNED, 0, 8)
        )

        return result_data_word.value(Representation.SIGNED)

    @staticmethod
    def mul_sl_sh(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        result_data_word.set_value(
            data_word1.bit_slice(Representation.SIGNED, 0, 8) * data_word2.bit_slice(Representation.SIGNED, 8, 16)
        )

        return result_data_word.value(Representation.SIGNED)

    @staticmethod
    def mul_sl_sl(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        result_data_word.set_value(
            data_word1.bit_slice(Representation.SIGNED, 0, 8) * data_word2.bit_slice(Representation.SIGNED, 0, 8)
        )

        return result_data_word.value(Representation.SIGNED)

    @staticmethod
    def mul_sl_uh(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        result_data_word.set_value(
            data_word1.bit_slice(Representation.SIGNED, 0, 8) * data_word2.bit_slice(Representation.UNSIGNED, 8, 16)
        )

        return result_data_word.value(Representation.SIGNED)

    @staticmethod
    def mul_sl_ul(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        result_data_word.set_value(
            data_word1.bit_slice(Representation.SIGNED, 0, 8) * data_word2.bit_slice(Representation.UNSIGNED, 0, 8)
        )

        return result_data_word.value(Representation.SIGNED)

    @staticmethod
    def mul_uh_uh(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        result_data_word.set_value(
            data_word1.bit_slice(Representation.UNSIGNED, 8, 16) * data_word2.bit_slice(Representation.UNSIGNED, 8, 16)
        )

        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def mul_uh_ul(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        result_data_word.set_value(
            data_word1.bit_slice(Representation.UNSIGNED, 8, 16) * data_word2.bit_slice(Representation.UNSIGNED, 0, 8)
        )

        return result_data_word.value(Representation.SIGNED)

    @staticmethod
    def mul_ul_uh(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        result_data_word.set_value(
            data_word1.bit_slice(Representation.UNSIGNED, 0, 8) * data_word2.bit_slice(Representation.UNSIGNED, 8, 16)
        )

        return result_data_word.value(Representation.SIGNED)

    @staticmethod
    def mul_ul_ul(operand1: int, operand2: int) -> int:
        data_word1 = DataWord()
        data_word1.set_value(operand1)

        data_word2 = DataWord()
        data_word2.set_value(operand2)

        result_data_word = DataWord()
        result_data_word.set_value(
            data_word1.bit_slice(Representation.UNSIGNED, 0, 8) * data_word2.bit_slice(Representation.UNSIGNED, 0, 8)
        )

        return result_data_word.value(Representation.SIGNED)

    @staticmethod
    def sats(operand: int) -> int:
        raise NotImplementedError

    @staticmethod
    def hash(operand1: int, operand2: int) -> int:
        raise NotImplementedError

    @staticmethod
    def tilde(operand: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)

        result_data_word = DataWord()
        for i in range(data_word.width()):
            if data_word.bit(i):
                result_data_word.clear_bit(i)
            else:
                result_data_word.set_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def signed_extension(operand: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)

        result_data_word = DoubleDataWord()
        result_data_word.set_bit_slice(0, data_word.width(), data_word.value(Representation.UNSIGNED))

        if data_word.sign_bit():
            for i in range(data_word.width(), 2 * data_word.width()):
                result_data_word.set_bit(i)
        return result_data_word.value(Representation.UNSIGNED)

    @staticmethod
    def unsigned_extension(operand: int) -> int:
        data_word = DataWord()
        data_word.set_value(operand)

        result_data_word = DoubleDataWord()
        result_data_word.set_bit_slice(0, data_word.width(), data_word.value(Representation.UNSIGNED))

        for i in range(data_word.width(), 2 * data_word.width()):
            result_data_word.clear_bit(i)

        return result_data_word.value(Representation.UNSIGNED)
