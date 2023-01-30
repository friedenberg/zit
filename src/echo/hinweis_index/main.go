package hinweis_index

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
	hinweis_index_v0 "github.com/friedenberg/zit/src/echo/hinweis_index/v0"
)

type HinweisStore interface {
	schnittstellen.Flusher
	CreateHinweis() (kennung.Hinweis, error)
}

type HinweisIndex interface {
	HinweisStore
	schnittstellen.Resetter
	AddHinweis(kennung.Hinweis) error
	PeekHinweisen(int) ([]kennung.Hinweis, error)
}

func MakeIndex(
	k schnittstellen.Konfig,
	s schnittstellen.Standort,
	su schnittstellen.VerzeichnisseFactory,
) (i HinweisIndex, err error) {
	switch v := k.GetStoreVersion().Int(); {
	case v >= 1:
		errors.TodoP0("fix hinweis_index_v1 and turn on by default")
		// if i, err = hinweis_index_v1.MakeIndex(
		// 	k,
		// 	s,
		// 	su,
		// ); err != nil {
		// 	err = errors.Wrap(err)
		// 	return
		// }
		fallthrough

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
