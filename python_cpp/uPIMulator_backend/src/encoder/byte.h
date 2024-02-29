#ifndef UPMEM_SIM_ENCODER_BYTE_H_
#define UPMEM_SIM_ENCODER_BYTE_H_

#include <cassert>

namespace upmem_sim::encoder {

class Byte {
 public:
  explicit Byte(int value) : value_(value) {
    assert(0 <= value and value < 256);
  }
  ~Byte() = default;

  int value() { return value_; }

 private:
  int value_;
};

}  // namespace upmem_sim::encoder

#endif
