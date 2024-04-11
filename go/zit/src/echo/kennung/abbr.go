package kennung

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/delta/sha"
)

type (
	// TODO use catgut.String
	FuncExpandString                                        func(string) (string, error)
	FuncAbbreviateString[V any, VPtr schnittstellen.Ptr[V]] func(VPtr) (string, error)

	Abbr struct {
		Sha struct {
			Expand     FuncExpandString
			Abbreviate FuncAbbreviateString[sha.Sha, *sha.Sha]
		}
		// TODO switch to Kennung2
		Hinweis abbrOne[Hinweis, *Hinweis]
	}

	abbrOne[V KennungLike[V], VPtr KennungLikePtr[V]] struct {
		Expand     FuncExpandString
		Abbreviate FuncAbbreviateString[V, VPtr]
	}
)

func DontExpandString(v string) (string, error) {
	return v, nil
}

func DontAbbreviateString[VPtr schnittstellen.Stringer](k VPtr) (string, error) {
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
	k Kennung,
) (v string, err error) {
	if ao.Abbreviate == nil {
		v = k.String()
		return
	}

	var ka1 V

	if ka1.GetGattung() != k.GetGattung() {
		err = gattung.ErrWrongType{
			ExpectedType: gattung.Must(ka1.GetGattung()),
			ActualType:   gattung.Must(k.GetGattung()),
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
	in *Kennung2,
) (kopf, schwanz int, err error) {
	if in.GetGattung() != gattung.Zettel || a.Hinweis.Abbreviate == nil {
		kopf, schwanz = in.LenKopfUndSchwanz()
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

	kopf = len(h.Kopf())
	schwanz = len(h.Schwanz())

	return
}

func (a Abbr) AbbreviateHinweisOnly(
	in *Kennung2,
) (err error) {
	if in.GetGattung() != gattung.Zettel {
		return
	}

	var getAbbr func(Kennung) (string, error)

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

	if err = in.SetWithGattung(abbr, h); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a Abbr) ExpandHinweisOnly(
	in *Kennung2,
) (err error) {
	if in.GetGattung() != gattung.Zettel || a.Hinweis.Expand == nil {
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

	if err = in.SetWithGattung(ex, h); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a Abbr) AbbreviateKennung(
	in *Kennung2,
	out *Kennung2,
) (err error) {
	var getAbbr func(Kennung) (string, error)

	switch in.GetGattung() {
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

	if err = out.SetWithGattung(abbr, in); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
