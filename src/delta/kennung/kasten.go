package kennung

import (
	"regexp"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/values"
)

type Kasten = Kennung[kasten, *kasten]

const KastenRegexString = `^(//)?[-a-z0-9_]+$`

var KastenRegex *regexp.Regexp

func init() {
	KastenRegex = regexp.MustCompile(KastenRegexString)
}

func MustKasten(v string) (e Kasten) {
	var err error

	if e, err = makeKennung[kasten, *kasten](v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func MakeKasten(v string) (e Kasten, err error) {
	if e, err = makeKennung[kasten, *kasten](v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type kasten string

func (e *kasten) Reset() {
	*e = kasten("")
}

func (e *kasten) ResetWith(e1 kasten) {
	*e = e1
}

func (a kasten) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a kasten) Equals(b kasten) bool {
	return a == b
}

func (o kasten) GetGattung() schnittstellen.Gattung {
	return gattung.Kasten
}

func (e kasten) String() string {
	return string(e)
}

// func (e etikett) GetQueryPrefix() string {
// 	return "//"
// }

func (e *kasten) Set(v string) (err error) {
	v = strings.TrimPrefix(v, "//")

	if !KastenRegex.Match([]byte(v)) {
		err = errors.Errorf("not a valid kasten: '%s'", v)
		return
	}

	*e = kasten(v)

	return
}
