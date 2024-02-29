#include "abi/word/_base_word.h"

#include <algorithm>
#include <cassert>
#include <cmath>

namespace upmem_sim::abi::word {

_BaseWord::_BaseWord(int width) {
  assert(width > 0);
  bits_.resize(width);
}

int64_t _BaseWord::bit_slice(Representation representation, int begin,
                             int end) {
  assert(0 <= begin and begin < end and end <= width());
  assert(end - begin < 64);

  int slice_width = end - begin;
  int64_t value = 0;
  for (int i = 0; i < slice_width; i++) {
    if (bit(begin + i)) {
      if (representation == SIGNED and i == slice_width - 1) {
        value -= static_cast<int64_t>(pow(2, i));
      } else {
        value += static_cast<int64_t>(pow(2, i));
      }
    }
  }
  return value;
}

void _BaseWord::set_bit_slice(int begin, int end, int64_t value) {
  assert(0 <= begin and begin < end and end <= width());
  assert(end - begin < 64);

  if (value >= 0) {
    set_positive_bit_slice(begin, end, value);
  } else {
    set_negative_bit_slice(begin, end, value);
  }
}

encoder::ByteStream *_BaseWord::to_byte_stream() {
  int num_bytes = ceil(width() / 8.0);
  auto byte_stream = new encoder::ByteStream();
  for (int i = 0; i < num_bytes; i++) {
    int begin = 8 * i;
    int end = std::min(begin + 8, width());

    auto byte = static_cast<int>(bit_slice(UNSIGNED, begin, end));
    byte_stream->append(byte);
  }
  return byte_stream;
}

void _BaseWord::from_byte_stream(encoder::ByteStream *byte_stream) {
  for (int i = 0; i < byte_stream->size(); i++) {
    int begin = 8 * i;
    int end = std::min(begin + 8, width());

    int64_t byte = byte_stream->byte(i);
    set_bit_slice(begin, end, byte);
  }
}

void _BaseWord::set_positive_bit_slice(int begin, int end, int64_t value) {
  assert(0 <= begin and begin < end and end <= width());
  assert(value >= 0);

  int slice_width = end - begin;
  for (int i = 0; i < slice_width; i++) {
    if (value % 2) {
      set_bit(begin + i);
    } else {
      clear_bit(begin + i);
    }

    value /= 2;
  }

  assert(value == 0);
}

void _BaseWord::set_negative_bit_slice(int begin, int end, int64_t value) {
  assert(0 <= begin and begin < end and end <= width());
  assert(value < 0);

  set_bit(end - 1);

  if (begin + 1 < end) {
    int slice_width = end - begin;
    value += static_cast<int64_t>(pow(2, slice_width - 1));
    set_positive_bit_slice(begin, end - 1, value);
  }
}

}  // namespace upmem_sim::abi::word
