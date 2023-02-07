package kennung

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type ProtoIdSet struct {
	types  []protoId
	always []protoId
}

func MakeProtoIdSet(types ...ProtoId) (ps ProtoIdSet) {
	ps.types = make([]protoId, 0, len(types))

	for _, t := range types {
		ps.Add(t)
	}

	ps.always = []protoId{
		makeProtoId(
			ProtoId{
				Setter: MakeSigil(SigilNone),
			},
		),
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

func (ps ProtoIdSet) Contains(i schnittstellen.Setter) (ok bool) {
	i2 := makeProtoId(ProtoId{Setter: i})
	for _, i1 := range ps.types {
		if i1.Type == i2.Type {
			ok = true
			break
		}
	}

	for _, i1 := range ps.always {
		if i1.Type == i2.Type {
			ok = true
			break
		}
	}

	return
}

func (ps ProtoIdSet) MakeOne(v string) (i schnittstellen.Value, err error) {
	for _, t := range ps.types {
		if i, err = t.Make(v); err == nil {
			break
		}
	}

	if i != nil && err == nil {
		return
	}

	for _, t := range ps.always {
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
	s = MakeSet()

	for _, v := range vs {
		var i schnittstellen.Value

		if i, err = ps.MakeOne(v); err != nil {
			err = errors.Wrap(err)
			return
		}

		s.Add(i)
	}

	return
}
