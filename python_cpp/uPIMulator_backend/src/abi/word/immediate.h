#ifndef UPMEM_SIM_ABI_WORD_IMMEDIATE_H_
#define UPMEM_SIM_ABI_WORD_IMMEDIATE_H_

#include <cstdint>

#include "abi/word/_base_word.h"
#include "abi/word/representation.h"
#include "encoder/byte_stream.h"

namespace upmem_sim::abi::word {

class Immediate {
 public:
  Immediate(Representation representation, int width, int64_t value)
      : representation_(representation), word_(new _BaseWord(width)) {
    word_->set_value(value);
  }
  ~Immediate() { delete word_; }

  Representation representation() { return representation_; }
  int width() { return word_->width(); }

  bool bit(int index) { return word_->bit(index); }
  int64_t bit_slice(int begin, int end) {
    return word_->bit_slice(representation_, begin, end);
  }
  int64_t value() { return word_->value(representation_); }
  encoder::ByteStream *to_byte_stream() { return word_->to_byte_stream(); }

 private:
  Representation representation_;
  _BaseWord *word_;
};

}  // namespace upmem_sim::abi::word

#endif
