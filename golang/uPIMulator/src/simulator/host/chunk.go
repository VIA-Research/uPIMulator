package host

import (
	"errors"
	"strconv"
	"strings"
	"uPIMulator/src/abi/encoding"
)

type ChunkType int

const (
	INPUT_DPU_HOST ChunkType = iota
	OUTPUT_DPU_HOST
	INPUT_DPU_MRAM_HEAP_POINTER_NAME
	OUTPUT_DPU_MRAM_HEAP_POINTER_NAME
)

type Chunk struct {
	chunk_type ChunkType

	name   *string
	offset *int64

	execution int
	dpu_id    int

	byte_stream *encoding.ByteStream
}

func (this *Chunk) Init(filename string, byte_stream *encoding.ByteStream) {
	if this.IsInputDpuHost(filename) {
		this.InitInputDpuHost(filename)
	} else if this.IsOutputDpuHost(filename) {
		this.InitOutputDpuHost(filename)
	} else if this.IsInputDpuMramHeapPointerName(filename) {
		this.InitInputDpuMramHeapPointerName(filename)
	} else if this.IsOutputDpuMramHeapPointerName(filename) {
		this.InitOutputDpuMramHeapPointerName(filename)
	} else {
		err := errors.New("filename cannot be parsed")
		panic(err)
	}

	this.byte_stream = byte_stream
}

func (this *Chunk) IsInputDpuHost(filename string) bool {
	words := strings.Split(strings.Split(filename, ".")[0], "_")

	if words[0] == "input" &&
		!(words[1] == "dpu" &&
			words[2] == "mram" &&
			words[3] == "heap" &&
			words[4] == "pointer" &&
			words[5] == "name") {
		return true
	} else {
		return false
	}
}

func (this *Chunk) IsOutputDpuHost(filename string) bool {
	words := strings.Split(strings.Split(filename, ".")[0], "_")

	if words[0] == "output" &&
		!(words[1] == "dpu" &&
			words[2] == "mram" &&
			words[3] == "heap" &&
			words[4] == "pointer" &&
			words[5] == "name") {
		return true
	} else {
		return false
	}
}

func (this *Chunk) IsInputDpuMramHeapPointerName(filename string) bool {
	words := strings.Split(strings.Split(filename, ".")[0], "_")

	if words[0] == "input" &&
		words[1] == "dpu" &&
		words[2] == "mram" &&
		words[3] == "heap" &&
		words[4] == "pointer" &&
		words[5] == "name" {
		return true
	} else {
		return false
	}
}

func (this *Chunk) IsOutputDpuMramHeapPointerName(filename string) bool {
	words := strings.Split(strings.Split(filename, ".")[0], "_")

	if words[0] == "output" &&
		words[1] == "dpu" &&
		words[2] == "mram" &&
		words[3] == "heap" &&
		words[4] == "pointer" &&
		words[5] == "name" {
		return true
	} else {
		return false
	}
}

func (this *Chunk) InitInputDpuHost(filename string) {
	if !this.IsInputDpuHost(filename) {
		err := errors.New("filename is not input DPU host")
		panic(err)
	}

	this.chunk_type = INPUT_DPU_HOST

	words := strings.Split(strings.Split(filename, ".")[0], "_")

	this.name = new(string)
	*this.name = strings.Join(words[1:len(words)-2], "_")

	var err error
	this.execution, err = strconv.Atoi(words[len(words)-2])
	if err != nil {
		panic(err)
	}

	this.dpu_id, err = strconv.Atoi(words[len(words)-1])
	if err != nil {
		panic(err)
	}
}

func (this *Chunk) InitOutputDpuHost(filename string) {
	if !this.IsOutputDpuHost(filename) {
		err := errors.New("filename is not output DPU host")
		panic(err)
	}

	this.chunk_type = OUTPUT_DPU_HOST

	words := strings.Split(strings.Split(filename, ".")[0], "_")

	this.name = new(string)
	*this.name = strings.Join(words[1:len(words)-2], "_")

	var err error
	this.execution, err = strconv.Atoi(words[len(words)-2])
	if err != nil {
		panic(err)
	}

	this.dpu_id, err = strconv.Atoi(words[len(words)-1])
	if err != nil {
		panic(err)
	}
}

func (this *Chunk) InitInputDpuMramHeapPointerName(filename string) {
	if !this.IsInputDpuMramHeapPointerName(filename) {
		err := errors.New("filename is not input DPU MRAM heap pointer name")
		panic(err)
	}

	this.chunk_type = INPUT_DPU_MRAM_HEAP_POINTER_NAME

	words := strings.Split(strings.Split(filename, ".")[0], "_")

	var err error
	this.offset = new(int64)
	*this.offset, err = strconv.ParseInt(words[len(words)-3], 10, 64)
	if err != nil {
		panic(err)
	}

	this.execution, err = strconv.Atoi(words[len(words)-2])
	if err != nil {
		panic(err)
	}

	this.dpu_id, err = strconv.Atoi(words[len(words)-1])
	if err != nil {
		panic(err)
	}
}

func (this *Chunk) InitOutputDpuMramHeapPointerName(filename string) {
	if !this.IsOutputDpuMramHeapPointerName(filename) {
		err := errors.New("filename is not output DPU MRAM heap pointer name")
		panic(err)
	}

	this.chunk_type = OUTPUT_DPU_MRAM_HEAP_POINTER_NAME

	words := strings.Split(strings.Split(filename, ".")[0], "_")

	var err error
	this.offset = new(int64)
	*this.offset, err = strconv.ParseInt(words[len(words)-3], 10, 64)
	if err != nil {
		panic(err)
	}

	this.execution, err = strconv.Atoi(words[len(words)-2])
	if err != nil {
		panic(err)
	}

	this.dpu_id, err = strconv.Atoi(words[len(words)-1])
	if err != nil {
		panic(err)
	}
}

func (this *Chunk) ChunkType() ChunkType {
	return this.chunk_type
}

func (this *Chunk) Name() string {
	if this.name == nil {
		err := errors.New("name == nil")
		panic(err)
	}

	return *this.name
}

func (this *Chunk) Offset() int64 {
	if this.offset == nil {
		err := errors.New("offset == nil")
		panic(err)
	}

	return *this.offset
}

func (this *Chunk) Execution() int {
	return this.execution
}

func (this *Chunk) DpuId() int {
	return this.dpu_id
}

func (this *Chunk) ByteStream() *encoding.ByteStream {
	return this.byte_stream
}
