package kennung_index

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
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
