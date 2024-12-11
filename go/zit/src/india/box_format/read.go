package box_format

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/token_types"
	"code.linenisgreat.com/zit/go/zit/src/alfa/unicorn"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/query_spec"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO make this error for invalid input
func (f *BoxTransacted) ReadStringFormat(
	rs io.RuneScanner,
	el *sku.Transacted,
) (n int64, err error) {
	var ts query_spec.TokenScanner
	ts.Reset(rs)

	if err = f.readStringFormatBox(&ts, el); err != nil {
		if errors.Is(err, errNotABox) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if ts.Error() != nil {
		err = errors.Wrap(ts.Error())
		return
	}

	n = ts.N()

	o := el.GetSku()

	if f.optionsPrint.DescriptionInBox {
		return
	}

	if err = o.Metadata.Description.ReadFromTokenScanner(&ts); err != nil {
		err = errors.Wrap(err)
		return
	}

	n = ts.N()

	return
}

var errNotABox = errors.New("not a box")

func (f *BoxTransacted) openBox(ts *query_spec.TokenScanner) (err error) {
	if !ts.ScanSkipSpace() {
		if ts.Error() != nil {
			err = errors.Wrap(ts.Error())
		} else {
			err = io.EOF
		}

		return
	}

	t, tokenType := ts.GetTokenAndType()

	if tokenType != token_types.TypeOperator || t.Bytes()[0] != '[' {
		err = errNotABox
		ts.Unscan()
		return
	}

	if !ts.ConsumeSpacesOrErrorOnFalse() {
		if ts.Error() != nil {
			err = errors.Wrap(ts.Error())
		} else {
			err = io.ErrUnexpectedEOF
		}

		return
	}

	return
}

func (f *BoxTransacted) readStringFormatBox(
	ts *query_spec.TokenScanner,
	el sku.ExternalLike,
) (err error) {
	o := el.GetSku()

	if err = f.openBox(ts); err != nil {
		err = errors.Wrap(err)
		return
	}

	{
		if !ts.ScanDotAllowedInIdentifiers() {
			if ts.Error() != nil {
				err = errors.Wrap(ts.Error())
			} else {
				err = io.ErrUnexpectedEOF
			}
		}

		t, _, parts := ts.GetTokenAndTypeAndParts()

		if t.Bytes()[0] == '/' {
			if err = o.ExternalObjectId.SetRaw(t.String()); err != nil {
				err = errors.Wrap(err)
				o.ExternalObjectId.Reset()
				return
			}
		} else if err = o.ObjectId.ReadFromTokenAndParts(t, parts); err != nil {
			err = errors.Wrap(err)
			o.ObjectId.Reset()
			return
		}

		if o.ObjectId.GetGenre() == genres.InventoryList {
			if err = o.Metadata.Tai.Set(t.String()); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	var k ids.ObjectId

LOOP_AFTER_OID:
	for ts.ScanDotAllowedInIdentifiers() {
		t, tokenType, tokenParts := ts.GetTokenAndTypeAndParts()

		if tokenType == token_types.TypeOperator {
			r := rune(t.Bytes()[0])

			switch {
			case r == ']':
				break LOOP_AFTER_OID

			case unicorn.IsSpace(r):
				continue LOOP_AFTER_OID
			}
		}

		switch tokenType {
		case token_types.TypeField:
			if len(tokenParts.Left) == 0 {
			} else {
				field := string_format_writer.Field{
					Key:   string(tokenParts.Left),
					Value: string(tokenParts.Right),
				}

				field.ColorType = string_format_writer.ColorTypeUserData
				o.Metadata.Fields = append(o.Metadata.Fields, field)
			}

			continue

		case token_types.TypeLiteral:
			if len(tokenParts.Left) == 0 {
			} else if len(tokenParts.Right) == 0 {
				if err = o.Metadata.Description.Set(string(tokenParts.Left)); err != nil {
					err = errors.Wrap(err)
					return
				}
			} else {
				field := string_format_writer.Field{
					Value: t.String(),
				}

				field.ColorType = string_format_writer.ColorTypeUserData
				o.Metadata.Fields = append(o.Metadata.Fields, field)
			}

			continue

		case token_types.TypeIdentifier:
			switch t.Bytes()[0] {
			case '/':
				if err = o.ExternalObjectId.SetRaw(t.String()); err != nil {
					o.ExternalObjectId.Reset()
					return
				}

				continue

			case '@':
				if err = o.Metadata.Blob.Set(t.String()); err != nil {
					err = errors.Wrap(err)
					return
				}

				continue
			}
		}

		if err = k.TodoSetBytes(t); err != nil {
			err = errors.Wrapf(err, "Type: %s", tokenType)
			return
		}

		g := k.GetGenre()

		switch g {
		case genres.InventoryList:
			// TODO make more performant
			if err = o.Metadata.Tai.Set(k.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

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

		case genres.Blob:
			field := string_format_writer.Field{
				Key:   string(tokenParts.Left),
				Value: string(tokenParts.Right),
			}

			field.ColorType = string_format_writer.ColorTypeUserData
			o.Metadata.Fields = append(o.Metadata.Fields, field)

		default:
			err = genres.MakeErrUnsupportedGenre(k.GetGenre())
			err = errors.Wrapf(err, "Token: %q", t)
			return
		}

		k.Reset()
	}

	if ts.Error() != nil {
		err = errors.Wrap(ts.Error())
		return
	}

	return
}
