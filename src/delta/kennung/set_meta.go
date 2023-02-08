package kennung

import (
	"encoding/gob"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/gattungen"
)

func init() {
	gob.Register(&metaSet{})
}

type MetaSet interface {
	Set(string) error
	SetMany(...string) error
	All(f func(gattung.Gattung, Set) error) error
}

type metaSet struct {
	expanders      Expanders
	defaultGattung gattung.Gattung
	Gattung        map[gattung.Gattung]Set
}

func MakeMetaSet(ex Expanders, dg gattung.Gattung) *metaSet {
	return &metaSet{
		expanders:      ex,
		defaultGattung: dg,
		Gattung:        make(map[gattung.Gattung]Set),
	}
}

func (s *metaSet) SetMany(vs ...string) (err error) {
	for _, v := range vs {
		if err = s.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (ms *metaSet) Set(v string) (err error) {
	sbs := [3]*strings.Builder{
		{},
		{},
		{},
	}

	sbIdx := 0

	for _, c := range v {
		isSigil := SigilFieldFunc(c)

		switch {
		case isSigil && sbIdx == 0:
			sbIdx += 1

		case isSigil && sbIdx > 1:
			err = errors.Errorf("invalid meta set: %q", v)
			return

		case !isSigil && sbIdx == 1:
			sbIdx += 1
		}

		sbs[sbIdx].WriteRune(c)
	}

	var sigil Sigil

	if err = sigil.Set(sbs[1].String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	before := sbs[0].String()
	after := sbs[2].String()

	var gs gattungen.Set

	if after != "" {
		if gs, err = gattungen.GattungFromString(after); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		gs = gattungen.MakeSet(ms.defaultGattung)
	}

	if err = gs.Each(
		func(g gattung.Gattung) (err error) {
			var ids Set
			ok := false

			if ids, ok = ms.Gattung[g]; !ok {
				ids = MakeSet(ms.expanders)
				ids.Sigil = sigil
			}

			if before != "" {
				if err = ids.Set(before); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			ms.Gattung[g] = ids
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// Runs in parallel
func (ms metaSet) All(f func(gattung.Gattung, Set) error) (err error) {
	errors.TodoP0("lock")
	chErr := make(chan error, len(ms.Gattung))

	for g, s := range ms.Gattung {
		go func(g gattung.Gattung, ids Set) {
			var err error

			defer func() {
				chErr <- err
			}()

			if err = f(g, ids); err != nil {
				if collections.IsStopIteration(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}
		}(g, s)
	}

	for i := 0; i < len(ms.Gattung); i++ {
		err = errors.Join(err, <-chErr)
	}

	return
}
