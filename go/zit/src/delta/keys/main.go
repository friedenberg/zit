package keys

import (
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
)

type Key byte

const (
	Unknown                      = Key(iota)
	ContentLength                = 'C'
	Sigil                        = 'S'
	Blob                         = 'A'
	Description                  = 'B'
	Tag                          = 'E'
	Genre                        = 'G'
	ObjectId                     = 'K'
	Comment                      = 'k'
	Tai                          = 'T'
	Type                         = 't'
	MutterMetadateiMutterKennung = 'M'
	ShaMetadateiMutterKennung    = 's'
	ShaMetadateiSansTai          = 'n'
	ShaMetadatei                 = 'm'
	VerzeichnisseArchiviert      = 'a'
	VerzeichnisseEtikettImplicit = 'I'
	VerzeichnisseEtikettExpanded = 'e'
	VerzeichnisseEtiketten       = 'x'
	VerzeichnisseEtiketten2      = 'y'
)

var ErrInvalid = errors.New("invalid key")

func (s Key) String() string {
	return fmt.Sprintf("%c", byte(s))
}

func (s *Key) Reset() {
	*s = 0
}

func (s *Key) ReadByte() (byte, error) {
	return byte(*s), nil
}

func (s *Key) WriteTo(w io.Writer) (n int64, err error) {
	b := [1]byte{byte(*s)}
	var n1 int
	n1, err = ohio.WriteAllOrDieTrying(w, b[:])
	n += int64(n1)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	return
}

func (s *Key) WriteByte(b byte) (err error) {
	*s = Key(b)

	return
}

func (s *Key) ReadFrom(r io.Reader) (n int64, err error) {
	var b [1]byte
	var n1 int
	n1, err = ohio.ReadAllOrDieTrying(r, b[:])
	n += int64(n1)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	err = s.WriteByte(b[0])

	return
}
