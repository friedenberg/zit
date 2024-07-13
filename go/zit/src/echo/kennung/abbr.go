package kennung

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type (
	// TODO use catgut.String
	FuncExpandString                                    func(string) (string, error)
	FuncAbbreviateString[V any, VPtr interfaces.Ptr[V]] func(VPtr) (string, error)

	Abbr struct {
		Sha struct {
			Expand     FuncExpandString
			Abbreviate FuncAbbreviateString[sha.Sha, *sha.Sha]
		}
		// TODO switch to Kennung2
		Hinweis abbrOne[Hinweis, *Hinweis]
	}

	abbrOne[V IdGeneric[V], VPtr IdGenericPtr[V]] struct {
		Expand     FuncExpandString
		Abbreviate FuncAbbreviateString[V, VPtr]
	}
)

func DontExpandString(v string) (string, error) {
	return v, nil
}

func DontAbbreviateString[VPtr interfaces.Stringer](k VPtr) (string, error) {
	return k.String(), nil
}

func (a Abbr) ExpanderFor(g gattung.Gattung) FuncExpandString {
	switch g {
	case gattung.Zettel:
		return a.Hinweis.Expand

	case gattung.Etikett, gattung.Typ, gattung.Kasten:
		return DontExpandString

	default:
		return nil
	}
}

func (ao abbrOne[V, VPtr]) AbbreviateKennung(
	k IdLike,
) (v string, err error) {
	if ao.Abbreviate == nil {
		v = k.String()
		return
	}

	var ka1 V

	if ka1.GetGenre() != k.GetGenre() {
		err = gattung.ErrWrongType{
			ExpectedType: gattung.Must(ka1.GetGenre()),
			ActualType:   gattung.Must(k.GetGenre()),
		}

		return
	}

	if err = VPtr(&ka1).Set(k.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if v, err = ao.Abbreviate(&ka1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a Abbr) LenKopfUndSchwanz(
	in *Id,
) (kopf, schwanz int, err error) {
	if in.GetGenre() != gattung.Zettel || a.Hinweis.Abbreviate == nil {
		kopf, schwanz = in.LenHeadAndTail()
		return
	}

	var h Hinweis

	if err = h.Set(in.String()); err != nil {
		err = nil
		return
	}

	var abbr string

	if abbr, err = a.Hinweis.AbbreviateKennung(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = h.Set(abbr); err != nil {
		err = errors.Wrap(err)
		return
	}

	kopf = len(h.GetHead())
	schwanz = len(h.GetTail())

	return
}

func (a Abbr) AbbreviateHinweisOnly(
	in *Id,
) (err error) {
	if in.GetGenre() != gattung.Zettel || in.IsVirtual() {
		return
	}

	var getAbbr func(IdLike) (string, error)

	var h Hinweis

	if err = h.Set(in.String()); err != nil {
		err = nil
		return
	}

	getAbbr = a.Hinweis.AbbreviateKennung

	var abbr string

	if abbr, err = getAbbr(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = in.SetWithGenre(abbr, h); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a Abbr) ExpandHinweisOnly(
	in *Id,
) (err error) {
	if in.GetGenre() != gattung.Zettel || a.Hinweis.Expand == nil {
		return
	}

	var h Hinweis

	if err = h.Set(in.String()); err != nil {
		err = nil
		return
	}

	var ex string

	if ex, err = a.Hinweis.Expand(h.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = in.SetWithGenre(ex, h); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a Abbr) AbbreviateKennung(
	in *Id,
	out *Id,
) (err error) {
	var getAbbr func(IdLike) (string, error)

	switch in.GetGenre() {
	case gattung.Zettel:
		getAbbr = a.Hinweis.AbbreviateKennung

	case gattung.Etikett, gattung.Typ, gattung.Kasten:
		getAbbr = DontAbbreviateString

	case gattung.Konfig:
		out = in
		return

	default:
		err = errors.Errorf("unsupported Kennung: %q, %T", in, in)
		return
	}

	var abbr string

	if abbr, err = getAbbr(in); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = out.SetWithGenre(abbr, in); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
