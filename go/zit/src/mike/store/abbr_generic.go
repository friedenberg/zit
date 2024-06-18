package store

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type AbbrStorePresenceGeneric[V any] interface {
	Exists([3]string) error
}

type AbbrStoreGeneric[V any, VPtr schnittstellen.Ptr[V]] interface {
	AbbrStorePresenceGeneric[V]
	ExpandStringString(string) (string, error)
	ExpandString(string) (VPtr, error)
	Expand(VPtr) (VPtr, error)
	Abbreviate(VPtr) (string, error)
}

type AbbrStoreMutableGeneric[V any, VPtr schnittstellen.Ptr[V]] interface {
	Add(VPtr) error
}

type AbbrStoreCompleteGeneric[V any, VPtr schnittstellen.Ptr[V]] interface {
	AbbrStoreGeneric[V, VPtr]
	AbbrStoreMutableGeneric[V, VPtr]
}

type indexNoAbbr[
	V schnittstellen.Stringer,
	VPtr schnittstellen.SetterPtr[V],
] struct {
	AbbrStoreGeneric[V, VPtr]
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

func (ih *indexHinweis) Add(h *kennung.Hinweis) (err error) {
	ih.Kopfen.Add(h.Kopf())
	ih.Schwanzen.Add(h.Schwanz())
	return
}

func (ih *indexHinweis) Exists(parts [3]string) (err error) {
	if err = ih.readFunc(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !ih.Kopfen.ContainsExpansion(parts[0]) {
		err = collections.MakeErrNotFoundString(parts[0])
		return
	}

	if !ih.Schwanzen.ContainsExpansion(parts[2]) {
		err = collections.MakeErrNotFoundString(parts[2])
		return
	}

	return
}

func (ih *indexHinweis) ExpandStringString(in string) (out string, err error) {
	var h *kennung.Hinweis

	if h, err = ih.ExpandString(in); err != nil {
		err = errors.Wrap(err)
		return
	}

	out = h.String()

	return
}

func (ih *indexHinweis) ExpandString(s string) (h *kennung.Hinweis, err error) {
	if err = ih.readFunc(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ha *kennung.Hinweis

	if ha, err = kennung.MakeHinweis(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	if h, err = ih.Expand(ha); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ih *indexHinweis) Expand(
	hAbbr *kennung.Hinweis,
) (h *kennung.Hinweis, err error) {
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

func (ih *indexHinweis) Abbreviate(h *kennung.Hinweis) (v string, err error) {
	if err = ih.readFunc(); err != nil {
		err = errors.Wrap(err)
		return
	}

	kopf := h.Kopf()
	schwanz := h.Schwanz()

	kopf = ih.Kopfen.Abbreviate(h.Kopf())
	schwanz = ih.Schwanzen.Abbreviate(h.Schwanz())

	if kopf == "" {
		v = h.String()
		return
	}

	if schwanz == "" {
		v = h.String()
		return
	}

	v = fmt.Sprintf("%s/%s", kopf, schwanz)

	return
}

type indexNotHinweis[
	K any,
	KPtr schnittstellen.StringerSetterPtr[K],
] struct {
	readFunc  func() error
	Kennungen schnittstellen.MutableTridex
}

func (ih *indexNotHinweis[K, KPtr]) Add(k KPtr) (err error) {
	ih.Kennungen.Add(k.String())
	return
}

func (ih *indexNotHinweis[K, KPtr]) Exists(parts [3]string) (err error) {
	if err = ih.readFunc(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !ih.Kennungen.ContainsExpansion(parts[2]) {
		err = collections.MakeErrNotFoundString(parts[2])
		return
	}

	return
}

func (ih *indexNotHinweis[K, KPtr]) ExpandStringString(
	in string,
) (out string, err error) {
	var k KPtr

	if k, err = ih.ExpandString(in); err != nil {
		err = errors.Wrap(err)
		return
	}

	out = k.String()

	return
}

func (ih *indexNotHinweis[K, KPtr]) ExpandString(s string) (k KPtr, err error) {
	if err = ih.readFunc(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var k1 K
	k = &k1

	if err = k.Set(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	if k, err = ih.Expand(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ih *indexNotHinweis[K, KPtr]) Expand(
	abbr KPtr,
) (exp KPtr, err error) {
	if err = ih.readFunc(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ex := ih.Kennungen.Expand(abbr.String())

	if ex == "" {
		// TODO-P4 should try to use the expansion if possible
		ex = abbr.String()
	}

	var k K
	exp = &k

	if err = exp.Set(ex); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ih *indexNotHinweis[K, KPtr]) Abbreviate(k KPtr) (v string, err error) {
	if err = ih.readFunc(); err != nil {
		err = errors.Wrap(err)
		return
	}

	v = ih.Kennungen.Abbreviate(k.String())

	return
}
