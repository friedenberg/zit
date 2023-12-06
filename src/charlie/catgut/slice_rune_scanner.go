package catgut

import (
	"unicode/utf8"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type SliceRuneScanner struct {
	slice               Slice
	locFirst, locSecond int
	overlap             [6]byte // three bytes from first, three bytes from second
}

var errSliceTooSmall = errors.New("slice too small")

func MakeSliceRuneScanner(slice Slice) (s *SliceRuneScanner, err error) {
	if slice.Len() < 8 {
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
}

func (s *SliceRuneScanner) ResetWith(slice Slice) {
	s.slice = slice
	s.overlap = slice.Overlap()
	s.locFirst = 0
	s.locSecond = 0
}

var errInvalidRune = errors.New("invalid rune")

func (s *SliceRuneScanner) ReadRune() (r rune, size int, err error) {
	ok := false
	r, size, ok = s.Scan()

	if !ok {
		err = errors.Wrap(errInvalidRune)
		return
	}

	return
}

func (s *SliceRuneScanner) Scan() (r rune, width int, ok bool) {
	firstRemaining := len(s.slice.First()) - s.locFirst

	switch {
	case firstRemaining >= 4:
		r, width = utf8.DecodeRune(s.slice.First()[s.locFirst:])
		s.locFirst += width
		return

	case firstRemaining > 0:
		r, width = utf8.DecodeRune(s.overlap[firstRemaining-3:])
		s.locFirst += width
		return

	case len(s.slice.Second())-s.locSecond > 0:
		r, width = utf8.DecodeRune(s.slice.Second()[s.locSecond:])
		s.locSecond += width
		return
	}

	return
}
