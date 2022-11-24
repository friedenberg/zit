package kennung

import (
	"regexp"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

const EtikettRegexString = `^[-a-z0-9_/]+$`

var EtikettRegex *regexp.Regexp

func init() {
	EtikettRegex = regexp.MustCompile(EtikettRegexString)
}

type Etikett = Kennung[etikett, *etikett]

func MustEtikett(v string) (e Etikett) {
	var err error

	if e, err = makeKennung[etikett, *etikett](v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func MakeEtikett(v string) (e Etikett, err error) {
	if e, err = makeKennung[etikett, *etikett](v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type etikett string

func (e etikett) String() string {
	return string(e)
}

func (e *etikett) Set(v string) (err error) {
	if !EtikettRegex.Match([]byte(v)) {
		err = errors.Errorf("not a valid tag: '%s'", v)
		return
	}

	*e = etikett(v)

	return
}

func IsDependentLeaf(a Etikett) (has bool) {
	has = strings.HasPrefix(strings.TrimSpace(a.String()), "-")
	return
}

func HasParentPrefix(a, b Etikett) (has bool) {
	has = strings.HasPrefix(strings.TrimSpace(a.String()), b.String())
	return
}
