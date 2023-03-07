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
	Add(schnittstellen.IdLike, Sigil) error
	Get(g gattung.Gattung) (s Set, ok bool)
	AddFDs(FDSet) error
	Set(string) error
	SetMany(...string) error
	All(f func(gattung.Gattung, Set) error) error
}

type metaSet struct {
	cwd                 Matcher
	fileExtensionGetter schnittstellen.FileExtensionGetter
	expanders           Expanders
	Etiketten           QuerySet[Etikett, *Etikett]
	DefaultGattungen    gattungen.Set
	Gattung             map[gattung.Gattung]Set
}

func MakeMetaSet(
	cwd Matcher,
	ex Expanders,
	etikett QuerySet[Etikett, *Etikett],
	feg schnittstellen.FileExtensionGetter,
	dg gattungen.Set,
) MetaSet {
	return &metaSet{
		cwd:                 cwd,
		fileExtensionGetter: feg,
		expanders:           ex,
		Etiketten:           etikett,
		DefaultGattungen:    dg.MutableClone(),
		Gattung:             make(map[gattung.Gattung]Set),
	}
}

func MakeMetaSetAll(
	cwd Matcher,
	ex Expanders,
	etikett QuerySet[Etikett, *Etikett],
	feg schnittstellen.FileExtensionGetter,
) MetaSet {
	errors.TodoP2("support allowed sigils")
	return &metaSet{
		cwd:                 cwd,
		fileExtensionGetter: feg,
		expanders:           ex,
		Etiketten:           etikett,
		DefaultGattungen:    gattungen.MakeSet(gattung.All()...),
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
		if err = s.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (ms *metaSet) Set(v string) (err error) {
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
					ids = MakeSet(ms.cwd, ms.expanders, ms.Etiketten)
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
				def := ms.Etiketten

				if sigil.IncludesHidden() {
					def = nil
				}

				ids = MakeSet(ms.cwd, ms.expanders, def)
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

func (ms *metaSet) Add(
	id schnittstellen.IdLike,
	sigil Sigil,
) (err error) {
	g := gattung.Make(id.GetGattung())

	var ids Set
	ok := false

	if ids, ok = ms.Gattung[g]; !ok {
		def := ms.Etiketten

		if sigil.IncludesHidden() {
			def = nil
		}

		ids = MakeSet(ms.cwd, ms.expanders, def)
		ids.Sigil = sigil
	}

	if err = ids.Add(id); err != nil {
		err = errors.Wrap(err)
		return
	}

	ms.Gattung[g] = ids

	return
}

func (ms metaSet) Get(g gattung.Gattung) (s Set, ok bool) {
	s, ok = ms.Gattung[g]
	return
}

func (ms metaSet) addFDs(fd FD) (err error) {
	ext := fd.ExtSansDot()

	g := gattung.MakeOrUnknown(ext)

	ok := false
	var ids Set

	if ids, ok = ms.Gattung[g]; !ok {
		ids = MakeSet(ms.cwd, ms.expanders, ms.Etiketten)
	}

	if err = ids.Add(fd); err != nil {
		err = errors.Wrap(err)
		return
	}

	ms.Gattung[g] = ids

	return
}

func (ms metaSet) AddFDs(fds FDSet) (err error) {
	return fds.Each(ms.addFDs)
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
