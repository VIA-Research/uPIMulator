#ifndef UPMEM_SIM_SIMULATOR_BASIC_TIMER_QUEUE_H_
#define UPMEM_SIM_SIMULATOR_BASIC_TIMER_QUEUE_H_

#include <cassert>
#include <tuple>
#include <vector>

#include "main.h"

namespace upmem_sim::simulator::basic {

template <typename T>
class TimerQueue {
 public:
  explicit TimerQueue(int size);
  explicit TimerQueue(int size, SimTime timer);
  ~TimerQueue() { assert(q_.empty()); }

  bool empty() { return q_.empty(); }
  int size() { return q_.size(); }

  bool can_push(int num_items);
  bool can_push() { return can_push(1); }
  void push(T *item);
  void push(T *item, SimTime timer);
  bool can_pop();
  T *pop();
  void cycle();

  std::tuple<T *, int> front();

 private:
  int size_;
  SimTime timer_;
  std::vector<std::tuple<T *, int>> q_;
};

template <typename T>
TimerQueue<T>::TimerQueue(int size) : size_(size), timer_(0) {
  assert(size != 0);
}

template <typename T>
TimerQueue<T>::TimerQueue(int size, SimTime timer) : size_(size), timer_(timer) {
  assert(size != 0);
  assert(timer > 0);
}

template <typename T>
bool TimerQueue<T>::can_push(int num_items) {
  if (size_ >= 0) {
    return size_ - q_.size() >= num_items;
  } else {
    return true;
  }
}

template <typename T>
void TimerQueue<T>::push(T *item) {
  assert(can_push());

  push(item, timer_);
}

template <typename T>
void TimerQueue<T>::push(T *item, SimTime timer) {
  assert(can_push());

  q_.push_back({item, timer});
}

template <typename T>
bool TimerQueue<T>::can_pop() {
  if (q_.empty()) {
    return false;
  } else {
    auto [item, timer] = q_[0];
    if (timer <= 0) {
      return true;
    } else {
      return false;
    }
  }
}

template <typename T>
T *TimerQueue<T>::pop() {
  assert(can_pop());
  auto [item, timer] = q_[0];
  q_.erase(q_.begin());
  return item;
}

template <typename T>
std::tuple<T *, int> TimerQueue<T>::front() {
  if (q_.empty()) {
    return {nullptr, 0};
  } else {
    return q_[0];
  }
}

template <typename T>
void TimerQueue<T>::cycle() {
  if (not q_.empty()) {
    auto &[item, timer] = q_[0];
    timer -= 1;
  }
}

}  // namespace upmem_sim::simulator::basic

#endif
