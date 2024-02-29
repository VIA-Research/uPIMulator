#ifndef UPMEM_SIM_ISA_SUFFIX_H_
#define UPMEM_SIM_ISA_SUFFIX_H_

namespace upmem_sim::abi::instruction {

enum Suffix {
  RICI,

  RRI,
  RRIC,
  RRICI,
  RRIF,
  RRR,
  RRRC,
  RRRCI,

  ZRI,
  ZRIC,
  ZRICI,
  ZRIF,
  ZRR,
  ZRRC,
  ZRRCI,

  S_RRI,
  S_RRIC,
  S_RRICI,
  S_RRIF,
  S_RRR,
  S_RRRC,
  S_RRRCI,

  U_RRI,
  U_RRIC,
  U_RRICI,
  U_RRIF,
  U_RRR,
  U_RRRC,
  U_RRRCI,

  RR,
  RRC,
  RRCI,

  ZR,
  ZRC,
  ZRCI,

  S_RR,
  S_RRC,
  S_RRCI,

  U_RR,
  U_RRC,
  U_RRCI,

  DRDICI,

  RRRI,
  RRRICI,

  ZRRI,
  ZRRICI,

  S_RRRI,
  S_RRRICI,

  U_RRRI,
  U_RRRICI,

  RIR,
  RIRC,
  RIRCI,

  ZIR,
  ZIRC,
  ZIRCI,

  S_RIRC,
  S_RIRCI,

  U_RIRC,
  U_RIRCI,

  R,
  RCI,

  Z,
  ZCI,

  S_R,
  S_RCI,

  U_R,
  U_RCI,

  CI,
  I,

  DDCI,

  ERRI,

  S_ERRI,
  U_ERRI,

  EDRI,

  ERII,
  ERIR,
  ERID,

  DMA_RRI
};
}

#endif
