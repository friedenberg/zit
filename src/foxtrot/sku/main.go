package sku

import (
	"flag"
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/kennung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/typ"
)

type Mutter [2]ts.Time

type Sku struct {
	gattung.Gattung
	Mutter
	Id flag.Value
	sha.Sha
}

func (a Sku) Equals(b Sku) (ok bool) {
	if a.Gattung != b.Gattung {
		return
	}

	if a.Mutter != b.Mutter {
		return
	}

	if a.Id.String() != b.Id.String() {
		return
	}

	if !a.Sha.Equals(b.Sha) {
		return
	}

	return true
}

func (o *Sku) Set(v string) (err error) {
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
		o.Id = &kennung.Etikett{}

	case gattung.Typ:
		o.Id = &typ.Kennung{}

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

func (o Sku) GetKey() string {
	return fmt.Sprintf("%s.%s", o.Gattung, o.Id)
}
