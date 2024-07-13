package test_metadatei_io

import (
	"bytes"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
)

type BlobIOFactory struct {
	contents      map[string]string
	currentBuffer *bytes.Buffer
}

func FixtureFactoryReadWriteCloser(
	contents map[string]string,
) *BlobIOFactory {
	return &BlobIOFactory{
		contents: contents,
	}
}

func (aw BlobIOFactory) BlobReader(
	sh sha.ShaLike,
) (rc sha.ReadCloser, err error) {
	if s, ok := aw.contents[sh.GetShaLike().String()]; ok {
		rc = sha.MakeNopReadCloser(io.NopCloser(strings.NewReader(s)))
	} else {
		err = errors.Errorf("not found: %s", sh)
		return
	}

	return
}

func (aw *BlobIOFactory) BlobWriter() (sha.WriteCloser, error) {
	aw.currentBuffer = bytes.NewBuffer(nil)
	wo := standort.WriteOptions{
		Writer: aw.currentBuffer,
	}

	return standort.NewWriter(wo)
}

func (aw BlobIOFactory) CurrentBufferString() string {
	if aw.currentBuffer == nil {
		return ""
	}

	sb := &strings.Builder{}

	io.Copy(sb, aw.currentBuffer)

	return sb.String()
}
