package etiketten

import (
	"log"
	"path"
	"strings"
)

type Etiketten interface {
	All() (e []_Etikett, err error)
	Add(e1 _Etikett)
	AddString(in string) (err error)
	Flush() (err error)
}

type etiketten struct {
	age    _Age
	path   string
	shard  _Shard
	store  _Store
	konfig _Konfig
}

func New(k _Konfig, a _Age, p string) (e *etiketten, err error) {
	p = path.Join(p, "Etiketten")

	e = &etiketten{
		konfig: k,
		age:    a,
		path:   p,
	}

	if e.shard, err = _NewShard(p, a, _ShardGeneric{}); err != nil {
		err = _Error(err)
		return
	}

	if e.store, err = _NewStore(p, e); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (zs etiketten) NewShard(p string, id string) (s _Shard, err error) {
	if s, err = _NewShard(path.Join(p, id), zs.age, _ShardLine{}); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (e *etiketten) AddString(in string) (err error) {
	if strings.TrimSpace(in) == "" {
		return
	}

	var e1 _Etikett

	if err = e1.Set(in); err != nil {
		err = _Error(err)
		return
	}

	e.Add(e1)

	return
}

func (e *etiketten) Add(e1 _Etikett) {
	e.shard.Set(e1.Sha().String(), e1.String())

	return
}

func (e etiketten) Flush() (err error) {
	return e.shard.Flush()
}

func (en *etiketten) All() (e []_Etikett, err error) {
	s := en.shard.All()
	e = make([]_Etikett, len(s))

OUTER:
	for i, v := range s {
	INNER:
		for tn, tv := range en.konfig.Tags {
			if !tv.Hide {
				log.Printf("hiding %s", tn)
				continue INNER
			}

			if tn == v.Value {
				continue OUTER
			}
		}

		if err = e[i].Set(v.Value); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}
