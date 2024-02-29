#ifndef UPMEM_SIM_SIMULATOR_DPU_ALU_H_
#define UPMEM_SIM_SIMULATOR_DPU_ALU_H_

#include <cstdint>
#include <tuple>

namespace upmem_sim::simulator::dpu {

class ALU {
 public:
  static int64_t atomic_address_hash(int64_t operand1, int64_t operand2);

  static std::tuple<int64_t, bool, bool> add(int64_t operand1,
                                             int64_t operand2);
  static std::tuple<int64_t, bool, bool> addc(int64_t operand1,
                                              int64_t operand2,
                                              bool carry_flag);

  static std::tuple<int64_t, bool, bool> sub(int64_t operand1,
                                             int64_t operand2);
  static std::tuple<int64_t, bool, bool> subc(int64_t operand1,
                                              int64_t operand2,
                                              bool carry_flag);

  static int64_t and_(int64_t operand1, int64_t operand2);
  static int64_t nand(int64_t operand1, int64_t operand2);
  static int64_t andn(int64_t operand1, int64_t operand2);
  static int64_t or_(int64_t operand1, int64_t operand2);
  static int64_t nor(int64_t operand1, int64_t operand2);
  static int64_t orn(int64_t operand1, int64_t operand2);
  static int64_t xor_(int64_t operand1, int64_t operand2);
  static int64_t nxor(int64_t operand1, int64_t operand2);

  static int64_t asr(int64_t operand, int64_t shift);
  static int64_t lsl(int64_t operand, int64_t shift);

  static std::tuple<int64_t, bool, bool> lsl_add(int64_t operand1,
                                                 int64_t operand2,
                                                 int64_t shift);
  static std::tuple<int64_t, bool, bool> lsl_sub(int64_t operand1,
                                                 int64_t operand2,
                                                 int64_t shift);

  static int64_t lsl1(int64_t operand, int64_t shift);
  static int64_t lsl1x(int64_t operand, int64_t shift);
  static int64_t lslx(int64_t operand, int64_t shift);

  static int64_t lsr(int64_t operand, int64_t shift);

  static std::tuple<int64_t, bool, bool> lsr_add(int64_t operand1,
                                                 int64_t operand2,
                                                 int64_t shift);

  static int64_t lsr1(int64_t operand, int64_t shift);
  static int64_t lsr1x(int64_t operand, int64_t shift);
  static int64_t lsrx(int64_t operand, int64_t shift);

  static int64_t rol(int64_t operand, int64_t shift);

  static std::tuple<int64_t, bool, bool> rol_add(int64_t operand1,
                                                 int64_t operand2,
                                                 int64_t shift);

  static int64_t ror(int64_t operand, int64_t shift);

  static int64_t cao(int64_t operand);
  static int64_t clo(int64_t operand);
  static int64_t cls(int64_t operand);
  static int64_t clz(int64_t operand);

  static int64_t cmpb4(int64_t operand1, int64_t operand2);

  static int64_t extsb(int64_t operand);
  static int64_t extsh(int64_t operand);
  static int64_t extub(int64_t operand);
  static int64_t extuh(int64_t operand);

  static int64_t mul_sh_sh(int64_t operand1, int64_t operand2);
  static int64_t mul_sh_sl(int64_t operand1, int64_t operand2);
  static int64_t mul_sh_uh(int64_t operand1, int64_t operand2);
  static int64_t mul_sh_ul(int64_t operand1, int64_t operand2);
  static int64_t mul_sl_sh(int64_t operand1, int64_t operand2);
  static int64_t mul_sl_sl(int64_t operand1, int64_t operand2);
  static int64_t mul_sl_uh(int64_t operand1, int64_t operand2);
  static int64_t mul_sl_ul(int64_t operand1, int64_t operand2);
  static int64_t mul_uh_uh(int64_t operand1, int64_t operand2);
  static int64_t mul_uh_ul(int64_t operand1, int64_t operand2);
  static int64_t mul_ul_uh(int64_t operand1, int64_t operand2);
  static int64_t mul_ul_ul(int64_t operand1, int64_t operand2);

  static int64_t sats(int64_t operand);
  static int64_t hash(int64_t operand1, int64_t operand2);

  static std::tuple<int64_t, int64_t> signed_extension(int64_t operand);
  static std::tuple<int64_t, int64_t> unsigned_extension(int64_t operand);
};

}  // namespace upmem_sim::simulator::dpu

#endif
