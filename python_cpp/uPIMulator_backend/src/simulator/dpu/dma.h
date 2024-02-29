#ifndef UPMEM_SIM_SIMULATOR_DPU_DMA_H_
#define UPMEM_SIM_SIMULATOR_DPU_DMA_H_

#include "simulator/dpu/operand_collector.h"
#include "simulator/dram/memory_controller.h"
#include "simulator/sram/atomic.h"
#include "simulator/sram/iram.h"
#include "simulator/sram/wram.h"

namespace upmem_sim::simulator::dpu {

class DMA {
 public:
  explicit DMA()
      : atomic_(nullptr),
        iram_(nullptr),
        operand_collector_(nullptr),
        memory_controller_(nullptr),
        input_q_(new basic::Queue<DMACommand>(
            util::ConfigLoader::max_num_tasklets())),
        ready_q_(new basic::Queue<DMACommand>(
            util::ConfigLoader::max_num_tasklets())) {}
  ~DMA();

  void connect_atomic(sram::Atomic *atomic);
  void connect_iram(sram::IRAM *iram);
  void connect_operand_collector(OperandCollector *operand_collector);
  void connect_memory_controller(dram::MemoryController *memory_controller);

  void transfer_to_atomic(Address address, encoder::ByteStream *byte_stream);

  void transfer_to_iram(Address address, encoder::ByteStream *byte_stream);

  encoder::ByteStream *transfer_from_wram(Address address, Address size);
  void transfer_to_wram(Address address, encoder::ByteStream *byte_stream);

  encoder::ByteStream *transfer_from_mram(Address address, Address size);
  void transfer_to_mram(Address address, encoder::ByteStream *byte_stream);

  void transfer_from_wram_to_mram(Address wram_address, Address mram_address,
                                  Address size,
                                  abi::instruction::Instruction *instruction);
  void transfer_from_mram_to_wram(Address wram_address, Address mram_address,
                                  Address size,
                                  abi::instruction::Instruction *instruction);

  bool can_push() { return input_q_->can_push(); }
  void push(DMACommand *dma_command) = delete;
  bool can_pop() { return ready_q_->can_pop(); }
  DMACommand *pop() { return ready_q_->pop(); }
  void cycle();

 protected:
  void service_input_q();
  void service_ready_q();

 private:
  sram::Atomic *atomic_;
  sram::IRAM *iram_;
  OperandCollector *operand_collector_;
  dram::MemoryController *memory_controller_;

  basic::Queue<DMACommand> *input_q_;
  basic::Queue<DMACommand> *ready_q_;
};

}  // namespace upmem_sim::simulator::dpu

#endif
