package test_metadatei_io

import (
	"bytes"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/echo/age_io"
)

type AkteIOFactory struct {
	contents      map[string]string
	currentBuffer *bytes.Buffer
}

func FixtureFactoryReadWriteCloser(
	contents map[string]string,
) *AkteIOFactory {
	return &AkteIOFactory{
		contents: contents,
	}
}

func (aw AkteIOFactory) AkteReader(
	sh sha.ShaLike,
) (rc sha.ReadCloser, err error) {
	if s, ok := aw.contents[sh.GetSha().String()]; ok {
		rc = sha.MakeNopReadCloser(io.NopCloser(strings.NewReader(s)))
	} else {
		err = errors.Errorf("not found: %s", sh)
		return
	}

	return
}

func (aw *AkteIOFactory) AkteWriter() (sha.WriteCloser, error) {
	aw.currentBuffer = bytes.NewBuffer(nil)
	wo := age_io.WriteOptions{
		Writer: aw.currentBuffer,
	}

	return age_io.NewWriter(wo)
}

func (aw AkteIOFactory) CurrentBufferString() string {
	if aw.currentBuffer == nil {
		return ""
	}

	sb := &strings.Builder{}

	io.Copy(sb, aw.currentBuffer)

	return sb.String()
}
