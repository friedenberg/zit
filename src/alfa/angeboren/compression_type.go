package angeboren

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/bravo/files"
)

type ErrUnsupportedCompression string

func (e ErrUnsupportedCompression) Error() string {
	return fmt.Sprintf("unsupported compression type: %q", string(e))
}

func (e ErrUnsupportedCompression) Is(err error) (ok bool) {
	_, ok = err.(ErrUnsupportedCompression)
	return
}

const (
	CompressionTypeNone    = CompressionType("")
	CompressionTypeGzip    = CompressionType("gzip")
	CompressionTypeDefault = CompressionTypeGzip
	// CompressionTypeZstd = "zstd"
)

type CompressionType string

func (ct *CompressionType) AddToFlagSet(f *flag.FlagSet) {
	f.Var(ct, "compression-type", "")
}

func (ct CompressionType) String() string {
	return string(ct)
}

func (ct *CompressionType) Set(v string) (err error) {
	v1 := CompressionType(strings.TrimSpace(strings.ToLower(v)))

	switch v1 {
	case CompressionTypeGzip, CompressionTypeNone:
		*ct = v1

	default:
		err = ErrUnsupportedCompression(v)
	}

	return
}

func (ct CompressionType) NewReader(
	r io.Reader,
) (out io.ReadCloser, err error) {
	switch ct {
	case CompressionTypeGzip:
		out, err = gzip.NewReader(r)

	default:
		out = io.NopCloser(r)
	}

	return
}

func (ct CompressionType) NewWriter(w io.Writer) io.WriteCloser {
	switch ct {
	case CompressionTypeGzip:
		return gzip.NewWriter(w)

	default:
		return files.NopWriteCloser{Writer: w}
	}
}
