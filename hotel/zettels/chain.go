package zettels

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type Chain struct {
	Hinweis hinweis.Hinweis
	//stored in reverse (latest is at 0)
	Zettels []stored_zettel.Stored
}

func (zs zettels) AllInChain(h hinweis.Hinweis) (c Chain, err error) {
	var s sha.Sha

	if s, c.Hinweis, err = zs.TailFromId(h); err != nil {
		err = errors.Error(err)
		return
	}

	shas := make(map[string]bool)

	for {
		if s.IsNull() {
			break
		}

		if _, ok := shas[s.String()]; ok {
			err = ErrHistoryLoopDetected{Sha: s}
			return
		}

		shas[s.String()] = true

		var sz stored_zettel.Named

		if sz, err = zs.Read(s); err != nil {
			err = errors.Error(err)
			return
		}

		c.Zettels = append(
			c.Zettels,
			sz.Stored,
		)

		s = sz.Mutter
	}

	return
}

func (zs zettels) TailFromId(id id.Id) (s sha.Sha, h hinweis.Hinweis, err error) {
	ok := false

	if s, ok = id.(sha.Sha); ok {
		if h, err = zs.hinweisen.ReadSha(s); err != nil {
			err = errors.Error(err)
			return
		}
	} else {
		if h, ok = id.(hinweis.Hinweis); !ok {
			err = errors.Errorf("unsupported id: '%q'", id)
			return
		}

		if s, err = zs.hinweisen.Read(h); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}
