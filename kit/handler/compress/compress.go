package compress

import (
	"compress/flate"
	"compress/gzip"
	"io"
)

type Compress interface {
	New(io.Writer) io.WriteCloser
	ContentEncoding() string
}

var (
	_ Compress = (*gZip)(nil)
	_ Compress = (*flateStd)(nil)
)

type gZip struct {
}

// ContentEncoding implements Compress.
func (g *gZip) ContentEncoding() string {
	return "gzip"
}

// New implements Compress.
func (g *gZip) New(w io.Writer) io.WriteCloser {
	return gzip.NewWriter(w)
}

type flateStd struct {
}

// ContentEncoding implements Compress.
func (f *flateStd) ContentEncoding() string {
	return "deflate"
}

// New implements Compress.
func (f *flateStd) New(w io.Writer) io.WriteCloser {
	fw, _ := flate.NewWriter(w, 6)
	return fw
}
