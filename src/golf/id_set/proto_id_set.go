package id_set

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/foxtrot/id"
)

type ProtoIdSet struct {
	types []protoId
}

func MakeProtoIdSet(types ...ProtoId) (ps ProtoIdSet) {
	ps.types = make([]protoId, 0, len(types))

	for _, t := range types {
		ps.Add(t)
	}

	return
}

func (ps *ProtoIdSet) Add(t ProtoId) (err error) {
	pid := makeProtoId(t)
	ps.types = append(ps.types, pid)

	return
}

func (ps *ProtoIdSet) AddMany(ts ...ProtoId) (err error) {
	for _, t := range ts {
		ps.Add(t)
	}

	return
}

func (ps ProtoIdSet) Len() int {
	return len(ps.types)
}

func (ps ProtoIdSet) Contains(i id.MutableId) (ok bool) {
	i2 := makeProtoId(ProtoId{MutableId: i})
	for _, i1 := range ps.types {
		if i1.Type == i2.Type {
			ok = true
			break
		}
	}

	return
}

func (ps ProtoIdSet) MakeOne(v string) (i gattung.IdLike, err error) {
	for _, t := range ps.types {
		if i, err = t.Make(v); err == nil {
			break
		}
	}

	switch {
	case err != nil && len(ps.types) == 1:
		return

	case err != nil:
		err = errors.Errorf("no proto id was able to parse: %q", v)
		return
	}

	return
}

func (ps ProtoIdSet) Make(vs ...string) (s Set, err error) {
	s = Make(len(vs))

	for _, v := range vs {
		var i gattung.IdLike

		if i, err = ps.MakeOne(v); err != nil {
			err = errors.Wrap(err)
			return
		}

		s.Add(i)
	}

	return
}
