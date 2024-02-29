#ifndef UPMEM_SIM_SIMULATOR_CPU_THREAD_H_
#define UPMEM_SIM_SIMULATOR_CPU_THREAD_H_

#include "simulator/dpu/dpu.h"

namespace upmem_sim::simulator::cpu {

class Thread {
 public:
  explicit Thread(util::ArgumentParser *argument_parser);
  ~Thread() = default;

  std::string benchmark() { return benchmark_; }
  int num_dpus() { return num_dpus_; }
  int num_tasklets() { return num_tasklets_; }

  Address sys_used_mram_end_pointer() { return sys_used_mram_end_pointer_; }
  Address dpu_input_arguments_pointer() { return dpu_input_arguments_pointer_; }
  Address dpu_results_pointer() { return dpu_results_pointer_; }
  Address sys_end_pointer() { return sys_end_pointer_; }

  int num_executions() { return num_executions_; }

 protected:
  encoder::ByteStream *load_byte_stream(std::string filename);

  void init_dpu_transfer_pointer();
  void init_num_executions();

 private:
  std::string bindir_;
  std::string benchmark_;

  int num_dpus_;
  int num_tasklets_;

  Address sys_used_mram_end_pointer_;
  Address dpu_input_arguments_pointer_;
  Address dpu_results_pointer_;
  Address sys_end_pointer_;

  int num_executions_;
};

}  // namespace upmem_sim::simulator::cpu

#endif
