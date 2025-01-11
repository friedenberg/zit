package repo_blobs

import (
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
)

type V0 struct {
	Uri values.Uri `toml:"uri"`
}

func (b V0) GetRepoBlob() Blob {
	return b
}

func (b V0) GetRepoType() {
}

func (a *V0) Reset() {
	a.Uri = values.Uri{}
}

func (a *V0) ResetWith(b V0) {
	a.Uri = b.Uri
}

func (a V0) Equals(b V0) bool {
	if a.Uri != b.Uri {
		return false
	}

	return true
}
