package store_verzeichnisse

import (
	"bufio"
	"encoding/gob"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/zettel_verzeichnisse"
)

type zettelenPage struct {
	pool        *zettel_verzeichnisse.Pool
	existing    []*zettel_verzeichnisse.Zettel
	added       []*zettel_verzeichnisse.Zettel
	flushFilter zettel_verzeichnisse.Writer
}

func (zp zettelenPage) Copy(
	r1 io.Reader,
	w zettel_verzeichnisse.Writer,
) (n int64, err error) {
	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	i := 0
	for {
		i += 1

		var tz *zettel_verzeichnisse.Zettel

		if zp.pool != nil {
			tz = zp.pool.Get()
		} else {
			tz = &zettel_verzeichnisse.Zettel{}
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
		zettel_verzeichnisse.MakeWriterMulti(
			zp.pool,
			zettel_verzeichnisse.MakeWriter(
				func(z *zettel_verzeichnisse.Zettel) (err error) {
					zp.existing = append(zp.existing, z)
					return
				},
			),
		),
	)
}

func (zp *zettelenPage) WriteTo(w1 io.Writer) (n int64, err error) {
	w := bufio.NewWriter(w1)

	defer errors.PanicIfError(w.Flush)

	enc := gob.NewEncoder(w)

	for _, z := range zp.existing {
		if err = zp.writeOne(enc, z); err != nil {
			err = errors.Wrapf(err, "failed to write zettel: %v", z)
			return
		}
	}

	for _, z := range zp.added {
		if err = zp.writeOne(enc, z); err != nil {
			err = errors.Wrapf(err, "failed to write zettel: %v", z)
			return
		}
	}

	return
}

func (zp *zettelenPage) writeOne(
	enc *gob.Encoder,
	z *zettel_verzeichnisse.Zettel,
) (err error) {
	if zp.flushFilter != nil {
		if err = zp.flushFilter.WriteZettelVerzeichnisse(z); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	if err = enc.Encode(z); err != nil {
		err = errors.Wrapf(err, "failed to write zettel: %v", z)
		return
	}

	return
}
