package immutable_config

import (
	"compress/gzip"
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"github.com/DataDog/zstd"
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
	CompressionTypeEmpty = CompressionType("")
	CompressionTypeNone  = CompressionType("none")
	CompressionTypeGzip  = CompressionType("gzip")
	CompressionTypeZlib  = CompressionType("zlib")
	CompressionTypeZstd  = CompressionType("zstd")

	CompressionTypeDefault = CompressionTypeGzip
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
	case CompressionTypeGzip,
		CompressionTypeNone,
		CompressionTypeEmpty,
		CompressionTypeZstd,
		CompressionTypeZlib:
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

	case CompressionTypeZlib:
		out, err = zlib.NewReader(r)

	case CompressionTypeZstd:
		out = zstd.NewReader(r)

	default:
		out = io.NopCloser(r)
	}

	return
}

func (ct CompressionType) NewWriter(w io.Writer) io.WriteCloser {
	switch ct {
	case CompressionTypeGzip:
		return gzip.NewWriter(w)

	case CompressionTypeZlib:
		return zlib.NewWriter(w)

	case CompressionTypeZstd:
		return zstd.NewWriter(w)

	default:
		return files.NopWriteCloser{Writer: w}
	}
}
