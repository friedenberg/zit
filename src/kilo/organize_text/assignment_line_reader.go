package organize_text

import (
	"io"
	"unicode"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/unicorn"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/kennung_fmt"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type assignmentLineReader struct {
	options            Options
	lineNo             int
	root               *assignment
	currentAssignment  *assignment
	stringFormatReader catgut.StringFormatReader[*sku.Transacted]
}

func (ar *assignmentLineReader) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := catgut.MakeRingBuffer(r1, 0)
	rbs := catgut.MakeRingBufferScanner(r)

	ar.root = newAssignment(0)
	ar.currentAssignment = ar.root

	for {
		var sl []byte
		sl, err = rbs.AdvanceToFirstMatch(unicorn.Not(unicode.IsSpace))
		n += int64(len(sl))

		if err == io.EOF && len(sl) <= 1 {
			err = nil
			break
		}

		if err == catgut.ErrBufferEmpty || err == catgut.ErrNoMatch {
			var n1 int64
			n1, err = r.Fill()

			if n1 == 0 && err == io.EOF {
				err = nil
				break
			} else {
				continue
			}
		}

		if err != nil && err != io.EOF {
			err = errors.Wrap(err)
			return
		}

		slen := len(sl)

		if slen >= 1 {
			pr := sl[0]

			switch pr {
			case '#':
				if err = ar.readOneHeading(r, sl); err != nil {
					err = errors.Wrap(err)
					return
				}

			case '-':
				if err = ar.readOneObj(r); err != nil {
					err = errors.Wrap(err)
					return
				}

			default:
				err = ErrorRead{
					error:  errors.Errorf("unsupported verb %q", pr),
					line:   ar.lineNo,
					column: 0,
				}

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
	match []byte,
) (err error) {
	depth := unicorn.CountRune(match, '#')

	currentEtiketten := kennung.MakeMutableEtikettSet()

	reader := kennung_fmt.MakeEtikettenReader()

	if _, err = reader.ReadStringFormat(rb, currentEtiketten); err != nil {
		err = errors.Wrap(err)
		return
	}

	var newAssignment *assignment

	if depth < ar.currentAssignment.depth {
		newAssignment, err = ar.readOneHeadingLesserDepth(
			depth,
			currentEtiketten,
		)
	} else if depth == ar.currentAssignment.depth {
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
	e kennung.EtikettSet,
) (newCurrent *assignment, err error) {
	depthDiff := d - ar.currentAssignment.Depth()

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
		assignment.etiketten = e.CloneSetPtrLike()
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
	e kennung.EtikettSet,
) (newCurrent *assignment, err error) {
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
		assignment.etiketten = e.CloneSetPtrLike()
		newCurrent.addChild(assignment)
		newCurrent = assignment
	}

	return
}

func (ar *assignmentLineReader) readOneHeadingGreaterDepth(
	d int,
	e kennung.EtikettSet,
) (newCurrent *assignment, err error) {
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
		assignment.etiketten = e.CloneSetPtrLike()
		newCurrent.addChild(assignment)
		// logz.Print("adding to parent")
		// logz.Print("child", assignment)
		// logz.Print("parent", newCurrent)
		newCurrent = assignment
	}

	return
}

func (ar *assignmentLineReader) readOneObj(r *catgut.RingBuffer) (err error) {
	// logz.Print("reading one zettel", l)

	var z obj

	if _, err = ar.stringFormatReader.ReadStringFormat(r, &z.Sku); err != nil {
		err = ErrorRead{
			error:  err,
			line:   ar.lineNo,
			column: 2,
		}

		return
	}

	if z.Sku.Kennung.IsEmpty() {
		ar.currentAssignment.unnamed.Add(&z)
	} else {
		ar.currentAssignment.named.Add(&z)
	}

	return
}
