package kennung

import "code.linenisgreat.com/zit/go/zit/src/alfa/errors"

type IndexedLike struct {
	Int            int
	Kennung        ObjectId
	SchwanzenCount int
	Count          int
}

func MakeIndexed(k IdLike) (i *IndexedLike) {
	i = &IndexedLike{}
	i.ResetWithKennung(k)
	return
}

func (i *IndexedLike) ResetWithKennung(k IdLike) {
	errors.PanicIfError(i.Kennung.SetWithIdLike(k))
}

func (z *IndexedLike) GetInt() int {
	return 0
}

func (z *IndexedLike) GetKennung() IdLike {
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
