package misc

import (
	"bufio"
	"os"
)

type FileScanner struct {
	path string
}

func (this *FileScanner) Init(path string) {
	this.path = path
}

func (this *FileScanner) ReadLines() []string {
	file, open_err := os.Open(this.path)

	if open_err != nil {
		panic(open_err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if scan_err := scanner.Err(); scan_err != nil {
		panic(scan_err)
	}

	return lines
}
