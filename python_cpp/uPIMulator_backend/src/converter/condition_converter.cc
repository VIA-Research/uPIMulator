#include "converter/condition_converter.h"

#include <stdexcept>

namespace upmem_sim::converter {

std::string ConditionConverter::to_string(abi::isa::Condition condition) {
  if (condition == abi::isa::TRUE) {
    return "true";
  } else if (condition == abi::isa::FALSE) {
    return "false";
  } else if (condition == abi::isa::Z) {
    return "z";
  } else if (condition == abi::isa::NZ) {
    return "nz";
  } else if (condition == abi::isa::E) {
    return "e";
  } else if (condition == abi::isa::O) {
    return "o";
  } else if (condition == abi::isa::PL) {
    return "pl";
  } else if (condition == abi::isa::MI) {
    return "mi";
  } else if (condition == abi::isa::OV) {
    return "ov";
  } else if (condition == abi::isa::NOV) {
    return "nov";
  } else if (condition == abi::isa::C) {
    return "c";
  } else if (condition == abi::isa::NC) {
    return "nc";
  } else if (condition == abi::isa::SZ) {
    return "sz";
  } else if (condition == abi::isa::SNZ) {
    return "snz";
  } else if (condition == abi::isa::SPL) {
    return "spl";
  } else if (condition == abi::isa::SMI) {
    return "smi";
  } else if (condition == abi::isa::SO) {
    return "so";
  } else if (condition == abi::isa::SE) {
    return "se";
  } else if (condition == abi::isa::NC5) {
    return "nc5";
  } else if (condition == abi::isa::NC6) {
    return "nc6";
  } else if (condition == abi::isa::NC7) {
    return "nc7";
  } else if (condition == abi::isa::NC8) {
    return "nc8";
  } else if (condition == abi::isa::NC9) {
    return "nc9";
  } else if (condition == abi::isa::NC10) {
    return "nc10";
  } else if (condition == abi::isa::NC11) {
    return "nc11";
  } else if (condition == abi::isa::NC12) {
    return "nc12";
  } else if (condition == abi::isa::NC13) {
    return "nc13";
  } else if (condition == abi::isa::NC14) {
    return "nc14";
  } else if (condition == abi::isa::MAX) {
    return "max";
  } else if (condition == abi::isa::NMAX) {
    return "nmax";
  } else if (condition == abi::isa::SH32) {
    return "sh32";
  } else if (condition == abi::isa::NSH32) {
    return "nsh32";
  } else if (condition == abi::isa::EQ) {
    return "eq";
  } else if (condition == abi::isa::NEQ) {
    return "neq";
  } else if (condition == abi::isa::LTU) {
    return "ltu";
  } else if (condition == abi::isa::LEU) {
    return "leu";
  } else if (condition == abi::isa::GTU) {
    return "gtu";
  } else if (condition == abi::isa::GEU) {
    return "geu";
  } else if (condition == abi::isa::LTS) {
    return "lts";
  } else if (condition == abi::isa::LES) {
    return "les";
  } else if (condition == abi::isa::GTS) {
    return "gts";
  } else if (condition == abi::isa::GES) {
    return "ges";
  } else if (condition == abi::isa::XZ) {
    return "xz";
  } else if (condition == abi::isa::XNZ) {
    return "xnz";
  } else if (condition == abi::isa::XLEU) {
    return "xleu";
  } else if (condition == abi::isa::XGTU) {
    return "xgtu";
  } else if (condition == abi::isa::XLES) {
    return "xles";
  } else if (condition == abi::isa::XGTS) {
    return "xgts";
  } else if (condition == abi::isa::SMALL) {
    return "small";
  } else if (condition == abi::isa::LARGE) {
    return "large";
  } else {
    throw std::invalid_argument("");
  }
}

}  // namespace upmem_sim::converter
