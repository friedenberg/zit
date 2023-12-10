package catgut

import (
	"io"
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

func (s *RingBufferScanner) AdvanceToFirstMatch(
	mf func(rune) bool,
) (match []byte, err error) {
	s.sliceScanner.ResetWith(s.rb.PeekReadable())

	if s.rb.PeekReadable().Len() == 0 {
		err = ErrBufferEmpty
		return
	}

	offset := 0
	startedMatch := false
	endedMatch := false
	startMatchOffset := -1

LOOP:
	for {
		r, w, okScan := s.sliceScanner.Scan()

		if !okScan {
			if err = s.sliceScanner.Error(); err != nil {
				if err == io.EOF {
					err = ErrBufferEmpty
				}
			}
		}

		offset += w
		currentMatch := mf(r)

		switch {
		case currentMatch:
			if startMatchOffset < 0 {
				startMatchOffset = offset - w
			}

			match = s.rb.data[s.rb.rIdx+startMatchOffset : s.rb.rIdx+offset]
			startedMatch = true

		case !currentMatch && startedMatch:
			endedMatch = true
			break LOOP
		}

		if !okScan {
			break
		}
	}

	if endedMatch {
		offset--
	}

	s.rb.AdvanceRead(offset)

	return
}
