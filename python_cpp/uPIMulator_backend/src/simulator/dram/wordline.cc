#include "simulator/dram/wordline.h"

namespace upmem_sim::simulator::dram {

Wordline::Wordline(util::ArgumentParser *argument_parser, Address address)
    : address_(new abi::word::DataAddressWord()),
      size_(argument_parser->get_int_parameter("wordline_size")) {
  assert(address >= util::ConfigLoader::mram_offset());
  assert(address + size_ <=
         util::ConfigLoader::mram_offset() + util::ConfigLoader::mram_size());
  assert(address % size_ == 0);
  assert(size_ % util::ConfigLoader::min_access_granularity() == 0);
  assert(size_ % abi::word::DataWord().size() == 0);

  address_->set_value(address);

  data_words_.resize(num_data_words());
  for (int i = 0; i < num_data_words(); i++) {
    data_words_[i] = new abi::word::DataWord();
  }
}

Wordline::~Wordline() {
  delete address_;

  for (int i = 0; i < num_data_words(); i++) {
    delete data_words_[i];
  }
}

std::vector<int> Wordline::read() {
  auto byte_stream = new encoder::ByteStream();

  for (int i = 0; i < num_data_words(); i++) {
    auto data_word_byte_stream = data_words_[i]->to_byte_stream();
    byte_stream->merge(data_word_byte_stream);
    delete data_word_byte_stream;
  }

  std::vector<int> bytes = byte_stream->bytes();

  delete byte_stream;

  return std::move(bytes);
}

void Wordline::write(std::vector<int> bytes) {
  assert(bytes.size() == size_);

  Address data_word_size = abi::word::DataWord().size();
  for (int i = 0; i < num_data_words(); i++) {
    for (int j = 0; j < data_word_size; j++) {
      int index = static_cast<int>(i * data_word_size + j);

      data_words_[i]->set_bit_slice(8 * j, 8 * (j + 1), bytes[index]);
    }
  }
}

void Wordline::write(encoder::ByteStream *byte_stream) {
  assert(byte_stream->size() == size_);

  for (int i = 0; i < num_data_words(); i++) {
    encoder::ByteStream *slice = byte_stream->slice(8 * i, 8 * (i + 1));
    data_words_[i]->from_byte_stream(slice);
    delete slice;
  }
}

}  // namespace upmem_sim::simulator::dram
