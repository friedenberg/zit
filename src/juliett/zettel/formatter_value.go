package zettel

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/collections_coding"
	"github.com/friedenberg/zit/src/delta/format"
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
		"log",
		"akte",
		"hinweis-text",
		"text",
		"objekte",
		"json",
		"toml",
		"action-names",
		"hinweis-akte":
		f.string = v1

	default:
		err = errors.Errorf("unsupported format type: %s", v)
		return
	}

	return
}

func (fv *FormatterValue) FuncFormatterVerzeichnisse(
	out io.Writer,
	af schnittstellen.AkteIOFactory,
	k konfig.Compiled,
	logFunc collections.WriterFunc[*Transacted],
) collections.WriterFunc[*Transacted] {
	return fv.FuncFormatter(
		out,
		af,
		k,
		logFunc,
	)
}

func (fv *FormatterValue) FuncFormatter(
	out io.Writer,
	af schnittstellen.AkteIOFactory,
	k konfig.Compiled,
	logFunc collections.WriterFunc[*Transacted],
) collections.WriterFunc[*Transacted] {
	errors.TodoP2("convert to verzeichnisse")

	switch fv.string {
	case "log":
		return logFunc

	case "formatters":
		return func(o *Transacted) (err error) {
			t := k.GetApproximatedTyp(o.Objekte.Typ)

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
			if _, err = io.WriteString(out, o.Objekte.Typ.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "typ-vim-syntax-type":
		return func(o *Transacted) (err error) {
			var t konfig.ApproximatedTyp

			if t = k.GetApproximatedTyp(o.Objekte.Typ); !t.HasValue() {
				return
			}

			if _, err = io.WriteString(
				out,
				t.Unwrap().Objekte.Akte.VimSyntaxType,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "objekte":
		f := FormatObjekte{}

		return func(o *Transacted) (err error) {
			if _, err = f.Format(out, &o.Objekte); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "text", "hinweis-text":
		f := MakeObjekteTextFormatterIncludeAkte(
			standort.Standort{},
			k,
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

			c := ObjekteFormatterContext{
				Zettel: o.Objekte,
			}

			if _, err = f.Format(out, &c); err != nil {
				err = errors.Wrapf(err, "Hinweis: %s", o.Sku.Kennung)

				if errors.IsNotExist(err) {
					errors.Err().Print(err)
					err = nil
				} else {
					return
				}
			}

			return
		}

	case "json":
		f := collections_coding.MakeEncoderJson[Transacted](out)

		return func(o *Transacted) (err error) {
			if _, err = f.Encode(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "akte":
		return func(o *Transacted) (err error) {
			var r sha.ReadCloser

			if r, err = af.AkteReader(o.Objekte.Akte); err != nil {
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

			if r, err = af.AkteReader(o.Objekte.Akte); err != nil {
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
			c := ObjekteFormatterContext{
				Zettel:      o.Objekte,
				IncludeAkte: true,
			}

			if _, err = f.Format(out, c); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "action-names":
		f := MakeFormatterTypActionNames(k, true)

		return func(o *Transacted) (err error) {
			c := ObjekteFormatterContext{
				Zettel:      o.Objekte,
				IncludeAkte: true,
			}

			if _, err = f.Format(out, c); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "hinweis-akte":
		return func(o *Transacted) (err error) {
			errors.TodoP3("convert into an option")
			if o.Objekte.Akte.IsNull() {
				return
			}

			if _, err = io.WriteString(
				out, fmt.Sprintf("%s %s\n", o.Sku.Kennung, o.Objekte.Akte),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	default:
		return func(_ *Transacted) (err error) {
			return errors.Errorf("unsupported format for typen: %s", fv.string)
		}
	}
}
