package objekten

import (
	"bytes"
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/zk_types"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/verzeichnisse"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
)

type indexReaderManyTail struct {
	chRow   chan verzeichnisse.Row
	chDone  chan struct{}
	rows    map[string][]byte
	zettels map[hinweis.Hinweis]stored_zettel.Named
}

func (r *indexReaderManyTail) Begin() (err error) {
	r.rows = make(map[string][]byte)
	r.zettels = make(map[hinweis.Hinweis]stored_zettel.Named)
	r.chRow = make(chan verzeichnisse.Row)
	r.chDone = make(chan struct{})

	go func() {
		defer func() {
			logz.Print("sending done")
			r.chDone <- struct{}{}
		}()

		for row := range r.chRow {
			logz.Print(row.Key)
			r.rows[row.Key] = row.Objekte
		}
	}()

	return
}

func (r *indexReaderManyTail) ReadRow(id string, row verzeichnisse.Row) (err error) {
	if row.Type != zk_types.TypeHinweis.String() {
		return
	}

	r.chRow <- row

	return
}

func (r *indexReaderManyTail) End() (err error) {
	close(r.chRow)

	logz.Print("waiting for done")
	<-r.chDone
	logz.Print("done")

	for _, row := range r.rows {
		b := bytes.NewBuffer(row)
		dec := gob.NewDecoder(b)

		var nz stored_zettel.Named

		if err = dec.Decode(&nz); err != nil {
			err = errors.Wrapped(err, "failed to decode zettel")
			return
		}

		r.zettels[nz.Hinweis] = nz
	}

	return
}
