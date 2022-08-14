package etiketten

import (
	"path"
	"strings"

	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/delta/age"
	"github.com/friedenberg/zit/delta/etikett"
	"github.com/friedenberg/zit/delta/konfig"
	"github.com/friedenberg/zit/foxtrot/sharded_store"
)

type Etiketten interface {
	All() (e []etikett.Etikett, err error)
	Add(e1 etikett.Etikett)
	AddString(in string) (err error)
	Flush() (err error)
}

type etiketten struct {
	age    age.Age
	path   string
	shard  sharded_store.Shard
	store  sharded_store.Store
	konfig konfig.Konfig
}

func New(k konfig.Konfig, a age.Age, p string) (e *etiketten, err error) {
	p = path.Join(p, "Etiketten")

	e = &etiketten{
		konfig: k,
		age:    a,
		path:   p,
	}

	if e.shard, err = sharded_store.NewShard(p, a, sharded_store.ShardGeneric{}); err != nil {
		err = errors.Error(err)
		return
	}

	if e.store, err = sharded_store.NewStore(p, e); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (zs etiketten) NewShard(p string, id string) (s sharded_store.Shard, err error) {
	if s, err = sharded_store.NewShard(path.Join(p, id), zs.age, sharded_store.ShardLine{}); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (e *etiketten) AddString(in string) (err error) {
	if strings.TrimSpace(in) == "" {
		return
	}

	var e1 etikett.Etikett

	if err = e1.Set(in); err != nil {
		err = errors.Error(err)
		return
	}

	e.Add(e1)

	return
}

func (e *etiketten) Add(e1 etikett.Etikett) {
	e.shard.Set(e1.Sha().String(), e1.String())

	return
}

func (e etiketten) Flush() (err error) {
	return e.shard.Flush()
}

func (en *etiketten) All() (e []etikett.Etikett, err error) {
	s := en.shard.All()
	e = make([]etikett.Etikett, len(s))

OUTER:
	for i, v := range s {
	INNER:
		for tn, tv := range en.konfig.Tags {
			if !tv.Hide {
				logz.Printf("hiding %s", tn)
				continue INNER
			}

			if tn == v.Value {
				continue OUTER
			}
		}

		if err = e[i].Set(v.Value); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}
