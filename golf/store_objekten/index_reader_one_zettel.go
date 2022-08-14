package objekten

import (
	"bytes"
	"encoding/gob"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/delta/id"
	"github.com/friedenberg/zit/foxtrot/verzeichnisse"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type indexReaderOneZettel struct {
	id.Id
	objekte []byte
	stored_zettel.Named
}

func (r *indexReaderOneZettel) Begin() (err error) {
	return
}

func (r *indexReaderOneZettel) ReadRow(id string, row verzeichnisse.Row) (err error) {
	if row.Key != r.Id.String() {
		return
	}

	r.objekte = row.Objekte

	return
}

func (r *indexReaderOneZettel) End() (err error) {
	if len(r.objekte) == 0 {
		err = ErrNotFound{Id: r.Id}
		return
	}

	b := bytes.NewBuffer(r.objekte)
	dec := gob.NewDecoder(b)

	if err = dec.Decode(&r.Named); err != nil {
		err = errors.Wrapped(err, "failed to decode zettel: %s", r.Id)
		return
	}

	return
}
