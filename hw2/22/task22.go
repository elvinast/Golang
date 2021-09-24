package main

import "golang.org/x/tour/reader"

type MyReader struct{}

func (r MyReader) Read(bytes []byte) (int, error) {
	for idx := range bytes {
		bytes[idx] = 'A'
	}
	return len(bytes), nil
}

func main() {
	reader.Validate(MyReader{})
}
