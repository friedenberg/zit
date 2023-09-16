package objekte

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/sku_formats"
)

type FormatterValue struct {
	string
}

func (f FormatterValue) String() string {
	return f.string
}

func (f *FormatterValue) Set(v string) (err error) {
	v1 := strings.TrimSpace(strings.ToLower(v))

	switch v1 {
	case
		// TODO-P3 add toml
		"bestandsaufnahme",
		"bestandsaufnahme-sans-tai",
		"bestandsaufnahme-verzeichnisse",
		"objekte",
		"kennung",
		"kennung-akte-sha",
		"bezeichnung",
		"akte",
		"akte-sku-prefix",
		"metadatei",
		"akte-sha",
		"debug",
		"etiketten",
		"etiketten-implicit",
		"json",
		"log",
		"sku",
		"sku-metadatei",
		"sku-metadatei-sans-tai",
		"text",
		"text-sku-prefix",
		"sku2":
		f.string = v1

	default:
		err = MakeErrUnsupportedFormatterValue(v1, gattung.Unknown)
		return
	}

	return
}

func (fv *FormatterValue) MakeFormatterObjekte(
	out io.Writer,
	af schnittstellen.AkteReaderFactory,
	k Konfig,
	logFunc schnittstellen.FuncIter[sku.SkuLikePtr],
	cliFmt schnittstellen.StringFormatWriter[sku.SkuLikePtr],
) schnittstellen.FuncIter[sku.SkuLikePtr] {
	switch fv.string {
	case "etiketten-implicit":
		return func(tl sku.SkuLikePtr) (err error) {
			esImp := tl.GetMetadateiPtr().Verzeichnisse.GetImplicitEtiketten()
			// TODO-P3 determine if empty sets should be printed or not

			if _, err = fmt.Fprintln(
				out,
				iter.StringCommaSeparated[kennung.Etikett](esImp),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "etiketten":
		return func(tl sku.SkuLikePtr) (err error) {
			if _, err = fmt.Fprintln(
				out,
				iter.StringCommaSeparated[kennung.Etikett](
					tl.GetMetadatei().GetEtiketten(),
				),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	case "bezeichnung":
		return func(tl sku.SkuLikePtr) (err error) {
			if _, err = fmt.Fprintln(out, tl.GetMetadatei().Bezeichnung); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "text":
		fMetadateiAndAkte := metadatei.MakeTextFormatterMetadateiInlineAkte(
			af,
			nil,
		)
		fMetadateiOnly := metadatei.MakeTextFormatterMetadateiOnly(af, nil)
		fAkteOnly := metadatei.MakeTextFormatterExcludeMetadatei(af, nil)

		return func(tl sku.SkuLikePtr) (err error) {
			if gattung.Konfig.EqualsGattung(tl.GetGattung()) {
				_, err = fAkteOnly.FormatMetadatei(out, tl)
			} else if k.IsInlineTyp(tl.GetTyp()) {
				_, err = fMetadateiAndAkte.FormatMetadatei(out, tl)
			} else {
				_, err = fMetadateiOnly.FormatMetadatei(out, tl)
			}

			return
		}

	case "objekte":
		f := objekte_format.FormatForVersion(k.GetStoreVersion())
		o := objekte_format.Options{
			IncludeTai: true,
		}

		return func(tl sku.SkuLikePtr) (err error) {
			if _, err = f.FormatPersistentMetadatei(out, tl, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "kennung-akte-sha":
		return func(tl sku.SkuLikePtr) (err error) {
			errors.TodoP3("convert into an option")

			sh := tl.GetAkteSha()

			if sh.IsNull() {
				return
			}

			if _, err = fmt.Fprintf(
				out,
				"%s %s\n",
				tl.GetSkuLike().GetKennungLike(),
				sh,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "kennung":
		return func(e sku.SkuLikePtr) (err error) {
			_, err = fmt.Fprintln(out, e.GetSkuLike().GetKennungLike())
			return
		}

	case "sku-metadatei-sans-tai":
		return func(e sku.SkuLikePtr) (err error) {
			_, err = fmt.Fprintln(
				out,
				sku_formats.StringMetadateiSansTai(e.GetSkuLike()),
			)
			return
		}

	case "sku-metadatei":
		return func(e sku.SkuLikePtr) (err error) {
			_, err = fmt.Fprintln(
				out,
				sku_formats.StringMetadatei(e.GetSkuLike()),
			)
			return
		}

	case "sku":
		return func(e sku.SkuLikePtr) (err error) {
			_, err = fmt.Fprintln(out, sku_formats.String(e.GetSkuLike()))
			return
		}

	case "metadatei":
		return func(e sku.SkuLikePtr) (err error) {
			_, err = fmt.Fprintf(out, "%#v\n", e.GetMetadatei())
			return
		}

	case "debug":
		return func(e sku.SkuLikePtr) (err error) {
			_, err = fmt.Fprintf(out, "%#v\n", e)
			return
		}

	case "log":
		return logFunc

	// case "objekte":
	// 	f := Format{}

	// 	return func(o TransactedLikePtr) (err error) {
	// 		if _, err = f.Format(out, &o.Objekte); err != nil {
	// 			err = errors.Wrap(err)
	// 			return
	// 		}

	// 		return
	// 	}
	case "json":
		enc := json.NewEncoder(out)

		return func(o sku.SkuLikePtr) (err error) {
			if err = enc.Encode(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "akte":
		return func(o sku.SkuLikePtr) (err error) {
			var r sha.ReadCloser

			if r, err = af.AkteReader(o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.DeferredCloser(&err, r)

			if _, err = io.Copy(out, r); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "text-sku-prefix":
		return func(o sku.SkuLikePtr) (err error) {
			var r sha.ReadCloser

			if r, err = af.AkteReader(o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.DeferredCloser(&err, r)

			if _, err = io.Copy(out, r); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "akte-sku-prefix":
		return func(o sku.SkuLikePtr) (err error) {
			var r sha.ReadCloser

			if r, err = af.AkteReader(o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.DeferredCloser(&err, r)

			sb := &strings.Builder{}

			if _, err = cliFmt.WriteStringFormat(sb, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, err = ohio.CopyWithPrefixOnDelim('\n', sb.String(), out, r); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "bestandsaufnahme-sans-tai":
		be := sku_formats.MakeFormatBestandsaufnahmePrinter(
			out,
			objekte_format.Default(),
			objekte_format.Options{},
		)

		return func(o sku.SkuLikePtr) (err error) {
			if _, err = be.Print(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "bestandsaufnahme":
		f := sku_formats.MakeFormatBestandsaufnahmePrinter(
			out,
			objekte_format.Default(),
			objekte_format.Options{IncludeTai: true},
		)

		return func(o sku.SkuLikePtr) (err error) {
			if _, err = f.Print(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "bestandsaufnahme-verzeichnisse":
		f := sku_formats.MakeFormatBestandsaufnahmePrinter(
			out,
			objekte_format.Default(),
			objekte_format.Options{
				IncludeTai:           true,
				IncludeVerzeichnisse: true,
			},
		)

		return func(o sku.SkuLikePtr) (err error) {
			if _, err = f.Print(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "akte-sha":
		return func(o sku.SkuLikePtr) (err error) {
			if _, err = fmt.Fprintln(out, o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	return func(e sku.SkuLikePtr) (err error) {
		return MakeErrUnsupportedFormatterValue(fv.string, e.GetGattung())
	}
}
