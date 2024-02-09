package sku_fmt

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type Json struct {
	Akte        string
	AkteSha     string
	Bezeichnung string
	Etiketten   []string
	Kennung     string
	Typ         string
}

type JsonWithUrl struct {
	Json
	TomlBookmark
}

func MakeJson(sk *sku.Transacted, s standort.Standort) (j Json, err error) {
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

	j = Json{
		Akte:        out.String(),
		AkteSha:     sk.GetAkteSha().String(),
		Bezeichnung: sk.Metadatei.Bezeichnung.String(),
		// Etiketten   []string
		Kennung: sk.Kennung.String(),
		Typ:     sk.Metadatei.Typ.String(),
	}

	return
}

func MakeJsonTomlBookmark(
	sk *sku.Transacted,
	s standort.Standort,
) (j JsonWithUrl, err error) {
	if j.Json, err = MakeJson(sk, s); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = toml.Unmarshal([]byte(j.Akte), &j.TomlBookmark); err != nil {
		err = errors.Wrapf(err, "%q", j.Akte)
		return
	}

	return
}

