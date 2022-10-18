package transaktion

import (
	"flag"
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/charlie/typ"
)

type Transaktion struct {
	ts.Time
	Objekten map[string]Objekte
}

func (t *Transaktion) AddObjekte(o Objekte) {
	k := o.GetKey()

	o1, ok := t.Objekten[k]

	if ok {
		//TODO migrate to a hard fail here
		errors.Printf(
			"Transaktion %s has duplicate entries: (%s %s %s) & (%s %s %s)",
			t.Time,
			o1.Gattung,
			o1.Id,
			o1.Sha,
			o.Gattung,
			o.Id,
			o.Sha,
		)
	}

	t.Objekten[k] = o
}

type Mutter [2]ts.Time

type Objekte struct {
	gattung.Gattung
	Mutter
	Id flag.Value
	sha.Sha
}

func (o *Objekte) Set(v string) (err error) {
	vs := strings.Split(v, " ")

	if len(vs) != 5 {
		err = errors.Errorf("expected 5 elements but got %d", len(vs))
		return
	}

	if err = o.Gattung.Set(vs[0]); err != nil {
		err = errors.Wrapf(err, "failed to set type: %s", vs[0])
		return
	}

	vs = vs[1:]

	if err = o.Mutter[0].Set(vs[0]); err != nil {
		err = errors.Wrapf(err, "failed to set mutter 0: %s", vs[0])
		return
	}

	vs = vs[1:]

	if err = o.Mutter[1].Set(vs[0]); err != nil {
		err = errors.Wrapf(err, "failed to set mutter 1: %s", vs[0])
		return
	}

	vs = vs[1:]

	switch o.Gattung {
	case gattung.Zettel:
		o.Id = &hinweis.Hinweis{}

	case gattung.Etikett:
		o.Id = &etikett.Etikett{}

	case gattung.Typ:
		o.Id = &typ.Typ{}

	default:
		err = errors.Errorf("unsupported gattung: %s", o.Gattung)
		return
	}

	if err = o.Id.Set(vs[0]); err != nil {
		err = errors.Wrapf(err, "failed to set id: %s", vs[1])
		return
	}

	vs = vs[1:]

	if err = o.Sha.Set(vs[0]); err != nil {
		err = errors.Wrapf(err, "failed to set sha: %s", vs[2])
		return
	}

	return
}

func (o Objekte) GetKey() string {
	return fmt.Sprintf("%s.%s", o.Gattung, o.Id)
}
