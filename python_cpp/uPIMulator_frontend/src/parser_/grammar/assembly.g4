grammar assembly;

document: (directive | instruction | label)* EOF;

ACQUIRE: '$acquire';
RELEASE: '$release';
BOOT: '$boot';
RESUME: '$resume';

ADD: '$add';
ADDC: '$addc';
AND: '$and';
ANDN: '$andn';
ASR: '$asr';
CMPB4: '$cmpb4';
LSL: '$lsl';
LSL1: '$lsl1';
LSL1X: '$lsl1x';
LSLX: '$lslx';
LSR: '$lsr';
LSR1: '$lsr1';
LSR1X: '$lsr1x';
LSRX: '$lsrx';
MUL_SH_SH: '$mul_sh_sh';
MUL_SH_SL: '$mul_sh_sl';
MUL_SH_UH: '$mul_sh_uh';
MUL_SH_UL: '$mul_sh_ul';
MUL_SL_SH: '$mul_sl_sh';
MUL_SL_SL: '$mul_sl_sl';
MUL_SL_UH: '$mul_sl_uh';
MUL_SL_UL: '$mul_sl_ul';
MUL_UH_UH: '$mul_uh_uh';
MUL_UH_UL: '$mul_uh_ul';
MUL_UL_UH: '$mul_ul_uh';
MUL_UL_UL: '$mul_ul_ul';
NAND: '$nand';
NOR: '$nor';
NXOR: '$nxor';
OR: '$or';
ORN: '$orn';
ROL: '$rol';
ROR: '$ror';
RSUB: '$rsub';
RSUBC: '$rsubc';
SUB: '$sub';
SUBC: '$subc';
XOR: '$xor';
CALL: '$call';
HASH: '$hash';

CAO: '$cao';
CLO: '$clo';
CLS: '$cls';
CLZ: '$clz';
EXTSB: '$extsb';
EXTSH: '$extsh';
EXTUB: '$extub';
EXTUH: '$extuh';
SATS: '$sats';
TIME_CFG: '$time_cfg';

DIV_STEP: '$div_step';
MUL_STEP: '$mul_step';

LSL_ADD: '$lsl_add';
LSL_SUB: '$lsl_sub';
LSR_ADD: '$lsr_add';
ROL_ADD: '$rol_add';
ROR_ADD: '$ror_add';

TIME: '$time';
NOP: '$nop';

STOP: '$stop';

FAULT: '$fault';

MOVD: '$movd';
SWAPD: '$swapd';

LBS: '$lbs';
LBU: '$lbu';
LD: '$ld';
LHS: '$lhs';
LHU: '$lhu';
LW: '$lw';

SB: '$sb';
SB_ID: '$sb_id';
SD: '$sd';
SD_ID: '$sd_id';
SH: '$sh';
SH_ID: '$sh_id';
SW: '$sw';
SW_ID: '$sw_id';

LDMA: '$ldma';
LDMAI: '$ldmai';
SDMA: '$sdma';

MOVE: '$move';
NEG: '$neg';
NOT: '$not';
BKP: '$bkp';

JEQ: '$jeq';
JNEQ: '$jneq';
JZ: '$jz';
JNZ: '$jnz';
JLTU: '$jltu';
JGTU: '$jgtu';
JLEU: '$jleu';
JGEU: '$jgeu';
JLTS: '$jlts';
JGTS: '$jgts';
JLES: '$jles';
JGES: '$jges';
JUMP: '$jump';

ATOMIC: '%atomic';
BSS: '%bss';
DATA: '%data';
DEBUG_ABBREV: '%debug_abbrev';
DEBUG_FRAME: '%debug_frame';
DEBUG_INFO: '%debug_info';
DEBUG_LINE: '%debug_line';
DEBUG_LOC: '%debug_loc';
DEBUG_RANGES: '%debug_ranges';
DEBUG_STR: '%debug_str';
DPU_HOST: '%dpu_host';
MRAM: '%mram';
RODATA: '%rodata';
STACK_SIZES: '%stack_sizes';
TEXT_SECTION: '%text';

PROGBITS: '@progbits';
NOBITS: '@nobits';

FUNCTION: '@function';
OBJECT: '@object';

TRUE: 'true';
FALSE: 'false';
Z: 'z';
NZ: 'nz';
E: 'e';
O: 'o';
PL: 'pl';
MI: 'mi';
OV: 'ov';
NOV: 'nov';
C: 'c';
NC: 'nc';
SZ: 'sz';
SNZ: 'snz';
SPL: 'spl';
SMI: 'smi';
SO: 'so';
SE: 'se';
NC5: 'nc5';
NC6: 'nc6';
NC7: 'nc7';
NC8: 'nc8';
NC9: 'nc9';
NC10: 'nc10';
NC11: 'nc11';
NC12: 'nc12';
NC13: 'nc13';
NC14: 'nc14';
MAX: 'max';
NMAX: 'nmax';
SH32: 'sh32';
NSH32: 'nsh32';
EQ: 'eq';
NEQ: 'neq';
LTU: 'ltu';
LEU: 'leu';
GTU: 'gtu';
GEU: 'geu';
LTS: 'lts';
LES: 'les';
GTS: 'gts';
GES: 'ges';
XZ: 'xz';
XNZ: 'xnz';
XLEU: 'xleu';
XGTU: 'xgtu';
XLES: 'xles';
XGTS: 'xgts';
SMALL: 'small';
LARGE: 'large';

LITTLE: '!little';
BIG: '!big';

ZERO_REGISTER: 'zero';
ONE: 'one';
ID: 'id';
ID2: 'id2';
ID4: 'id4';
ID8: 'id8';
LNEG: 'lneg';
MNEG: 'mneg';

ADDRSIG: '$addrsig';
ADDRSIG_SYM: '$addrsig_sym';
ASCII: '$ascii';
ASCIZ: '$asciz';
BYTE: '$byte';
CFI_DEF_CFA_OFFSET: '$cfi_def_cfa_offset';
CFI_ENDPROC: '$cfi_endproc';
CFI_OFFSET: '$cfi_offset';
CFI_SECTIONS: '$cfi_sections';
CFI_STARTPROC: '$cfi_startproc';
FILE: '$file';
GLOBL: '$globl';
LOC: '$loc';
LONG: '$long';
P2ALIGN: '$p2align';
QUAD: '$quad';
SECTION: '$section';
SET: '$set';
SHORT: '$short';
SIZE: '$size';
TEXT_DIRECTIVE: '$text';
TYPE: '$type';
WEAK: '$weak';
ZERO_DIRECTIVE: '$zero';

IS_STMT: 'is_stmt';
PROLOGUE_END: 'prologue_end';

S_SUFFIX: '.s';
U_SUFFIX: '.u';

PositiveNumber: [0-9]+;
negative_number: '-' PositiveNumber;
hex_number: '0x' PositiveNumber;
number
	: PositiveNumber
	| negative_number
	| hex_number
	;

rici_op_code
	: ACQUIRE
	| RELEASE
	| BOOT
	| RESUME
	;

rri_op_code
	: ADD
	| ADDC
	| AND
	| ANDN
	| ASR
	| CMPB4
	| LSL
	| LSL1
	| LSL1X
	| LSLX
	| LSR
	| LSR1
	| LSR1X
	| LSRX
	| MUL_SH_SH
	| MUL_SH_SL
	| MUL_SH_UH
	| MUL_SH_UL
	| MUL_SL_SH
	| MUL_SL_SL
	| MUL_SL_UH
	| MUL_SL_UL
	| MUL_UH_UH
	| MUL_UH_UL
	| MUL_UL_UH
	| MUL_UL_UL
	| NAND
	| NOR
	| NXOR
	| OR
	| ORN
	| ROL
	| ROR
	| RSUB
	| RSUBC
	| SUB
	| SUBC
	| XOR
	| CALL
	| HASH
	;

rr_op_code
	: CAO
	| CLO
	| CLS
	| CLZ
	| EXTSB
	| EXTSH
	| EXTUB
	| EXTUH
	| SATS
	| TIME_CFG
	;

drdici_op_code 
	: DIV_STEP
	| MUL_STEP
	;

rrri_op_code
	: LSL_ADD
	| LSL_SUB
	| LSR_ADD
	| ROL_ADD
	| ROR_ADD
	;

r_op_code
	: TIME
	| NOP
	;

ci_op_code
	: STOP
	;

i_op_code
	: FAULT
	;

ddci_op_code
	: MOVD
	| SWAPD
	;

load_op_code
	: LBS
	| LBU
	| LD
	| LHS
	| LHU
	| LW
	;

store_op_code
	: SB
	| SB_ID
	| SD
	| SD_ID
	| SH
	| SH_ID
	| SW
	| SW_ID
	;

dma_op_code
	: LDMA
	| LDMAI
	| SDMA
	;

section_name
	: ATOMIC
	| BSS
	| DATA
	| DEBUG_ABBREV
	| DEBUG_FRAME
	| DEBUG_INFO
	| DEBUG_LINE
	| DEBUG_LOC
	| DEBUG_RANGES
	| DEBUG_STR
	| DPU_HOST
	| MRAM
	| RODATA
	| STACK_SIZES
	| TEXT_SECTION
	;

section_types: PROGBITS | NOBITS;

symbol_type: FUNCTION | OBJECT;

condition
	: TRUE
	| FALSE
	| Z
	| NZ
	| E
	| O
	| PL
	| MI
	| OV
	| NOV
	| C
	| NC
	| SZ
	| SNZ
	| SPL
	| SMI
	| SO
	| SE
	| NC5
	| NC6
	| NC7
	| NC8
	| NC9
	| NC10
	| NC11
	| NC12
	| NC13
	| NC14
	| MAX
	| NMAX
	| SH32
	| NSH32
	| EQ
	| NEQ
	| LTU
	| LEU
	| GTU
	| GEU
	| LTS
	| LES
	| GTS
	| GES
	| XZ
	| XNZ
	| XLEU
	| XGTU
	| XLES
	| XGTS
	| SMALL
	| LARGE
	;

endian: LITTLE | BIG ;

GPRegister: 'r' [0-9]+;

sp_register
	: ZERO_REGISTER
	| ONE
	| ID
	| ID2
	| ID4
	| ID8
	| LNEG
	| MNEG
	;

PairRegister: 'd' [0-9]+;


Identifier: ([a-zA-Z] | '.' | '_') ([a-zA-Z] | [0-9] | '.' | '_')*;

StringLiteral
	: '"'
		( [a-zA-Z]
		| [0-9]
		| '.'
		| ','
		| '_'
		| '+'
		| '-'
		| '/'
		| ':'
		| ';'
		| '('
		| ')'
		| '['
		| ']'
		| '{'
		| '}'
		| '?'
		| '!'
		| '\\'
		| WHITE_SPACE
		)*
	'"'
	;

src_register: GPRegister | sp_register;

program_counter
	: primary_expression
	| add_expression
	| sub_expression
	;
add_expression: primary_expression '+' primary_expression;
sub_expression: primary_expression '-' primary_expression;
primary_expression: number | Identifier | section_name;

directive
	: addrsig_directive
	| addrsig_sym_directive
	| ascii_directive
	| asciz_directive
	| byte_directive
	| cfi_def_cfa_offset_directive
	| cfi_endproc_directive
	| cfi_offset_directive
	| cfi_sections_directive
	| cfi_startproc_directive
	| file_directive
	| global_directive
	| loc_directive
	| long_directive
	| p2align_directive
	| quad_directive
	| section_directive
	| set_directive
	| short_directive
	| size_directive
	| stack_sizes_directive
	| text_directive
	| type_directive
	| weak_directive
	| zero_directive
	;

addrsig_directive: ADDRSIG ',';
addrsig_sym_directive: ADDRSIG_SYM ',' Identifier;
ascii_directive: ASCII ',' StringLiteral;
asciz_directive: ASCIZ ',' StringLiteral;
byte_directive: BYTE ',' program_counter;
cfi_def_cfa_offset_directive: CFI_DEF_CFA_OFFSET ',' number;
cfi_endproc_directive: CFI_ENDPROC ',';
cfi_offset_directive: CFI_OFFSET ',' number ',' number;
cfi_sections_directive: CFI_SECTIONS ',' section_name;
cfi_startproc_directive: CFI_STARTPROC ',';
file_directive
	: FILE ',' StringLiteral
	| FILE ',' number StringLiteral StringLiteral
	;
global_directive: GLOBL ',' Identifier;
loc_directive
	: LOC ',' number number number
	| LOC ',' number number number IS_STMT number
	| LOC ',' number number number PROLOGUE_END
	;
long_directive: LONG ',' program_counter;
p2align_directive: P2ALIGN ',' number;
quad_directive: QUAD ',' program_counter;
section_directive
	: SECTION ',' section_name ',' StringLiteral ',' section_types
	| SECTION ',' section_name ',' StringLiteral ',' section_types ',' number
	| SECTION ',' section_name ',' Identifier ',' StringLiteral ',' section_types
	| SECTION ',' section_name ',' Identifier ',' StringLiteral ',' section_types ',' number
	;
set_directive: SET ',' Identifier ',' Identifier;
short_directive: SHORT ',' program_counter;
size_directive: SIZE ',' Identifier ',' program_counter;
stack_sizes_directive: SECTION ',' STACK_SIZES ',' StringLiteral ',' section_types ',' section_name ',' Identifier;
text_directive: TEXT_DIRECTIVE ',';
type_directive: TYPE ',' Identifier ',' symbol_type;
weak_directive: WEAK ',' Identifier;
zero_directive
	: ZERO_DIRECTIVE ',' number
	| ZERO_DIRECTIVE ',' number ',' number
	;

instruction
	: rici_instruction
	| rri_instruction
	| rric_instruction
	| rrici_instruction
	| rrr_instruction
	| rrrc_instruction
	| rrrci_instruction
	| zri_instruction
	| zric_instruction
	| zrici_instruction
	| zrr_instruction
	| zrrc_instruction
	| zrrci_instruction
	| s_rri_instruction
	| s_rric_instruction
	| s_rrici_instruction
	| s_rrr_instruction
	| s_rrrc_instruction
	| s_rrrci_instruction
	| u_rri_instruction
	| u_rric_instruction
	| u_rrici_instruction
	| u_rrr_instruction
	| u_rrrc_instruction
	| u_rrrci_instruction
	| rr_instruction
	| rrc_instruction
	| rrci_instruction
	| zr_instruction
	| zrc_instruction
	| zrci_instruction
	| s_rr_instruction
	| s_rrc_instruction
	| s_rrci_instruction
	| u_rr_instruction
	| u_rrc_instruction
	| u_rrci_instruction
	| drdici_instruction
	| rrri_instruction
	| rrrici_instruction
	| zrri_instruction
	| zrrici_instruction
	| s_rrri_instruction
	| s_rrrici_instruction
	| u_rrri_instruction
	| u_rrrici_instruction
	| rir_instruction
	| rirc_instruction
	| rirci_instruction
	| zir_instruction
	| zirc_instruction
	| zrici_instruction
	| s_rirc_instruction
	| s_rirci_instruction
	| u_rric_instruction
	| u_rrici_instruction
	| r_instruction
	| rci_instruction
	| z_instruction
	| zci_instruction
	| s_r_instruction
	| s_rci_instruction
	| u_r_instruction
	| u_rci_instruction
	| ci_instruction
	| i_instruction
	| ddci_instruction
	| erri_instruction
	| edri_instruction
	| s_erri_instruction
	| u_erri_instruction
	| erii_instruction
	| erir_instruction
	| erid_instruction
	| dma_rri_instruction
	| synthetic_sugar_instruction
	;

rici_instruction: rici_op_code ',' src_register ',' program_counter ',' condition ',' program_counter;

rri_instruction: rri_op_code ',' GPRegister ',' src_register ',' number;
rric_instruction: rri_op_code ',' GPRegister ',' src_register ',' number ',' condition;
rrici_instruction: rri_op_code ',' GPRegister ',' src_register ',' number ',' condition ',' program_counter;
rrr_instruction: rri_op_code ',' GPRegister ',' src_register ',' src_register;
rrrc_instruction: rri_op_code ',' GPRegister ',' src_register ',' src_register ',' condition;
rrrci_instruction: rri_op_code ',' GPRegister ',' src_register ',' src_register ',' condition ',' program_counter;

zri_instruction: rri_op_code ',' ZERO_REGISTER ',' src_register ',' program_counter;
zric_instruction: rri_op_code ',' ZERO_REGISTER ',' src_register ',' number ',' condition;
zrici_instruction: rri_op_code ',' ZERO_REGISTER ',' src_register ',' number ',' condition ',' program_counter;
zrr_instruction: rri_op_code ',' ZERO_REGISTER ',' src_register ',' src_register;
zrrc_instruction: rri_op_code ',' ZERO_REGISTER ',' src_register ',' src_register ',' condition;
zrrci_instruction: rri_op_code ',' ZERO_REGISTER ',' src_register ',' src_register ',' condition ',' program_counter;

s_rri_instruction: rri_op_code S_SUFFIX ',' PairRegister ',' src_register ',' number;
s_rric_instruction: rri_op_code S_SUFFIX ',' PairRegister ',' src_register ',' number ',' condition;
s_rrici_instruction: rri_op_code S_SUFFIX ',' PairRegister ',' src_register ',' number ',' condition ',' program_counter;
s_rrr_instruction: rri_op_code S_SUFFIX ',' PairRegister ',' src_register ',' src_register;
s_rrrc_instruction: rri_op_code S_SUFFIX ',' PairRegister ',' src_register ',' src_register ',' condition;
s_rrrci_instruction: rri_op_code S_SUFFIX ',' PairRegister ',' src_register ',' src_register ',' condition ',' program_counter;

u_rri_instruction: rri_op_code U_SUFFIX ',' PairRegister ',' src_register ',' number;
u_rric_instruction: rri_op_code U_SUFFIX ',' PairRegister ',' src_register ',' number ',' condition;
u_rrici_instruction: rri_op_code U_SUFFIX ',' PairRegister ',' src_register ',' number ',' condition ',' program_counter;
u_rrr_instruction: rri_op_code U_SUFFIX ',' PairRegister ',' src_register ',' src_register;
u_rrrc_instruction: rri_op_code U_SUFFIX ',' PairRegister ',' src_register ',' src_register ',' condition;
u_rrrci_instruction: rri_op_code U_SUFFIX ',' PairRegister ',' src_register ',' src_register ',' condition ',' program_counter;

rr_instruction: rr_op_code ',' GPRegister ',' src_register;
rrc_instruction: rr_op_code ',' GPRegister ',' src_register ',' condition;
rrci_instruction: rr_op_code ',' GPRegister ',' src_register ',' condition ',' program_counter;

zr_instruction: rr_op_code ',' ZERO_REGISTER ',' src_register;
zrc_instruction: rr_op_code ',' ZERO_REGISTER ',' src_register ',' condition;
zrci_instruction: rr_op_code ',' ZERO_REGISTER ',' src_register ',' condition ',' program_counter;

s_rr_instruction: rr_op_code S_SUFFIX ',' PairRegister ',' src_register;
s_rrc_instruction: rr_op_code S_SUFFIX ',' PairRegister ',' src_register ',' condition;
s_rrci_instruction: rr_op_code S_SUFFIX ',' PairRegister ',' src_register ',' condition ',' program_counter;

u_rr_instruction: rr_op_code U_SUFFIX ',' PairRegister ',' src_register;
u_rrc_instruction: rr_op_code U_SUFFIX ',' PairRegister ',' src_register ',' condition;
u_rrci_instruction: rr_op_code U_SUFFIX ',' PairRegister ',' src_register ',' condition ',' program_counter;

drdici_instruction: drdici_op_code ',' PairRegister ',' src_register ',' PairRegister ',' number ',' condition ',' program_counter;

rrri_instruction: rrri_op_code ',' GPRegister ',' src_register ',' src_register ',' number;
rrrici_instruction: rrri_op_code ',' GPRegister ',' src_register ',' src_register ',' number ',' condition ',' program_counter;

zrri_instruction: rrri_op_code ',' ZERO_REGISTER ',' src_register ',' src_register ',' number;
zrrici_instruction: rrri_op_code ',' ZERO_REGISTER ',' src_register ',' src_register ',' number ',' condition ',' program_counter;

s_rrri_instruction: rrri_op_code S_SUFFIX ',' PairRegister ',' src_register ',' src_register ',' number;
s_rrrici_instruction: rrri_op_code S_SUFFIX ',' PairRegister ',' src_register ',' src_register ',' number ',' condition ',' program_counter;

u_rrri_instruction: rrri_op_code U_SUFFIX ',' PairRegister ',' src_register ',' src_register ',' number;
u_rrrici_instruction: rrri_op_code U_SUFFIX ',' PairRegister ',' src_register ',' src_register ',' number ',' condition ',' program_counter;

rir_instruction: rri_op_code ',' GPRegister ',' number ',' src_register;
rirc_instruction: rri_op_code ',' GPRegister ',' number ',' src_register ',' condition;
rirci_instruction: rri_op_code ',' GPRegister ',' number ',' src_register ',' condition ',' program_counter;

zir_instruction: rri_op_code ',' ZERO_REGISTER ',' number ',' src_register;
zirc_instruction: rri_op_code ',' ZERO_REGISTER ',' number ',' src_register ',' condition;
zirci_instruction: rri_op_code ',' ZERO_REGISTER ',' number ',' src_register ',' condition ',' program_counter;

s_rirc_instruction: rri_op_code S_SUFFIX ',' PairRegister ',' number ',' src_register;
s_rirci_instruction: rri_op_code S_SUFFIX ',' PairRegister ',' number ',' src_register ',' program_counter;

u_rirc_instruction: rri_op_code U_SUFFIX ',' PairRegister ',' number ',' src_register;
u_rirci_instruction: rri_op_code U_SUFFIX ',' PairRegister ',' number ',' src_register ',' program_counter;

r_instruction: r_op_code ',' GPRegister;
rci_instruction: r_op_code ',' condition ',' condition ',' program_counter;

z_instruction
	: r_op_code ','
	| r_op_code ',' ZERO_REGISTER
	;
zci_instruction: r_op_code ',' ZERO_REGISTER ',' condition ',' program_counter;

s_r_instruction: r_op_code S_SUFFIX ',' PairRegister;
s_rci_instruction: r_op_code S_SUFFIX ',' PairRegister ',' condition ',' program_counter;

u_r_instruction: r_op_code U_SUFFIX ',' PairRegister;
u_rci_instruction: r_op_code U_SUFFIX ',' PairRegister ',' condition ',' program_counter;

ci_instruction: ci_op_code ',' condition ',' program_counter;
i_instruction: i_op_code ',' number;

ddci_instruction: ddci_op_code ',' PairRegister ',' PairRegister ',' condition ',' program_counter;

erri_instruction: load_op_code ',' endian ',' GPRegister ',' src_register ',' program_counter;
edri_instruction: load_op_code ',' endian ',' PairRegister ',' src_register ',' program_counter;
s_erri_instruction: load_op_code S_SUFFIX ',' endian ',' PairRegister ',' src_register ',' program_counter;
u_erri_instruction: load_op_code U_SUFFIX ',' endian ',' PairRegister ',' src_register ',' program_counter;

erii_instruction: store_op_code ',' endian ',' src_register ',' number ',' number;
erir_instruction: store_op_code ',' endian ',' src_register ',' program_counter ',' src_register;
erid_instruction: store_op_code ',' endian ',' src_register ',' program_counter ',' PairRegister;

dma_rri_instruction: dma_op_code ',' src_register ',' src_register ',' number;

synthetic_sugar_instruction
	: rrif_instruction
	| move_instruction
	| neg_instruction
	| not_instruction
	| jump_instruction
	| shortcut_instruction
	;

rrif_instruction
	: andn_rrif_instruction
	| nand_rrif_instruction
	| nor_rrif_instruction
	| nxor_rrif_instruction
	| orn_rrif_instruction
	| hash_rrif_instruction
	;
andn_rrif_instruction: ANDN ',' GPRegister ',' src_register ',' number;
nand_rrif_instruction: NAND ',' GPRegister ',' src_register ',' number;
nor_rrif_instruction: NOR ',' GPRegister ',' src_register ',' number;
nxor_rrif_instruction: NXOR ',' GPRegister ',' src_register ',' number;
orn_rrif_instruction: ORN ',' GPRegister ',' src_register ',' number;
hash_rrif_instruction: HASH ',' GPRegister ',' src_register ',' number;

move_instruction
	: move_ri_instruction
	| move_rici_instruction
	| move_rr_instruction
	| move_rrci_instruction
	| move_s_ri_instruction
	| move_s_rici_instruction
	| move_s_rr_instruction
	| move_s_rrci_instruction
	| move_u_ri_instruction
	| move_u_rici_instruction
	| move_u_rr_instruction
	| move_u_rrci_instruction
	;
move_ri_instruction: MOVE ',' GPRegister ',' program_counter;
move_rici_instruction: MOVE ',' GPRegister ',' number ',' condition ',' program_counter;
move_rr_instruction: MOVE ',' GPRegister ',' src_register;
move_rrci_instruction: MOVE ',' GPRegister ',' src_register ',' condition ',' program_counter;
move_s_ri_instruction: MOVE S_SUFFIX ',' PairRegister ',' number;
move_s_rici_instruction: MOVE S_SUFFIX ',' PairRegister ',' number ',' condition ',' program_counter;
move_s_rr_instruction: MOVE S_SUFFIX ',' PairRegister ',' src_register;
move_s_rrci_instruction: MOVE S_SUFFIX ',' PairRegister ',' src_register ',' condition ',' program_counter;
move_u_ri_instruction: MOVE U_SUFFIX ',' PairRegister ',' number;
move_u_rici_instruction: MOVE U_SUFFIX ',' PairRegister ',' number ',' condition ',' program_counter;
move_u_rr_instruction: MOVE U_SUFFIX ',' PairRegister ',' src_register;
move_u_rrci_instruction: MOVE U_SUFFIX ',' PairRegister ',' src_register ',' condition ',' program_counter;

neg_instruction
	: neg_rr_instruction
	| neg_rrci_instruction
	;
neg_rr_instruction: NEG ',' GPRegister ',' src_register;
neg_rrci_instruction: NEG ',' GPRegister ',' src_register ',' condition ',' program_counter;

not_instruction
	: not_rr_instruction
	| not_rrci_instruction
	| not_zrci_instruction
	;
not_rr_instruction: NOT ',' GPRegister ',' src_register;
not_rrci_instruction: NOT ',' GPRegister ',' src_register ',' condition ',' program_counter;
not_zrci_instruction: NOT ',' src_register ',' condition ',' program_counter;

jump_instruction
	: jeq_rii_instruction
	| jeq_rri_instruction
	| jneq_rii_instruction
	| jneq_rri_instruction
	| jz_ri_instruction
	| jnz_ri_instruction
	| jltu_rii_instruction
	| jltu_rri_instruction
	| jgtu_rii_instruction
	| jgtu_rri_instruction
	| jleu_rii_instruction
	| jleu_rri_instruction
	| jgeu_rii_instruction
	| jgeu_rri_instruction
	| jlts_rii_instruction
	| jlts_rri_instruction
	| jgts_rii_instruction
	| jgts_rri_instruction
	| jles_rii_instruction
	| jles_rri_instruction
	| jges_rii_instruction
	| jges_rri_instruction
	| jump_ri_instruction
	| jump_i_instruction
	| jump_r_instruction
	;
jeq_rii_instruction: JEQ ',' src_register ',' program_counter ',' program_counter;
jeq_rri_instruction: JEQ ',' src_register ',' src_register ',' program_counter;
jneq_rii_instruction: JNEQ ',' src_register ',' number ',' program_counter;
jneq_rri_instruction: JNEQ ',' src_register ',' src_register ',' program_counter;
jz_ri_instruction: JZ ',' src_register ',' program_counter;
jnz_ri_instruction: JNZ ',' src_register ',' program_counter;
jltu_rii_instruction: JLTU ',' src_register ',' number ',' program_counter;
jltu_rri_instruction: JLTU ',' src_register ',' src_register ',' program_counter;
jgtu_rii_instruction: JGTU ',' src_register ',' number ',' program_counter;
jgtu_rri_instruction: JGTU ',' src_register ',' src_register ',' program_counter;
jleu_rii_instruction: JLEU ',' src_register ',' number ',' program_counter;
jleu_rri_instruction: JLEU ',' src_register ',' src_register ',' program_counter;
jgeu_rii_instruction: JGEU ',' src_register ',' number ',' program_counter;
jgeu_rri_instruction: JGEU ',' src_register ',' src_register ',' program_counter;
jlts_rii_instruction: JLTS ',' src_register ',' number ',' program_counter;
jlts_rri_instruction: JLTS ',' src_register ',' src_register ',' program_counter;
jgts_rii_instruction: JGTS ',' src_register ',' number ',' program_counter;
jgts_rri_instruction: JGTS ',' src_register ',' src_register ',' program_counter;
jles_rii_instruction: JLES ',' src_register ',' number ',' program_counter;
jles_rri_instruction: JLES ',' src_register ',' src_register ',' program_counter;
jges_rii_instruction: JGES ',' src_register ',' number ',' program_counter;
jges_rri_instruction: JGES ',' src_register ',' src_register ',' program_counter;
jump_ri_instruction: JUMP ',' src_register ',' program_counter;
jump_i_instruction: JUMP ',' program_counter;
jump_r_instruction: JUMP ',' src_register;

shortcut_instruction
	: div_step_drdici_instruction
	| mul_step_drdici_instruction
	| boot_rici_instruction
	| resume_rici_instruction
	| stop_ci_instruction
	| call_ri_instruction
	| call_rr_instruction
	| bkp_instruction
	| movd_ddci_instruction
	| swapd_ddci_instruction
	| time_cfg_zr_instruction
	| lbs_erri_instruction
	| lbs_s_erri_instruction
	| lbu_erri_instruction
	| lbu_u_erri_instruction
	| ld_edri_instruction
	| lhs_erri_instruction
	| lhs_s_erri_instruction
	| lhu_erri_instruction
	| lhu_u_erri_instruction
	| lw_erri_instruction
	| lw_s_erri_instruction
	| lw_u_erri_instruction
	| sb_erii_instruction
	| sb_erir_instruction
	| sb_id_rii_instruction
	| sb_id_ri_instruction
	| sd_erii_instruction
	| sd_erid_instruction
	| sd_id_rii_instruction
	| sd_id_ri_instruction
	| sh_erii_instruction
	| sh_erir_instruction
	| sh_id_rii_instruction
	| sh_id_ri_instruction
	| sw_erii_instruction
	| sw_erir_instruction
	| sw_id_rii_instruction
	| sw_id_ri_instruction
	;
div_step_drdici_instruction: DIV_STEP ',' PairRegister ',' src_register ',' PairRegister ',' number;
mul_step_drdici_instruction: MUL_STEP ',' PairRegister ',' src_register ',' PairRegister ',' number;
boot_rici_instruction: BOOT ',' src_register ',' number;
resume_rici_instruction: RESUME ',' src_register ',' number;
stop_ci_instruction: STOP ',';
call_ri_instruction: CALL ',' GPRegister ',' program_counter;
call_rr_instruction: CALL ',' GPRegister ',' src_register;
bkp_instruction: BKP ',';
movd_ddci_instruction: MOVD ',' PairRegister ',' PairRegister;
swapd_ddci_instruction: SWAPD ',' PairRegister ',' PairRegister;
time_cfg_zr_instruction: TIME_CFG ',' src_register;
lbs_erri_instruction: LBS ',' GPRegister ',' src_register ',' program_counter;
lbs_s_erri_instruction: LBS S_SUFFIX ',' PairRegister ',' src_register ',' program_counter;
lbu_erri_instruction: LBU ',' GPRegister ',' src_register ',' program_counter;
lbu_u_erri_instruction: LBU U_SUFFIX ',' PairRegister ',' src_register ',' program_counter;
ld_edri_instruction: LD ',' PairRegister ',' src_register ',' program_counter;
lhs_erri_instruction: LHS ',' GPRegister ',' src_register ',' program_counter;
lhs_s_erri_instruction: LHS S_SUFFIX ',' PairRegister ',' src_register ',' program_counter;
lhu_erri_instruction: LHU ',' GPRegister ',' src_register ',' program_counter;
lhu_u_erri_instruction: LHU U_SUFFIX ',' PairRegister ',' src_register ',' program_counter;
lw_erri_instruction: LW ',' GPRegister ',' src_register ',' program_counter;
lw_s_erri_instruction: LW S_SUFFIX ',' PairRegister ',' src_register ',' program_counter;
lw_u_erri_instruction: LW U_SUFFIX ',' PairRegister ',' src_register ',' program_counter;
sb_erii_instruction: SB ',' src_register ',' number ',' program_counter;
sb_erir_instruction: SB ',' src_register ',' program_counter ',' src_register;
sb_id_rii_instruction: SB ',' src_register ',' number ',' number;
sb_id_ri_instruction: SB ',' src_register ',' number;
sd_erii_instruction: SD ',' src_register ',' program_counter ',' program_counter;
sd_erid_instruction: SD ',' src_register ',' program_counter ',' PairRegister;
sd_id_rii_instruction: SD ',' src_register ',' number ',' number;
sd_id_ri_instruction: SD ',' src_register ',' number;
sh_erii_instruction: SH ',' src_register ',' number ',' program_counter;
sh_erir_instruction: SH ',' src_register ',' program_counter ',' src_register;
sh_id_rii_instruction: SH ',' src_register ',' number ',' number;
sh_id_ri_instruction: SH ',' src_register ',' number;
sw_erii_instruction: SW ',' src_register ',' number ',' program_counter;
sw_erir_instruction: SW ',' src_register ',' program_counter ',' src_register;
sw_id_rii_instruction: SW ',' src_register ',' number ',' number;
sw_id_ri_instruction: SW ',' src_register ',' number;

label: Identifier ':';

COMMENT: '//' ~[\n\r]* -> skip;
WHITE_SPACE: [ \n\t\r]+ -> skip;
