#ifndef UPMEM_SIM_SIMULATOR_BASIC_QUEUE_H_
#define UPMEM_SIM_SIMULATOR_BASIC_QUEUE_H_

#include <cassert>
#include <iostream>
#include <queue>

namespace upmem_sim::simulator::basic {

template <typename T>
class Queue {
 public:
  explicit Queue(int size) : size_(size) { assert(size != 0); }
  ~Queue() { assert(q_.empty()); }

  bool empty() { return q_.empty(); }
  int size() { return q_.size(); }

  bool can_push(int num_items);
  bool can_push() { return can_push(1); }
  void push(T *item);
  bool can_pop() { return not q_.empty(); }
  T *pop();
  void cycle() = delete;

  T *front();

 private:
  int size_;
  std::queue<T *> q_;
};

template <typename T>
bool Queue<T>::can_push(int num_items) {
  if (size_ >= 0) {
    return size_ - q_.size() >= num_items;
  } else {
    return true;
  }
}

template <typename T>
void Queue<T>::push(T *item) {
  assert(can_push());

  q_.push(item);
}

template <typename T>
T *Queue<T>::pop() {
  assert(can_pop());

  T *item = q_.front();
  q_.pop();
  return item;
}

template <typename T>
T *Queue<T>::front() {
  if (q_.empty()) {
    return nullptr;
  } else {
    return q_.front();
  }
}

}  // namespace upmem_sim::simulator::basic

#endif
