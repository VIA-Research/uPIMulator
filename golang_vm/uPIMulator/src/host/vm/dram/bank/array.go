package bank

import (
	"errors"
	"uPIMulator/src/encoding"
	"uPIMulator/src/misc"
)

type Array struct {
	channel_id int
	rank_id    int
	bank_id    int

	address int64
	size    int64

	wordline_size int64
	wordlines     []*Wordline
}

func (this *Array) Init(command_line_parser *misc.CommandLineParser) {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.address = 0
	this.size = config_loader.VmBankSize()
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

func (this *Array) Fini() {
	for _, wordline := range this.wordlines {
		wordline.Fini()
	}
}

func (this *Array) Address() int64 {
	return this.address
}

func (this *Array) Size() int64 {
	return this.size
}

func (this *Array) Read(address int64) *encoding.ByteStream {
	return this.wordlines[this.Index(address)].Read()
}

func (this *Array) Write(address int64, byte_stream *encoding.ByteStream) {
	this.wordlines[this.Index(address)].Write(byte_stream)
}

func (this *Array) Index(address int64) int {
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
