#ifndef UPMEM_SIM_SIMULATOR_DRAM_WORDLINE_H_
#define UPMEM_SIM_SIMULATOR_DRAM_WORDLINE_H_

#include <vector>

#include "abi/word/data_address_word.h"
#include "abi/word/data_word.h"
#include "util/argument_parser.h"

namespace upmem_sim::simulator::dram {

class Wordline {
 public:
  explicit Wordline(util::ArgumentParser *argument_parser, Address address);
  ~Wordline();

  Address address() { return address_->address(); }
  Address size() { return size_; }

  std::vector<int> read();

  void write(std::vector<int> bytes);
  void write(encoder::ByteStream *byte_stream);

  void cycle() = delete;

 protected:
  int num_data_words() {
    return static_cast<int>(size_ / abi::word::DataWord().size());
  }

 private:
  abi::word::DataAddressWord *address_;
  Address size_;
  std::vector<abi::word::DataWord *> data_words_;
};

}  // namespace upmem_sim::simulator::dram

#endif
