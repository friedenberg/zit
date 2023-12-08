package catgut

import (
	"unicode/utf8"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type SliceRuneScanner struct {
	slice               Slice
	locFirst, locSecond int
	overlap             [6]byte // three bytes from first, three bytes from second
	err                 error
}

var errSliceTooSmall = errors.New("slice too small")

func MakeSliceRuneScanner(slice Slice) (s *SliceRuneScanner, err error) {
	if slice.Len() < 8 && len(slice.Second()) < 4 && len(slice.Second()) > 0 {
		err = errors.Wrap(errSliceTooSmall)
		return
	}

	s = &SliceRuneScanner{}
	s.ResetWith(slice)

	return
}

func (s *SliceRuneScanner) Reset() {
	s.slice = Slice{}
	s.overlap = [6]byte{}
	s.locFirst = 0
	s.locSecond = 0
	s.err = nil
}

func (s *SliceRuneScanner) ResetWith(slice Slice) {
	s.slice = slice
	s.overlap = slice.Overlap()
	s.locFirst = 0
	s.locSecond = 0
	s.err = nil
}

var errInvalidRune = errors.New("invalid rune")

func (s *SliceRuneScanner) ReadRune() (r rune, size int, err error) {
	ok := false
	r, size, ok = s.Scan()

	if !ok {
		if s.err == nil {
			panic("expected error got nil")
		}

		err = s.err

		return
	}

	return
}

func (s *SliceRuneScanner) Error() error {
	return s.err
}

func (s *SliceRuneScanner) Scan() (r rune, width int, ok bool) {
	if s.err != nil {
		return
	}

	firstRemaining := len(s.slice.First()) - s.locFirst

	idxChanged := false

	switch {
	case firstRemaining > 0 && len(s.slice.Second()) == 0:
		fallthrough

	case firstRemaining >= 4:
		r, width = utf8.DecodeRune(s.slice.First()[s.locFirst:])
		s.locFirst += width
		idxChanged = true

	case firstRemaining > 0 && len(s.slice.Second()) > 0:
		r, width = utf8.DecodeRune(s.overlap[firstRemaining-3:])
		s.locFirst += width
		idxChanged = true

	case len(s.slice.Second())-s.locSecond > 0:
		r, width = utf8.DecodeRune(s.slice.Second()[s.locSecond:])
		s.locSecond += width
		idxChanged = true

	default:
		s.err = ErrBufferEmpty
		return
	}

	if r == utf8.RuneError {
		s.err = errInvalidRune
		return
	}

	ok = idxChanged

	return
}
