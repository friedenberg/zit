package box

import (
	"io"
	"unicode/utf8"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type SeqRuneScanner struct {
	Seq
	token_index int
	byte_index  int
}

func (scanner *SeqRuneScanner) Reset() {
	scanner.token_index = 0
	scanner.byte_index = 0
}

func (scanner *SeqRuneScanner) IsFull() bool {
	if scanner == nil {
		return false
	}

	return scanner.token_index == 0 && scanner.byte_index == 0
}

func (scanner *SeqRuneScanner) ReadRune() (r rune, size int, err error) {
	if scanner.token_index == scanner.Len() {
		err = io.EOF
		return
	}

	token := scanner.At(scanner.token_index)
	contents := token.Contents[scanner.byte_index:]

	r, size = utf8.DecodeRune(contents)

	if r == utf8.RuneError && size > 0 {
		err = errors.Errorf("invalid utf8 byte sequence: %q", contents)
		return
	} else if r == utf8.RuneError {
		// should never happen
		err = errors.Errorf("tried to read past end of tokens: %q", contents)
		return
	}

	scanner.byte_index += size

	if len(contents) == size {
		scanner.byte_index = 0
		scanner.token_index += 1
	}

	return
}

func (scanner *SeqRuneScanner) UnreadRune() (err error) {
	if scanner.token_index == 0 && scanner.byte_index == 0 {
		err = errors.Errorf("seq rune scanner fully unread")
		return
	}

	token := scanner.At(scanner.token_index - 1)
	bytes := token.Contents[:scanner.byte_index+1]

	_, size := utf8.DecodeLastRune(bytes)
	scanner.byte_index -= size

	if size == len(bytes) {
		scanner.token_index -= 1
		scanner.byte_index = len(scanner.At(scanner.token_index).Contents) - 1
	}

	return
}
