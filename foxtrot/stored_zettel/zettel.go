package stored_zettel

type Stored struct {
	Mutter _Sha
	Kinder _Sha
	Sha    _Sha
	Zettel _Zettel
}

type Named struct {
	Stored
	Hinweis _Hinweis
}
