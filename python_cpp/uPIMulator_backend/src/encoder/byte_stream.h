#ifndef UPMEM_SIM_ENCODER_BYTE_STREAM_H_
#define UPMEM_SIM_ENCODER_BYTE_STREAM_H_

#include <string>
#include <vector>

#include "encoder/byte.h"
#include "main.h"

namespace upmem_sim::encoder {

class ByteStream {
 public:
  explicit ByteStream() = default;
  explicit ByteStream(std::vector<int> bytes);
  explicit ByteStream(std::string filename);

  ~ByteStream();

  Address size() { return static_cast<int>(bytes_.size()); }
  int byte(int index) { return bytes_[index]->value(); }
  std::vector<int> bytes();

  void append(int value) { bytes_.push_back(new Byte(value)); }
  void merge(ByteStream *byte_stream);

  ByteStream *slice(int begin, int end);

 private:
  std::vector<Byte *> bytes_;
};

}  // namespace upmem_sim::encoder

#endif
