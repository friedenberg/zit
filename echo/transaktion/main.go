package transaktion

import (
	"flag"
	"strings"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/bravo/zk_types"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/charlie/ts"
)

type Transaktion struct {
	ts.Time
	Objekten []Objekte
}

type Objekte struct {
	zk_types.Type
	Id flag.Value
	sha.Sha
}

func (o *Objekte) Set(v string) (err error) {
	vs := strings.Split(v, " ")

	if len(vs) != 3 {
		err = errors.Errorf("expected 3 elements but got %d", len(vs))
		return
	}

	if err = o.Type.Set(vs[0]); err != nil {
		err = errors.Wrapped(err, "failed to set type: %s", vs[0])
		return
	}

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

	if err = o.Id.Set(vs[1]); err != nil {
		err = errors.Wrapped(err, "failed to set id: %s", vs[1])
		return
	}

	if err = o.Sha.Set(vs[2]); err != nil {
		err = errors.Wrapped(err, "failed to set sha: %s", vs[2])
		return
	}

	return
}
