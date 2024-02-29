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

	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range lines {
		_, write_err := writer.WriteString(line + "\n")

		if write_err != nil {
			panic(write_err)
		}
	}

	writer.Flush()
}
