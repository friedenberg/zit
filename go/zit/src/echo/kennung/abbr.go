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
		Etikett abbrOne[Etikett, *Etikett]
		Typ     abbrOne[Typ, *Typ]
		Hinweis abbrOne[Hinweis, *Hinweis]
		Kasten  abbrOne[Kasten, *Kasten]
	}

	abbrOne[V KennungLike[V], VPtr KennungLikePtr[V]] struct {
		Expand     FuncExpandString
		Abbreviate FuncAbbreviateString[V, VPtr]
	}
)

func (a Abbr) ExpanderFor(g gattung.Gattung) FuncExpandString {
	switch g {
	case gattung.Zettel:
		return a.Hinweis.Expand

	case gattung.Etikett:
		return a.Etikett.Expand

	case gattung.Typ:
		return a.Typ.Expand

	case gattung.Kasten:
		return a.Kasten.Expand

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

	case gattung.Etikett:
		getAbbr = a.Etikett.AbbreviateKennung

	case gattung.Typ:
		getAbbr = a.Typ.AbbreviateKennung

	case gattung.Kasten:
		getAbbr = a.Kasten.AbbreviateKennung

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
