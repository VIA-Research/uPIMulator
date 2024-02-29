#include "encoder/byte_stream.h"

#include <fstream>

namespace upmem_sim::encoder {

ByteStream::ByteStream(std::vector<int> bytes) {
  bytes_.resize(bytes.size());
  for (int i = 0; i < bytes_.size(); i++) {
    bytes_[i] = new Byte(bytes[i]);
  }
}

ByteStream::ByteStream(std::string filename) {
  std::ifstream ifs(filename);
  int byte;
  while (ifs >> byte) {
    append(byte);
  }
}

ByteStream::~ByteStream() {
  for (auto &byte : bytes_) {
    delete byte;
  }
}

std::vector<int> ByteStream::bytes() {
  std::vector<int> bytes;
  bytes.resize(bytes_.size());
  for (int i = 0; i < size(); i++) {
    bytes[i] = byte(i);
  }
  return std::move(bytes);
}

void ByteStream::merge(ByteStream *byte_stream) {
  int size = static_cast<int>(bytes_.size());
  bytes_.resize(bytes_.size() + byte_stream->bytes_.size());
  for (int i = 0; i < byte_stream->size(); i++) {
    bytes_[i + size] = new Byte(static_cast<int>(byte_stream->byte(i)));
  }
}

ByteStream *ByteStream::slice(int begin, int end) {
  assert(0 <= begin and begin < end and end <= size());

  auto *byte_stream = new ByteStream();
  for (int i = begin; i < end; i++) {
    byte_stream->append(byte(i));
  }
  return byte_stream;
}

}  // namespace upmem_sim::encoder
