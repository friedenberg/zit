package etikett

import (
	"crypto/sha256"
	"io"
	"regexp"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
)

const EtikettRegexString = `^[-a-z0-9_/]+$`

var EtikettRegex *regexp.Regexp

func init() {
	EtikettRegex = regexp.MustCompile(EtikettRegexString)
}

type Etikett struct {
	Value string
}

func (e Etikett) Sha() sha.Sha {
	hash := sha256.New()
	sr := strings.NewReader(e.String())

	if _, err := io.Copy(hash, sr); err != nil {
		errors.PanicIfError(err)
	}

	return sha.FromHash(hash)
}

func (e Etikett) String() string {
	return e.Value
}

func (e *Etikett) Set(v string) (err error) {
	v1 := strings.ToLower(v)
	v3 := strings.TrimSpace(v1)

	if !EtikettRegex.Match([]byte(v3)) {
		err = errors.Errorf("not a valid tag: '%s'", v)
		return
	}

	e.Value = v3

	return
}

func (a Etikett) Equals(b Etikett) bool {
	return a.Value == b.Value
}

func (a Etikett) LeftSubtract(b Etikett) (c Etikett) {
	return Etikett{Value: strings.TrimPrefix(a.String(), b.String())}
}

func (a Etikett) IsEmpty() bool {
	return a.Value == ""
}

func (a Etikett) IsDependentLeaf() (has bool) {
	has = strings.HasPrefix(strings.TrimSpace(a.Value), "-")
	return
}

func (a Etikett) HasParentPrefix(b Etikett) (has bool) {
	has = strings.HasPrefix(strings.TrimSpace(a.Value), b.Value)
	return
}

func (e Etikett) Expanded(exes ...Expander) (expanded *Set) {
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
