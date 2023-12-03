package kennung

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type Abbr struct {
	Sha struct {
		Expand     func(string) (string, error)
		Abbreviate func(*sha.Sha) (string, error)
	}
	Etikett abbrOne[Etikett, *Etikett]
	Typ     abbrOne[Typ, *Typ]
	Hinweis abbrOne[Hinweis, *Hinweis]
	Kasten  abbrOne[Kasten, *Kasten]
}

type abbrOne[V KennungLike[V], VPtr KennungLikePtr[V]] struct {
	Expand     func(string) (string, error)
	Abbreviate func(VPtr) (string, error)
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
