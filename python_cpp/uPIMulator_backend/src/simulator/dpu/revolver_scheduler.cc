#include "simulator/dpu/revolver_scheduler.h"

namespace upmem_sim::simulator::dpu {

RevolverScheduler::RevolverScheduler(util::ArgumentParser *argument_parser,
                                     std::vector<Thread *> threads)
    : thread_q_(new basic::Queue<Thread>(threads.size())),
      stat_factory_(new util::StatFactory("RevolverScheduler")) { 
  num_revolver_scheduling_cycles_ = static_cast<int>(
      argument_parser->get_int_parameter("num_revolver_scheduling_cycles"));

  assert(num_revolver_scheduling_cycles_ > 0);

  threads_ = threads;

  for (auto &thread : threads) {
    thread_q_->push(thread);
  }
}

RevolverScheduler::~RevolverScheduler() {
  while (thread_q_->can_pop()) {
    thread_q_->pop();
  }

  delete stat_factory_;
}

util::StatFactory *RevolverScheduler::stat_factory() { 
  auto stat_factory = new util::StatFactory("");
  stat_factory->merge(stat_factory_);
  return stat_factory;
}

Thread *RevolverScheduler::schedule() {
  bool is_blocked = false;
  for (int i = 0; i < thread_q_->size(); i++) {
    Thread *thread = thread_q_->pop();
    thread_q_->push(thread);

    if (thread->issue_cycle() >= num_revolver_scheduling_cycles_) {
      if (thread->state() == Thread::RUNNABLE) {
        thread->reset_issue_cycle();

        stat_factory_->increment("breakdown_run");

        return thread;
      } else if (thread->state() == Thread::BLOCK) {
        is_blocked = true;
      }
    }
  }

  if (is_blocked) {
    stat_factory_->increment("breakdown_dma");
  } else {
    stat_factory_->increment("breakdown_etc");
  }

  return nullptr;
}

bool RevolverScheduler::boot(ThreadID id) {
  Thread *thread = threads_[id];
  assert(thread->id() == id);

  if (thread->state() == Thread::EMBRYO) {
    thread->set_state(Thread::RUNNABLE);
    return true;
  } else if (thread->state() == Thread::ZOMBIE) {
    thread->set_state(Thread::RUNNABLE);
    return true;
  } else {
    throw std::invalid_argument("");
  }
}

bool RevolverScheduler::sleep(ThreadID id) {
  Thread *thread = threads_[id];
  assert(thread->id() == id);

  if (thread->state() == Thread::RUNNABLE) {
    thread->set_state(Thread::SLEEP);
    return true;
  } else {
    throw std::invalid_argument("");
  }
}

bool RevolverScheduler::block(ThreadID id) {
  Thread *thread = threads_[id];
  assert(thread->id() == id);

  if (thread->state() == Thread::RUNNABLE) {
    thread->set_state(Thread::BLOCK);
    return true;
  } else {
    throw std::invalid_argument("");
  }
}

bool RevolverScheduler::awake(ThreadID id) {
  Thread *thread = threads_[id];
  assert(thread->id() == id);

  if (thread->state() == Thread::EMBRYO) {
    thread->set_state(Thread::RUNNABLE);
    return true;
  } else if (thread->state() == Thread::SLEEP) {
    thread->set_state(Thread::RUNNABLE);
    return true;
  } else if (thread->state() == Thread::BLOCK) {
    thread->set_state(Thread::RUNNABLE);
    return true;
  } else {
    throw std::invalid_argument("");
  }
}

bool RevolverScheduler::shutdown(ThreadID id) {
  Thread *thread = threads_[id];
  assert(thread->id() == id);

  if (thread->state() == Thread::SLEEP) {
    thread->set_state(Thread::ZOMBIE);
    return true;
  } else {
    throw std::invalid_argument("");
  }
}

void RevolverScheduler::cycle() {
  int num_active_tasklets = 0, num_embryo = 0, num_sleep = 0, num_block = 0,
      num_zombie = 0;

  for (auto &thread : threads_) {
    if (thread->state() == Thread::RUNNABLE and
        thread->issue_cycle() < num_revolver_scheduling_cycles_) {
      stat_factory_->increment("revolver_wait");                 
      stat_factory_->increment(std::to_string(thread->id()) +
                               "_revolver_wait");
    }

    thread->increment_issue_cycle();

    if (thread->state() == Thread::RUNNABLE)
      num_active_tasklets++;
    else if (thread->state() == Thread::EMBRYO)
      num_embryo++;
    else if (thread->state() == Thread::SLEEP)
      num_sleep++;
    else if (thread->state() == Thread::BLOCK)
      num_block++;
    else if (thread->state() == Thread::ZOMBIE)
      num_zombie++;
    else
      assert(0);
  }
  assert(num_active_tasklets <= 16);
  stat_factory_->overwrite("current_active_tasklets",
                           num_active_tasklets); 

  stat_factory_->increment("active_tasklets_" +
                           std::to_string(num_active_tasklets));

  stat_factory_->increment("total_EMBRYO", num_embryo);
  stat_factory_->increment("total_RUNNABLE", num_active_tasklets);
  stat_factory_->increment("total_SLEEP", num_sleep);
  stat_factory_->increment("total_BLOCK", num_block);
  stat_factory_->increment("total_ZOMBIE", num_zombie);

  issuable_threads_ = num_active_tasklets;
}

}  // namespace upmem_sim::simulator::dpu