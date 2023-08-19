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
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/hotel/transacted"
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
	tagp schnittstellen.AkteGetterPutter[*typ.Akte],
) schnittstellen.FuncIter[*transacted.Zettel] {
	return fv.FuncFormatter(
		out,
		af,
		k,
		tagp,
	)
}

func (fv *FormatterValue) FuncFormatter(
	out io.Writer,
	af schnittstellen.AkteIOFactory,
	k konfig.Compiled,
	tagp schnittstellen.AkteGetterPutter[*typ.Akte],
) schnittstellen.FuncIter[*transacted.Zettel] {
	errors.TodoP2("convert to verzeichnisse")

	switch fv.string {
	case "formatters":
		return func(o *transacted.Zettel) (err error) {
			t := k.GetApproximatedTyp(o.GetTyp())

			if !t.HasValue() {
				return
			}

			tt := t.ActualOrNil()

			var ta *typ.Akte

			if ta, err = tagp.GetAkte(tt.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer tagp.PutAkte(ta)

			lw := format.MakeLineWriter()

			for fn, f := range ta.Formatters {
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
		return func(o *transacted.Zettel) (err error) {
			if _, err = io.WriteString(out, o.GetTyp().String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "typ-vim-syntax-type":
		return func(o *transacted.Zettel) (err error) {
			var t *transacted.Typ

			if t = k.GetApproximatedTyp(
				o.GetTyp(),
			).ApproximatedOrActual(); t == nil {
				return
			}

			var ta *typ.Akte

			if ta, err = tagp.GetAkte(t.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer tagp.PutAkte(ta)

			if _, err = fmt.Fprintln(
				out,
				ta.VimSyntaxType,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "objekte":
		f := objekte_format.FormatForVersion(k.GetStoreVersion())

		return func(o *transacted.Zettel) (err error) {
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

		return func(o *transacted.Zettel) (err error) {
			if fv.string == "hinweis-text" {
				if _, err = io.WriteString(
					out,
					fmt.Sprintf("= %s\n", o.GetKennung()),
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			if _, err = f.FormatMetadatei(
				out,
				o,
			); err != nil {
				err = errors.Wrapf(err, "Hinweis: %s", o.GetKennung())

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
		return func(o *transacted.Zettel) (err error) {
			if _, err = io.WriteString(
				out, fmt.Sprintf("['%s']\n", o.GetKennung()),
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
		f := MakeFormatterTypFormatterUTIGroups(k, tagp)

		return func(o *transacted.Zettel) (err error) {
			if _, err = f.Format(out, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "action-names":
		f := MakeFormatterTypActionNames(k, true, tagp)

		return func(o *transacted.Zettel) (err error) {
			if _, err = f.Format(out, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	default:
		return func(_ *transacted.Zettel) (err error) {
			err = objekte.MakeErrUnsupportedFormatterValue(
				fv.string,
				gattung.Zettel,
			)
			return
		}
	}
}
