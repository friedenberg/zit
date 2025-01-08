package interfaces

type ObjectId interface {
	GenreGetter
	Stringer
	GetObjectIdString() string
	Parts() [3]string
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
	GenreGetter
	Stringer
	EqualsGenre(GenreGetter) bool
	GetGenreBitInt() byte
	GetGenreString() string
	GetGenreStringPlural(StoreVersion) string
}

type GenreGetter interface {
	GetGenre() Genre
}
