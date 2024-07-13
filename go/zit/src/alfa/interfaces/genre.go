package interfaces

type StringerGenreGetter interface {
	GenreGetter
	Stringer
}

type StringerGenreRepoIdGetter interface {
	GenreGetter
	RepoIdGetter
	Stringer
}

type RepoId interface {
	Stringer
	EqualsRepoId(RepoIdGetter) bool
	GetRepoIdString() string
}

type RepoIdGetter interface {
	GetRepoId() RepoId
}

type Genre interface {
	StringerGenreGetter
	EqualsGenre(GenreGetter) bool
	GetGenreBitInt() byte
	GetGenreString() string
	GetGenreStringPlural() string
}

type GenreGetter interface {
	GetGenre() Genre
}
