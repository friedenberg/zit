package organize_text

import (
	"io"
	"unicode"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/unicorn"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/etiketten_path"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/kennung_fmt"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type assignmentLineReader struct {
	options            Options
	lineNo             int
	root               *Assignment
	currentAssignment  *Assignment
	stringFormatReader catgut.StringFormatReader[*sku.Transacted]
}

func (ar *assignmentLineReader) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := catgut.MakeRingBuffer(r1, 0)
	rbs := catgut.MakeRingBufferScanner(r)

	ar.root = newAssignment(0)
	ar.currentAssignment = ar.root

LOOP:
	for {
		var sl catgut.Slice
		var offsetPlusMatch int

		sl, offsetPlusMatch, err = rbs.FirstMatch(unicorn.Not(unicode.IsSpace))

		if err == io.EOF && sl.Len() == 0 {
			err = nil
			break
		}

		switch err {
		case catgut.ErrBufferEmpty, catgut.ErrNoMatch:
			var n1 int64
			n1, err = r.Fill()

			if n1 == 0 && err == io.EOF {
				err = nil
				break LOOP
			} else {
				err = nil
				continue
			}
		}

		if err != nil && err != io.EOF {
			err = errors.Wrap(err)
			return
		}

		r.AdvanceRead(offsetPlusMatch)
		n += int64(sl.Len())
		sb := sl.SliceBytes()

		slen := sl.Len()

		if slen >= 1 {
			pr := sl.FirstByte()

			switch pr {
			case '#':
				if err = ar.readOneHeading(r, sb); err != nil {
					err = errors.Wrap(err)
					return
				}

			case '%':
				if err = ar.readOneObj(r, etiketten_path.TypeUnknown); err != nil {
					if err == io.EOF {
						err = nil
					} else {
						err = errors.Wrap(err)
						return
					}
				}

			case '-':
				if err = ar.readOneObj(r, etiketten_path.TypeDirect); err != nil {
					if err == io.EOF {
						err = nil
					} else {
						err = errors.Wrap(err)
						return
					}
				}

			default:
				err = errors.Errorf("unsupported verb. slice: %q", sl)
				return
			}
		}

		ar.lineNo++

		if err == io.EOF {
			err = nil
			break
		} else {
			continue
		}
	}

	return
}

func (ar *assignmentLineReader) readOneHeading(
	rb *catgut.RingBuffer,
	match catgut.SliceBytes,
) (err error) {
	depth := unicorn.CountRune(match.Bytes, '#')

	currentEtiketten := kennung.MakeMutableTagSet()

	reader := kennung_fmt.MakeEtikettenReader()

	if _, err = reader.ReadStringFormat(rb, currentEtiketten); err != nil {
		err = errors.Wrap(err)
		return
	}

	var newAssignment *Assignment

	if depth < ar.currentAssignment.Depth {
		newAssignment, err = ar.readOneHeadingLesserDepth(
			depth,
			currentEtiketten,
		)
	} else if depth == ar.currentAssignment.Depth {
		newAssignment, err = ar.readOneHeadingEqualDepth(depth, currentEtiketten)
	} else {
		// always use currentEtiketten.depth + 1 because it corrects movements
		newAssignment, err = ar.readOneHeadingGreaterDepth(depth, currentEtiketten)
	}

	if err != nil {
		err = ErrorRead{
			error:  err,
			line:   ar.lineNo,
			column: 2,
		}

		return
	}

	if newAssignment == nil {
		err = errors.Errorf("read heading function return nil new assignment")
		return
	}

	ar.currentAssignment = newAssignment

	return
}

func (ar *assignmentLineReader) readOneHeadingLesserDepth(
	d int,
	e kennung.TagSet,
) (newCurrent *Assignment, err error) {
	depthDiff := d - ar.currentAssignment.GetDepth()

	if newCurrent, err = ar.currentAssignment.nthParent(depthDiff - 1); err != nil {
		err = errors.Wrap(err)
		return
	}

	if e.Len() == 0 {
		// `
		// # task-todo
		// ## priority-1
		// - wow
		// #
		// `
		// logz.Print("new set is empty")
	} else {
		// `
		// # task-todo
		// ## priority-1
		// - wow
		// # zz-inbox
		// `
		assignment := newAssignment(d)
		assignment.Etiketten = e.CloneSetPtrLike()
		newCurrent.addChild(assignment)
		// logz.Print("adding to parent")
		// logz.Print("child", assignment.etiketten)
		// logz.Print("parent", newCurrent.etiketten)
		newCurrent = assignment
	}

	return
}

func (ar *assignmentLineReader) readOneHeadingEqualDepth(
	d int,
	e kennung.TagSet,
) (newCurrent *Assignment, err error) {
	// logz.Print("depth count is ==")

	if newCurrent, err = ar.currentAssignment.nthParent(1); err != nil {
		err = errors.Wrap(err)
		return
	}

	if e.Len() == 0 {
		// `
		// # task-todo
		// ## priority-1
		// - wow
		// ##
		// `
		// logz.Print("new set is empty")
	} else {
		// `
		// # task-todo
		// ## priority-1
		// - wow
		// ## priority-2
		// `
		assignment := newAssignment(d)
		assignment.Etiketten = e.CloneSetPtrLike()
		newCurrent.addChild(assignment)
		newCurrent = assignment
	}

	return
}

func (ar *assignmentLineReader) readOneHeadingGreaterDepth(
	d int,
	e kennung.TagSet,
) (newCurrent *Assignment, err error) {
	// logz.Print("depth count is >")
	// logz.Print(e)

	newCurrent = ar.currentAssignment

	if e.Len() == 0 {
		// `
		// # task-todo
		// ## priority-1
		// - wow
		// ###
		// `
		// logz.Print("new set is empty")
	} else {
		// `
		// # task-todo
		// ## priority-1
		// - wow
		// ### priority-2
		// `
		assignment := newAssignment(d)
		assignment.Etiketten = e.CloneSetPtrLike()
		newCurrent.addChild(assignment)
		// logz.Print("adding to parent")
		// logz.Print("child", assignment)
		// logz.Print("parent", newCurrent)
		newCurrent = assignment
	}

	return
}

func (ar *assignmentLineReader) readOneObj(
	r *catgut.RingBuffer,
	t etiketten_path.Type,
) (err error) {
	// logz.Print("reading one zettel", l)

	var z obj
	z.Type = t

	// {
	// 	var sl catgut.Slice

	// 	if sl, err = r.PeekUpto('['); err != nil {
	// 		if collections.IsErrNotFound(err) {
	// 			err = nil
	// 		} else {
	// 			err = errors.Wrap(err)
	// 			return
	// 		}
	// 	} else if sl.LastByte() == '%' {
	// 		z.Type = etiketten_path.TypeUnknown
	// 		r.AdvanceRead(sl.Len())
	// 	}
	// }

	if _, err = ar.stringFormatReader.ReadStringFormat(r, &z.Transacted); err != nil {
		err = ErrorRead{
			error:  err,
			line:   ar.lineNo,
			column: 2,
		}

		return
	}

	if z.Kennung.IsEmpty() {
		// set empty hinweis to ensure middle is '/'
		if err = z.Kennung.SetWithIdLike(kennung.Hinweis{}); err != nil {
			err = errors.Wrap(err)
			return
		}

		ar.currentAssignment.AddObjekte(&z)

		return
	}

	if err = ar.options.Abbr.ExpandHinweisOnly(&z.Kennung); err != nil {
		err = errors.Wrap(err)
		return
	}

	ar.currentAssignment.AddObjekte(&z)

	return
}
