package sku

import (
	"github.com/friedenberg/zit/src/alfa/errors"
)

type (
	FuncMakeSkuLike func(string) (*Transacted, error)
)

func TryMakeSkuWithFormats(fms ...FuncMakeSkuLike) FuncMakeSkuLike {
	return func(line string) (sk *Transacted, err error) {
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

// func MakeSkuFromLineGattungFirst(line string) (sk *Transacted2, err error) {
// 	var (
// 		m  metadatei.Metadatei
// 		k  kennung.Kennung
// 		os sha.Sha
// 		g  gattung.Gattung
// 	)

// 	r := strings.NewReader(line)

// 	if _, err = format.ReadSep(
// 		' ',
// 		r,
// 		ohio.MakeLineReaderIterateStrict(
// 			g.Set,
// 			m.Tai.Set,
// 			func(v string) (err error) {
// 				if k, err = kennung.MakeWithGattung(g, v); err != nil {
// 					err = errors.Wrap(err)
// 					return
// 				}

// 				return
// 			},
// 			os.Set,
// 			m.AkteSha.Set,
// 		),
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if !g.EqualsGattung(k.GetGattung()) {
// 		err = errors.Errorf(
// 			"sku gattung does not match kennung gattung: %q (sku), %q (kennung)",
// 			g,
// 			k.GetGattung(),
// 		)

// 		err = errors.Wrapf(err, "Line: %q", line)

// 		return
// 	}

// 	return MakeSkuLike(m, k, os)
// }

// func MakeSkuFromLineTaiFirst(line string) (sk *Transacted2, err error) {
// 	var (
// 		m  metadatei.Metadatei
// 		k  kennung.Kennung
// 		os sha.Sha
// 		g  gattung.Gattung
// 	)

// 	r := strings.NewReader(line)

// 	if _, err = format.ReadSep(
// 		' ',
// 		r,
// 		ohio.MakeLineReaderIterateStrict(
// 			m.Tai.Set,
// 			g.Set,
// 			func(v string) (err error) {
// 				if k, err = kennung.MakeWithGattung(g, v); err != nil {
// 					err = errors.Wrap(err)
// 					return
// 				}

// 				return
// 			},
// 			os.Set,
// 			m.AkteSha.Set,
// 		),
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if !g.EqualsGattung(k.GetGattung()) {
// 		err = errors.Errorf(
// 			"sku gattung does not match kennung gattung: %q (sku), %q (kennung)",
// 			g,
// 			k.GetGattung(),
// 		)

// 		err = errors.Wrapf(err, "Line: %q", line)

// 		return
// 	}

// 	return MakeSkuLike(m, k, os)
// }
