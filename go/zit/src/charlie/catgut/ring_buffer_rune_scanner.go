package catgut

import (
	"io"
	"unicode/utf8"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
)

func MakeRingBufferRuneScanner(rb *RingBuffer) (s *RingBufferRuneScanner) {
	s = &RingBufferRuneScanner{}
	s.ResetWith(rb)

	return
}

type RingBufferRuneScanner struct {
	rb *RingBuffer
	SliceRuneScanner
}

func (s *RingBufferRuneScanner) Reset() {
	s.SliceRuneScanner.Reset()
}

func (s *RingBufferRuneScanner) ResetWith(rb *RingBuffer) {
	s.rb = rb
	s.SliceRuneScanner.ResetWith(rb.PeekReadable())
}

func (s *RingBufferRuneScanner) UnreadRune() (err error) {
	return s.SliceRuneScanner.UnreadRune()
}

func (s *RingBufferRuneScanner) ReadRune() (r rune, size int, err error) {
	if s.Remaining() <= utf8.UTFMax {
		var n int64

		if n, err = s.rb.Fill(); err != nil {
			if err == io.EOF {
				if n > 0 || s.Remaining() > 0 {
					err = nil
				} else {
					return
				}
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		s.ResetSliceOnly(s.rb.PeekReadable())
	}

	return s.SliceRuneScanner.ReadRune()
}
