package output

import "io"

// NopCloser will return a new NopCloser for writers without a Close method
func NopCloser(w io.Writer) io.WriteCloser {
	var c closer
	c.Writer = w
	return &c
}

type closer struct {
	io.Writer
}

func (c *closer) Close() (err error) {
	return
}
