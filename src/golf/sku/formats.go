package sku

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/echo/format"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type FuncMakeSkuLike func(string) (SkuLike, error)

func TryMakeSkuWithFormats(fms ...FuncMakeSkuLike) FuncMakeSkuLike {
	return func(line string) (sk SkuLike, err error) {
		em := errors.MakeMulti()

		for _, f := range fms {
			if sk, err = f(line); err == nil {
				return
			}

			em.Add(err)
		}

		return nil, em
	}
}

func MakeSkuFromLineGattungFirst(line string) (sk SkuLike, err error) {
	var (
		m  metadatei.Metadatei
		k  kennung.Kennung
		os sha.Sha
		g  gattung.Gattung
	)

	r := strings.NewReader(line)

	if _, err = format.ReadSep(
		' ',
		r,
		ohio.MakeLineReaderIterateStrict(
			g.Set,
			m.Tai.Set,
			func(v string) (err error) {
				if k, err = kennung.MakeWithGattung(g, v); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
			os.Set,
			m.AkteSha.Set,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !g.EqualsGattung(k.GetGattung()) {
		err = errors.Errorf(
			"sku gattung does not match kennung gattung: %q (sku), %q (kennung)",
			g,
			k.GetGattung(),
		)

		err = errors.Wrapf(err, "Line: %q", line)

		return
	}

	return MakeSkuLike(m, k, os)
}

func MakeSkuFromLineTaiFirst(line string) (sk SkuLike, err error) {
	var (
		m  metadatei.Metadatei
		k  kennung.Kennung
		os sha.Sha
		g  gattung.Gattung
	)

	r := strings.NewReader(line)

	if _, err = format.ReadSep(
		' ',
		r,
		ohio.MakeLineReaderIterateStrict(
			m.Tai.Set,
			g.Set,
			func(v string) (err error) {
				if k, err = kennung.MakeWithGattung(g, v); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
			os.Set,
			m.AkteSha.Set,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !g.EqualsGattung(k.GetGattung()) {
		err = errors.Errorf(
			"sku gattung does not match kennung gattung: %q (sku), %q (kennung)",
			g,
			k.GetGattung(),
		)

		err = errors.Wrapf(err, "Line: %q", line)

		return
	}

	return MakeSkuLike(m, k, os)
}

func MakeOldSkuFromTaiAndLine(
	t kennung.Tai,
	line string,
) (out1 SkuLikePtr, err error) {
	fields := strings.Fields(line)
	var g gattung.Gattung

	if err = g.Set(fields[0]); err != nil {
		err = errors.Wrapf(err, "failed to set type: %s", fields[0])
		return
	}

	switch g {
	case gattung.Zettel:
		out := &Transacted[kennung.Hinweis, *kennung.Hinweis]{}
		out1 = out
		err = setTimeAndFields(out, t, fields[1:]...)

	case gattung.Typ:
		out := &Transacted[kennung.Typ, *kennung.Typ]{}
		out1 = out
		err = setTimeAndFields(out, t, fields[1:]...)

	case gattung.Etikett:
		out := &Transacted[kennung.Etikett, *kennung.Etikett]{}
		out1 = out
		err = setTimeAndFields(out, t, fields[1:]...)

	case gattung.Konfig:
		out := &Transacted[kennung.Konfig, *kennung.Konfig]{}
		out1 = out
		return

	default:
		err = errors.Errorf("unsupported gattung: %s", g)
		return
	}

	return
}

func setTimeAndFields[
	K kennung.KennungLike[K], KPtr kennung.KennungLikePtr[K],
](
	o *Transacted[K, KPtr],
	t kennung.Tai,
	vs ...string,
) (err error) {
	o.GetMetadateiPtr().Tai = t

	if len(vs) != 4 {
		err = errors.Errorf("expected 4 elements but got %d", len(vs))
		return
	}

	// Mutter[0] used to be here

	vs = vs[1:]

	// Mutter[1] used to be here

	vs = vs[1:]

	if err = KPtr(&o.Kennung).Set(vs[0]); err != nil {
		err = errors.Wrapf(err, "failed to set id: %s", vs[1])
		return
	}

	vs = vs[1:]

	if err = o.ObjekteSha.Set(vs[0]); err != nil {
		err = errors.Wrapf(err, "failed to set sha: %s", vs[2])
		return
	}

	return
}
