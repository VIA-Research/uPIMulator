#include "simulator/dpu/dpu.h"

#include <cmath>
#include <iostream>

namespace upmem_sim::simulator::dpu {

DPU::DPU(DPUID dpu_id, util::ArgumentParser *argument_parser)
    : dpu_id_(dpu_id),
      atomic_(new sram::Atomic()),
      iram_(new sram::IRAM()),
      wram_(new sram::WRAM()),
      mram_(new dram::MRAM(argument_parser)),
      logic_(new Logic(dpu_id, argument_parser)),
      dma_(new DMA()),
      operand_collector_(new OperandCollector()),
      memory_controller_(new dram::MemoryController(argument_parser)),
      stat_factory_(new util::StatFactory("DPU#" + std::to_string(dpu_id))) {
  int num_threads =
      static_cast<int>(argument_parser->get_int_parameter("num_tasklets"));

  assert(0 <= num_threads and
         num_threads <= util::ConfigLoader::max_num_tasklets());

  threads_.resize(num_threads);
  for (ThreadID id = 0; id < num_threads; id++) {
    threads_[id] = new Thread(id);
  }

  scheduler_ = new RevolverScheduler(argument_parser, threads_);

  logic_->connect_scheduler(scheduler_);
  logic_->connect_atomic(atomic_);
  logic_->connect_iram(iram_);
  logic_->connect_operand_collector(operand_collector_);
  logic_->connect_dma(dma_);

  dma_->connect_atomic(atomic_);
  dma_->connect_iram(iram_);
  dma_->connect_operand_collector(operand_collector_);
  dma_->connect_memory_controller(memory_controller_);

  operand_collector_->connect_wram(wram_);

  memory_controller_->connect_mram(mram_);

  logic_frequency_ =
      static_cast<int>(argument_parser->get_int_parameter("logic_frequency"));
  memory_frequency_ =
      static_cast<int>(argument_parser->get_int_parameter("memory_frequency"));
  frequency_ratio_ = static_cast<double>(memory_frequency_) /
                     static_cast<double>(logic_frequency_);
}

DPU::~DPU() {
  for (auto &thread : threads_) {
    delete thread;
  }

  delete scheduler_;
  delete atomic_;
  delete iram_;
  delete wram_;
  delete mram_;
  delete logic_;
  delete dma_;
  delete operand_collector_;
  delete memory_controller_;

  delete stat_factory_;
}

util::StatFactory *DPU::stat_factory() {
  auto stat_factory = new util::StatFactory("");

  util::StatFactory *logic_stat_factory = logic_->stat_factory();
  util::StatFactory *memory_stat_factory = memory_controller_->stat_factory();

  for (auto &thread : threads_) {
    std::map<ThreadStatus, int64_t> status_tracker = thread->status_tracker();
    for (const auto &stat : status_tracker) {
      stat_factory_->increment(
          std::to_string(thread->id()) + "_latency_breakdown_" + stat.first,
          stat.second);
    }
  }

  stat_factory->merge(stat_factory_);
  stat_factory->merge(logic_stat_factory);
  stat_factory->merge(memory_stat_factory);

  delete logic_stat_factory;
  delete memory_stat_factory;

  return stat_factory;
}

bool DPU::is_zombie() {
  for (auto &thread : threads_) {
    if (thread->state() != Thread::ZOMBIE) {
      return false;
    }
  }

  return logic_->empty() and memory_controller_->empty();
}

void DPU::cycle() {
  scheduler_->cycle();
  logic_->cycle();
  dma_->cycle();
  int num_memory_cycles = static_cast<int>(
      floor(frequency_ratio_ *
            static_cast<double>(stat_factory_->value("cycle"))) -
      floor(frequency_ratio_ *
            static_cast<double>(stat_factory_->value("cycle") - 1)));
  for (int i = 0; i < num_memory_cycles; i++) {
    memory_controller_->cycle();
  }

  stat_factory_->increment("cycle");
}

}  // namespace upmem_sim::simulator::dpu
