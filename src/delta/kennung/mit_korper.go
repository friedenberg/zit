package kennung

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
)

type MitKorperLike[T any] interface {
	KennungLike[T]
	Kopf() string
	Schwanz() string
}

type MitKorper[T MitKorperLike[T], T1 KennungLikePtr[T]] Kennung[T, T1]

func (mk MitKorper[T, T1]) AlignedParts(kopf, schwanz int) (string, string) {
	parts := mk.Parts()

	diffKopf := kopf - len(parts[0])
	if diffKopf > 0 {
		parts[0] = strings.Repeat(" ", diffKopf) + parts[0]
	}

	diffSchwanz := schwanz - len(parts[1])
	if diffSchwanz > 0 {
		parts[1] = parts[1] + strings.Repeat(" ", diffSchwanz)
	}

	return parts[0], parts[1]
}

func (mk MitKorper[T, T1]) Aligned(kopf, schwanz int) string {
	p1, p2 := mk.AlignedParts(kopf, schwanz)
	return fmt.Sprintf("%s/%s", p1, p2)
}

func (mk MitKorper[T, T1]) Parts() [2]string {
	return [2]string{mk.Kopf(), mk.Schwanz()}
}

func (mk MitKorper[T, T1]) Kopf() string {
	return mk.value.Kopf()
}

func (mk MitKorper[T, T1]) Schwanz() string {
	return mk.value.Schwanz()
}

func (mk MitKorper[T, T1]) String() string {
	v := Kennung[T, T1](mk).String()
	// if v == "/" {
	// 	errors.Err().Caller(1, "empty hinweis")
	// }
	return v
}

func (mk *MitKorper[T, T1]) Set(v string) (err error) {
	if err = (*Kennung[T, T1])(mk).Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (mk MitKorper[T, T1]) Equals(b MitKorper[T, T1]) bool {
	return Kennung[T, T1](mk).Equals(Kennung[T, T1](b))
}

func (mk MitKorper[T, T1]) GetGattung() schnittstellen.Gattung {
	return Kennung[T, T1](mk).GetGattung()
}

func (mk *MitKorper[T, T1]) Reset() {
	(*Kennung[T, T1])(mk).Reset()
}

func (mk *MitKorper[T, T1]) ResetWith(b MitKorper[T, T1]) {
	(*Kennung[T, T1])(mk).ResetWith((Kennung[T, T1])(b))
}

func (mk MitKorper[T, T1]) Less(b MitKorper[T, T1]) bool {
	return (Kennung[T, T1])(mk).Less((Kennung[T, T1])(b))
}

func (mk MitKorper[T, T1]) GetSha() sha.Sha {
	return (Kennung[T, T1])(mk).GetSha()
}

func (mk MitKorper[T, T1]) GetSigil() Sigil {
	return (Kennung[T, T1])(mk).GetSigil()
}

func (t MitKorper[T, T1]) GobEncode() (by []byte, err error) {
	return (Kennung[T, T1])(t).GobEncode()
}

func (t *MitKorper[T, T1]) GobDecode(b []byte) (err error) {
	return (*Kennung[T, T1])(t).GobDecode(b)
}
