package dram

import (
	"errors"
	"uPIMulator/src/encoding"
	"uPIMulator/src/misc"
)

type Mram struct {
	address int64
	size    int64

	wordline_size int64
	wordlines     []*Wordline
}

func (this *Mram) Init(command_line_parser *misc.CommandLineParser) {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.address = config_loader.MramOffset()
	this.size = config_loader.MramSize()
	this.wordline_size = command_line_parser.IntParameter("wordline_size")

	if this.wordline_size <= 0 {
		err := errors.New("wordline size <= 0")
		panic(err)
	} else if this.address%this.wordline_size != 0 {
		err := errors.New("address is not aligned with wordline size")
		panic(err)
	} else if this.size%this.wordline_size != 0 {
		err := errors.New("size is not aligned with wordline size")
		panic(err)
	}

	this.wordlines = make([]*Wordline, 0)
	num_wordlines := int(this.size / this.wordline_size)
	for i := 0; i < num_wordlines; i++ {
		wordline := new(Wordline)
		wordline.Init(this.address+int64(i)*this.wordline_size, this.wordline_size)
		this.wordlines = append(this.wordlines, wordline)
	}
}

func (this *Mram) Fini() {
	for _, wordline := range this.wordlines {
		wordline.Fini()
	}
}

func (this *Mram) Address() int64 {
	return this.address
}

func (this *Mram) Size() int64 {
	return this.size
}

func (this *Mram) Read(address int64) *encoding.ByteStream {
	return this.wordlines[this.Index(address)].Read()
}

func (this *Mram) Write(address int64, size int64, byte_stream *encoding.ByteStream) {
	if size != byte_stream.Size() {
		err := errors.New("size != byte stream's size")
		panic(err)
	}

	this.wordlines[this.Index(address)].Write(byte_stream)
}

func (this *Mram) Index(address int64) int {
	if address < this.address {
		err := errors.New("address < MRAM offset")
		panic(err)
	} else if address+this.wordline_size > this.address+this.size {
		err := errors.New("address + wordline size > MRAM offset + MRAM size")
		panic(err)
	} else if address%this.wordline_size != 0 {
		err := errors.New("address is not aligned with wordline size")
		panic(err)
	}

	return int((address - this.address) / this.wordline_size)
}
