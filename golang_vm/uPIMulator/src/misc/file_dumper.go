package misc

import (
	"bufio"
	"os"
)

type FileDumper struct {
	path string
}

func (this *FileDumper) Init(path string) {
	this.path = path
}

func (this *FileDumper) WriteLines(lines []string) {
	file, create_err := os.Create(this.path)

	if create_err != nil {
		panic(create_err)
	}

	writer := bufio.NewWriter(file)

	for _, line := range lines {
		_, write_err := writer.WriteString(line + "\n")

		if write_err != nil {
			panic(write_err)
		}
	}

	flush_err := writer.Flush()
	if flush_err != nil {
		panic(flush_err)
	}

	close_err := file.Close()
	if close_err != nil {
		panic(close_err)
	}
}
