package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
)

type MatcherBuilder struct {
	implicitEtikettenGetter ImplicitEtikettenGetter
	cwd                     Matcher
	fileExtensionGetter     schnittstellen.FileExtensionGetter
	expanders               Expanders
	hidden                  Matcher
	defaultGattungen        gattungen.Set
	gattung                 map[gattung.Gattung]Set
}

func (mb *MatcherBuilder) WithCwd(
	cwd Matcher,
) *MatcherBuilder {
	mb.cwd = cwd
	return mb
}

func (mb *MatcherBuilder) WithFileExtensionGetter(
	feg schnittstellen.FileExtensionGetter,
) *MatcherBuilder {
	mb.fileExtensionGetter = feg
	return mb
}

func (mb *MatcherBuilder) WithExpanders(
	expanders Expanders,
) *MatcherBuilder {
	mb.expanders = expanders
	return mb
}

func (mb *MatcherBuilder) WithDefaultGattungen(
	defaultGattungen gattungen.Set,
) *MatcherBuilder {
	mb.defaultGattungen = defaultGattungen
	return mb
}

func (mb *MatcherBuilder) WithHidden(
	hidden Matcher,
) *MatcherBuilder {
	mb.hidden = hidden
	return mb
}

func (mb *MatcherBuilder) WithImplicitEtikettenGetter(
	ieg ImplicitEtikettenGetter,
) *MatcherBuilder {
	mb.implicitEtikettenGetter = ieg
	return mb
}

func (mb MatcherBuilder) Build(v string) (m Matcher, err error) {
	// v = strings.TrimSpace(v)

	// sbs := [3]*strings.Builder{
	// 	{},
	// 	{},
	// 	{},
	// }

	// sbIdx := 0

	// for _, c := range v {
	// 	isSigil := SigilFieldFunc(c)

	// 	switch {
	// 	case isSigil && sbIdx == 0:
	// 		sbIdx += 1

	// 	case isSigil && sbIdx > 1:
	// 		err = errors.Errorf("invalid meta set: %q", v)
	// 		return

	// 	case !isSigil && sbIdx == 1:
	// 		sbIdx += 1
	// 	}

	// 	sbs[sbIdx].WriteRune(c)
	// }

	// var sigil Sigil

	// if err = sigil.Set(sbs[1].String()); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// before := sbs[0].String()
	// after := sbs[2].String()

	// var gs gattungen.Set

	// if after != "" {
	// 	if gs, err = gattungen.GattungFromString(after); err != nil {
	// 		if gattung.IsErrUnrecognizedGattung(err) {
	// 			err = nil

	// 			if err = collections.AddString[FD, *FD](
	// 				ms.FDs,
	// 				v,
	// 			); err != nil {
	// 				err = errors.Wrap(err)
	// 				return
	// 			}

	// 		} else {
	// 			err = errors.Wrap(err)
	// 		}

	// 		return
	// 	}
	// } else {
	// 	gs = ms.DefaultGattungen.ImmutableClone()
	// }

	// if err = gs.Each(
	// 	func(g gattung.Gattung) (err error) {
	// 		var ids Set
	// 		ok := false

	// 		if ids, ok = ms.Gattung[g]; !ok {
	// 			ids = ms.MakeSet()
	// 			ids.AddSigil(sigil)
	// 		}

	// 		switch {
	// 		case before == "":
	// 			break

	// 		case ids.Sigil.IncludesCwd():
	// 			fp := fmt.Sprintf("%s.%s", before, after)

	// 			var fd FD

	// 			if fd, err = FDFromPath(fp); err == nil {
	// 				ids.Add(fd)
	// 				break
	// 			}

	// 			err = nil

	// 			fallthrough

	// 		default:
	// 			if err = ids.Set(before); err != nil {
	// 				err = errors.Wrap(err)
	// 				return
	// 			}
	// 		}

	// 		if g.Equals(gattung.Konfig) {
	// 			ids.Add(Konfig{})
	// 		}

	// 		ms.Gattung[g] = ids

	// 		return
	// 	},
	// ); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	return
}