#include "simulator/rank/rank.h"

namespace upmem_sim::simulator::rank {

Rank::Rank(util::ArgumentParser* argument_parser)
    : read_bandwidth_(
          argument_parser->get_int_parameter("rank_read_bandwidth")),
      write_bandwidth_(
          argument_parser->get_int_parameter("rank_write_bandwidth")),
      stat_factory_(new util::StatFactory("Rank")) {
  assert(read_bandwidth_ > 0);
  assert(write_bandwidth_ > 0);

  int num_dpus =
      static_cast<int>(argument_parser->get_int_parameter("num_dpus"));
  dpus_.resize(num_dpus);
  communication_qs_.resize(num_dpus);
  for (DPUID dpu_id = 0; dpu_id < num_dpus; dpu_id++) {
    dpus_[dpu_id] = new dpu::DPU(dpu_id, argument_parser);
    communication_qs_[dpu_id] = new basic::TimerQueue<RankMessage>(-1);
  }
}

Rank::~Rank() {
  for (auto& dpu : dpus_) {
    delete dpu;
  }

  for (auto & communication_q : communication_qs_) {
    delete communication_q;
  }

  delete stat_factory_;
}

util::StatFactory* Rank::stat_factory() {
  auto stat_factory = new util::StatFactory("");

  stat_factory->merge(stat_factory_);

  for (auto& dpu : dpus_) {
    util::StatFactory* dpu_stat_factory = dpu->stat_factory();

    stat_factory->merge(dpu_stat_factory);

    delete dpu_stat_factory;
  }

  return stat_factory;
}

void Rank::launch() {
  for (auto& dpu : dpus_) {
    for (auto& thread : dpu->scheduler()->threads()) {
      Address bootstrap = util::ConfigLoader::iram_offset();
      thread->reg_file()->write_pc_reg(bootstrap);
    }
    dpu->boot();
  }
}

bool Rank::is_zombie() {
  for (auto& dpu : dpus_) {
    if (not dpu->is_zombie()) {
      return false;
    }
  }
  return true;
}

void Rank::read(upmem_sim::simulator::rank::RankMessage*rank_message) {
  assert(rank_message->operation() == RankMessage::READ);
  assert(0 <= rank_message->dpu_id() and rank_message->dpu_id() < dpus_.size());

  basic::TimerQueue<RankMessage>* communication_q = communication_qs_[rank_message->dpu_id()];
  assert(communication_q->can_push());

  SimTime latency = static_cast<SimTime>(10 * rank_message->size() / read_bandwidth_);
  communication_q->push(rank_message, latency);

  stat_factory_->increment("num_reads");
  stat_factory_->increment("read_bytes", rank_message->size());
}

void Rank::write(RankMessage*rank_message) {
  assert(rank_message->operation() == RankMessage::WRITE);
  assert(0 <= rank_message->dpu_id() and rank_message->dpu_id() < dpus_.size());


  basic::TimerQueue<RankMessage>* communication_q = communication_qs_[rank_message->dpu_id()];
  assert(communication_q->can_push());

  SimTime latency = static_cast<SimTime>(10 * rank_message->size() / write_bandwidth_);
  communication_q->push(rank_message, latency);

  stat_factory_->increment("num_writes");
  stat_factory_->increment("write_bytes", rank_message->size());
}

void Rank::cycle() {
  service_sequence_q();

  for (auto& dpu : dpus_) {
    dpu->cycle();
  }

  bool is_communication_q_empty = true;
  for (auto & communication_q : communication_qs_) {
    communication_q->cycle();

    if (not communication_q->empty()) {
      is_communication_q_empty = false;
    }
  }

  if (not is_communication_q_empty) {
    stat_factory_->increment("communication_cycle");
  }

  stat_factory_->increment("rank_cycle");
}

void Rank::service_sequence_q() {
  for (auto &communication_q : communication_qs_) {
    if (communication_q->can_pop()) {
      RankMessage * rank_message = communication_q->pop();
      rank_message->set_ack();
    }
  }
}

}  // namespace upmem_sim::simulator::rank
