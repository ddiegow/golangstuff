package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (r13 rot13Reader) Read(b []byte) (int, error) {
	n, e := r13.r.Read(b)

	for i := range b {
		if b[i] >= 'A' && b[i] < 'Z' {
			b[i] = 65 + (((b[i] - 65) + 13) % 26)
		} else if b[i] >= 'a' && b[i] <= 'z' {
			b[i] = 97 + (((b[i] - 97) + 13) % 26)
		}
	}
	return n, e
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
