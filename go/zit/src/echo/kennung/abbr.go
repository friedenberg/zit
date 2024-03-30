package kennung

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/charlie/sha"
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

	switch ka := k.(type) {
	case VPtr:
		if err = VPtr(&ka1).Set(ka.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

	case V:
		ka1 = ka

	default:
		err = errors.Errorf("expected kennung type %T but got %T", ka, k)
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
	in Kennung,
) (out Kennung, err error) {
	var getAbbr func(Kennung) (string, error)

	switch in.(type) {
	case Hinweis, *Hinweis:
		getAbbr = a.Hinweis.AbbreviateKennung

	case Etikett, *Etikett:
		getAbbr = a.Etikett.AbbreviateKennung

	case Typ, *Typ:
		getAbbr = a.Typ.AbbreviateKennung

	case Kasten, *Kasten:
		getAbbr = a.Kasten.AbbreviateKennung

	case Konfig, *Konfig:
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

	outPtr := &Kennung2{}

	if err = outPtr.SetWithGattung(abbr, in); err != nil {
		err = errors.Wrap(err)
		return
	}

	out = outPtr

	return
}
