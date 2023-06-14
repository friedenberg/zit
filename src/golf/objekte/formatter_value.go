package objekte

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte_format"
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
		"objekte",
		"kennung",
		"kennung-akte-sha",
		"bezeichnung",
		"akte",
		"metadatei",
		"akte-sha",
		"debug",
		"etiketten",
		"etiketten-implicit",
		"json",
		"log",
		"sku",
		"text",
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
	logFunc schnittstellen.FuncIter[TransactedLikePtr],
) schnittstellen.FuncIter[TransactedLikePtr] {
	switch fv.string {
	case "etiketten-implicit":
		return func(tl TransactedLikePtr) (err error) {
			ets := tl.GetMetadatei().GetEtiketten().MutableClone()

			ets.Each(
				func(e kennung.Etikett) (err error) {
					impl := k.GetImplicitEtiketten(e)
					return impl.Each(ets.Add)
				},
			)

			if _, err = fmt.Fprintln(
				out,
				collections.StringCommaSeparated[kennung.Etikett](
					ets,
				),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "etiketten":
		return func(tl TransactedLikePtr) (err error) {
			if _, err = fmt.Fprintln(
				out,
				collections.StringCommaSeparated[kennung.Etikett](
					tl.GetMetadatei().GetEtiketten(),
				),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	case "bezeichnung":
		return func(tl TransactedLikePtr) (err error) {
			if _, err = fmt.Fprintln(out, tl.GetMetadatei().Bezeichnung); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "text":
		fInlineAkte := metadatei.MakeTextFormatterMetadateiInlineAkte(af, nil)
		fOmitMetadatei := metadatei.MakeTextFormatterExcludeMetadatei(af, nil)
		// f := MakeSavedAkteFormatter(af)

		return func(tl TransactedLikePtr) (err error) {
			if tl.GetGattung() == gattung.Zettel {
				_, err = fInlineAkte.FormatMetadatei(out, tl)
			} else {
				_, err = fOmitMetadatei.FormatMetadatei(out, tl)
			}

			return
		}

	case "objekte":
		f := objekte_format.FormatForVersion(k.GetStoreVersion())

		return func(tl TransactedLikePtr) (err error) {
			if _, err = f.FormatPersistentMetadatei(out, tl); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "kennung-akte-akte":
		return func(tl TransactedLikePtr) (err error) {
			errors.TodoP3("convert into an option")

			sh := tl.GetAkteSha()

			if sh.IsNull() {
				return
			}

			if _, err = fmt.Fprintf(
				out,
				"%s %s\n",
				tl.GetDataIdentity().GetId(),
				sh,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "kennung":
		return func(e TransactedLikePtr) (err error) {
			_, err = fmt.Fprintln(out, e.GetDataIdentity().GetId())
			return
		}

	case "sku":
		return func(e TransactedLikePtr) (err error) {
			_, err = fmt.Fprintln(out, e.GetSku().String())
			return
		}

	case "metadatei":
		return func(e TransactedLikePtr) (err error) {
			_, err = fmt.Fprintf(out, "%#v\n", e.GetMetadatei())
			return
		}

	case "debug":
		return func(e TransactedLikePtr) (err error) {
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

		return func(o TransactedLikePtr) (err error) {
			if err = enc.Encode(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "akte":
		return func(o TransactedLikePtr) (err error) {
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

	case "bestandsaufnahme":
		f := MakeFormatBestandsaufnahme(
			out,
			objekte_format.BestandsaufnahmeFormat(),
		)

		return func(o TransactedLikePtr) (err error) {
			if _, err = f.PrintOne(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "akte-sha":
		return func(o TransactedLikePtr) (err error) {
			if _, err = fmt.Fprintln(out, o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	return func(e TransactedLikePtr) (err error) {
		return MakeErrUnsupportedFormatterValue(fv.string, e.GetGattung())
	}
}
