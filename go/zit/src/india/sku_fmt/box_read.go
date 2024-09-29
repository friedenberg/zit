package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/token_types"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/query_spec"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (f *Box) ReadStringFormat(
	rb *catgut.RingBuffer,
	el *sku.Transacted,
) (n int64, err error) {
	var ts query_spec.TokenScanner
	ts.Reset(catgut.MakeRingBufferRuneScanner(rb))

	if err = f.readStringFormatBox(&ts, el); err != nil {
		err = errors.Wrap(err)
		return
	}

	n = ts.N()

	o := el.GetSku()

	if err = o.Metadata.Description.ReadFromRuneScanner(&ts); err != nil {
		err = errors.Wrap(err)
		return
	}

	n = ts.N()

	return
}

func (f *Box) readStringFormatBox(
	ts *query_spec.TokenScanner,
	el sku.ExternalLike,
) (err error) {
	o := el.GetSku()

	state := 0
	var k ids.ObjectId

LOOP:
	for ts.ScanIdentifierLikeSkipSpaces() {
		t, tokenType, tokenParts := ts.GetTokenAndTypeAndParts()

		if t.EqualsString(" ") || t.EqualsString("\n") {
			continue
		}

		switch state {
		case 0:
			if !t.EqualsString("[") {
				ts.Unscan()
				return
			}

			state++

		case 1:
			if t.Bytes()[0] == '/' {
				if err = o.ExternalObjectId.SetRaw(t.String()); err != nil {
					o.ExternalObjectId.Reset()
					return
				}
			} else if err = o.ObjectId.ReadFromToken(t); err != nil {
				o.ObjectId.Reset()
				return
			}

			state++

		case 2:
			if t.EqualsString("]") {
				break LOOP
			}

			switch tokenType {
			case token_types.TypeLiteral:
				if len(tokenParts.Left) == 0 {
				} else {
					field := string_format_writer.Field{
						Key:   string(tokenParts.Left),
						Value: string(tokenParts.Right),
					}

					o.Metadata.Fields = append(o.Metadata.Fields, field)
				}

				continue LOOP
			case token_types.TypeIdentifier:
				if t.Bytes()[0] == '/' {
					if err = o.ExternalObjectId.SetRaw(t.String()); err != nil {
						o.ExternalObjectId.Reset()
						return
					}

					continue LOOP
				}
			}

			if err = k.TodoSetBytes(t); err != nil {
				err = errors.Wrap(err)
				return
			}

			g := k.GetGenre()

			switch g {
			case genres.Type:
				if err = o.Metadata.Type.TodoSetFromObjectId(&k); err != nil {
					err = errors.Wrap(err)
					return
				}

			case genres.Tag:
				var e ids.Tag

				if err = e.TodoSetFromObjectId(&k); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = o.AddTagPtr(&e); err != nil {
					err = errors.Wrap(err)
					return
				}

			default:
				err = genres.MakeErrUnsupportedGenre(k.GetGenre())
				err = errors.Wrapf(err, "Token: %q", t)
				return
			}

			k.Reset()

		default:
			err = errors.Errorf("invalid state: %d", state)
			return
		}
	}

	if ts.Error() != nil {
		err = errors.Wrap(ts.Error())
		return
	}

	return
}
