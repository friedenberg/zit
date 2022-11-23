package test_metadatei_io

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/charlie/sha"
)

type nopReadWriteCloser struct {
	io.ReadWriter
}

func NopReadWriteCloser(rw io.ReadWriter) *nopReadWriteCloser {
	return &nopReadWriteCloser{
		ReadWriter: rw,
	}
}

func (nrwc *nopReadWriteCloser) Close() (err error) {
	return
}

type nopAkteIoFactory struct {
	io.ReadWriteCloser
}

func NopFactoryReadWriter(rw io.ReadWriter) *nopAkteIoFactory {
	return NopFactoryReadWriteCloser(
		NopReadWriteCloser(rw),
	)
}

func NopFactoryReadWriteCloser(rwc io.ReadWriteCloser) *nopAkteIoFactory {
	return &nopAkteIoFactory{
		ReadWriteCloser: rwc,
	}
}

func (b nopAkteIoFactory) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = io.Copy(b.ReadWriteCloser, r)
	return
}

func (b nopAkteIoFactory) WriteTo(w io.Writer) (n int64, err error) {
	n, err = io.Copy(w, b.ReadWriteCloser)
	return
}

func (b nopAkteIoFactory) Sha() sha.Sha {
	return sha.Sha{}
}

func (aw nopAkteIoFactory) AkteReader(sh sha.Sha) (sha.ReadCloser, error) {
	return aw, nil
}

func (aw nopAkteIoFactory) AkteWriter() (sha.WriteCloser, error) {
	return aw, nil
}

func (aw nopAkteIoFactory) String() string {
	sb := &strings.Builder{}

	io.Copy(sb, aw)

	return sb.String()
}
