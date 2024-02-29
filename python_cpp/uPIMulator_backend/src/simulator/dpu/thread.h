#ifndef UPMEM_SIM_SIMULATOR_DPU_THREAD_H_
#define UPMEM_SIM_SIMULATOR_DPU_THREAD_H_

#include <cassert>
#include <map>

#include "simulator/reg/reg_file.h"
#include "util/config_loader.h"

namespace upmem_sim::simulator::dpu {

using ThreadStatus = std::string;

class Thread {
 public:
  enum State { EMBRYO = 0, RUNNABLE, SLEEP, BLOCK, ZOMBIE };

  explicit Thread(ThreadID id)
      : id_(id),
        state_(EMBRYO),
        reg_file_(new reg::RegFile(id_)),
        issue_cycle_(0) {
    assert(0 <= id and id < upmem_sim::util::ConfigLoader::max_num_tasklets());
    status_tracker_.emplace("WAIT_DATA", 0);
    status_tracker_.emplace("WAIT_SYNC", 0);
    status_tracker_.emplace("ARITHMETIC", 0);
    status_tracker_.emplace("SPM_ACCESS", 0);
    status_tracker_.emplace("WAIT_SCHEDULE", 0);
  }
  ~Thread();

  ThreadID id() { return id_; }
  State state() { return state_; }
  void set_state(State state) { state_ = state; }
  reg::RegFile *reg_file() { return reg_file_; }
  int issue_cycle() { return issue_cycle_; }
  void increment_issue_cycle() { issue_cycle_ += 1; }
  void reset_issue_cycle() { issue_cycle_ = 0; }

  void update_thread_status(ThreadStatus status, int64_t value) {
    status_tracker_[std::move(status)] += value;
  }
  std::map<ThreadStatus, int64_t> &status_tracker() { return status_tracker_; }

 private:
  ThreadID id_;
  State state_;
  std::map<ThreadStatus, int64_t> status_tracker_;
  reg::RegFile *reg_file_;
  int issue_cycle_;
};

}  // namespace upmem_sim::simulator::dpu

#endif
