package zettel_external

import "github.com/friedenberg/zit/src/charlie/gattung"

type fileEncoder struct {
	arf gattung.AkteReaderFactory
}

func MakeFileEncoder(
	arf gattung.AkteReaderFactory,
) fileEncoder {
	return fileEncoder{
		arf: arf,
	}
}

func (e *fileEncoder) Encode(z *Zettel) (err error) {
	return
}
