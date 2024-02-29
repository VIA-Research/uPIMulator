#include "simulator/sram/iram.h"

#include "encoder/instruction_encoder.h"

namespace upmem_sim::simulator::sram {

IRAM::IRAM() {
  address_ = new abi::word::InstructionAddressWord();
  address_->set_value(util::ConfigLoader::iram_offset());

  size_ = util::ConfigLoader::iram_size();

  assert(address() % abi::word::InstructionWord().size() == 0);
  assert(size_ % abi::word::InstructionWord().size() == 0);

  cells_.resize(num_instruction_words());
  for (int i = 0; i < num_instruction_words(); i++) {
    cells_[i] = new abi::word::InstructionWord();
  }
}

IRAM::~IRAM() {
  delete address_;

  for (int i = 0; i < size_ / num_instruction_words(); i++) {
    delete cells_[i];
  }
}

abi::instruction::Instruction *IRAM::read(Address address) {
  encoder::ByteStream *byte_stream = cells_[index(address)]->to_byte_stream();
  abi::instruction::Instruction *instruction =
      encoder::InstructionEncoder::decode(byte_stream);
  delete byte_stream;
  return instruction;
}

void IRAM::write(Address address, encoder::ByteStream *byte_stream) {
  cells_[index(address)]->from_byte_stream(byte_stream);
}

int IRAM::index(Address address) {
  assert(address >= this->address());
  assert(address + abi::word::InstructionWord().size() <=
         this->address() + size_);
  assert((address - this->address()) % abi::word::InstructionWord().size() ==
         0);

  return static_cast<int>((address - this->address()) /
                          abi::word::InstructionWord().size());
}

}  // namespace upmem_sim::simulator::sram
