package stored_zettel

import "github.com/friedenberg/zit/charlie/hinweis"

type CollectionNamed interface {
	Hinweisen() (h []hinweis.Hinweis)
	HinweisStrings() (h []string)
}
