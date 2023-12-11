package catgut

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type RingBufferScanner struct {
	sliceScanner SliceRuneScanner
	rb           *RingBuffer
}

func MakeRingBufferScanner(rb *RingBuffer) *RingBufferScanner {
	return &RingBufferScanner{
		rb: rb,
	}
}

var ErrNoMatch = errors.New("no match")

func (s *RingBufferScanner) AdvanceToFirstMatch(
	mf func(rune) bool,
) (match []byte, err error) {
	s.sliceScanner.ResetWith(s.rb.PeekReadable())

	if s.rb.PeekReadable().Len() == 0 {
		err = ErrBufferEmpty
		return
	}

	startMatchOffset := -1

	var n int
	var r rune
	var w int
	var ok bool

LOOP:
	for {
		r, w, ok = s.sliceScanner.Scan()

		if !ok {
			if err = s.sliceScanner.Error(); err != nil {
				if err == io.EOF {
					err = ErrBufferEmpty
				}
			}

			break
		}

		currentMatch := mf(r)

		switch {
		case currentMatch:
			if startMatchOffset < 0 {
				startMatchOffset = n
			}

		case !currentMatch && startMatchOffset >= 0:
			break LOOP
		}

		n += w
	}

	if startMatchOffset == -1 {
		s.rb.AdvanceRead(n)
		err = ErrNoMatch
		return
	}

	match = s.rb.PeekReadable().Slice(startMatchOffset, n).Bytes()
	s.rb.AdvanceRead(n)

	return
}
