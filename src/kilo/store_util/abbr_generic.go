package store_util

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
)

type AbbrStorePresenceGeneric[V any] interface {
	Exists(V) error
}

type AbbrStoreGeneric[V any] interface {
	AbbrStorePresenceGeneric[V]
	ExpandStringString(string) (string, error)
	ExpandString(string) (V, error)
	Expand(V) (V, error)
	Abbreviate(V) (string, error)
}

type AbbrStoreMutableGeneric[V any] interface {
	Add(V) error
}

type AbbrStoreCompleteGeneric[V any] interface {
	AbbrStoreGeneric[V]
	AbbrStoreMutableGeneric[V]
}

type indexNoAbbr[
	V schnittstellen.Stringer,
	VPtr schnittstellen.SetterPtr[V],
] struct {
	AbbrStoreGeneric[V]
}

func (ih indexNoAbbr[V, VPtr]) Abbreviate(h V) (v string, err error) {
	v = h.String()
	return
}

type indexHinweis struct {
	readFunc  func() error
	Kopfen    schnittstellen.MutableTridex
	Schwanzen schnittstellen.MutableTridex
}

func (ih *indexHinweis) Add(h kennung.Hinweis) (err error) {
	ih.Kopfen.Add(h.Kopf())
	ih.Schwanzen.Add(h.Schwanz())
	return
}

func (ih *indexHinweis) Exists(h kennung.Hinweis) (err error) {
	if err = ih.readFunc(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !ih.Kopfen.ContainsExpansion(h.Kopf()) {
		err = objekte_store.ErrNotFound{Id: h}
		return
	}

	if !ih.Schwanzen.ContainsExpansion(h.Schwanz()) {
		err = objekte_store.ErrNotFound{Id: h}
		return
	}

	return
}

func (ih *indexHinweis) ExpandStringString(in string) (out string, err error) {
	var h kennung.Hinweis

	if h, err = ih.ExpandString(in); err != nil {
		err = errors.Wrap(err)
		return
	}

	out = h.String()

	return
}

func (ih *indexHinweis) ExpandString(s string) (h kennung.Hinweis, err error) {
	if err = ih.readFunc(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ha kennung.Hinweis

	if ha, err = kennung.MakeHinweis(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	return ih.Expand(ha)
}

func (ih *indexHinweis) Expand(
	hAbbr kennung.Hinweis,
) (h kennung.Hinweis, err error) {
	if err = ih.readFunc(); err != nil {
		err = errors.Wrap(err)
		return
	}

	kopf := ih.Kopfen.Expand(hAbbr.Kopf())
	schwanz := ih.Schwanzen.Expand(hAbbr.Schwanz())

	if h, err = kennung.MakeHinweisKopfUndSchwanz(kopf, schwanz); err != nil {
		err = errors.Wrapf(err, "{Abbreviated: '%s'}", hAbbr)
		return
	}

	return
}

func (ih *indexHinweis) Abbreviate(h kennung.Hinweis) (v string, err error) {
	if err = ih.readFunc(); err != nil {
		err = errors.Wrap(err)
		return
	}

	kopf := h.Kopf()
	schwanz := h.Schwanz()

	kopf = ih.Kopfen.Abbreviate(h.Kopf())
	schwanz = ih.Schwanzen.Abbreviate(h.Schwanz())

	if kopf == "" {
		err = errors.Errorf("abbreviated kopf would be empty for %s", h)
		return
	}

	if schwanz == "" {
		err = errors.Errorf("abbreviated schwanz would be empty for %s", h)
		return
	}

	v = fmt.Sprintf("%s/%s", kopf, schwanz)

	return
}

type indexNotHinweis[
	K schnittstellen.Stringer,
	KPtr schnittstellen.SetterPtr[K],
] struct {
	readFunc  func() error
	Kennungen schnittstellen.MutableTridex
}

func (ih *indexNotHinweis[K, KPtr]) Add(k K) (err error) {
	ih.Kennungen.Add(k.String())
	return
}

func (ih *indexNotHinweis[K, KPtr]) Exists(k K) (err error) {
	if err = ih.readFunc(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !ih.Kennungen.ContainsExpansion(k.String()) {
		err = objekte_store.ErrNotFound{Id: k}
		return
	}

	return
}

func (ih *indexNotHinweis[K, KPtr]) ExpandStringString(
	in string,
) (out string, err error) {
	var k K

	if k, err = ih.ExpandString(in); err != nil {
		err = errors.Wrap(err)
		return
	}

	out = k.String()

	return
}

func (ih *indexNotHinweis[K, KPtr]) ExpandString(s string) (k K, err error) {
	if err = ih.readFunc(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = KPtr(&k).Set(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	return ih.Expand(k)
}

func (ih *indexNotHinweis[K, KPtr]) Expand(
	abbr K,
) (exp K, err error) {
	if err = ih.readFunc(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ex := ih.Kennungen.Expand(abbr.String())

	if ex == "" {
		// TODO-P4 should try to use the expansion if possible
		ex = abbr.String()
	}

	if err = KPtr(&exp).Set(ex); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ih *indexNotHinweis[K, KPtr]) Abbreviate(k K) (v string, err error) {
	if err = ih.readFunc(); err != nil {
		err = errors.Wrap(err)
		return
	}

	v = ih.Kennungen.Abbreviate(k.String())

	return
}
