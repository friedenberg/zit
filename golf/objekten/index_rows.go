package objekten

import (
	"github.com/friedenberg/zit/bravo/zk_types"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/verzeichnisse"
)

func (s Store) indexRowsForZettel(nz stored_zettel.Named) (rs []verzeichnisse.Row, err error) {
	rs = make([]verzeichnisse.Row, 0)

	for _, e := range nz.Zettel.Etiketten.Expanded(etikett.ExpanderAll{}) {
		rs = append(
			rs,
			verzeichnisse.Row{
				Sha:     e.Sha(),
				Key:     e.String(),
				Type:    zk_types.TypeEtikett.String(),
				Objekte: nz,
			},
		)
	}

	rs = append(
		rs,
		verzeichnisse.Row{
			Sha:     nz.Hinweis.Sha(),
			Key:     nz.Hinweis.String(),
			Type:    zk_types.TypeHinweis.String(),
			Objekte: nz,
		},
	)

	if !nz.Zettel.Akte.IsNull() {
		rs = append(
			rs,
			verzeichnisse.Row{
				Sha:     nz.Zettel.Akte,
				Key:     nz.Zettel.Akte.String(),
				Type:    zk_types.TypeAkte.String(),
				Objekte: nz,
			},
		)
	}

	//TODO add akte typ
	//TODO add beziechnung?

	return
}
