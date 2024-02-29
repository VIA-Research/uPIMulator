#ifndef UPMEM_SIM_SIMULATOR_DRAM_MRAM_H_
#define UPMEM_SIM_SIMULATOR_DRAM_MRAM_H_

#include <vector>

#include "abi/word/data_address_word.h"
#include "simulator/dram/wordline.h"
#include "util/argument_parser.h"

namespace upmem_sim::simulator::dram {

class MRAM {
 public:
  explicit MRAM(util::ArgumentParser *argument_parser);
  ~MRAM();

  Address address() { return address_->address(); }
  Address size() { return size_; }

  std::vector<int> read(Address address);

  void write(Address address, std::vector<int> bytes);
  void write(Address address, encoder::ByteStream *byte_stream);

  void cycle() = delete;

 protected:
  int index(Address address);

 private:
  abi::word::DataAddressWord *address_;
  Address size_;
  std::vector<Wordline *> wordlines_;

  Address wordline_size_;
  int num_wordlines_;
};

}  // namespace upmem_sim::simulator::dram

#endif
