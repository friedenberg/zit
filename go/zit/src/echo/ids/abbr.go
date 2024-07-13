package ids

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
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
		Hinweis abbrOne[ZettelId, *ZettelId]
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

func (a Abbr) ExpanderFor(g genres.Genre) FuncExpandString {
	switch g {
	case genres.Zettel:
		return a.Hinweis.Expand

	case genres.Tag, genres.Type, genres.Repo:
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
		err = genres.ErrWrongType{
			ExpectedType: genres.Must(ka1.GetGenre()),
			ActualType:   genres.Must(k.GetGenre()),
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
	in *ObjectId,
) (kopf, schwanz int, err error) {
	if in.GetGenre() != genres.Zettel || a.Hinweis.Abbreviate == nil {
		kopf, schwanz = in.LenHeadAndTail()
		return
	}

	var h ZettelId

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
	in *ObjectId,
) (err error) {
	if in.GetGenre() != genres.Zettel || in.IsVirtual() {
		return
	}

	var getAbbr func(IdLike) (string, error)

	var h ZettelId

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
	in *ObjectId,
) (err error) {
	if in.GetGenre() != genres.Zettel || a.Hinweis.Expand == nil {
		return
	}

	var h ZettelId

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
	in *ObjectId,
	out *ObjectId,
) (err error) {
	var getAbbr func(IdLike) (string, error)

	switch in.GetGenre() {
	case genres.Zettel:
		getAbbr = a.Hinweis.AbbreviateKennung

	case genres.Tag, genres.Type, genres.Repo:
		getAbbr = DontAbbreviateString

	case genres.Config:
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
