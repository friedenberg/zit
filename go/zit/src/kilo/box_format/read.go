package box_format

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/unicorn"
	"code.linenisgreat.com/zit/go/zit/src/charlie/box"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

// TODO make this error for invalid input
func (f *BoxTransacted) ReadStringFormat(
	el *sku.Transacted,
	rs io.RuneScanner,
) (n int64, err error) {
	var ts box.Scanner
	ts.Reset(rs)

	if err = f.readStringFormatBox(&ts, el); err != nil {
		if errors.Is(err, errNotABox) {
			err = nil
		} else {
			err = errors.WrapExcept(err, io.EOF)
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

	if err = o.Metadata.Description.ReadFromBoxScanner(&ts); err != nil {
		err = errors.Wrap(err)
		return
	}

	n = ts.N()

	return
}

var errNotABox = errors.New("not a box")

func (f *BoxTransacted) openBox(scanner *box.Scanner) (err error) {
	if !scanner.ScanSkipSpace() {
		if scanner.Error() != nil {
			err = errors.Wrap(scanner.Error())
		} else {
			err = io.EOF
		}

		return
	}

	seq := scanner.GetSeq()

	if !seq.MatchAll(box.TokenMatcherOp(box.OpGroupOpen)) {
		err = errNotABox
		scanner.Unscan()
		return
	}

	if !scanner.ConsumeSpacesOrErrorOnFalse() {
		if scanner.Error() != nil {
			err = errors.Wrap(scanner.Error())
		} else {
			err = io.ErrUnexpectedEOF
		}

		return
	}

	return
}

func (f *BoxTransacted) readStringFormatBox(
	scanner *box.Scanner,
	el sku.ExternalLike,
) (err error) {
	o := el.GetSku()

	if err = f.openBox(scanner); err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	{
		if !scanner.ScanDotAllowedInIdentifiers() {
			if scanner.Error() != nil {
				err = errors.Wrap(scanner.Error())
			} else {
				err = io.ErrUnexpectedEOF
			}

			return
		}

		seq := scanner.GetSeq()

		if err = o.ObjectId.ReadFromSeq(seq); err != nil {
			err = nil
			o.ObjectId.Reset()

			if seq.MatchAll(box.TokenTypeLiteral) {
				if err = o.ExternalObjectId.Set(seq.String()); err != nil {
					err = errors.Wrap(err)
					return
				}
			} else if ok, left, _, _ := seq.PartitionFavoringLeft(
				box.TokenMatcherOp(box.OpPathSeparator),
			); ok && left.Len() == 0 {
				if err = o.ExternalObjectId.Set(seq.String()); err != nil {
					err = errors.Wrap(err)
					return
				}
			} else if ok, left, right := seq.MatchEnd(
				box.TokenMatcherOp(box.OpSigilExternal),
				box.TokenTypeIdentifier,
			); ok {
				var g genres.Genre

				// left: one/uno, right: .zettel
				if err = g.Set(right.At(1).String()); err != nil {
					err = nil
				} else {
					if err = o.ObjectId.SetWithGenre(left.String(), g); err != nil {
						o.ObjectId.Reset()
						err = errors.Wrap(err)
						return
					}
				}

				if err = o.ExternalObjectId.Set(seq.String()); err != nil {
					err = errors.Wrap(err)
					return
				}

			} else {
				err = errors.ErrorWithStackf("unsupported seq: %q", seq)
				return
			}
		}

		if o.ObjectId.GetGenre() == genres.InventoryList {
			if err = o.Metadata.Tai.Set(o.ObjectId.String()); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	var k ids.ObjectId

LOOP_AFTER_OID:
	for scanner.ScanDotAllowedInIdentifiers() {
		seq := scanner.GetSeq()

		switch {
		// ] ' '
		case seq.MatchAll(box.TokenTypeOperator):
			r := rune(seq.At(0).Contents[0])

			switch {
			case r == ']':
				break LOOP_AFTER_OID

			case unicorn.IsSpace(r):
				continue
			}

			// "value"
		case seq.MatchAll(box.TokenTypeLiteral):
			if err = o.Metadata.Description.Set(
				seq.String(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			continue

			// @abcd
		case seq.MatchAll(box.TokenMatcherOp('@'), box.TokenTypeIdentifier):
			if err = o.Metadata.Blob.Set(
				string(seq.At(1).Contents),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			continue

			// "value"
		case seq.MatchAll(
			box.TokenTypeLiteral,
		):
			if err = o.Metadata.Description.Set(seq.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			continue

			// key=value key="value"
		case seq.MatchStart(
			box.TokenTypeIdentifier,
			box.TokenMatcherOp(box.OpExact),
		) || seq.MatchStart(
			box.TokenTypeIdentifier,
			box.TokenMatcherOp(box.OpExact),
			box.TokenTypeLiteral,
		):

			value := seq[2:]

			field := string_format_writer.Field{
				Key:   string(seq.At(0).Contents),
				Value: value.String(),
			}

			field.ColorType = string_format_writer.ColorTypeUserData
			o.Metadata.Fields = append(o.Metadata.Fields, field)
			continue
		}

		if err = k.ReadFromSeq(seq); err != nil {
			err = nil
			scanner.Unscan()
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
			err = errors.Wrapf(err, "Seq: %q", seq)
			return
		}

		k.Reset()
	}

	if scanner.Error() != nil {
		err = errors.Wrap(scanner.Error())
		return
	}

	return
}
