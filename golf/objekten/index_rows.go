package objekten

import (
	"bytes"
	"encoding/gob"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/zk_types"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/charlie/verzeichnisse"
)

func (s Store) indexRowsForZettel(nz stored_zettel.Named) (rs []verzeichnisse.Row, err error) {
	rs = make([]verzeichnisse.Row, 0)

	b := &bytes.Buffer{}
	enc := gob.NewEncoder(b)

	if err = enc.Encode(nz); err != nil {
		err = errors.Wrapped(err, "failed to encode zettel: %s", nz.Hinweis)
		return
	}

	for _, e := range nz.Zettel.Etiketten.Expanded(etikett.ExpanderAll{}) {
		rs = append(
			rs,
			verzeichnisse.Row{
				Sha:     e.Sha(),
				Key:     e.String(),
				Type:    zk_types.TypeEtikett.String(),
				Objekte: b.Bytes(),
			},
		)
	}

	rs = append(
		rs,
		verzeichnisse.Row{
			Sha:     nz.Hinweis.Sha(),
			Key:     nz.Hinweis.String(),
			Type:    zk_types.TypeHinweis.String(),
			Objekte: b.Bytes(),
		},
	)

	rs = append(
		rs,
		verzeichnisse.Row{
			Sha:     nz.Stored.Sha,
			Key:     nz.Stored.Sha.String(),
			Type:    zk_types.TypeZettel.String(),
			Objekte: b.Bytes(),
		},
	)

	if !nz.Zettel.Akte.IsNull() {
		rs = append(
			rs,
			verzeichnisse.Row{
				Sha:     nz.Zettel.Akte,
				Key:     nz.Zettel.Akte.String(),
				Type:    zk_types.TypeAkte.String(),
				Objekte: b.Bytes(),
			},
		)
	}

	//TODO add akte typ
	//TODO add beziechnung?

	return
}
