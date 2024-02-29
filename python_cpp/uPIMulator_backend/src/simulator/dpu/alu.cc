#include "simulator/dpu/alu.h"

#include <cassert>
#include <cmath>
#include <functional>

#include "abi/word/data_word.h"

namespace upmem_sim::simulator::dpu {

int64_t ALU::atomic_address_hash(int64_t operand1, int64_t operand2) {
  assert(operand1 + operand2 < 256);
  auto [result, carry, overflow] = add(operand1, operand2);
  return result;
}

std::tuple<int64_t, bool, bool> ALU::add(int64_t operand1, int64_t operand2) {
  return addc(operand1, operand2, false);
}

std::tuple<int64_t, bool, bool> ALU::addc(int64_t operand1, int64_t operand2,
                                          bool carry_flag) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  int64_t result = data_word1->value(abi::word::UNSIGNED) +
                   data_word2->value(abi::word::UNSIGNED) + carry_flag;

  int64_t max_unsigned_value =
      static_cast<int64_t>(pow(2, abi::word::DataWord().width())) - 1;
  bool carry;
  if (result > max_unsigned_value) {
    result %= static_cast<int64_t>(pow(2, abi::word::DataWord().width()));
    carry = true;
  } else {
    carry = false;
  }

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(result);

  bool overflow;
  if (data_word1->sign_bit() and data_word2->sign_bit() and
      not result_data_word->sign_bit()) {
    overflow = true;
  } else if (not data_word1->sign_bit() and not data_word2->sign_bit() and
             result_data_word->sign_bit()) {
    overflow = true;
  } else {
    overflow = false;
  }

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return {result, carry, overflow};
}

std::tuple<int64_t, bool, bool> ALU::sub(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  int64_t result;
  bool carry;
  if (data_word1->value(abi::word::UNSIGNED) >=
      data_word2->value(abi::word::UNSIGNED)) {
    result = data_word1->value(abi::word::UNSIGNED) -
             data_word2->value(abi::word::UNSIGNED);
    carry = false;
  } else {
    result = static_cast<int64_t>(pow(2, abi::word::DataWord().width())) +
             data_word1->value(abi::word::UNSIGNED) -
             data_word2->value(abi::word::UNSIGNED);
    carry = true;
  }

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(result);

  bool overflow;
  if (data_word1->sign_bit() and not data_word2->sign_bit() and
      result_data_word->sign_bit()) {
    overflow = true;
  } else if (not data_word1->sign_bit() and data_word2->sign_bit() and
             not result_data_word->sign_bit()) {
    overflow = true;
  } else {
    overflow = false;
  }

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return {result, carry, overflow};
}

std::tuple<int64_t, bool, bool> ALU::subc(int64_t operand1, int64_t operand2,
                                          bool carry_flag) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  int64_t result;
  bool carry;
  if (data_word1->value(abi::word::UNSIGNED) + carry_flag >=
      data_word2->value(abi::word::UNSIGNED)) {
    result = data_word1->value(abi::word::UNSIGNED) -
             data_word2->value(abi::word::UNSIGNED) - carry_flag;
    carry = false;
  } else {
    result = static_cast<int64_t>(pow(2, abi::word::DataWord().width())) +
             data_word1->value(abi::word::UNSIGNED) -
             data_word2->value(abi::word::UNSIGNED) - carry_flag;
    carry = true;
  }

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(result);

  bool overflow;
  if (data_word1->sign_bit() and not data_word2->sign_bit() and
      result_data_word->sign_bit()) {
    overflow = true;
  } else if (not data_word1->sign_bit() and data_word2->sign_bit() and
             not result_data_word->sign_bit()) {
    overflow = true;
  } else {
    overflow = false;
  }

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return {result, carry, overflow};
}

int64_t ALU::and_(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < result_data_word->width(); i++) {
    if (data_word1->bit(i) and data_word2->bit(i)) {
      result_data_word->set_bit(i);
    } else {
      result_data_word->clear_bit(i);
    }
  }

  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::nand(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < result_data_word->width(); i++) {
    if (data_word1->bit(i) and data_word2->bit(i)) {
      result_data_word->clear_bit(i);
    } else {
      result_data_word->set_bit(i);
    }
  }

  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::andn(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < result_data_word->width(); i++) {
    if (not data_word1->bit(i) and data_word2->bit(i)) {
      result_data_word->set_bit(i);
    } else {
      result_data_word->clear_bit(i);
    }
  }

  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::or_(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < result_data_word->width(); i++) {
    if (data_word1->bit(i) or data_word2->bit(i)) {
      result_data_word->set_bit(i);
    } else {
      result_data_word->clear_bit(i);
    }
  }

  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::nor(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < result_data_word->width(); i++) {
    if (data_word1->bit(i) or data_word2->bit(i)) {
      result_data_word->clear_bit(i);
    } else {
      result_data_word->set_bit(i);
    }
  }

  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::orn(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < result_data_word->width(); i++) {
    if (not data_word1->bit(i) or data_word2->bit(i)) {
      result_data_word->set_bit(i);
    } else {
      result_data_word->clear_bit(i);
    }
  }

  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::xor_(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < result_data_word->width(); i++) {
    if (not data_word1->bit(i) and data_word2->bit(i)) {
      result_data_word->set_bit(i);
    } else if (data_word1->bit(i) and not data_word2->bit(i)) {
      result_data_word->set_bit(i);
    } else {
      result_data_word->clear_bit(i);
    }
  }

  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::nxor(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < result_data_word->width(); i++) {
    if (not data_word1->bit(i) and data_word2->bit(i)) {
      result_data_word->clear_bit(i);
    } else if (data_word1->bit(i) and not data_word2->bit(i)) {
      result_data_word->clear_bit(i);
    } else {
      result_data_word->set_bit(i);
    }
  }

  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::asr(int64_t operand, int64_t shift) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);
  bool msb = data_word->sign_bit();

  auto shift_data_word = new abi::word::DataWord();
  shift_data_word->set_value(shift);
  auto shift_value = shift_data_word->bit_slice(abi::word::UNSIGNED, 0, 5);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < result_data_word->width(); i++) {
    if (i + shift_value >= result_data_word->width()) {
      if (msb) {
        result_data_word->set_bit(i);
      } else {
        result_data_word->clear_bit(i);
      }
    } else {
      if (data_word->bit(static_cast<int>(i + shift_value))) {
        result_data_word->set_bit(i);
      } else {
        result_data_word->clear_bit(i);
      }
    }
  }

  int64_t result = data_word->value(abi::word::UNSIGNED);

  delete data_word;
  delete result_data_word;

  return result;
}

int64_t ALU::lsl(int64_t operand, int64_t shift) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);

  auto shift_data_word = new abi::word::DataWord();
  shift_data_word->set_value(shift);
  auto shift_value = shift_data_word->bit_slice(abi::word::UNSIGNED, 0, 5);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < result_data_word->width(); i++) {
    if (i < shift_value) {
      result_data_word->clear_bit(i);
    } else {
      if (data_word->bit(static_cast<int>(i - shift_value))) {
        result_data_word->set_bit(i);
      } else {
        result_data_word->clear_bit(i);
      }
    }
  }

  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word;
  delete result_data_word;

  return result;
}

std::tuple<int64_t, bool, bool> ALU::lsl_add(int64_t operand1, int64_t operand2,
                                             int64_t shift) {
  return add(operand1, lsl(operand2, shift));
}

std::tuple<int64_t, bool, bool> ALU::lsl_sub(int64_t operand1, int64_t operand2,
                                             int64_t shift) {
  return sub(operand1, lsl(operand2, shift));
}

int64_t ALU::lsl1(int64_t operand, int64_t shift) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);

  auto shift_data_word = new abi::word::DataWord();
  shift_data_word->set_value(shift);
  auto shift_value = shift_data_word->bit_slice(abi::word::UNSIGNED, 0, 5);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < result_data_word->width(); i++) {
    if (i < shift_value) {
      result_data_word->set_bit(i);
    } else {
      if (data_word->bit(static_cast<int>(i - shift_value))) {
        result_data_word->set_bit(i);
      } else {
        result_data_word->clear_bit(i);
      }
    }
  }

  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word;
  delete result_data_word;

  return result;
}

int64_t ALU::lsl1x(int64_t operand, int64_t shift) {
  throw std::bad_function_call();
}

int64_t ALU::lslx(int64_t operand, int64_t shift) {
  auto shift_data_word = new abi::word::DataWord();
  shift_data_word->set_value(shift);
  auto shift_value = shift_data_word->bit_slice(abi::word::UNSIGNED, 0, 5);

  if (shift_value == 0) {
    return 0;
  } else {
    return lsr(operand, abi::word::DataWord().width() - shift_value);
  }
}

int64_t ALU::lsr(int64_t operand, int64_t shift) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);

  auto shift_data_word = new abi::word::DataWord();
  shift_data_word->set_value(shift);
  auto shift_value = shift_data_word->bit_slice(abi::word::UNSIGNED, 0, 5);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < result_data_word->width(); i++) {
    if (i + shift_value >= result_data_word->width()) {
      result_data_word->clear_bit(i);
    } else {
      if (data_word->bit(static_cast<int>(i + shift_value))) {
        result_data_word->set_bit(i);
      } else {
        result_data_word->clear_bit(i);
      }
    }
  }
  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word;
  delete result_data_word;

  return result;
}

std::tuple<int64_t, bool, bool> ALU::lsr_add(int64_t operand1, int64_t operand2,
                                             int64_t shift) {
  return add(operand1, lsr(operand2, shift));
}

int64_t ALU::lsr1(int64_t operand, int64_t shift) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);

  auto shift_data_word = new abi::word::DataWord();
  shift_data_word->set_value(shift);
  auto shift_value = shift_data_word->bit_slice(abi::word::UNSIGNED, 0, 5);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < result_data_word->width(); i++) {
    if (i + shift_value >= result_data_word->width()) {
      result_data_word->set_bit(i);
    } else {
      if (data_word->bit(static_cast<int>(i + shift_value))) {
        result_data_word->set_bit(i);
      } else {
        result_data_word->clear_bit(i);
      }
    }
  }

  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word;
  delete result_data_word;

  return result;
}

int64_t ALU::lsr1x(int64_t operand, int64_t shift) {
  throw std::bad_function_call();
}

int64_t ALU::lsrx(int64_t operand, int64_t shift) {
  auto shift_data_word = new abi::word::DataWord();
  shift_data_word->set_value(shift);
  auto shift_value = shift_data_word->bit_slice(abi::word::UNSIGNED, 0, 5);

  if (shift_value == 0) {
    return 0;
  } else {
    return lsl(operand, abi::word::DataWord().width() - shift_value);
  }
}

int64_t ALU::rol(int64_t operand, int64_t shift) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);

  auto shift_data_word = new abi::word::DataWord();
  shift_data_word->set_value(shift);
  auto shift_value = shift_data_word->bit_slice(abi::word::UNSIGNED, 0, 5);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < result_data_word->width(); i++) {
    if (i < shift_value) {
      if (data_word->bit(
              static_cast<int>(i + data_word->width() - shift_value))) {
        result_data_word->set_bit(i);
      } else {
        result_data_word->clear_bit(i);
      }
    } else {
      if (data_word->bit(static_cast<int>(i - shift_value))) {
        result_data_word->set_bit(i);
      } else {
        result_data_word->clear_bit(i);
      }
    }
  }

  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word;
  delete result_data_word;

  return result;
}

std::tuple<int64_t, bool, bool> ALU::rol_add(int64_t operand1, int64_t operand2,
                                             int64_t shift) {
  return add(operand1, rol(operand2, shift));
}

int64_t ALU::ror(int64_t operand, int64_t shift) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);

  auto shift_data_word = new abi::word::DataWord();
  shift_data_word->set_value(shift);
  auto shift_value = shift_data_word->bit_slice(abi::word::UNSIGNED, 0, 5);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < result_data_word->width(); i++) {
    if (i + shift_value >= result_data_word->width()) {
      if (data_word->bit(
              static_cast<int>((i + shift_value) % data_word->width()))) {
        result_data_word->set_bit(i);
      } else {
        result_data_word->clear_bit(i);
      }
    } else {
      if (data_word->bit(static_cast<int>(i + shift_value))) {
        result_data_word->set_bit(i);
      } else {
        result_data_word->clear_bit(i);
      }
    }
  }

  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word;
  delete result_data_word;

  return result;
}

int64_t ALU::cao(int64_t operand) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);

  int64_t ones = 0;
  for (int i = 0; i < data_word->width(); i++) {
    if (data_word->bit(i)) {
      ones += 1;
    }
  }

  delete data_word;

  return ones;
}

int64_t ALU::clo(int64_t operand) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);

  int64_t leading_ones = 0;
  for (int i = 0; i < data_word->width(); i++) {
    if (data_word->bit(data_word->width() - 1 - i)) {
      leading_ones += 1;
    } else {
      break;
    }
  }
  delete data_word;

  return leading_ones;
}

int64_t ALU::cls(int64_t operand) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);

  bool msb = data_word->sign_bit();
  int64_t leading_sign_bits = 0;
  for (int i = 0; i < data_word->width(); i++) {
    if (data_word->bit(data_word->width() - 1 - i) == msb) {
      leading_sign_bits += 1;
    } else {
      break;
    }
  }

  delete data_word;
  return leading_sign_bits;
}

int64_t ALU::clz(int64_t operand) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);

  int64_t leading_zeros = 0;
  for (int i = 0; i < data_word->width(); i++) {
    if (not data_word->bit(data_word->width() - 1 - i)) {
      leading_zeros += 1;
    } else {
      break;
    }
  }

  delete data_word;

  return leading_zeros;
}

int64_t ALU::cmpb4(int64_t operand1, int64_t operand2) {
  assert(abi::word::DataWord().width() == 4 * 8);

  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  for (int i = 0; i < 4; i++) {
    int begin = 8 * i;
    int end = 8 * (i + 1);

    int64_t byte1 = data_word1->bit_slice(abi::word::UNSIGNED, begin, end);
    int64_t byte2 = data_word2->bit_slice(abi::word::UNSIGNED, begin, end);

    if (byte1 == byte2) {
      result_data_word->set_bit_slice(begin, end, 1);
    } else {
      result_data_word->set_bit_slice(begin, end, 0);
    }
  }

  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete result_data_word;

  return result;
}

int64_t ALU::extsb(int64_t operand) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);
  int64_t result = data_word->bit_slice(abi::word::SIGNED, 0, 8);
  delete data_word;
  return result;
}

int64_t ALU::extsh(int64_t operand) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);
  int64_t result = data_word->bit_slice(abi::word::SIGNED, 0, 16);
  delete data_word;
  return result;
}

int64_t ALU::extub(int64_t operand) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);
  int64_t result = data_word->bit_slice(abi::word::UNSIGNED, 0, 8);
  delete data_word;
  return result;
}

int64_t ALU::extuh(int64_t operand) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);
  int64_t result = data_word->bit_slice(abi::word::UNSIGNED, 0, 16);
  delete data_word;
  return result;
}

int64_t ALU::mul_sh_sh(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(data_word1->bit_slice(abi::word::SIGNED, 8, 16) *
                              data_word2->bit_slice(abi::word::SIGNED, 8, 16));
  int64_t result = result_data_word->value(abi::word::SIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::mul_sh_sl(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(data_word1->bit_slice(abi::word::SIGNED, 8, 16) *
                              data_word2->bit_slice(abi::word::UNSIGNED, 0, 8));
  int64_t result = result_data_word->value(abi::word::SIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::mul_sh_uh(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(
      data_word1->bit_slice(abi::word::SIGNED, 8, 16) *
      data_word2->bit_slice(abi::word::UNSIGNED, 8, 16));
  int64_t result = result_data_word->value(abi::word::SIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::mul_sh_ul(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(data_word1->bit_slice(abi::word::SIGNED, 8, 16) *
                              data_word2->bit_slice(abi::word::UNSIGNED, 0, 8));
  int64_t result = result_data_word->value(abi::word::SIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::mul_sl_sh(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(data_word1->bit_slice(abi::word::SIGNED, 0, 8) *
                              data_word2->bit_slice(abi::word::SIGNED, 8, 16));
  int64_t result = result_data_word->value(abi::word::SIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::mul_sl_sl(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(data_word1->bit_slice(abi::word::SIGNED, 0, 8) *
                              data_word2->bit_slice(abi::word::SIGNED, 0, 8));
  int64_t result = result_data_word->value(abi::word::SIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::mul_sl_uh(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(
      data_word1->bit_slice(abi::word::SIGNED, 0, 8) *
      data_word2->bit_slice(abi::word::UNSIGNED, 8, 16));
  int64_t result = result_data_word->value(abi::word::SIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::mul_sl_ul(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(data_word1->bit_slice(abi::word::SIGNED, 0, 8) *
                              data_word2->bit_slice(abi::word::UNSIGNED, 0, 8));
  int64_t result = result_data_word->value(abi::word::SIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::mul_uh_uh(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(
      data_word1->bit_slice(abi::word::UNSIGNED, 8, 16) *
      data_word2->bit_slice(abi::word::UNSIGNED, 8, 16));
  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::mul_uh_ul(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(
      data_word1->bit_slice(abi::word::UNSIGNED, 8, 16) *
      data_word2->bit_slice(abi::word::UNSIGNED, 0, 8));
  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::mul_ul_uh(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(
      data_word1->bit_slice(abi::word::UNSIGNED, 0, 8) *
      data_word2->bit_slice(abi::word::UNSIGNED, 8, 16));
  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::mul_ul_ul(int64_t operand1, int64_t operand2) {
  auto data_word1 = new abi::word::DataWord();
  data_word1->set_value(operand1);

  auto data_word2 = new abi::word::DataWord();
  data_word2->set_value(operand2);

  auto result_data_word = new abi::word::DataWord();
  result_data_word->set_value(data_word1->bit_slice(abi::word::UNSIGNED, 0, 8) *
                              data_word2->bit_slice(abi::word::UNSIGNED, 0, 8));
  int64_t result = result_data_word->value(abi::word::UNSIGNED);

  delete data_word1;
  delete data_word2;
  delete result_data_word;

  return result;
}

int64_t ALU::sats(int64_t operand) { throw std::bad_function_call(); }

int64_t ALU::hash(int64_t operand1, int64_t operand2) {
  throw std::bad_function_call();
}

std::tuple<int64_t, int64_t> ALU::signed_extension(int64_t operand) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);

  auto even_data_word = new abi::word::DataWord();
  if (data_word->sign_bit()) {
    even_data_word->set_value(-1);
  }
  int64_t even = even_data_word->value(abi::word::UNSIGNED);

  auto odd_data_word = new abi::word::DataWord();
  odd_data_word->set_value(data_word->value(abi::word::UNSIGNED));
  int64_t odd = odd_data_word->value(abi::word::UNSIGNED);

  delete data_word;
  delete even_data_word;
  delete odd_data_word;

  return {even, odd};
}

std::tuple<int64_t, int64_t> ALU::unsigned_extension(int64_t operand) {
  auto data_word = new abi::word::DataWord();
  data_word->set_value(operand);

  auto even_data_word = new abi::word::DataWord();
  even_data_word->set_value(0);
  int64_t even = even_data_word->value(abi::word::UNSIGNED);

  auto odd_data_word = new abi::word::DataWord();
  odd_data_word->set_value(data_word->value(abi::word::UNSIGNED));
  int64_t odd = odd_data_word->value(abi::word::UNSIGNED);

  delete data_word;
  delete even_data_word;
  delete odd_data_word;

  return {even, odd};
}

}  // namespace upmem_sim::simulator::dpu
