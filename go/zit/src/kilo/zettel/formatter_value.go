package zettel

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/src/echo/format"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/src/juliett/objekte"
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
		"json-blob",
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
	k *konfig.Compiled,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) schnittstellen.FuncIter[*sku.Transacted] {
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
	k *konfig.Compiled,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) schnittstellen.FuncIter[*sku.Transacted] {
	errors.TodoP2("convert to verzeichnisse")

	switch fv.string {
	case "formatters":
		return func(o *sku.Transacted) (err error) {
			t := k.GetApproximatedTyp(o.GetTyp())

			if !t.HasValue {
				return
			}

			tt := t.ActualOrNil()

			var ta *typ_akte.V0

			if ta, err = tagp.GetAkte(tt.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer tagp.PutAkte(ta)

			lw := format.MakeLineWriter()

			for fn, f := range ta.Formatters {
				fe := f.FileExtension

				if fe == "" {
					fe = fn
				}

				lw.WriteFormat("%s\t%s", fn, fe)
			}

			if _, err = lw.WriteTo(out); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "typ-vim-syntax-type":
		return func(o *sku.Transacted) (err error) {
			var t *sku.Transacted

			if t = k.GetApproximatedTyp(
				o.GetTyp(),
			).ApproximatedOrActual(); t == nil {
				return
			}

			var ta *typ_akte.V0

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
		op := objekte_format.Options{}

		return func(o *sku.Transacted) (err error) {
			if _, err = f.FormatPersistentMetadatei(out, o, op); err != nil {
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

		return func(o *sku.Transacted) (err error) {
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

	case "json-blob":
		e := json.NewEncoder(out)

		return func(o *sku.Transacted) (err error) {
			var a map[string]interface{}

			var r sha.ReadCloser

			if r, err = af.AkteReader(o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.Deferred(&err, r.Close)

			d := toml.NewDecoder(r)

			if err = d.Decode(&a); err != nil {
				log.Err().Printf("%s: %s", o, err)
				err = nil
				return
			}

			a["description"] = o.Metadatei.Bezeichnung.String()
			a["identifier"] = o.Kennung.String()

			if err = e.Encode(&a); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "toml":
		errors.TodoP3("limit to only zettels supporting toml")
		return func(o *sku.Transacted) (err error) {
			var a map[string]interface{}

			var r sha.ReadCloser

			if r, err = af.AkteReader(o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.Deferred(&err, r.Close)

			d := toml.NewDecoder(r)

			if err = d.Decode(&a); err != nil {
				log.Err().Printf("%s: %s", o, err)
				err = nil
				return
			}

			a["description"] = o.Metadatei.Bezeichnung.String()
			a["identifier"] = o.Kennung.String()

			e := toml.NewEncoder(out)

			if err = e.Encode(&a); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, err = out.Write([]byte("\x00")); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "typ-formatter-uti-groups":
		f := MakeFormatterTypFormatterUTIGroups(k, tagp)

		return func(o *sku.Transacted) (err error) {
			if _, err = f.Format(out, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "action-names":
		f := MakeFormatterTypActionNames(k, true, tagp)

		return func(o *sku.Transacted) (err error) {
			if _, err = f.Format(out, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	default:
		return func(_ *sku.Transacted) (err error) {
			err = objekte.MakeErrUnsupportedFormatterValue(
				fv.string,
				gattung.Zettel,
			)
			return
		}
	}
}
