package objekte_formatter

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/konfig"
	"github.com/friedenberg/zit/src/kilo/typ"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

type funcFormat = schnittstellen.FuncIter[sku.SkuLikePtr]

type FormatterFactory interface {
	MakeFormatterObjekte(
		out io.Writer,
		af schnittstellen.AkteIOFactory,
		k konfig.Compiled,
		logFunc schnittstellen.FuncIter[sku.SkuLikePtr],
	) funcFormat
}

type formatter struct {
	formatters map[gattung.Gattung]funcFormat
}

func makeFuncFormatter[T kennung.Matchable](
	f schnittstellen.FuncIter[T],
) funcFormat {
	return func(e sku.SkuLikePtr) (err error) {
		if e1, ok := e.(T); ok {
			return f(e1)
		}

		var e1 T
		return errors.Errorf("could not convert %T into %T", e, e1)
	}
}

type Formatter interface {
	MakeFormatFunc() funcFormat
}

func MakeFormatter(
	ms kennung.MetaSet,
	v string,
	out io.Writer,
	af schnittstellen.AkteIOFactory,
	k konfig.Compiled,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) (fo Formatter, err error) {
	f := formatter{
		formatters: make(map[gattung.Gattung]funcFormat),
	}

	if _, ok := ms.Get(gattung.Zettel); ok {
		var zv zettel.FormatterValue

		if err = zv.Set(v); err != nil {
			err = errors.Normal(err)
			return
		}

		zvf := zv.FuncFormatterVerzeichnisse(
			out,
			af,
			k,
			tagp,
		)

		f.formatters[gattung.Zettel] = makeFuncFormatter(zvf)
	}

	if _, ok := ms.Get(gattung.Typ); ok {
		var tv typ.FormatterValue

		if err = tv.Set(v); err != nil {
			err = errors.Normal(err)
			return
		}

		f.formatters[gattung.Typ] = makeFuncFormatter(
			tv.FuncFormatter(
				out,
				af,
				tagp,
			),
		)
	}

	// if _, ok := ms.Get(gattung.Etikett); ok {
	// 	var ev etikett.FormatterValue

	// 	if err = ev.Set(v); err != nil {
	// 		err = errors.Normal(err)
	// 		return
	// 	}

	// 	f.formatters[gattung.Etikett] = makeFuncFormatter(
	// 		ev.FuncFormatter(
	// 			out,
	// 			af,
	// 		),
	// 	)
	// }

	// if _, ok := ms.Get(gattung.Kasten); ok {
	// 	var kv kasten.FormatterValue

	// 	if err = kv.Set(v); err != nil {
	// 		err = errors.Normal(err)
	// 		return
	// 	}

	// 	f.formatters[gattung.Kasten] = makeFuncFormatter(
	// 		kv.FuncFormatter(
	// 			out,
	// 			af,
	// 		),
	// 	)
	// }

	// if _, ok := ms.Get(gattung.Konfig); ok {
	// 	var kv erworben.FormatterValue

	// 	if err = kv.Set(v); err != nil {
	// 		err = errors.Normal(err)
	// 		return
	// 	}

	// 	f.formatters[gattung.Konfig] = makeFuncFormatter(
	// 		kv.FuncFormatter(
	// 			out,
	// 			af,
	// 		),
	// 	)
	// }

	fo = f

	return
}

func (f formatter) MakeFormatFunc() funcFormat {
	return func(tl sku.SkuLikePtr) (err error) {
		g := gattung.Must(tl.GetGattung())

		if f1, ok := f.formatters[g]; ok {
			return f1(tl)
		}

		return gattung.MakeErrUnsupportedGattung(g)
	}
}
