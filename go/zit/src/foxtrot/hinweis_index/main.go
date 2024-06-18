package hinweis_index

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	hinweis_index_v0 "code.linenisgreat.com/zit/go/zit/src/foxtrot/hinweis_index/v0"
	hinweis_index_v1 "code.linenisgreat.com/zit/go/zit/src/foxtrot/hinweis_index/v1"
)

type HinweisStore interface {
	schnittstellen.Flusher
	CreateHinweis() (*kennung.Hinweis, error)
}

type HinweisIndex interface {
	HinweisStore
	schnittstellen.ResetterWithError
	AddHinweis(kennung.Kennung) error
	PeekHinweisen(int) ([]*kennung.Hinweis, error)
}

func MakeIndex(
	k schnittstellen.Konfig,
	s schnittstellen.Standort,
	su schnittstellen.VerzeichnisseFactory,
) (i HinweisIndex, err error) {
	switch v := k.GetStoreVersion().GetInt(); {
	case v >= 1 && false:
		errors.TodoP3("investigate using bitsets")
		if i, err = hinweis_index_v1.MakeIndex(
			k,
			s,
			su,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		if i, err = hinweis_index_v0.MakeIndex(
			k,
			s,
			su,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
