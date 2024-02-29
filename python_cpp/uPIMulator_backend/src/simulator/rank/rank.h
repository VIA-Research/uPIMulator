#ifndef UPMEM_SIM_SIMULATOR_RANK_RANK_H_
#define UPMEM_SIM_SIMULATOR_RANK_RANK_H_

#include <vector>

#include "main.h"
#include "simulator/basic/timer_queue.h"
#include "simulator/dpu/dpu.h"
#include "simulator/rank/rank_message.h"

namespace upmem_sim::simulator::rank {

class Rank {
 public:
  explicit Rank(util::ArgumentParser *argument_parser);
  ~Rank();

  util::StatFactory *stat_factory();

  std::vector<dpu::DPU *> dpus() { return dpus_; }

  void launch();
  bool is_zombie();

  void read(RankMessage *rank_message);
  void write(RankMessage *rank_message);

  void cycle();

 protected:
  void service_sequence_q();

 private:
  Address read_bandwidth_;
  Address write_bandwidth_;

  std::vector<dpu::DPU *> dpus_;
  std::vector<basic::TimerQueue<RankMessage>*> communication_qs_;

  util::StatFactory *stat_factory_;
};

}  // namespace upmem_sim::simulator::rank

#endif
