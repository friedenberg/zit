package verzeichnisse

import (
	"bufio"
	"encoding/gob"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type zettelenPageIndex struct {
	self map[string]string
}

func (zpi *zettelenPageIndex) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	for {
		if err = dec.Decode(&zpi.self); err != nil {
			if errors.IsEOF(err) {
				err = nil
				break
			} else {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (zpi zettelenPageIndex) WriteTo(w1 io.Writer) (n int64, err error) {
	w := bufio.NewWriter(w1)

	defer errors.PanicIfError(w.Flush)

	enc := gob.NewEncoder(w)

	if err = enc.Encode(zpi.self); err != nil {
		err = errors.Wrapf(err, "failed to write page index")
		return
	}

	return
}
