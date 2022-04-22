package etikett

import (
	"crypto/sha256"
	"io"
	"regexp"
	"strings"
)

const EtikettRegexString = `^[-a-z0-9_/]+$`

var EtikettRegex *regexp.Regexp

func init() {
	EtikettRegex = regexp.MustCompile(EtikettRegexString)
}

type Etikett struct {
	Value string
	sha   _Sha
}

func (e Etikett) Sha() _Sha {
	return e.sha
}

func (e Etikett) String() string {
	return e.Value
}

func (e *Etikett) Set(v string) (err error) {
	v1 := strings.ToLower(v)
	v2 := strings.TrimPrefix(v1, "-")
	v3 := strings.TrimSpace(v2)

	if !EtikettRegex.Match([]byte(v3)) {
		err = _Errorf("not a valid tag: '%s'", v)
		return
	}

	hash := sha256.New()
	sr := strings.NewReader(v3)

	if _, err = io.Copy(hash, sr); err != nil {
		err = _Error(err)
		return
	}

	e.sha = _MakeShaFromHash(hash)

	e.Value = v3

	return
}

func (a Etikett) Equals(b Etikett) bool {
	return a.Value == b.Value
}

func (e Etikett) Expanded(exes ...Expander) (expanded Set) {
	expanded = NewSet()

	if len(exes) == 0 {
		exes = []Expander{ExpanderAll{}}
	}

	for _, ex := range exes {
		for _, e := range ex.Expand(e) {
			expanded.addOnlyExact(e)
		}
	}

	return
}
