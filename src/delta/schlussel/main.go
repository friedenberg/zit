package schlussel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/ohio"
)

//go:generate stringer -type=Schlussel
type Schlussel byte

const (
	Unknown                      = Schlussel(iota)
	ContentLength                = 'C'
	Sigil                        = 'S'
	Akte                         = 'A'
	Bezeichnung                  = 'B'
	Etikett                      = 'E'
	Gattung                      = 'G'
	Kennung                      = 'K'
	Komment                      = 'k'
	Tai                          = 'T'
	Typ                          = 't'
	Mutter                       = 'M'
	Sha                          = 's'
	VerzeichnisseArchiviert      = 'a'
	VerzeichnisseEtikettImplicit = 'I'
	VerzeichnisseEtikettExpanded = 'e'
)

var ErrInvalid = errors.New("invalid schlussel")

func (s *Schlussel) Reset() {
	*s = 0
}

func (s *Schlussel) ReadByte() (byte, error) {
	return byte(*s), nil
}

func (s *Schlussel) WriteTo(w io.Writer) (n int64, err error) {
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

func (s *Schlussel) WriteByte(b byte) (err error) {
	*s = Schlussel(b)

	return
}

func (s *Schlussel) ReadFrom(r io.Reader) (n int64, err error) {
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
