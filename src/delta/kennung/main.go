package kennung

import (
	"crypto/sha256"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type KennungLike[T any] interface {
	schnittstellen.GattungGetter
	schnittstellen.Value
	schnittstellen.Equatable[T]
}

type KennungLikePtr[T schnittstellen.Value] interface {
	schnittstellen.ValuePtr[T]
	schnittstellen.Resetable[T]
}

type Kennung[T KennungLike[T], T1 KennungLikePtr[T]] struct {
	value T
}

func makeKennung[T KennungLike[T], T1 KennungLikePtr[T]](
	v string,
) (k Kennung[T, T1], err error) {
	k.value = *T1(new(T))

	if err = k.Set(v); err != nil {
		err = errors.Wrap(err)
	}

	return
}

func (e Kennung[T, T1]) GetSha() sha.Sha {
	hash := sha256.New()
	sr := strings.NewReader(e.String())

	if _, err := io.Copy(hash, sr); err != nil {
		errors.PanicIfError(err)
	}

	return sha.FromHash(hash)
}

func (e Kennung[T, T1]) GetGattung() schnittstellen.Gattung {
	return e.value.GetGattung()
}

func (e Kennung[T, T1]) String() string {
	return e.value.String()
}

func (e *Kennung[T, T1]) Set(v string) (err error) {
	v1 := strings.ToLower(v)
	v3 := strings.TrimSpace(v1)

	if v3 == "" {
		err = errors.Wrap(gattung.ErrEmptyKennung{})
		return
	}

	if err = T1(&e.value).Set(v3); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e Kennung[T, T1]) Len() int {
	return len(e.value.String())
}

func (a Kennung[T, T1]) Includes(b Kennung[T, T1]) bool {
	return b.Contains(a)
}

func (a Kennung[T, T1]) Contains(b Kennung[T, T1]) bool {
	if b.Len() > a.Len() {
		return false
	}

	return strings.HasPrefix(a.value.String(), b.value.String())
}

func (a *Kennung[T, T1]) Reset() {
	var a1 T
	a.value = a1
}

func (a *Kennung[T, T1]) ResetWith(b Kennung[T, T1]) {
	a.value = b.value
}

func (a Kennung[T, T1]) Equals(b Kennung[T, T1]) bool {
	return a.value.Equals(b.value)
}

func (a Kennung[T, T1]) Less(b Kennung[T, T1]) bool {
	return a.value.String() < b.value.String()
}

func (a *Kennung[T, T1]) LeftSubtract(b Kennung[T, T1]) (c Kennung[T, T1], err error) {
	c.value = *T1(new(T))

	if err = c.Set(strings.TrimPrefix(a.String(), b.String())); err != nil {
		err = errors.Wrapf(err, "'%s' - '%s'", a, b)
		return
	}

	return
}

func (a Kennung[T, T1]) IsEmpty() bool {
	return a.Len() == 0
}

func (e Kennung[T, T1]) Expanded(
	exes ...Expander,
) (out collections.ValueSet[Kennung[T, T1], *Kennung[T, T1]]) {
	expanded := collections.MakeMutableValueSet[Kennung[T, T1], *Kennung[T, T1]]()

	for _, ex := range exes {
		ex.Expand(collections.MakeFuncSetString[Kennung[T, T1], *Kennung[T, T1]](expanded), e.String())
	}

	out = expanded.Copy()

	return
}

func (t Kennung[T, T1]) MarshalText() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Kennung[T, T1]) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t Kennung[T, T1]) GobEncode() (by []byte, err error) {
	by = []byte(t.String())
	return
}

func (t *Kennung[T, T1]) GobDecode(b []byte) (err error) {
	if err = t.Set(string(b)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
