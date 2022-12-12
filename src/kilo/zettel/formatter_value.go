package zettel

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/collections_coding"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/juliett/konfig_compiled"
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
	case "text", "objekte", "json", "toml", "action-names":
		f.string = v1

	default:
		err = errors.Errorf("unsupported format type: %s", v)
		return
	}

	return
}

func (f *FormatterValue) FuncFormatter(
	out io.Writer,
	af gattung.AkteIOFactory,
	k konfig_compiled.Compiled,
) collections.WriterFunc[*Transacted] {
	switch f.string {
	case "objekte":
		f := objekte.MakeFormatter[*Transacted](af)

		return func(o *Transacted) (err error) {
			if _, err = f.WriteFormat(out, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "text":
		f := Text{}

		return func(o *Transacted) (err error) {
			c := FormatContextWrite{
				Out:               out,
				Zettel:            o.Objekte,
				AkteReaderFactory: af,
				IncludeAkte:       true,
			}

			if _, err = f.WriteTo(c); err != nil {
				err = errors.Wrap(err)
				return
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

	case "toml":
    //TODO-P0 limit to just zettels that support toml
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

	case "action-names":
		f := &EncoderTypActionNames{
			konfig: k,
		}

		return func(o *Transacted) (err error) {
			c := FormatContextWrite{
				Out:               out,
				Zettel:            o.Objekte,
				AkteReaderFactory: af,
				IncludeAkte:       true,
			}

			if _, err = f.WriteTo(c); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	default:
		return func(_ *Transacted) (err error) {
			return errors.Errorf("unsupported format for typen: %s", f.string)
		}
	}
}
