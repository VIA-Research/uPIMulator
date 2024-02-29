#ifndef UPMEM_SIM_SIMULATOR_RANK_SEQUENCE_MESSAGE_H_
#define UPMEM_SIM_SIMULATOR_RANK_SEQUENCE_MESSAGE_H_

#include "encoder/byte_stream.h"
#include "main.h"

namespace upmem_sim::simulator::rank {

class RankMessage {
 public:
  enum Operation { READ = 0, WRITE };

  explicit RankMessage(Operation operation, DPUID dpu_id, Address address,
                       Address size)
      : operation_(operation),
        dpu_id_(dpu_id),
        address_(address),
        size_(size),
        ack_(false) {}
  explicit RankMessage(Operation operation, DPUID dpu_id, Address address,
                       Address size, encoder::ByteStream* byte_stream)
      : operation_(operation),
        dpu_id_(dpu_id),
        address_(address),
        size_(size),
        byte_stream_(byte_stream),
        ack_(false) {
    assert(byte_stream->size() == size);
  }

  ~RankMessage() = default;

  Operation operation() { return operation_; }

  DPUID dpu_id() { return dpu_id_; }

  Address address() { return address_; }
  Address size() { return size_; }

  encoder::ByteStream* byte_stream() { return byte_stream_; }

  bool ack() { return ack_; }
  void set_ack() {
    assert(not ack_);
    ack_ = true;
  }

 private:
  Operation operation_;
  DPUID dpu_id_;
  Address address_;
  Address size_;
  encoder::ByteStream* byte_stream_;
  bool ack_;
};

}  // namespace upmem_sim::simulator::rank

#endif
