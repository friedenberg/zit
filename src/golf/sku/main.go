package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/int_value"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/ts"
)

type Mutter [2]ts.Time

type IdLike = fmt.Stringer

type SkuLike interface {
	SetFields(...string) error
	GetKey() string
	SetTransactionIndex(int)
	GetGattung() gattung.Gattung
	GetId() IdLike
	GetMutter() Mutter
	GetObjekteSha() sha.Sha
	GetTransactionIndex() int_value.IntValue
	GetKopf() ts.Time
	GetSchwanz() ts.Time
}

// type Sku struct {
// 	Gattung          gattung.Gattung
// 	Mutter           Mutter
// 	Id               fmt.Stringer
// 	Sha              sha.Sha
// 	TransactionIndex int_value.IntValue
// }

// func (a Sku) Equals(b Sku) (ok bool) {
// 	if a.Gattung != b.Gattung {
// 		return
// 	}

// 	if a.Mutter != b.Mutter {
// 		return
// 	}

// 	if a.Id.String() != b.Id.String() {
// 		return
// 	}

// 	if !a.Sha.Equals(b.Sha) {
// 		return
// 	}

// 	return true
// }

// func (o *Sku) Set(v string) (err error) {
// 	vs := strings.Split(v, " ")

// 	if len(vs) != 5 {
// 		err = errors.Errorf("expected 5 elements but got %d", len(vs))
// 		return
// 	}

// 	if err = o.Gattung.Set(vs[0]); err != nil {
// 		err = errors.Wrapf(err, "failed to set type: %s", vs[0])
// 		return
// 	}

// 	vs = vs[1:]

// 	if err = o.Mutter[0].Set(vs[0]); err != nil {
// 		err = errors.Wrapf(err, "failed to set mutter 0: %s", vs[0])
// 		return
// 	}

// 	vs = vs[1:]

// 	if err = o.Mutter[1].Set(vs[0]); err != nil {
// 		err = errors.Wrapf(err, "failed to set mutter 1: %s", vs[0])
// 		return
// 	}

// 	vs = vs[1:]

// 	var id flag.Value

// 	switch o.Gattung {
// 	case gattung.Zettel:
// 		id = &hinweis.Hinweis{}

// 	case gattung.Etikett:
// 		id = &kennung.Etikett{}

// 	case gattung.Typ:
// 		id = &kennung.Typ{}

// 	case gattung.Konfig:
// 		id = &kennung.Konfig{}

// 	default:
// 		err = errors.Errorf("unsupported gattung: %s", o.Gattung)
// 		return
// 	}

// 	if err = id.Set(vs[0]); err != nil {
// 		err = errors.Wrapf(err, "failed to set id: %s", vs[1])
// 		return
// 	}

// 	o.Id = id

// 	vs = vs[1:]

// 	if err = o.Sha.Set(vs[0]); err != nil {
// 		err = errors.Wrapf(err, "failed to set sha: %s", vs[2])
// 		return
// 	}

// 	return
// }

// func (o Sku) GetKey() string {
// 	return fmt.Sprintf("%s.%s", o.Gattung, o.Id)
// }
