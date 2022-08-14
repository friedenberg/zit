package stored_zettel

import "github.com/friedenberg/zit/src/delta/hinweis"

type CollectionNamed interface {
	Hinweisen() (h []hinweis.Hinweis)
	HinweisStrings() (h []string)
}
