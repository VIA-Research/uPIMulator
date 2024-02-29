#ifndef UPMEM_SIM_ABI_ISA_INSTRUCTION_CONDITION_H_
#define UPMEM_SIM_ABI_ISA_INSTRUCTION_CONDITION_H_

namespace upmem_sim::abi::isa {

enum Condition {
  TRUE,
  FALSE,

  Z,
  NZ,

  E,
  O,

  PL,
  MI,

  OV,
  NOV,

  C,
  NC,

  SZ,
  SNZ,

  SPL,
  SMI,

  SO,
  SE,

  NC5,
  NC6,
  NC7,
  NC8,
  NC9,
  NC10,
  NC11,
  NC12,
  NC13,
  NC14,

  MAX,
  NMAX,

  SH32,
  NSH32,

  EQ,
  NEQ,

  LTU,
  LEU,
  GTU,
  GEU,

  LTS,
  LES,
  GTS,
  GES,

  XZ,
  XNZ,

  XLEU,
  XGTU,

  XLES,
  XGTS,

  SMALL,
  LARGE
};
}

#endif
