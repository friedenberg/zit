package sku_fmt

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/standort"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Json struct {
	Akte        string   `json:"akte"`
	AkteSha     string   `json:"akte-sha"`
	Bezeichnung string   `json:"bezeichnung"`
	Etiketten   []string `json:"etiketten"`
	Kennung     string   `json:"kennung"`
	Typ         string   `json:"typ"`
	Tai         string   `json:"tai"`
}

type JsonWithUrl struct {
	Json
	TomlBookmark
}

func (j *Json) FromTransacted(sk *sku.Transacted, s standort.Standort) (err error) {
	var r sha.ReadCloser

	if r, err = s.AkteReader(sk.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, r)

	var out strings.Builder

	if _, err = io.Copy(&out, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	j.Akte = out.String()
	j.AkteSha = sk.GetAkteSha().String()
	j.Bezeichnung = sk.Metadatei.Bezeichnung.String()
	j.Etiketten = iter.Strings(sk.GetEtiketten())
	j.Kennung = sk.Kennung.String()
	j.Tai = sk.Metadatei.Tai.String()
	j.Typ = sk.Metadatei.Typ.String()

	return
}

func (j *Json) ToTransacted(sk *sku.Transacted, s standort.Standort) (err error) {
	var w sha.WriteCloser

	if w, err = s.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w)

	if _, err = io.Copy(w, strings.NewReader(j.Akte)); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO-P1 support states of akte vs akte sha
	sk.SetAkteSha(w.GetShaLike())

	// if err = sk.Metadatei.Tai.Set(j.Tai); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	if err = sk.Kennung.Set(j.Kennung); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = sk.Metadatei.Typ.Set(j.Typ); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = sk.Metadatei.Bezeichnung.Set(j.Bezeichnung); err != nil {
		err = errors.Wrap(err)
		return
	}

	var es kennung.EtikettSet

	if es, err = kennung.MakeEtikettSetStrings(j.Etiketten...); err != nil {
		err = errors.Wrap(err)
		return
	}

	sk.Metadatei.SetEtiketten(es)
	sk.Metadatei.GenerateExpandedEtiketten()

	return
}

func MakeJsonTomlBookmark(
	sk *sku.Transacted,
	s standort.Standort,
) (j JsonWithUrl, err error) {
	if err = j.FromTransacted(sk, s); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = toml.Unmarshal([]byte(j.Akte), &j.TomlBookmark); err != nil {
		err = errors.Wrapf(err, "%q", j.Akte)
		return
	}

	return
}
