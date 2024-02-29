#ifndef UPMEM_SIM_SIMULATOR_DPU_DPU_H_
#define UPMEM_SIM_SIMULATOR_DPU_DPU_H_

#include "simulator/dpu/dma.h"
#include "simulator/dpu/logic.h"
#include "simulator/dpu/operand_collector.h"
#include "simulator/dram/memory_controller.h"
#include "simulator/dram/mram.h"
#include "simulator/sram/atomic.h"
#include "simulator/sram/iram.h"
#include "simulator/sram/wram.h"
#include "util/argument_parser.h"
#include "util/stat_factory.h"

namespace upmem_sim::simulator::dpu {

class DPU {
 public:
  explicit DPU(DPUID dpu_id, util::ArgumentParser *argument_parser);
  ~DPU();

  DPUID dpu_id() { return dpu_id_; }

  RevolverScheduler *scheduler() { return scheduler_; }
  DMA *dma() { return dma_; }

  util::StatFactory *stat_factory();

  bool is_zombie();
  void boot() { scheduler_->boot(0); }
  void cycle();

 private:
  DPUID dpu_id_;

  std::vector<Thread *> threads_;

  RevolverScheduler *scheduler_;
  sram::Atomic *atomic_;
  sram::IRAM *iram_;
  sram::WRAM *wram_;
  dram::MRAM *mram_;
  Logic *logic_;
  DMA *dma_;
  OperandCollector *operand_collector_;
  dram::MemoryController *memory_controller_;

  int logic_frequency_;
  int memory_frequency_;
  double frequency_ratio_;

  util::StatFactory *stat_factory_;
};

}  // namespace upmem_sim::simulator::dpu

#endif
