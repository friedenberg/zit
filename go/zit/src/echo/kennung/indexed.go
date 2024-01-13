package kennung

import "github.com/friedenberg/zit/src/alfa/errors"

type IndexedLike struct {
	Int            int
	Kennung        Kennung2
	SchwanzenCount int
	Count          int
}

func MakeIndexed(k Kennung) (i *IndexedLike) {
	i = &IndexedLike{}
	i.ResetWithKennung(k)
	return
}

func (i *IndexedLike) ResetWithKennung(k Kennung) {
	errors.PanicIfError(i.Kennung.SetWithKennung(k))
}

func (z *IndexedLike) GetInt() int {
	return 0
}

func (z *IndexedLike) GetKennung() Kennung {
	return &z.Kennung
}

func (k *IndexedLike) GetSchwanzenCount() int {
	return k.SchwanzenCount
}

func (k *IndexedLike) GetCount() int {
	return k.Count
}

func (z *IndexedLike) Reset() {
	z.Kennung.Reset()
	z.SchwanzenCount = 0
	z.Count = 0
}
