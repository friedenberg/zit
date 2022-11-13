package sha

type Abbr interface {
	AbbreviateSha(Sha) (string, error)
}
