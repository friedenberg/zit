package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/schnittstellen"
)

type Typ = Kennung[typ, *typ]

func MustTyp(v string) (e Typ) {
	var err error

	if e, err = makeKennung[typ, *typ](v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func MakeTyp(v string) (e Typ, err error) {
	if e, err = makeKennung[typ, *typ](v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type typ string

func (e *typ) Reset(e1 *typ) {
	if e1 == nil {
		*e = typ("")
	} else {
		*e = *e1
	}

	return
}

func (e typ) Equals(e1 *typ) bool {
	if e1 == nil {
		return false
	}

	return e == *e1
}

func (o typ) GetGattung() schnittstellen.Gattung {
	return gattung.Typ
}

func (e typ) String() string {
	return string(e)
}

func (e *typ) Set(v string) (err error) {
	v = strings.TrimSpace(strings.Trim(v, ".! "))

	if !EtikettRegex.Match([]byte(v)) {
		err = errors.Errorf("not a valid typ: '%s'", v)
		return
	}

	*e = typ(v)

	return
}
