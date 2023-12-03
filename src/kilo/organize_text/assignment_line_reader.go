package organize_text

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type assignmentLineReader struct {
	options           Options
	lineNo            int
	root              *assignment
	currentAssignment *assignment
	ex                kennung.Abbr
}

func (ar *assignmentLineReader) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	ar.root = newAssignment(0)
	ar.currentAssignment = ar.root

	for {
		var s string
		s, err = r.ReadString('\n')

		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n += int64(len(s))

		s = strings.TrimSpace(s)
		slen := len(s)

		if slen < 1 {
			continue
		}

		var l line

		if err = l.Set(s); err != nil {
			err = ErrorRead{
				error: err,
				line:  ar.lineNo,
			}
		}

		if err = ar.readOne(l); err != nil {
			err = ErrorRead{
				error: err,
				line:  ar.lineNo,
			}
			return
		}

		ar.lineNo++
	}

	return
}

func (ar *assignmentLineReader) readOne(l line) (err error) {
	switch l.PrefixRune() {
	case '#':
		return ar.readOneHeading(l)

	case '-':
		err = ar.readOneObj(l)
		// logz.Print(len(ar.currentAssignment.named))
		return err

	default:
		err = ErrorRead{
			error:  errors.Errorf("unsupported verb %q, %q", l.PrefixRune(), l),
			line:   ar.lineNo,
			column: 0,
		}

		return
	}
}

func (ar *assignmentLineReader) readOneHeading(l line) (err error) {
	var depth int

	// logz.Print("getting depth count")

	if depth, err = l.Depth('#'); err != nil {
		err = errors.Wrap(err)
		return
	}

	currentEtiketten := kennung.MakeMutableEtikettSet()

	flag := collections_ptr.MakeFlagCommasFromExisting(
		collections_ptr.SetterPolicyAppend,
		currentEtiketten,
	)

	if l.value != "" {
		if err = flag.Set(l.value); err != nil {
			err = ErrorRead{
				error:  err,
				line:   ar.lineNo,
				column: 2,
			}

			return
		}

		errors.Log().Print(flag.String())
		errors.Log().Print(l.value)
		errors.Log().Print(currentEtiketten.Len())
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

func (ar *assignmentLineReader) readOneObj(l line) (err error) {
	// logz.Print("reading one zettel", l)

	var z obj

	if err = z.setExistingObj(ar.options.PrintOptions, l.String(), ar.ex); err == nil {
		// logz.Print("added to named zettels")
		ar.currentAssignment.named.Add(&z)
		// logz.Print(len(ar.currentAssignment.named))
		return
	}

	if len(l.String()) < 2 && l.String()[:2] != "- " {
		err = ErrorRead{
			error:  err,
			line:   ar.lineNo,
			column: 2,
		}

		return
	}

	var nz obj

	if err = nz.setNewObj(l.String()); err == nil {
		// logz.Print("added to unnamed zettels")
		ar.currentAssignment.unnamed.Add(&nz)
		return
	}

	// logz.Print("failed to read zettel")

	err = ErrorRead{
		error:  err,
		line:   ar.lineNo,
		column: 2,
	}

	return
}
