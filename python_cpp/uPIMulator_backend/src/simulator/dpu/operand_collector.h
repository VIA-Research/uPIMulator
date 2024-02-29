#ifndef UPMEM_SIM_SIMULATOR_DPU_OPERAND_COLLECTOR_H_
#define UPMEM_SIM_SIMULATOR_DPU_OPERAND_COLLECTOR_H_

#include "simulator/sram/wram.h"

namespace upmem_sim::simulator::dpu {

class OperandCollector {
 public:
  explicit OperandCollector() : wram_(nullptr) {}
  ~OperandCollector() = default;

  void connect_wram(sram::WRAM *wram);

  int64_t lbs(Address address);
  int64_t lbu(Address address);

  int64_t lhs(Address address);
  int64_t lhu(Address address);

  int64_t lw(Address address);
  std::tuple<int64_t, int64_t> ld(Address address);

  void sb(Address address, int64_t value);
  void sh(Address address, int64_t value);
  void sw(Address address, int64_t value);
  void sd(Address address, int64_t even, int64_t odd);

  void cycle() = delete;

 private:
  sram::WRAM *wram_;
};

}  // namespace upmem_sim::simulator::dpu

#endif
