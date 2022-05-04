package user_ops

import (
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type ZettelResults struct {
	Zettelen map[string]stored_zettel.Named
}

func (r ZettelResults) Hinweisen() (h []hinweis.Hinweis) {
	h = make([]hinweis.Hinweis, 0, len(r.Zettelen))

	for _, z := range r.Zettelen {
		h = append(h, z.Hinweis)
	}

	return
}

func (r ZettelResults) HinweisStrings() (h []string) {
	h = make([]string, 0, len(r.Zettelen))

	for i, _ := range r.Zettelen {
		h = append(h, i)
	}

	return
}
