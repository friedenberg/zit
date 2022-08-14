package organize_text

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/delta/etikett"
)

type assignmentLineReader struct {
	lineNo            int
	root              *assignment
	currentAssignment *assignment
}

type line struct {
	prefix string
	value  string
}

func (l line) String() string {
	return fmt.Sprintf("%s %s", l.prefix, l.value)
}

func (l *line) Set(v string) (err error) {
	v = strings.TrimSpace(v)

	if len(v) == 0 {
		err = errors.Errorf("line not long enough")
		return
	}

	firstSpace := strings.Index(v, " ")

	if firstSpace == -1 {
		l.prefix = v
		return
	}

	l.prefix = strings.TrimSpace(v[:firstSpace])
	l.value = strings.TrimSpace(v[firstSpace:])

	return
}

func (l line) PrefixRune() rune {
	if len(l.prefix) == 0 {
		panic(errors.Errorf("cannot find prefix in line: %q", l.value))
	}

	return rune(l.prefix[0])
}

func (l line) Depth(r rune) (depth int, err error) {
	for i, c := range l.prefix {
		if c != r {
			err = errors.Errorf("rune at index %d is %c and not %c", i, c, r)
			return
		}

		depth++
	}

	return
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
			err = errors.Error(err)
			return
		}

		n += int64(len(s))

		s = strings.TrimSuffix(s, "\n")
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
		err = ar.readOneZettel(l)
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
		err = errors.Error(err)
		return
	}

	currentEtiketten := etikett.NewSet()

	if l.value != "" {
		if err = currentEtiketten.Set(l.value); err != nil {
			err = ErrorRead{
				error:  err,
				line:   ar.lineNo,
				column: 2,
			}

			return
		}
	}

	var newAssignment *assignment

	// logz.Printf("line: %q\n", l)
	// logz.Printf("previous: %d, current: %d\n", ar.currentAssignment.depth, depth)
	if depth < ar.currentAssignment.depth {
		newAssignment, err = ar.readOneHeadingLesserDepth(depth, currentEtiketten)
	} else if depth == ar.currentAssignment.depth {
		newAssignment, err = ar.readOneHeadingEqualDepth(depth, currentEtiketten)
	} else {
		//always use currentEtiketten.depth + 1 because it corrects movements
		newAssignment, err = ar.readOneHeadingGreaterDepth(ar.currentAssignment.depth+1, currentEtiketten)
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

func (ar *assignmentLineReader) readOneHeadingLesserDepth(d int, e *etikett.Set) (newCurrent *assignment, err error) {
	// logz.Print("depth count is <")

	depthDiff := d - ar.currentAssignment.depth

	if newCurrent, err = ar.currentAssignment.nthParent(depthDiff - 1); err != nil {
		err = errors.Error(err)
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
		assignment.etiketten = *e
		newCurrent.addChild(assignment)
		// logz.Print("adding to parent")
		// logz.Print("child", assignment.etiketten)
		// logz.Print("parent", newCurrent.etiketten)
		newCurrent = assignment
	}

	return
}

func (ar *assignmentLineReader) readOneHeadingEqualDepth(d int, e *etikett.Set) (newCurrent *assignment, err error) {
	// logz.Print("depth count is ==")

	if newCurrent, err = ar.currentAssignment.nthParent(1); err != nil {
		err = errors.Error(err)
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
		assignment.etiketten = *e
		newCurrent.addChild(assignment)
		newCurrent = assignment
	}

	return
}

func (ar *assignmentLineReader) readOneHeadingGreaterDepth(d int, e *etikett.Set) (newCurrent *assignment, err error) {
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
		assignment.etiketten = *e
		newCurrent.addChild(assignment)
		// logz.Print("adding to parent")
		// logz.Print("child", assignment)
		// logz.Print("parent", newCurrent)
		newCurrent = assignment
	}

	return
}

func (ar *assignmentLineReader) readOneZettel(l line) (err error) {
	// logz.Print("reading one zettel", l)

	var z zettel

	if err = z.Set(l.String()); err == nil {
		// logz.Print("added to named zettels")
		ar.currentAssignment.named.Add(z)
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

	var nz newZettel

	if err = nz.Set(l.String()); err == nil {
		// logz.Print("added to unnamed zettels")
		ar.currentAssignment.unnamed.Add(nz)
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
