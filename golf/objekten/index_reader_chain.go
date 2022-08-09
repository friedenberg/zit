package objekten

import (
	"bytes"
	"encoding/gob"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/zk_types"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/verzeichnisse"
)

type indexReaderChain struct {
	hinweis.Hinweis
	zettels []stored_zettel.Stored
}

func (r *indexReaderChain) ReadRow(id string, row verzeichnisse.Row) (err error) {
	if row.Type != zk_types.TypeHinweis.String() {
		return
	}

	if row.Key != r.Hinweis.String() {
		return
	}

	b := bytes.NewBuffer(row.Objekte)
	dec := gob.NewDecoder(b)

	var nz stored_zettel.Named

	if err = dec.Decode(&nz); err != nil {
		err = errors.Wrapped(err, "failed to decode zettel: %s", r.Hinweis)
		return
	}

	r.zettels = append(r.zettels, nz.Stored)

	return
}

func (r *indexReaderChain) Done() (err error) {
	if len(r.zettels) == 0 {
		err = ErrNotFound{Id: r.Hinweis}
		return
	}

	return
}
