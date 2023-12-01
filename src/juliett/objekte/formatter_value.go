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
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	to_merge "github.com/friedenberg/zit/src/india/sku_fmt"
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
		"etiketten-all",
		"etiketten-implicit",
		"etiketten-expanded",
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
	logFunc schnittstellen.FuncIter[*sku.Transacted],
	cliFmt schnittstellen.StringFormatWriter[*sku.Transacted],
) schnittstellen.FuncIter[*sku.Transacted] {
	switch fv.string {
	case "etiketten-all":
		return func(tl *sku.Transacted) (err error) {
			esImp := tl.GetMetadatei().Verzeichnisse.GetExpandedEtiketten()
			esEx := tl.GetMetadatei().Verzeichnisse.GetImplicitEtiketten()
			// TODO-P3 determine if empty sets should be printed or not

			if _, err = fmt.Fprintln(
				out,
				iter.StringCommaSeparated[kennung.Etikett](esImp, esEx),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "etiketten-expanded":
		return func(tl *sku.Transacted) (err error) {
			esImp := tl.GetMetadatei().Verzeichnisse.GetExpandedEtiketten()
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

	case "etiketten-implicit":
		return func(tl *sku.Transacted) (err error) {
			esImp := tl.GetMetadatei().Verzeichnisse.GetImplicitEtiketten()
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
		return func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				out,
				iter.StringCommaSeparated[kennung.Etikett](
					tl.Metadatei.GetEtiketten(),
				),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	case "bezeichnung":
		return func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(out, tl.GetMetadatei().Bezeichnung); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "text":
		f := MakeTextFormatter(
			af,
			k,
		)

		return func(tl *sku.Transacted) (err error) {
			_, err = f.WriteStringFormat(out, tl)
			return
		}

	case "objekte":
		f := objekte_format.FormatForVersion(k.GetStoreVersion())
		o := objekte_format.Options{
			IncludeTai: true,
		}

		return func(tl *sku.Transacted) (err error) {
			if _, err = f.FormatPersistentMetadatei(out, tl, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "kennung-akte-sha":
		return func(tl *sku.Transacted) (err error) {
			errors.TodoP3("convert into an option")

			sh := tl.GetAkteSha()

			if sh.IsNull() {
				return
			}

			if _, err = fmt.Fprintf(
				out,
				"%s %s\n",
				tl.Kennung,
				sh,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "kennung":
		return func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, e.Kennung)
			return
		}

	case "sku-metadatei-sans-tai":
		return func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(
				out,
				to_merge.StringMetadateiSansTai(e),
			)
			return
		}

	case "sku-metadatei":
		return func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(
				out,
				to_merge.StringMetadatei(e),
			)
			return
		}

	case "sku":
		return func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, to_merge.String(e))
			return
		}

	case "metadatei":
		return func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintf(out, "%#v\n", e.GetMetadatei())
			return
		}

	case "debug":
		return func(e *sku.Transacted) (err error) {
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

		return func(o *sku.Transacted) (err error) {
			if err = enc.Encode(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "akte":
		return func(o *sku.Transacted) (err error) {
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
		return func(o *sku.Transacted) (err error) {
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
		return func(o *sku.Transacted) (err error) {
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
		be := to_merge.MakeFormatBestandsaufnahmePrinter(
			out,
			objekte_format.Default(),
			objekte_format.Options{},
		)

		return func(o *sku.Transacted) (err error) {
			if _, err = be.Print(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "bestandsaufnahme":
		f := to_merge.MakeFormatBestandsaufnahmePrinter(
			out,
			objekte_format.Default(),
			objekte_format.Options{IncludeTai: true},
		)

		return func(o *sku.Transacted) (err error) {
			if _, err = f.Print(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "bestandsaufnahme-verzeichnisse":
		f := to_merge.MakeFormatBestandsaufnahmePrinter(
			out,
			objekte_format.Default(),
			objekte_format.Options{
				IncludeTai:           true,
				IncludeVerzeichnisse: true,
			},
		)

		return func(o *sku.Transacted) (err error) {
			if _, err = f.Print(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "akte-sha":
		return func(o *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(out, o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	return func(e *sku.Transacted) (err error) {
		return MakeErrUnsupportedFormatterValue(fv.string, e.GetGattung())
	}
}
