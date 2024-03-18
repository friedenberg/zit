package objekte_formatter

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/src/hotel/matcher_proto"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/src/kilo/typ"
	"code.linenisgreat.com/zit/src/kilo/zettel"
)

type funcFormat = schnittstellen.FuncIter[*sku.Transacted]

type FormatterFactory interface {
	MakeFormatterObjekte(
		out io.Writer,
		af schnittstellen.AkteIOFactory,
		k konfig.Compiled,
		logFunc schnittstellen.FuncIter[*sku.Transacted],
	) funcFormat
}

type formatter struct {
	formatters map[gattung.Gattung]funcFormat
}

type Formatter interface {
	MakeFormatFunc() funcFormat
}

func MakeFormatter(
	ms matcher_proto.QueryGroup,
	v string,
	out io.Writer,
	af schnittstellen.AkteIOFactory,
	k *konfig.Compiled,
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

		f.formatters[gattung.Zettel] = zvf
	}

	if _, ok := ms.Get(gattung.Typ); ok {
		var tv typ.FormatterValue

		if err = tv.Set(v); err != nil {
			err = errors.Normal(err)
			return
		}

		f.formatters[gattung.Typ] = tv.FuncFormatter(
			out,
			af,
			tagp,
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
	return func(tl *sku.Transacted) (err error) {
		g := gattung.Must(tl.GetGattung())

		if f1, ok := f.formatters[g]; ok {
			return f1(tl)
		}

		return gattung.MakeErrUnsupportedGattung(g)
	}
}
