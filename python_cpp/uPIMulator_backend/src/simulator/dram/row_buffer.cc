#include "simulator/dram/row_buffer.h"

namespace upmem_sim::simulator::dram {

RowBuffer::RowBuffer(util::ArgumentParser *argument_parser)
    : mram_(nullptr),
      row_address_(nullptr),
      input_q_(new basic::Queue<MemoryCommand>(1)),
      ready_q_(new basic::Queue<MemoryCommand>(-1)),
      stat_factory_(new util::StatFactory("row_buffer")) {
  timing_parameters_["t_ras"] =
      static_cast<int>(argument_parser->get_int_parameter("t_ras"));
  timing_parameters_["t_rcd"] =
      static_cast<int>(argument_parser->get_int_parameter("t_rcd"));
  timing_parameters_["t_cl"] =
      static_cast<int>(argument_parser->get_int_parameter("t_cl"));
  timing_parameters_["t_bl"] =
      static_cast<int>(argument_parser->get_int_parameter("t_bl"));
  timing_parameters_["t_rp"] =
      static_cast<int>(argument_parser->get_int_parameter("t_rp"));

  assert(timing_parameters_["t_ras"] > 0);
  assert(timing_parameters_["t_rcd"] > 0);
  assert(timing_parameters_["t_cl"] > 0);
  assert(timing_parameters_["t_bl"] > 0);
  assert(timing_parameters_["t_rp"] > 0);

  wordline_size_ = argument_parser->get_int_parameter("wordline_size");

  assert(wordline_size_ > 0);

  activation_q_ =
      new basic::TimerQueue<MemoryCommand>(1, timing_parameters_["t_ras"]);
  io_q_ = new basic::TimerQueue<MemoryCommand>(1, timing_parameters_["t_cl"]);
  bus_q_ = new basic::TimerQueue<MemoryCommand>(1, timing_parameters_["t_bl"]);
  precharge_q_ =
      new basic::TimerQueue<MemoryCommand>(1, timing_parameters_["t_rp"]);
}

RowBuffer::~RowBuffer() {
  delete row_address_;
  delete input_q_;
  delete ready_q_;
  delete activation_q_;
  delete io_q_;
  delete bus_q_;
  delete precharge_q_;

  delete stat_factory_;
}

util::StatFactory *RowBuffer::stat_factory() {
  auto stat_factory = new util::StatFactory("");

  stat_factory->merge(stat_factory_);

  return stat_factory;
}

void RowBuffer::connect_mram(MRAM *mram) {
  assert(mram != nullptr);
  assert(mram_ == nullptr);

  mram_ = mram;
}

void RowBuffer::push(MemoryCommand *memory_command) {
  assert(memory_command != nullptr);
  assert(can_push());
  input_q_->push(memory_command);
}

MemoryCommand *RowBuffer::pop() {
  assert(can_pop());
  return ready_q_->pop();
}

void RowBuffer::flush() {
  if (row_address_ != nullptr) {
    write_to_mram();
    delete row_address_;
    row_address_ = nullptr;
  }
}

void RowBuffer::cycle() {
  service_input_q();
  service_activation_q();
  service_io_q();
  service_bus_q();
  service_precharge_q();

  activation_q_->cycle();
  io_q_->cycle();
  bus_q_->cycle();
  precharge_q_->cycle();
}

void RowBuffer::service_input_q() {
  if (input_q_->can_pop()) {
    MemoryCommand *memory_command = input_q_->front();

    if (memory_command->operation() == MemoryCommand::ACTIVATION) {
      if (activation_q_->empty() and row_address_ == nullptr) {
        activation_q_->push(memory_command);
        input_q_->pop();
      }
    } else if (memory_command->operation() == MemoryCommand::READ) {
      if (io_q_->can_push() and row_address_ != nullptr) {
        io_q_->push(memory_command);
        input_q_->pop();
      }
    } else if (memory_command->operation() == MemoryCommand::WRITE) {
      if (io_q_->can_push() and row_address_ != nullptr) {
        io_q_->push(memory_command);
        input_q_->pop();
      }
    } else if (memory_command->operation() == MemoryCommand::PRECHARGE) {
      if (activation_q_->empty() and io_q_->empty() and bus_q_->empty() and
          precharge_q_->empty()) {
        precharge_q_->push(memory_command);
        input_q_->pop();
      }
    } else {
      throw std::invalid_argument("");
    }
  }
}

void RowBuffer::service_activation_q() {
  auto [memory_command, cycle] = activation_q_->front();
  if (cycle == timing_parameters_["t_ras"] - timing_parameters_["t_rcd"]) {
    assert(row_address_ == nullptr);

    assert(memory_command->address() % wordline_size_ == 0);
    row_address_ = new abi::word::DataAddressWord();
    row_address_->set_value(memory_command->address());

    row_buffer_ = read_from_mram();
  }

  if (activation_q_->can_pop() and ready_q_->can_push()) {
    activation_q_->pop();
    ready_q_->push(memory_command);

    stat_factory_->increment("num_activations");
  }
}

void RowBuffer::service_io_q() {
  if (io_q_->can_pop() and bus_q_->can_push()) {
    MemoryCommand *memory_command = io_q_->pop();
    bus_q_->push(memory_command);
  }
}

void RowBuffer::service_bus_q() {
  if (bus_q_->can_pop() and ready_q_->can_push()) {
    MemoryCommand *memory_command = bus_q_->pop();
    ready_q_->push(memory_command);

    if (memory_command->operation() == MemoryCommand::READ) {
      std::vector<int> bytes = read_from_row_buffer(memory_command->address(),
                                                    memory_command->size());
      memory_command->set_bytes(bytes);

      stat_factory_->increment("num_reads");
      if (memory_command->dma_command()->has_instruction()) {
        stat_factory_->increment(
            std::to_string(
                memory_command->dma_command()->instruction()->thread()->id()) +
            "_num_reads");
      }

      stat_factory_->increment("read_bytes", memory_command->size());
      if (memory_command->dma_command()->has_instruction()) {
        stat_factory_->increment(
            std::to_string(
                memory_command->dma_command()->instruction()->thread()->id()) +
                "_read_bytes",
            memory_command->size());
      }
    } else if (memory_command->operation() == MemoryCommand::WRITE) {
      write_to_row_buffer(memory_command->address(), memory_command->size(),
                          memory_command->bytes());

      stat_factory_->increment("num_writes");
      if (memory_command->dma_command()->has_instruction()) {
        stat_factory_->increment(
            std::to_string(
                memory_command->dma_command()->instruction()->thread()->id()) +
            "_num_writes");
      }

      stat_factory_->increment("write_bytes", memory_command->size());
      if (memory_command->dma_command()->has_instruction()) {
        stat_factory_->increment(
            std::to_string(
                memory_command->dma_command()->instruction()->thread()->id()) +
                "_write_bytes",
            memory_command->size());
      }
    } else {
      throw std::invalid_argument("");
    }
  }
}

void RowBuffer::service_precharge_q() {
  if (precharge_q_->can_pop() and ready_q_->can_push()) {
    MemoryCommand *memory_command = precharge_q_->pop();

    assert(memory_command->address() % wordline_size_ == 0);
    assert(memory_command->address() == row_address_->address());

    write_to_mram();
    delete row_address_;
    row_address_ = nullptr;
    ready_q_->push(memory_command);

    stat_factory_->increment("num_precharges");
  }
}

std::vector<int> RowBuffer::read_from_mram() {
  assert(row_address_ != nullptr);
  return std::move(mram_->read(row_address_->address()));
}

std::vector<int> RowBuffer::read_from_row_buffer(Address address,
                                                 Address size) {
  assert(row_address_ != nullptr);

  std::vector<int> bytes;
  bytes.resize(size);
  std::copy(row_buffer_.begin() + index(address),
            row_buffer_.begin() + index(address + size), bytes.begin());
  return std::move(bytes);
}

void RowBuffer::write_to_mram() {
  assert(row_address_ != nullptr);

  mram_->write(row_address_->address(), std::move(row_buffer_));
}

void RowBuffer::write_to_row_buffer(Address address, Address size,
                                    std::vector<int> bytes) {
  assert(size == bytes.size());

  std::copy(bytes.begin(), bytes.end(), row_buffer_.begin() + index(address));
}

int RowBuffer::index(Address address) {
  assert(row_address_->address() <= address and
         address <= row_address_->address() + wordline_size_);

  return static_cast<int>(address - row_address_->address());
}

}  // namespace upmem_sim::simulator::dram
