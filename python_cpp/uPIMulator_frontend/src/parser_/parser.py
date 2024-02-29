from typing import Set

import antlr4

from parser_.grammar.assemblyLexer import assemblyLexer
from parser_.grammar.assemblyParser import assemblyParser


class Parser:
    def __init__(self):
        pass

    @staticmethod
    def parse_lines(lines: str) -> assemblyParser.DocumentContext:
        input_stream = antlr4.InputStream(lines)
        lexer = assemblyLexer(input_stream)
        common_token_stream = antlr4.CommonTokenStream(lexer)
        parser = assemblyParser(common_token_stream)
        return parser.document()

    @staticmethod
    def preprocess(lines: str) -> str:
        lines = Parser._preprocess_directive(lines)
        lines = Parser._preprocess_instruction(lines)
        lines = Parser._preprocess_section(lines)
        lines = Parser._preprocess_operator(lines)
        return lines

    @staticmethod
    def _preprocess_directive(lines: str) -> str:
        for directive in Parser._directives():
            lines = lines.replace(f".{directive} ", f"${directive}, ")
            lines = lines.replace(f".{directive}\t", f"${directive},\t")
            lines = lines.replace(f".{directive}\n", f"${directive},\n")
        return lines

    @staticmethod
    def _preprocess_instruction(lines: str) -> str:
        for op_code in Parser._op_codes():
            for suffix in Parser._suffixes():
                lines = lines.replace(f"\t{op_code}{suffix} ", f"\t${op_code} {suffix}, ")
                lines = lines.replace(f"\t{op_code}{suffix}\t", f"\t${op_code} {suffix},\t")
                lines = lines.replace(f"\t{op_code}{suffix}\n", f"\t${op_code} {suffix},\n")
        return lines

    @staticmethod
    def _preprocess_section(lines: str) -> str:
        for section_name in Parser._section_names():
            lines = lines.replace(f".{section_name}", f"%{section_name}")
            lines = lines.replace(f"%{section_name}.", f"%{section_name}, ")
        return lines

    @staticmethod
    def _preprocess_operator(lines: str) -> str:
        for operator in Parser._operators():
            lines = lines.replace(f"{operator}", f" {operator} ")
        return lines

    @staticmethod
    def _directives() -> Set[str]:
        return {
            "addrsig",
            "addrsig_sym",
            "ascii",
            "asciz",
            "bss",
            "byte",
            "cfi_def_cfa_offset",
            "cfi_endproc",
            "cfi_offset",
            "cfi_sections",
            "cfi_startproc",
            "file",
            "globl",
            "loc",
            "long",
            "p2align",
            "quad",
            "section",
            "set",
            "short",
            "size",
            "text",
            "type",
            "weak",
            "zero",
        }

    @staticmethod
    def _op_codes() -> Set[str]:
        return {
            "acquire",
            "release",
            "add",
            "addc",
            "and",
            "andn",
            "asr",
            "cao",
            "clo",
            "cls",
            "clz",
            "cmpb4",
            "div_step",
            "extsb",
            "extsh",
            "extub",
            "extuh",
            "lsl",
            "lsl1",
            "lsl1x",
            "lsl_add",
            "lsl_sub",
            "lslx",
            "lsr",
            "lsr1",
            "lsr1x",
            "lsr_add",
            "lsrx",
            "mul_sh_sh",
            "mul_sh_sl",
            "mul_sh_uh",
            "mul_sh_ul",
            "mul_sl_sh",
            "mul_sl_sl",
            "mul_sl_uh",
            "mul_sl_ul",
            "mul_step",
            "mul_uh_uh",
            "mul_uh_ul",
            "mul_ul_uh",
            "mul_ul_ul",
            "nand",
            "nor",
            "nxor",
            "or",
            "orn",
            "rol",
            "rol_add",
            "ror",
            "rsub",
            "rsubc",
            "sub",
            "subc",
            "xor",
            "boot",
            "resume",
            "stop",
            "call",
            "fault",
            "nop",
            "sats",
            "hash",
            "movd",
            "swapd",
            "time",
            "time_cfg",
            "lbs",
            "lbu",
            "ld",
            "lhs",
            "lhu",
            "lw",
            "sb",
            "sb_id",
            "sd",
            "sd_id",
            "sh",
            "sh_id",
            "sw",
            "sw_id",
            "ldma",
            "ldmai",
            "sdma",
            "adds",
            "move",
            "neg",
            "subs",
            "jump",
            "jeq",
            "jneq",
            "jz",
            "jnz",
            "jltu",
            "jgtu",
            "jleu",
            "jgeu",
            "jlts",
            "jgts",
            "jles",
            "jges",
            "not",
            "bkp",
            "lbss",
            "lbus",
            "lds",
            "lhss",
            "lhus",
            "lws",
            "sbs",
            "sds",
            "shs",
            "sws",
        }

    @staticmethod
    def _suffixes() -> Set[str]:
        return {"", ".s", ".u"}

    @staticmethod
    def _section_names() -> Set[str]:
        return {
            "atomic",
            "bss",
            "data",
            "debug_abbrev",
            "debug_frame",
            "debug_info",
            "debug_line",
            "debug_loc",
            "debug_ranges",
            "debug_str",
            "dpu_host",
            "mram",
            "rodata",
            "stack_sizes",
            "text",
        }

    @staticmethod
    def _operators() -> Set[str]:
        return {"+", "-"}
