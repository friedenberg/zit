package kennung

import (
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/gattungen"
)

func init() {
	gob.Register(&metaSet{})
}

// TODO rename to QueryGattungGroup
type MetaSet interface {
	Get(g gattung.Gattung) (s Matcher, ok bool)
	GetIds(g gattung.Gattung) (s Set, ok bool)
	Set(string) error
	SetMany(...string) error
	All(f func(gattung.Gattung, Set) error) error
}

type metaSet struct {
	cwd                 Matcher
	fileExtensionGetter schnittstellen.FileExtensionGetter
	expanders           Expanders
	Hidden              Matcher
	DefaultGattungen    gattungen.Set
	Gattung             map[gattung.Gattung]Set
}

func MakeMetaSet(
	cwd Matcher,
	ex Expanders,
	hidden Matcher,
	feg schnittstellen.FileExtensionGetter,
	dg gattungen.Set,
) MetaSet {
	return &metaSet{
		cwd:                 cwd,
		fileExtensionGetter: feg,
		expanders:           ex,
		Hidden:              hidden,
		DefaultGattungen:    dg.MutableClone(),
		Gattung:             make(map[gattung.Gattung]Set),
	}
}

func MakeMetaSetAll(
	cwd Matcher,
	ex Expanders,
	hidden Matcher,
	feg schnittstellen.FileExtensionGetter,
) MetaSet {
	errors.TodoP2("support allowed sigils")
	return &metaSet{
		cwd:                 cwd,
		fileExtensionGetter: feg,
		expanders:           ex,
		Hidden:              hidden,
		DefaultGattungen:    gattungen.MakeSet(gattung.TrueGattung()...),
		Gattung:             make(map[gattung.Gattung]Set),
	}
}

func (s metaSet) String() string {
	sb := &strings.Builder{}

	for g, ids := range s.Gattung {
		sb.WriteString(fmt.Sprintf("%s%s", ids, g))
	}

	return sb.String()
}

func (s *metaSet) SetMany(vs ...string) (err error) {
	for _, v := range vs {
		if err = s.set(v); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (ms *metaSet) Set(v string) (err error) {
	return ms.set(v)
}

func (ms *metaSet) set(v string) (err error) {
	v = strings.TrimSpace(v)

	// if v != "." {
	// 	if err = collections.AddString[FD, *FD](
	// 		ms.FDs,
	// 		v,
	// 	); err == nil {
	// 		ms.Gattung = make(map[gattung.Gattung]Set)
	// 		return
	// 	}

	// 	err = nil
	// }

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
			if gattung.IsErrUnrecognizedGattung(err) {
				err = nil

				var ids Set
				ok := false

				if ids, ok = ms.Gattung[gattung.Unknown]; !ok {
					ids = ms.MakeSet()
				}

				if err = collections.AddString[FD, *FD](
					ids.FDs,
					v,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				ms.Gattung[gattung.Unknown] = ids

			} else {
				err = errors.Wrap(err)
			}

			return
		}
	} else {
		gs = ms.DefaultGattungen.ImmutableClone()
	}

	if err = gs.Each(
		func(g gattung.Gattung) (err error) {
			var ids Set
			ok := false

			if ids, ok = ms.Gattung[g]; !ok {
				ids = ms.MakeSet()
				ids.Sigil = sigil
			}

			switch {
			case before == "":
				break

			case ids.Sigil.IncludesCwd():
				fp := fmt.Sprintf("%s.%s", before, after)

				var fd FD

				if fd, err = FDFromPath(fp); err == nil {
					ids.Add(fd)
					break
				}

				err = nil

				fallthrough

			default:
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

func (ms metaSet) Get(g gattung.Gattung) (s Matcher, ok bool) {
	s, ok = ms.Gattung[g]
	return
}

func (ms metaSet) GetIds(g gattung.Gattung) (s Set, ok bool) {
	s, ok = ms.Gattung[g]
	return
}

func (ms metaSet) MakeSet() Set {
	return MakeSet(ms.cwd, ms.expanders, ms.Hidden)
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

// func (s metaSet) MarshalBinary() (bs []byte, err error) {
// 	b := bytes.NewBuffer(bs)
// 	enc := gob.NewEncoder(b)

// 	if err = enc.Encode(s.Gattung); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	bs = b.Bytes()

// 	return
// }

// func (s *metaSet) UnmarshalBinary(bs []byte) (err error) {
// 	err = errors.New("wow")
// 	return

// 	// 	b := bytes.NewBuffer(bs)
// 	// 	dec := gob.NewDecoder(b)

// 	// 	if err = dec.Decode(&s.Gattung); err != nil {
// 	// 		err = errors.Wrap(err)
// 	// 		return
// 	// 	}

// 	// return
// }
