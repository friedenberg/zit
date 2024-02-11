package catgut

import (
	"io"
	"unicode/utf8"

	"code.linenisgreat.com/zit/src/alfa/errors"
)

func MakeSliceRuneScanner(slice Slice) (s *SliceRuneScanner) {
	s = &SliceRuneScanner{}
	s.ResetWith(slice)

	return
}

type SliceRuneScanner struct {
	slice                       Slice
	overlap                     [6]byte // three bytes from first, three bytes from second
	overlapFirst, overlapSecond int

	locFirst, locSecond         int
	lastLocFirst, lastLocSecond int
	err                         error
}

func (s *SliceRuneScanner) Offset() int {
	return s.locFirst + s.locSecond
}

func (s *SliceRuneScanner) Remaining() int {
	return s.slice.Len() - (s.locFirst + s.locSecond)
}

func (s *SliceRuneScanner) Reset() {
	s.slice = Slice{}
	s.overlap = [6]byte{}
	s.overlapFirst = 0
	s.overlapSecond = 0
	s.locFirst = 0
	s.locSecond = 0
	s.lastLocFirst = 0
	s.lastLocSecond = 0
	s.err = nil
}

func (s *SliceRuneScanner) ResetSliceOnly(slice Slice) {
	s.slice = slice
	s.overlap, s.overlapFirst, s.overlapSecond = slice.Overlap()
}

func (s *SliceRuneScanner) ResetWith(slice Slice) {
	s.ResetSliceOnly(slice)
	s.locFirst = 0
	s.locSecond = 0
	s.lastLocFirst = 0
	s.lastLocSecond = 0
	s.err = nil
}

func (s *SliceRuneScanner) UnreadRune() (err error) {
	if s.locFirst <= 0 && s.locSecond <= 0 {
		err = errors.Errorf("nothing to unread")
		return
	}

	if s.lastLocFirst == s.locFirst && s.lastLocSecond == s.locSecond {
		err = errors.New("already unread")
		return
	}

	s.locFirst = s.lastLocFirst
	s.locSecond = s.lastLocSecond

	return
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
	lastLocFirst := s.locFirst
	lastLocSecond := s.locSecond

	switch {
	case firstRemaining > 0:
		if firstRemaining <= s.overlapFirst {
			r, width = utf8.DecodeRune(s.overlap[s.overlapFirst-firstRemaining:])
			s.locFirst += width
			diff := width - firstRemaining

			if diff > 0 {
				s.locSecond += diff
			}
		} else {
			r, width = utf8.DecodeRune(s.slice.First()[s.locFirst:])
			s.locFirst += width
		}

		idxChanged = true

	case len(s.slice.Second())-s.locSecond > 0:
		r, width = utf8.DecodeRune(s.slice.Second()[s.locSecond:])
		s.locSecond += width
		idxChanged = true

	default:
		s.err = io.EOF
		return
	}

	if r == utf8.RuneError {
		s.err = errInvalidRune
		return
	}

	if idxChanged {
		s.lastLocFirst = lastLocFirst
		s.lastLocSecond = lastLocSecond
	}

	ok = idxChanged

	return
}
