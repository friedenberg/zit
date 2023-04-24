package zettel

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/persisted_metadatei_format"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/india/konfig"
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
		"formatters",
		"typ-vim-syntax-type",
		"typ",
		"typ-formatter-uti-groups",
		"hinweis-text",
		"text",
		"objekte",
		"toml",
		"action-names",
		"hinweis-akte":
		f.string = v1

	default:
		err = objekte.MakeErrUnsupportedFormatterValue(v, gattung.Zettel)
		return
	}

	return
}

func (fv *FormatterValue) FuncFormatterVerzeichnisse(
	out io.Writer,
	af schnittstellen.AkteIOFactory,
	k konfig.Compiled,
) schnittstellen.FuncIter[*Transacted] {
	return fv.FuncFormatter(
		out,
		af,
		k,
	)
}

func (fv *FormatterValue) FuncFormatter(
	out io.Writer,
	af schnittstellen.AkteIOFactory,
	k konfig.Compiled,
) schnittstellen.FuncIter[*Transacted] {
	errors.TodoP2("convert to verzeichnisse")

	switch fv.string {
	case "formatters":
		return func(o *Transacted) (err error) {
			t := k.GetApproximatedTyp(o.GetTyp())

			if !t.HasValue() {
				return
			}

			tt := t.ActualOrNil()

			lw := format.MakeLineWriter()

			for fn, f := range tt.Objekte.Akte.Formatters {
				if f.FileExtension != "" {
					lw.WriteFormat("%s %s", fn, f.FileExtension)
				} else {
					lw.WriteFormat("%s", fn)
				}
			}

			if _, err = lw.WriteTo(out); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "typ":
		return func(o *Transacted) (err error) {
			if _, err = io.WriteString(out, o.GetTyp().String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "typ-vim-syntax-type":
		return func(o *Transacted) (err error) {
			var t *typ.Transacted

			if t = k.GetApproximatedTyp(
				o.GetTyp(),
			).ApproximatedOrActual(); t == nil {
				return
			}

			if _, err = fmt.Fprintln(
				out,
				t.Objekte.Akte.VimSyntaxType,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "objekte":
		f := persisted_metadatei_format.FormatForVersion(k.GetStoreVersion())

		return func(o *Transacted) (err error) {
			if _, err = f.FormatPersistentMetadatei(out, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "text", "hinweis-text":
		f := metadatei.MakeTextFormatterMetadateiInlineAkte(
			af,
			nil,
		)

		return func(o *Transacted) (err error) {
			if fv.string == "hinweis-text" {
				if _, err = io.WriteString(
					out,
					fmt.Sprintf("= %s\n", o.Sku.Kennung),
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			if _, err = f.FormatMetadatei(
				out,
				o,
			); err != nil {
				err = errors.Wrapf(err, "Hinweis: %s", o.Sku.Kennung)

				if errors.IsNotExist(err) {
					err = nil
				} else {
					return
				}
			}

			return
		}

	case "toml":
		errors.TodoP3("limit to only zettels supporting toml")
		return func(o *Transacted) (err error) {
			if _, err = io.WriteString(
				out, fmt.Sprintf("['%s']\n", o.Sku.Kennung),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			var r sha.ReadCloser

			if r, err = af.AkteReader(o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.Deferred(&err, r.Close)

			if _, err = io.Copy(out, r); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "typ-formatter-uti-groups":
		f := MakeFormatterTypFormatterUTIGroups(k)

		return func(o *Transacted) (err error) {
			if _, err = f.Format(out, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "action-names":
		f := MakeFormatterTypActionNames(k, true)

		return func(o *Transacted) (err error) {
			if _, err = f.Format(out, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	default:
		return func(_ *Transacted) (err error) {
			err = objekte.MakeErrUnsupportedFormatterValue(fv.string, gattung.Zettel)
			return
		}
	}
}
