package box_format

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

	if err = o.Metadata.Description.ReadFromTokenScanner(&ts); err != nil {
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

	{
		if !ts.ScanSkipSpace() {
			if ts.Error() != nil {
				err = errors.Wrap(ts.Error())
			}

			return
		}

		t := ts.GetToken()

		if !t.EqualsString("[") {
			ts.Unscan()
			return
		}
	}

	if !ts.ConsumeSpaces() {
		if ts.Error() != nil {
			err = errors.Wrap(ts.Error())
		}

		return
	}

	{
		var toid catgut.String

	LOOP:
		for ts.Scan() {
			t, tokenType := ts.GetTokenAndType()
			first := t.Bytes()[0]

			if tokenType == token_types.TypeOperator {
				switch first {
				case '.':
					// fall through to append

				default:
					ts.Unscan()
					break LOOP
				}
			}

			if _, err = toid.Append(t); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if ts.Error() != nil {
			err = errors.Wrap(ts.Error())
			return
		}

		if toid.Bytes()[0] == '/' {
			if err = o.ExternalObjectId.SetRaw(toid.String()); err != nil {
				err = errors.Wrap(err)
				o.ExternalObjectId.Reset()
				return
			}
		} else if err = o.ObjectId.ReadFromToken(&toid); err != nil {
			err = errors.Wrap(err)
			o.ObjectId.Reset()
			return
		}
	}

	var k ids.ObjectId

LOOP_AFTER_OID:
	for ts.Scan() {
		t, tokenType, tokenParts := ts.GetTokenAndTypeAndParts()

		if tokenType == token_types.TypeOperator {
			switch t.Bytes()[0] {
			case ']':
				break LOOP_AFTER_OID

			case ' ':
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
