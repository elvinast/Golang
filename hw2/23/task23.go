package main

import (
	// "fmt"
	"io"
	"os"
	"strings"
)

//rotate by 13
type rot13Reader struct {
	r io.Reader
}

func (rot *rot13Reader) Read(bytes []byte) (n int, e error) {
	n, e = rot.r.Read(bytes)
	// fmt.Println(len(bytes))
	for i := 0; i < len(bytes); i++ {
		if ((bytes[i] >= 'a' && bytes[i] < 'n') || (bytes[i] >= 'A' && bytes[i] < 'N')) {
			bytes[i] += 13
		} else if ((bytes[i] > 'm' && bytes[i] <= 'z') || (bytes[i] > 'M' && bytes[i] <= 'Z')) {
			bytes[i] -= 13
		}
	}
	return n, e
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
