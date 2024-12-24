package interfaces

type ObjectId interface {
	GenreGetter
	Stringer
}

type ObjectIdWithRepoId interface {
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
	ObjectId
	EqualsGenre(GenreGetter) bool
	GetGenreBitInt() byte
	GetGenreString() string
	GetGenreStringPlural(StoreVersion) string
}

type GenreGetter interface {
	GetGenre() Genre
}
