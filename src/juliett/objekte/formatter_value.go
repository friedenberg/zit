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
	"github.com/friedenberg/zit/src/golf/ennui"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/sku_fmt"
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
		"akte",
		"akte-sha",
		"akte-sku-prefix",
		"bestandsaufnahme",
		"bestandsaufnahme-sans-tai",
		"bestandsaufnahme-sha",
		"bestandsaufnahme-shas",
		"bestandsaufnahme-verzeichnisse",
		"bezeichnung",
		"debug",
		"etiketten",
		"etiketten-all",
		"etiketten-expanded",
		"etiketten-implicit",
		"json",
		"kennung",
		"kennung-akte-sha",
		"kennung-sha",
		"kennung-tai",
		"log",
		"metadatei",
		"metadatei-plus-mutter",
		"mutter",
		"mutter-sha",
		"objekte",
		"sha",
		"sku",
		"sku-metadatei",
		"sku-metadatei-sans-tai",
		"sku2",
		"tai",
		"text",
		"text-sku-prefix",
		"typ":
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
	enn ennui.Ennui,
	rob func(*sha.Sha) (*sku.Transacted, error),
) schnittstellen.FuncIter[*sku.Transacted] {
	switch fv.string {
	case "sha":
		return func(tl *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, &tl.Metadatei.Sha)
			return
		}

	case "etiketten-all":
		return func(tl *sku.Transacted) (err error) {
			for _, es := range tl.Metadatei.Verzeichnisse.Etiketten {
				if _, err = fmt.Fprintln(out, tl.GetKennung(), "->", es); err != nil {
					err = errors.Wrap(err)
					return
				}
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
			Tai: true,
		}

		return func(tl *sku.Transacted) (err error) {
			if _, err = f.FormatPersistentMetadatei(out, tl, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "kennung-sha":
		return func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintf(
				out,
				"%s@%s\n",
				&tl.Kennung,
				tl.GetObjekteSha(),
			); err != nil {
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
				&tl.Kennung,
				sh,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "kennung":
		return func(e *sku.Transacted) (err error) {
			_, err = e.Kennung.WriteTo(out)
			return
		}

	case "kennung-tai":
		return func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, e.StringKennungTai())
			return
		}

	case "sku-metadatei-sans-tai":
		return func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(
				out,
				sku_fmt.StringMetadateiSansTai(e),
			)
			return
		}

	case "sku-metadatei":
		return func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(
				out,
				sku_fmt.StringMetadatei(e),
			)
			return
		}

	case "sku":
		return func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, sku_fmt.String(e))
			return
		}

	case "metadatei":
		f, err := objekte_format.FormatForKeyError("Metadatei")
		errors.PanicIfError(err)

		return func(e *sku.Transacted) (err error) {
			_, err = f.WriteMetadateiTo(out, e)
			return
		}

	case "metadatei-plus-mutter":
		f, err := objekte_format.FormatForKeyError("MetadateiMutter")
		errors.PanicIfError(err)

		return func(e *sku.Transacted) (err error) {
			_, err = f.WriteMetadateiTo(out, e)
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

	case "tai":
		return func(o *sku.Transacted) (err error) {
			fmt.Fprintln(out, o.GetTai())
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
		be := sku_fmt.MakeFormatBestandsaufnahmePrinter(
			out,
			objekte_format.Default(),
			objekte_format.Options{ExcludeMutter: true},
		)

		return func(o *sku.Transacted) (err error) {
			if _, err = be.Print(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "mutter-sha":
		return func(z *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, &z.Metadatei.Mutter)
			return
		}

	case "mutter":
		return func(z *sku.Transacted) (err error) {
			if z.Metadatei.Mutter.IsNull() {
				return
			}

			if z, err = rob(&z.GetMetadatei().Mutter); err != nil {
				fmt.Fprintln(out, err)
				err = nil
				return
			}

			return logFunc(z)
		}

	case "bestandsaufnahme-sha":
		return func(z *sku.Transacted) (err error) {
			// var loc ennui.Loc

			// if loc, err = enn.ReadOne("MetadateiMutter", z.GetMetadatei()); err != nil {
			// 	err = errors.Wrapf(err, "Kennung: %s", &z.Kennung)
			// 	return
			// }

			// fmt.Fprintf(out, "%s\n", loc)

			return
		}

	case "bestandsaufnahme-shas":
		return func(z *sku.Transacted) (err error) {
			// if err = enn.ReadAll(z.GetMetadatei(), &locs); err != nil {
			// 	err = errors.Wrapf(err, "Kennung: %s", &z.Kennung)
			// 	return
			// }

			// for _, loc := range locs {
			// 	fmt.Fprintf(out, "%d %d\n", loc.Page, loc.Offset)
			// }

			// locs = locs[:0]

			return
		}

	case "bestandsaufnahme":
		f := sku_fmt.MakeFormatBestandsaufnahmePrinter(
			out,
			objekte_format.Default(),
			objekte_format.Options{Tai: true},
		)

		return func(o *sku.Transacted) (err error) {
			if _, err = f.Print(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "bestandsaufnahme-verzeichnisse":
		f := sku_fmt.MakeFormatBestandsaufnahmePrinter(
			out,
			objekte_format.Default(),
			objekte_format.Options{
				Tai:           true,
				Verzeichnisse: true,
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

	case "typ":
		return func(o *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(out, o.GetTyp().String()); err != nil {
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
