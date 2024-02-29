#include "simulator/dram/fifo_scheduler.h"

namespace upmem_sim::simulator::dram {

void FIFOScheduler::cycle() {
  service_input_q();
  service_output_q();
}

void FIFOScheduler::service_input_q() {
  if (input_q_->can_pop()) {
    dpu::DMACommand *dma_command = input_q_->pop();
    service_dma_command(dma_command);
  }
}

void FIFOScheduler::service_dma_command(dpu::DMACommand *dma_command) {
  Address begin_address = dma_command->mram_address();
  Address end_address = dma_command->mram_address() + dma_command->size();
  Address address = begin_address;
  while (address < end_address) {
    Address min_access_granularity =
        util::ConfigLoader::min_access_granularity();
    Address wordline_address =
        (address / wordline_size_) * wordline_size_ + wordline_size_;
    Address size =
        std::min(std::min(address + min_access_granularity, wordline_address),
                 end_address) -
        address;

    reorder_buffer_.push_back({dma_command, address, size});

    address += size;
  }
}

void FIFOScheduler::service_output_q() {
  if (not reorder_buffer_.empty()) {
    service_fcfs();
  }
}

bool FIFOScheduler::service_fcfs() {
  auto [dma_command, address, size] = reorder_buffer_[0];
  Address wordline_address = (address / wordline_size_) * wordline_size_;

  if (row_address_ == nullptr and ready_q_->can_push(2)) {
    reorder_buffer_.erase(reorder_buffer_.begin());

    ready_q_->push(
        new MemoryCommand(MemoryCommand::ACTIVATION, wordline_address));

    row_address_ = new abi::word::DataAddressWord();
    row_address_->set_value(wordline_address);

    if (dma_command->operation() == dpu::DMACommand::READ) {
      ready_q_->push(
          new MemoryCommand(MemoryCommand::READ, address, size, dma_command));
    } else if (dma_command->operation() == dpu::DMACommand::WRITE) {
      ready_q_->push(new MemoryCommand(MemoryCommand::WRITE, address, size,
                                       dma_command->bytes(address, size),
                                       dma_command));
    } else {
      throw std::invalid_argument("");
    }

    stat_factory_->increment("row_buffer_miss");

    return true;
  } else if (row_address_ != nullptr and
             row_address_->address() == wordline_address and
             ready_q_->can_push(1)) {
    reorder_buffer_.erase(reorder_buffer_.begin());

    if (dma_command->operation() == dpu::DMACommand::READ) {
      ready_q_->push(
          new MemoryCommand(MemoryCommand::READ, address, size, dma_command));
    } else if (dma_command->operation() == dpu::DMACommand::WRITE) {
      ready_q_->push(new MemoryCommand(MemoryCommand::WRITE, address, size,
                                       dma_command->bytes(address, size),
                                       dma_command));
    } else {
      throw std::invalid_argument("");
    }

    stat_factory_->increment("row_buffer_hit");

    return true;
  } else if (row_address_ != nullptr and
             row_address_->address() != wordline_address and
             ready_q_->can_push(3)) {
    reorder_buffer_.erase(reorder_buffer_.begin());

    ready_q_->push(
        new MemoryCommand(MemoryCommand::PRECHARGE, row_address_->address()));
    ready_q_->push(
        new MemoryCommand(MemoryCommand::ACTIVATION, wordline_address));

    row_address_->set_value(wordline_address);

    if (dma_command->operation() == dpu::DMACommand::READ) {
      ready_q_->push(
          new MemoryCommand(MemoryCommand::READ, address, size, dma_command));
    } else if (dma_command->operation() == dpu::DMACommand::WRITE) {
      ready_q_->push(new MemoryCommand(MemoryCommand::WRITE, address, size,
                                       dma_command->bytes(address, size),
                                       dma_command));
    } else {
      throw std::invalid_argument("");
    }

    stat_factory_->increment("row_buffer_miss");

    return true;
  } else {
    return false;
  }
}

}  // namespace upmem_sim::simulator::dram
