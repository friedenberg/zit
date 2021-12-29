package organize_text

type Grouper interface {
	GroupZettel(_NamedZettel) []_EtikettSet
}

type Sorter interface {
	SortGroups(_EtikettSet, _EtikettSet) bool
	SortZettels(_NamedZettel, _NamedZettel) bool
}

type OrganizeOptions struct {
	RootEtiketten _EtikettSet
	Grouper
	Sorter
}
