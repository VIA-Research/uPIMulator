#ifndef UPMEM_SIM_ABI_WORD__BASE_WORD_H_
#define UPMEM_SIM_ABI_WORD__BASE_WORD_H_

#include <cstdint>
#include <vector>

#include "abi/word/representation.h"
#include "encoder/byte_stream.h"
#include "main.h"

namespace upmem_sim::abi::word {

class _BaseWord {
 public:
  explicit _BaseWord(int width);
  ~_BaseWord() = default;

  int width() { return static_cast<int>(bits_.size()); }
  Address size() { return width() / 8; }

  bool sign_bit() { return bit(width() - 1); }
  bool bit(int index) { return bits_[index]; }
  void set_bit(int index) { bits_[index] = true; }
  void clear_bit(int index) { bits_[index] = false; }

  int64_t bit_slice(Representation representation, int begin, int end);
  void set_bit_slice(int begin, int end, int64_t value);

  int64_t value(Representation representation) {
    return bit_slice(representation, 0, width());
  }
  void set_value(int64_t value) { set_bit_slice(0, width(), value); }

  encoder::ByteStream *to_byte_stream();
  void from_byte_stream(encoder::ByteStream *byte_stream);

 protected:
  void set_positive_bit_slice(int begin, int end, int64_t value);
  void set_negative_bit_slice(int begin, int end, int64_t value);

 private:
  std::vector<bool> bits_;
};

}  // namespace upmem_sim::abi::word

#endif
