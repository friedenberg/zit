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
		// TODO-P3 add text, objekte, toml, json
		"kennung",
		"akte",
		"akte-sha",
		"debug",
		"etiketten",
		"json",
		"log",
		"sku",
		"sku-transacted",
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
	k schnittstellen.Konfig,
	logFunc schnittstellen.FuncIter[TransactedLike],
) schnittstellen.FuncIter[TransactedLike] {
	switch fv.string {
	case "kennung":
		return func(e TransactedLike) (err error) {
			_, err = fmt.Fprintln(out, e.GetDataIdentity().GetId())
			return
		}

	case "sku-transacted":
		return func(e TransactedLike) (err error) {
			_, err = fmt.Fprintln(out, e.GetSku())
			return
		}

	case "sku":
		return func(e TransactedLike) (err error) {
			_, err = fmt.Fprintln(out, e.GetSku().String())
			return
		}

	case "sku2":
		return func(e TransactedLike) (err error) {
			_, err = fmt.Fprintln(out, e.GetSku2().String())
			return
		}

	case "debug":
		return func(e TransactedLike) (err error) {
			_, err = fmt.Fprintf(out, "%#v\n", e)
			return
		}

	case "log":
		return logFunc

	// case "objekte":
	// 	f := Format{}

	// 	return func(o TransactedLike) (err error) {
	// 		if _, err = f.Format(out, &o.Objekte); err != nil {
	// 			err = errors.Wrap(err)
	// 			return
	// 		}

	// 		return
	// 	}
	case "json":
		enc := json.NewEncoder(out)

		return func(o TransactedLike) (err error) {
			if err = enc.Encode(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "akte":
		return func(o TransactedLike) (err error) {
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

	case "akte-sha":
		return func(o TransactedLike) (err error) {
			if _, err = fmt.Fprintln(out, o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	return func(e TransactedLike) (err error) {
		return MakeErrUnsupportedFormatterValue(fv.string, e.GetGattung())
	}
}
