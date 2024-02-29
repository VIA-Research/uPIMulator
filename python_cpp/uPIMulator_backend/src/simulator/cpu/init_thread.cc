#include "simulator/cpu/init_thread.h"

#include <iostream>

namespace upmem_sim::simulator::cpu {

void InitThread::connect_rank(rank::Rank *rank) {
  assert(rank != nullptr);
  assert(rank_ == nullptr);

  rank_ = rank;
}

void InitThread::init() {
  dma_transfer_to_atomic();
  dma_transfer_to_iram();
  dma_transfer_to_wram();
  dma_transfer_to_mram();

  std::cout << "init completed..." << std::endl;
}

void InitThread::launch() {
  rank_->launch();

  std::cout << "launch completed..." << std::endl;
}

void InitThread::dma_transfer_to_atomic() {
  auto byte_stream = load_byte_stream("atomic");
  for (auto &dpu : rank_->dpus()) {
    dpu->dma()->transfer_to_atomic(util::ConfigLoader::atomic_offset(),
                                   byte_stream);
  }
  delete byte_stream;

  std::cout << "DMA to atomic completed..." << std::endl;
}

void InitThread::dma_transfer_to_iram() {
  auto byte_stream = load_byte_stream("iram");
  for (auto &dpu : rank_->dpus()) {
    dpu->dma()->transfer_to_iram(util::ConfigLoader::iram_offset(),
                                 byte_stream);
  }
  delete byte_stream;

  std::cout << "DMA to IRAM completed..." << std::endl;
}

void InitThread::dma_transfer_to_wram() {
  auto byte_stream = load_byte_stream("wram");
  for (auto &dpu : rank_->dpus()) {
    dpu->dma()->transfer_to_wram(util::ConfigLoader::wram_offset(),
                                 byte_stream);
  }
  delete byte_stream;

  std::cout << "DMA to WRAM completed..." << std::endl;
}

void InitThread::dma_transfer_to_mram() {
  auto byte_stream = load_byte_stream("mram");
  for (auto &dpu : rank_->dpus()) {
    dpu->dma()->transfer_to_mram(util::ConfigLoader::mram_offset(),
                                 byte_stream);
  }
  delete byte_stream;

  std::cout << "DMA to MRAM completed..." << std::endl;
}

}  // namespace upmem_sim::simulator::cpu
