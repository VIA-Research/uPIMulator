package assembler

import (
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/misc"
)

type Assemblable interface {
	Init(command_line_parser *misc.CommandLineParser)

	InputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream
	OutputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream

	InputDpuMramHeapPointerName(execution int, dpu_id int) (int64, *encoding.ByteStream)
	OutputDpuMramHeapPointerName(execution int, dpu_id int) (int64, *encoding.ByteStream)

	NumExecutions() int
}
