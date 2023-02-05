package kasten

import "net/url"

type Akte struct {
	Uri url.URL `toml:"uri"`
}

func (a *Akte) Reset() {
	a.Uri = url.URL{}
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
