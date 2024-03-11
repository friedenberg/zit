package sku_fmt

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/standort"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
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

func (j *Json) FromStringAndMetadatei(
	k string,
	m *metadatei.Metadatei,
	s standort.Standort,
) (err error) {
	var r sha.ReadCloser

	if r, err = s.AkteReader(&m.Akte); err != nil {
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
	j.AkteSha = m.Akte.String()
	j.Bezeichnung = m.Bezeichnung.String()
	j.Etiketten = iter.Strings(m.GetEtiketten())
	j.Kennung = k
	j.Tai = m.Tai.String()
	j.Typ = m.Typ.String()

	return
}

func (j *Json) FromTransacted(sk *sku.Transacted, s standort.Standort) (err error) {
	return j.FromStringAndMetadatei(sk.Kennung.String(), sk.GetMetadatei(), s)
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
