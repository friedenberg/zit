package verzeichnisse

import (
	"bufio"
	"encoding/gob"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type zettelenPage struct {
	existing []*Zettel
	added    []*Zettel
}

func (zp zettelenPage) Copy(
	r1 io.Reader,
	w writer,
) (n int64, err error) {
	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	i := 0
	for {
		i += 1

		var tz *Zettel

		if w.ZettelPool == nil {
			tz = &Zettel{}
		} else {
			tz = w.Get()
		}

		if err = dec.Decode(tz); err != nil {
			if errors.IsEOF(err) {
				err = nil
				break
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		if err = w.WriteZettelVerzeichnisse(tz); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (zp *zettelenPage) ReadFrom(r1 io.Reader) (n int64, err error) {
	return zp.Copy(
		r1,
		writer{
			writers: []Writer{
				MakeWriter(
					func(z *Zettel) (err error) {
						zp.existing = append(zp.existing, z)
						return
					},
				),
			},
		},
	)
}

func (zp *zettelenPage) WriteTo(w1 io.Writer) (n int64, err error) {
	w := bufio.NewWriter(w1)

	defer errors.PanicIfError(w.Flush)

	enc := gob.NewEncoder(w)

	for _, z := range zp.existing {
		if err = enc.Encode(z); err != nil {
			err = errors.Wrapf(err, "failed to write zettel: %v", z)
			return
		}
	}

	for _, z := range zp.added {
		if err = enc.Encode(z); err != nil {
			err = errors.Wrapf(err, "failed to write zettel: %v", z)
			return
		}
	}

	return
}
