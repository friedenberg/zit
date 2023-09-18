package iter2

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

func AddPtrOrReplaceIfGreater[T any, TPtr schnittstellen.Ptr[T]](
	c schnittstellen.MutableSetPtrLike[T, TPtr],
	l schnittstellen.Lessor2[T, TPtr],
	b TPtr,
) (err error) {
	a, ok := c.GetPtr(c.KeyPtr(b))

	if !ok || l.LessPtr(a, b) {
		return c.AddPtr(b)
	}

	return
}
