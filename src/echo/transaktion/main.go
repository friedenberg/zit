package transaktion

import (
	"flag"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/zk_types"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/ts"
)

type Transaktion struct {
	ts.Time
	Objekten []Objekte
}

type Mutter [2]ts.Time

type Objekte struct {
	zk_types.Type
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

	if err = o.Type.Set(vs[0]); err != nil {
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

	switch o.Type {
	case zk_types.TypeZettel:
		o.Id = &hinweis.Hinweis{}
	case zk_types.TypeEtikett:
		o.Id = &etikett.Etikett{}
	case zk_types.TypeAkteTyp:
		//TODO
		fallthrough
	default:
		err = errors.Errorf("unsupported type: %s", o.Type)
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
