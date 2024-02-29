#include "simulator/dram/memory_controller.h"

#include "simulator/dram/fifo_scheduler.h"
#include "simulator/dram/frfcfs_scheduler.h"

namespace upmem_sim::simulator::dram {

MemoryController::MemoryController(util::ArgumentParser *argument_parser)
    : wordline_size_(argument_parser->get_int_parameter("wordline_size")),
      row_buffer_(new RowBuffer(argument_parser)),
      mram_(nullptr),
      input_q_(new basic::Queue<dpu::DMACommand>(-1)),
      wait_q_(new basic::Queue<dpu::DMACommand>(-1)),
      memory_command_q_(new basic::Queue<MemoryCommand>(1)),
      ready_q_(new basic::Queue<dpu::DMACommand>(-1)),
      stat_factory_(new util::StatFactory("MemoryController")) {
  std::string memory_scheduling_policy =
      argument_parser->get_string_parameter("memory_scheduling_policy");
  if (memory_scheduling_policy == "fifo") {
    scheduler_ = new FIFOScheduler(argument_parser);
  } else if (memory_scheduling_policy == "frfcfs") {
    scheduler_ = new FRFCFSScheduler(argument_parser);
  } else {
    throw std::invalid_argument("");
  }
}

MemoryController::~MemoryController() {
  delete scheduler_;
  delete row_buffer_;
  delete input_q_;
  delete wait_q_;
  delete memory_command_q_;
  delete ready_q_;

  delete stat_factory_;
}

util::StatFactory *MemoryController::stat_factory() {
  auto stat_factory = new util::StatFactory("");

  util::StatFactory *scheduler_stat_factory = scheduler_->stat_factory();
  util::StatFactory *row_buffer_stat_factory = row_buffer_->stat_factory();

  stat_factory->merge(stat_factory_);
  stat_factory->merge(scheduler_stat_factory);
  stat_factory->merge(row_buffer_stat_factory);

  delete scheduler_stat_factory;
  delete row_buffer_stat_factory;

  return stat_factory;
}

void MemoryController::connect_mram(MRAM *mram) {
  assert(mram != nullptr);
  assert(mram_ == nullptr);

  mram_ = mram;
  row_buffer_->connect_mram(mram);
}

void MemoryController::push(dpu::DMACommand *dma_command) {
  assert(dma_command != nullptr);
  input_q_->push(dma_command);
}

dpu::DMACommand *MemoryController::pop() {
  assert(can_pop());
  return ready_q_->pop();
}

dpu::DMACommand *MemoryController::front() {
  assert(can_pop());
  return ready_q_->front();
}

std::vector<int> MemoryController::read(Address address, Address size) {
  Address end_address = address + size;

  Address cur_address = address;

  std::vector<int> bytes;
  while (cur_address < end_address) {
    Address cur_wordline_address =
        (cur_address / wordline_size_) * wordline_size_;
    Address cur_size =
        std::min(cur_wordline_address + wordline_size_, end_address) -
        cur_address;
    int cur_offset = cur_address % wordline_size_;

    std::vector<int> mram_bytes = mram_->read(cur_wordline_address);

    bytes.insert(bytes.end(), mram_bytes.begin() + cur_offset,
                 mram_bytes.begin() + cur_offset + cur_size);

    cur_address += cur_size;
  }

  return std::move(bytes);
}

void MemoryController::write(Address address, Address size,
                             std::vector<int> bytes) {
  assert(bytes.size() == size);

  Address end_address = address + size;

  Address cur_address = address;
  int cur_bytes_offset = 0;
  while (cur_address < end_address) {
    Address cur_wordline_address =
        (cur_address / wordline_size_) * wordline_size_;
    Address cur_size =
        std::min(cur_wordline_address + wordline_size_, end_address) -
        cur_address;
    int cur_offset = cur_address % wordline_size_;

    std::vector<int> mram_bytes = mram_->read(cur_wordline_address);

    std::copy(bytes.begin() + cur_bytes_offset,
              bytes.begin() + cur_bytes_offset + cur_size,
              mram_bytes.begin() + cur_offset);

    mram_->write(cur_wordline_address, mram_bytes);

    cur_address += cur_size;
    cur_bytes_offset += cur_size;
  }
}

void MemoryController::write(Address address, Address size,
                             encoder::ByteStream *byte_stream) {
  assert(byte_stream->size() == size);

  std::vector<int> bytes = byte_stream->bytes();

  Address end_address = address + size;

  Address cur_address = address;
  int cur_bytes_offset = 0;
  while (cur_address < end_address) {
    Address cur_wordline_address =
        (cur_address / wordline_size_) * wordline_size_;
    Address cur_size =
        std::min(cur_wordline_address + wordline_size_, end_address) -
        cur_address;
    int cur_offset = cur_address % wordline_size_;

    std::vector<int> mram_bytes = mram_->read(cur_wordline_address);

    std::copy(bytes.begin() + cur_bytes_offset,
              bytes.begin() + cur_bytes_offset + cur_size,
              mram_bytes.begin() + cur_offset);

    mram_->write(cur_wordline_address, mram_bytes);

    cur_address += cur_size;
    cur_bytes_offset += cur_size;
  }
}

void MemoryController::flush() {
  scheduler_->flush();
  row_buffer_->flush();
}

void MemoryController::cycle() {
  service_input_q();
  service_scheduler();
  service_memory_command_q();
  service_row_buffer();
  service_wait_q();

  scheduler_->cycle();
  row_buffer_->cycle();

  stat_factory_->increment("mem_cycle");
}

void MemoryController::service_input_q() {
  if (input_q_->can_pop() and scheduler_->can_push() and wait_q_->can_push()) {
    dpu::DMACommand *dma_command = input_q_->pop();
    scheduler_->push(dma_command);
    wait_q_->push(dma_command);
  }
}

void MemoryController::service_scheduler() {
  if (scheduler_->can_pop() and memory_command_q_->can_push()) {
    MemoryCommand *memory_command = scheduler_->pop();
    memory_command_q_->push(memory_command);
  }
}

void MemoryController::service_memory_command_q() {
  if (memory_command_q_->can_pop() and row_buffer_->can_push()) {
    MemoryCommand *memory_command = memory_command_q_->pop();
    row_buffer_->push(memory_command);
  }
}

void MemoryController::service_row_buffer() {
  if (row_buffer_->can_pop()) {
    MemoryCommand *memory_command = row_buffer_->pop();

    if (memory_command->operation() == MemoryCommand::ACTIVATION or
        memory_command->operation() == MemoryCommand::PRECHARGE) {
      delete memory_command;
    } else if (memory_command->operation() == MemoryCommand::READ) {
      memory_command->dma_command()->set_bytes(memory_command->address(),
                                               memory_command->size(),
                                               memory_command->bytes());
      memory_command->dma_command()->ack_bytes(memory_command->address(),
                                               memory_command->size());
      delete memory_command;
    } else if (memory_command->operation() == MemoryCommand::WRITE) {
      memory_command->dma_command()->ack_bytes(memory_command->address(),
                                               memory_command->size());
      delete memory_command;
    } else {
      throw std::invalid_argument("");
    }
  }
}

void MemoryController::service_wait_q() {
  if (not wait_q_->empty()) {
    dpu::DMACommand *dma_command = wait_q_->front();

    if (dma_command->is_ready() and ready_q_->can_push()) {
      wait_q_->pop();
      ready_q_->push(dma_command);
    }
  }
}

}  // namespace upmem_sim::simulator::dram
