package kennung

import (
	"crypto/sha256"
	"io"
	"regexp"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
)

const KennungRegexString = `^[-a-z0-9_/]+$`

var KennungRegex *regexp.Regexp

func init() {
	KennungRegex = regexp.MustCompile(KennungRegexString)
}

type Kennung struct {
	Value string
}

func Make(v string) Kennung {
	return Kennung{Value: v}
}

func (e Kennung) Sha() sha.Sha {
	hash := sha256.New()
	sr := strings.NewReader(e.String())

	if _, err := io.Copy(hash, sr); err != nil {
		errors.PanicIfError(err)
	}

	return sha.FromHash(hash)
}

func (e Kennung) String() string {
	return e.Value
}

func (e *Kennung) Set(v string) (err error) {
	v1 := strings.ToLower(v)
	v3 := strings.TrimSpace(v1)

	if v3 == "" {
		err = errors.Errorf("etikett cannot be empty")
		return
	}

	if !KennungRegex.Match([]byte(v3)) {
		err = errors.Errorf("not a valid tag: '%s'", v)
		return
	}

	e.Value = v3

	return
}

func (e Kennung) Len() int {
	return len(e.Value)
}

func (a Kennung) Includes(b Kennung) bool {
	return b.Contains(a)
}

func (a Kennung) Contains(b Kennung) bool {
	if b.Len() > a.Len() {
		return false
	}

	return strings.HasPrefix(a.Value, b.Value)
}

func (a Kennung) Equals(b Kennung) bool {
	return a.Value == b.Value
}

func (a Kennung) Less(b Kennung) bool {
	return a.Value < b.Value
}

func (a Kennung) LeftSubtract(b Kennung) (c Kennung) {
	return Kennung{Value: strings.TrimPrefix(a.String(), b.String())}
}

func (a Kennung) IsEmpty() bool {
	return a.Value == ""
}

func (a Kennung) IsDependentLeaf() (has bool) {
	has = strings.HasPrefix(strings.TrimSpace(a.Value), "-")
	return
}

func (a Kennung) HasParentPrefix(b Kennung) (has bool) {
	has = strings.HasPrefix(strings.TrimSpace(a.Value), b.Value)
	return
}

// func (e Kennung) Expanded(exes ...ExpanderKennung) (out Set) {
// 	expanded := MakeMutableSet()

// 	if len(exes) == 0 {
// 		exes = []ExpanderKennung{ExpanderKennungAll}
// 	}

// 	for _, ex := range exes {
// 		ex.Expand(e.String()).Each(expanded.Add)
// 	}

// 	out = Set(expanded.Copy())

// 	return
// }
