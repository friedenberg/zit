package hinweis

type Abbr interface {
	AbbreviateHinweis(h Hinweis) (ha Hinweis, err error)
}
