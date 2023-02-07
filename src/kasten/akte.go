package kasten

import "github.com/friedenberg/zit/src/uri"

type Akte struct {
	Uri uri.Uri `toml:"uri"`
}

func (a *Akte) Reset() {
	a.Uri = uri.Uri{}
}

func (a *Akte) ResetWith(b Akte) {
	a.Uri = b.Uri
}

func (a Akte) Equals(b Akte) bool {
	if a.Uri != b.Uri {
		return false
	}

	return true
}
