#include "simulator/cpu/thread.h"

#include <filesystem>
#include <fstream>
#include <sstream>

namespace upmem_sim::simulator::cpu {

Thread::Thread(util::ArgumentParser *argument_parser)
    : bindir_(argument_parser->get_string_parameter("bindir")),
      benchmark_(argument_parser->get_string_parameter("benchmark")),
      num_dpus_(
          static_cast<int>(argument_parser->get_int_parameter("num_dpus"))),
      num_tasklets_(static_cast<int>(
          argument_parser->get_int_parameter("num_tasklets"))) {
  assert(0 < num_dpus_);
  assert(0 < num_tasklets_ and
         num_tasklets_ <= util::ConfigLoader::max_num_tasklets());

  init_dpu_transfer_pointer();
  init_num_executions();
}

encoder::ByteStream *Thread::load_byte_stream(std::string filename) {
  std::string bin_filepath = bindir_ + "/" + benchmark_ + "." +
                             std::to_string(num_tasklets_) + "/" + filename +
                             ".bin";
  if (std::filesystem::exists(bin_filepath)) {
    auto byte_stream = new encoder::ByteStream(bin_filepath);
    return byte_stream;
  } else {
    return nullptr;
  }
}

void Thread::init_dpu_transfer_pointer() {
  std::string bin_filepath = bindir_ + "/" + benchmark_ + "." +
                             std::to_string(num_tasklets_) +
                             "/dpu_transfer_pointer.bin";
  std::ifstream ifs(bin_filepath);

  ifs >> sys_used_mram_end_pointer_ >> dpu_input_arguments_pointer_ >>
      dpu_results_pointer_ >> sys_end_pointer_;
}

void Thread::init_num_executions() {
  std::string bin_filepath = bindir_ + "/" + benchmark_ + "." +
                             std::to_string(num_tasklets_) +
                             "/num_executions.bin";
  std::ifstream ifs(bin_filepath);

  ifs >> num_executions_;
}

}  // namespace upmem_sim::simulator::cpu
