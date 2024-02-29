#include "converter/suffix_converter.h"

#include <stdexcept>

namespace upmem_sim::converter {

std::string SuffixConverter::to_string(abi::instruction::Suffix suffix) {
  if (suffix == abi::instruction::RICI) {
    return "rici";
  } else if (suffix == abi::instruction::RRI) {
    return "rri";
  } else if (suffix == abi::instruction::RRIC) {
    return "rric";
  } else if (suffix == abi::instruction::RRICI) {
    return "rrici";
  } else if (suffix == abi::instruction::RRIF) {
    return "rrif";
  } else if (suffix == abi::instruction::RRR) {
    return "rrr";
  } else if (suffix == abi::instruction::RRRC) {
    return "rrrc";
  } else if (suffix == abi::instruction::RRRCI) {
    return "rrrci";
  } else if (suffix == abi::instruction::ZRI) {
    return "zri";
  } else if (suffix == abi::instruction::ZRIC) {
    return "zric";
  } else if (suffix == abi::instruction::ZRICI) {
    return "zrici";
  } else if (suffix == abi::instruction::ZRIF) {
    return "zrif";
  } else if (suffix == abi::instruction::ZRR) {
    return "zrr";
  } else if (suffix == abi::instruction::ZRRC) {
    return "zrrc";
  } else if (suffix == abi::instruction::ZRRCI) {
    return "zrrci";
  } else if (suffix == abi::instruction::S_RRI) {
    return "s_rri";
  } else if (suffix == abi::instruction::S_RRIC) {
    return "s_rric";
  } else if (suffix == abi::instruction::S_RRICI) {
    return "s_rrici";
  } else if (suffix == abi::instruction::S_RRIF) {
    return "s_rrif";
  } else if (suffix == abi::instruction::S_RRR) {
    return "s_rrr";
  } else if (suffix == abi::instruction::S_RRRC) {
    return "s_rrrc";
  } else if (suffix == abi::instruction::S_RRRCI) {
    return "s_rrrci";
  } else if (suffix == abi::instruction::U_RRI) {
    return "u_rri";
  } else if (suffix == abi::instruction::U_RRIC) {
    return "u_rric";
  } else if (suffix == abi::instruction::U_RRICI) {
    return "u_rrici";
  } else if (suffix == abi::instruction::U_RRIF) {
    return "u_rrif";
  } else if (suffix == abi::instruction::U_RRR) {
    return "u_rrr";
  } else if (suffix == abi::instruction::U_RRRC) {
    return "u_rrrc";
  } else if (suffix == abi::instruction::U_RRRCI) {
    return "u_rrrci";
  } else if (suffix == abi::instruction::RR) {
    return "rr";
  } else if (suffix == abi::instruction::RRC) {
    return "rrc";
  } else if (suffix == abi::instruction::RRCI) {
    return "rrci";
  } else if (suffix == abi::instruction::ZR) {
    return "zr";
  } else if (suffix == abi::instruction::ZRC) {
    return "zrc";
  } else if (suffix == abi::instruction::ZRCI) {
    return "zrci";
  } else if (suffix == abi::instruction::S_RR) {
    return "s_rr";
  } else if (suffix == abi::instruction::S_RRC) {
    return "s_rrc";
  } else if (suffix == abi::instruction::S_RRCI) {
    return "s_rrci";
  } else if (suffix == abi::instruction::U_RR) {
    return "u_rr";
  } else if (suffix == abi::instruction::U_RRC) {
    return "u_rrc";
  } else if (suffix == abi::instruction::U_RRCI) {
    return "u_rrci";
  } else if (suffix == abi::instruction::DRDICI) {
    return "drdici";
  } else if (suffix == abi::instruction::RRRI) {
    return "rrri";
  } else if (suffix == abi::instruction::RRRICI) {
    return "rrrici";
  } else if (suffix == abi::instruction::ZRRI) {
    return "zrri";
  } else if (suffix == abi::instruction::ZRRICI) {
    return "zrrici";
  } else if (suffix == abi::instruction::S_RRRI) {
    return "s_rrri";
  } else if (suffix == abi::instruction::S_RRRICI) {
    return "s_rrrici";
  } else if (suffix == abi::instruction::U_RRRI) {
    return "u_rrri";
  } else if (suffix == abi::instruction::U_RRRICI) {
    return "u_rrrici";
  } else if (suffix == abi::instruction::RIR) {
    return "rir";
  } else if (suffix == abi::instruction::RIRC) {
    return "rirc";
  } else if (suffix == abi::instruction::RIRCI) {
    return "rirci";
  } else if (suffix == abi::instruction::ZIR) {
    return "zir";
  } else if (suffix == abi::instruction::ZIRC) {
    return "zirc";
  } else if (suffix == abi::instruction::ZIRCI) {
    return "zirci";
  } else if (suffix == abi::instruction::S_RIRC) {
    return "s_rirc";
  } else if (suffix == abi::instruction::S_RIRCI) {
    return "s_rirci";
  } else if (suffix == abi::instruction::U_RIRC) {
    return "u_rirc";
  } else if (suffix == abi::instruction::U_RIRCI) {
    return "u_rirci";
  } else if (suffix == abi::instruction::R) {
    return "r";
  } else if (suffix == abi::instruction::RCI) {
    return "rci";
  } else if (suffix == abi::instruction::Z) {
    return "z";
  } else if (suffix == abi::instruction::ZCI) {
    return "zci";
  } else if (suffix == abi::instruction::S_R) {
    return "s_r";
  } else if (suffix == abi::instruction::S_RCI) {
    return "s_rci";
  } else if (suffix == abi::instruction::U_R) {
    return "u_r";
  } else if (suffix == abi::instruction::U_RCI) {
    return "u_rci";
  } else if (suffix == abi::instruction::CI) {
    return "ci";
  } else if (suffix == abi::instruction::I) {
    return "i";
  } else if (suffix == abi::instruction::DDCI) {
    return "ddci";
  } else if (suffix == abi::instruction::ERRI) {
    return "erri";
  } else if (suffix == abi::instruction::S_ERRI) {
    return "s_erri";
  } else if (suffix == abi::instruction::U_ERRI) {
    return "u_erri";
  } else if (suffix == abi::instruction::EDRI) {
    return "edri";
  } else if (suffix == abi::instruction::ERII) {
    return "erii";
  } else if (suffix == abi::instruction::ERIR) {
    return "erir";
  } else if (suffix == abi::instruction::ERID) {
    return "erid";
  } else if (suffix == abi::instruction::DMA_RRI) {
    return "dma_rri";
  } else {
    throw std::invalid_argument("");
  }
}

}  // namespace upmem_sim::converter
