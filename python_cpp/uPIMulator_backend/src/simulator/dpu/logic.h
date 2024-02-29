#ifndef UPMEM_SIM_SIMULATOR_DPU_LOGIC_H_
#define UPMEM_SIM_SIMULATOR_DPU_LOGIC_H_

#include "simulator/dpu/cycle_rule.h"
#include "simulator/dpu/dma.h"
#include "simulator/dpu/operand_collector.h"
#include "simulator/dpu/pipeline.h"
#include "simulator/dpu/revolver_scheduler.h"
#include "simulator/dram/memory_controller.h"
#include "simulator/sram/atomic.h"
#include "simulator/sram/iram.h"
#include "simulator/sram/wram.h"
#include "util/argument_parser.h"
#include "util/stat_factory.h"

namespace upmem_sim::simulator::dpu {

class Logic {
 public:
  explicit Logic(DPUID dpu_id, util::ArgumentParser *argument_parser)
      : verbose_(argument_parser->get_int_parameter("verbose")),
        dpu_id_(dpu_id),
        scheduler_(nullptr),
        atomic_(nullptr),
        iram_(nullptr),
        dma_(nullptr),
        pipeline_(new Pipeline(argument_parser)),
        cycle_rule_(new CycleRule(argument_parser)),
        operand_collector_(nullptr),
        wait_instruction_q_(new basic::Queue<abi::instruction::Instruction>(
            util::ConfigLoader::max_num_tasklets())),
        stat_factory_(new util::StatFactory("Logic")),
        num_pipeline_stages_(
            argument_parser->get_int_parameter("num_pipeline_stages")) {}
  ~Logic();

  DPUID dpu_id() { return dpu_id_; }

  util::StatFactory *stat_factory();

  void connect_scheduler(RevolverScheduler *scheduler);
  void connect_atomic(sram::Atomic *atomic);
  void connect_iram(sram::IRAM *iram);
  void connect_operand_collector(OperandCollector *operand_collector);
  void connect_dma(DMA *dma);

  bool empty() {
    return pipeline_->empty() and cycle_rule_->empty() and
           wait_instruction_q_->empty();
  }
  void cycle();

 protected:
  void service_scheduler();
  void service_pipeline();
  void service_cycle_rule();
  void service_logic();
  void service_dma();

  void execute_instruction(abi::instruction::Instruction *instruction);

  void execute_rici(abi::instruction::Instruction *instruction);
  void execute_acquire_rici(abi::instruction::Instruction *instruction);
  void execute_release_rici(abi::instruction::Instruction *instruction);
  void execute_boot_rici(abi::instruction::Instruction *instruction);

  void execute_rri(abi::instruction::Instruction *instruction);
  void execute_add_rri(abi::instruction::Instruction *instruction);
  void execute_asr_rri(abi::instruction::Instruction *instruction);
  void execute_call_rri(abi::instruction::Instruction *instruction);

  void execute_rric(abi::instruction::Instruction *instruction);
  void execute_add_rric(abi::instruction::Instruction *instruction);
  void execute_asr_rric(abi::instruction::Instruction *instruction);
  void execute_sub_rric(abi::instruction::Instruction *instruction);

  void execute_rrici(abi::instruction::Instruction *instruction);
  void execute_add_rrici(abi::instruction::Instruction *instruction);
  void execute_and_rrici(abi::instruction::Instruction *instruction);
  void execute_asr_rrici(abi::instruction::Instruction *instruction);
  void execute_sub_rrici(abi::instruction::Instruction *instruction);

  void execute_rrif(abi::instruction::Instruction *instruction);

  void execute_rrr(abi::instruction::Instruction *instruction);

  void execute_rrrc(abi::instruction::Instruction *instruction);
  void execute_add_rrrc(abi::instruction::Instruction *instruction);
  void execute_rsub_rrrc(abi::instruction::Instruction *instruction);
  void execute_sub_rrrc(abi::instruction::Instruction *instruction);

  void execute_rrrci(abi::instruction::Instruction *instruction);
  void execute_add_rrrci(abi::instruction::Instruction *instruction);
  void execute_and_rrrci(abi::instruction::Instruction *instruction);
  void execute_asr_rrrci(abi::instruction::Instruction *instruction);
  void execute_mul_rrrci(abi::instruction::Instruction *instruction);
  void execute_rsub_rrrci(abi::instruction::Instruction *instruction);

  void execute_zri(abi::instruction::Instruction *instruction);
  void execute_add_zri(abi::instruction::Instruction *instruction);
  void execute_asr_zri(abi::instruction::Instruction *instruction);
  void execute_call_zri(abi::instruction::Instruction *instruction);

  void execute_zric(abi::instruction::Instruction *instruction);
  void execute_add_zric(abi::instruction::Instruction *instruction);
  void execute_asr_zric(abi::instruction::Instruction *instruction);
  void execute_sub_zric(abi::instruction::Instruction *instruction);

  void execute_zrici(abi::instruction::Instruction *instruction);
  void execute_add_zrici(abi::instruction::Instruction *instruction);
  void execute_and_zrici(abi::instruction::Instruction *instruction);
  void execute_asr_zrici(abi::instruction::Instruction *instruction);
  void execute_sub_zrici(abi::instruction::Instruction *instruction);

  void execute_zrif(abi::instruction::Instruction *instruction);

  void execute_zrr(abi::instruction::Instruction *instruction);

  void execute_zrrc(abi::instruction::Instruction *instruction);
  void execute_add_zrrc(abi::instruction::Instruction *instruction);
  void execute_rsub_zrrc(abi::instruction::Instruction *instruction);
  void execute_sub_zrrc(abi::instruction::Instruction *instruction);

  void execute_zrrci(abi::instruction::Instruction *instruction);
  void execute_add_zrrci(abi::instruction::Instruction *instruction);
  void execute_and_zrrci(abi::instruction::Instruction *instruction);
  void execute_asr_zrrci(abi::instruction::Instruction *instruction);
  void execute_mul_zrrci(abi::instruction::Instruction *instruction);
  void execute_rsub_zrrci(abi::instruction::Instruction *instruction);

  void execute_s_rri(abi::instruction::Instruction *instruction);
  void execute_add_s_rri(abi::instruction::Instruction *instruction);
  void execute_asr_s_rri(abi::instruction::Instruction *instruction);

  void execute_s_rric(abi::instruction::Instruction *instruction);

  void execute_s_rrici(abi::instruction::Instruction *instruction);
  void execute_add_s_rrici(abi::instruction::Instruction *instruction);
  void execute_and_s_rrici(abi::instruction::Instruction *instruction);
  void execute_asr_s_rrici(abi::instruction::Instruction *instruction);
  void execute_sub_s_rrici(abi::instruction::Instruction *instruction);

  void execute_s_rrif(abi::instruction::Instruction *instruction);

  void execute_s_rrr(abi::instruction::Instruction *instruction);

  void execute_s_rrrc(abi::instruction::Instruction *instruction);

  void execute_s_rrrci(abi::instruction::Instruction *instruction);

  void execute_u_rri(abi::instruction::Instruction *instruction);
  void execute_add_u_rri(abi::instruction::Instruction *instruction);
  void execute_asr_u_rri(abi::instruction::Instruction *instruction);

  void execute_u_rric(abi::instruction::Instruction *instruction);

  void execute_u_rrici(abi::instruction::Instruction *instruction);
  void execute_add_u_rrici(abi::instruction::Instruction *instruction);
  void execute_and_u_rrici(abi::instruction::Instruction *instruction);
  void execute_asr_u_rrici(abi::instruction::Instruction *instruction);
  void execute_sub_u_rrici(abi::instruction::Instruction *instruction);

  void execute_u_rrif(abi::instruction::Instruction *instruction);

  void execute_u_rrr(abi::instruction::Instruction *instruction);

  void execute_u_rrrc(abi::instruction::Instruction *instruction);

  void execute_u_rrrci(abi::instruction::Instruction *instruction);

  void execute_rr(abi::instruction::Instruction *instruction);

  void execute_rrc(abi::instruction::Instruction *instruction);

  void execute_rrci(abi::instruction::Instruction *instruction);
  void execute_cao_rrci(abi::instruction::Instruction *instruction);
  void execute_extsb_rrci(abi::instruction::Instruction *instruction);
  void execute_time_cfg_rrci(abi::instruction::Instruction *instruction);

  void execute_zr(abi::instruction::Instruction *instruction);

  void execute_zrc(abi::instruction::Instruction *instruction);

  void execute_zrci(abi::instruction::Instruction *instruction);
  void execute_cao_zrci(abi::instruction::Instruction *instruction);
  void execute_extsb_zrci(abi::instruction::Instruction *instruction);
  void execute_time_cfg_zrci(abi::instruction::Instruction *instruction);

  void execute_s_rr(abi::instruction::Instruction *instruction);
  void execute_s_rrc(abi::instruction::Instruction *instruction);
  void execute_s_rrci(abi::instruction::Instruction *instruction);

  void execute_u_rr(abi::instruction::Instruction *instruction);
  void execute_u_rrc(abi::instruction::Instruction *instruction);
  void execute_u_rrci(abi::instruction::Instruction *instruction);

  void execute_drdici(abi::instruction::Instruction *instruction);
  void execute_div_step_drdici(abi::instruction::Instruction *instruction);
  void execute_mul_step_drdici(abi::instruction::Instruction *instruction);

  void execute_rrri(abi::instruction::Instruction *instruction);
  void execute_rrrici(abi::instruction::Instruction *instruction);

  void execute_zrri(abi::instruction::Instruction *instruction);
  void execute_zrrici(abi::instruction::Instruction *instruction);

  void execute_s_rrri(abi::instruction::Instruction *instruction);
  void execute_s_rrrici(abi::instruction::Instruction *instruction);

  void execute_u_rrri(abi::instruction::Instruction *instruction);
  void execute_u_rrrici(abi::instruction::Instruction *instruction);

  void execute_rir(abi::instruction::Instruction *instruction);
  void execute_rirc(abi::instruction::Instruction *instruction);
  void execute_rirci(abi::instruction::Instruction *instruction);

  void execute_zir(abi::instruction::Instruction *instruction);
  void execute_zirc(abi::instruction::Instruction *instruction);
  void execute_zirci(abi::instruction::Instruction *instruction);

  void execute_s_rirc(abi::instruction::Instruction *instruction);
  void execute_s_rirci(abi::instruction::Instruction *instruction);

  void execute_u_rirc(abi::instruction::Instruction *instruction);
  void execute_u_rirci(abi::instruction::Instruction *instruction);

  void execute_r(abi::instruction::Instruction *instruction);
  void execute_rci(abi::instruction::Instruction *instruction);

  void execute_z(abi::instruction::Instruction *instruction);
  void execute_zci(abi::instruction::Instruction *instruction);

  void execute_s_r(abi::instruction::Instruction *instruction);
  void execute_s_rci(abi::instruction::Instruction *instruction);

  void execute_u_r(abi::instruction::Instruction *instruction);
  void execute_u_rci(abi::instruction::Instruction *instruction);

  void execute_ci(abi::instruction::Instruction *instruction);
  void execute_i(abi::instruction::Instruction *instruction);

  void execute_ddci(abi::instruction::Instruction *instruction);
  void execute_movd_ddci(abi::instruction::Instruction *instruction);
  void execute_swapd_ddci(abi::instruction::Instruction *instruction);

  void execute_erri(abi::instruction::Instruction *instruction);
  void execute_s_erri(abi::instruction::Instruction *instruction);
  void execute_u_erri(abi::instruction::Instruction *instruction);
  void execute_edri(abi::instruction::Instruction *instruction);

  void execute_erii(abi::instruction::Instruction *instruction);
  void execute_erir(abi::instruction::Instruction *instruction);
  void execute_erid(abi::instruction::Instruction *instruction);

  void execute_dma_rri(abi::instruction::Instruction *instruction);
  void execute_ldma(abi::instruction::Instruction *instruction);
  void execute_ldmai(abi::instruction::Instruction *instruction);
  void execute_sdma(abi::instruction::Instruction *instruction);

  void set_acquire_cc(abi::instruction::Instruction *instruction,
                      int64_t result);
  void set_add_nz_cc(abi::instruction::Instruction *instruction,
                     int64_t operand1, int64_t result, bool carry,
                     bool overflow);
  void set_boot_cc(abi::instruction::Instruction *instruction, int64_t operand1,
                   int64_t result);
  void set_count_nz_cc(abi::instruction::Instruction *instruction,
                       int64_t operand1, int64_t result);
  void set_div_cc(abi::instruction::Instruction *instruction, int64_t operand1);
  void set_div_nz_cc(abi::instruction::Instruction *instruction,
                     int64_t operand1);
  void set_ext_sub_set_cc(abi::instruction::Instruction *instruction,
                          int64_t operand1, int64_t operand2, int64_t result,
                          bool carry, bool overflow);
  void set_imm_shift_nz_cc(abi::instruction::Instruction *instruction,
                           int64_t operand1, int64_t result);
  void set_log_nz_cc(abi::instruction::Instruction *instruction,
                     int64_t operand1, int64_t result);
  void set_log_set_cc(abi::instruction::Instruction *instruction,
                      int64_t result);
  void set_mul_nz_cc(abi::instruction::Instruction *instruction,
                     int64_t operand1, int64_t result);
  void set_sub_nz_cc(abi::instruction::Instruction *instruction,
                     int64_t operand1, int64_t operand2, int64_t result,
                     bool carry, bool overflow);
  void set_sub_set_cc(abi::instruction::Instruction *instruction,
                      int64_t operand1, int64_t operand2, int64_t result);

  void set_flags(abi::instruction::Instruction *instruction, int64_t result,
                 bool carry);

 private:
  DPUID dpu_id_;
  int verbose_;

  RevolverScheduler *scheduler_;
  sram::Atomic *atomic_;
  sram::IRAM *iram_;
  DMA *dma_;

  Pipeline *pipeline_;
  int num_pipeline_stages_;
  CycleRule *cycle_rule_;
  OperandCollector *operand_collector_;

  basic::Queue<abi::instruction::Instruction> *wait_instruction_q_;

  util::StatFactory *stat_factory_;
};

}  // namespace upmem_sim::simulator::dpu

#endif
