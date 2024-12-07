package quiter

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func AddClone[E any, EPtr interface {
	*E
	ResetWithPtr(*E)
}](
	c interfaces.Adder[EPtr],
) interfaces.FuncIter[EPtr] {
	return func(e EPtr) (err error) {
		var e1 E
		EPtr(&e1).ResetWithPtr((*E)(e))
		c.Add(&e1)
		return
	}
}

func AddClonePool[E any, EPtr interfaces.Ptr[E]](
	s interfaces.AdderPtr[E, EPtr],
	p interfaces.Pool[E, EPtr],
	r interfaces.Resetter2[E, EPtr],
	b EPtr,
) (err error) {
	a := p.Get()
	r.ResetWith(a, b)
	return s.AddPtr(a)
}

func MakeAddClonePoolFunc[E any, EPtr interfaces.Ptr[E]](
	s interfaces.AdderPtr[E, EPtr],
	p interfaces.Pool[E, EPtr],
	r interfaces.Resetter2[E, EPtr],
) interfaces.FuncIter[EPtr] {
	return MakeSyncSerializer(func(e EPtr) (err error) {
		return AddClonePool(s, p, r, e)
	})
}

func ExpandAndAddString[E any, EPtr interfaces.SetterPtr[E]](
	c interfaces.Adder[E],
	expander func(string) (string, error),
	v string,
) (err error) {
	if expander != nil {
		v1 := v

		if v1, err = expander(v); err != nil {
			err = nil
			v1 = v
		}

		v = v1
	}

	if err = AddString[E, EPtr](c, v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type AddGetKeyer[E interfaces.Lessor[E]] interface {
	interfaces.Adder[E]
	Get(string) (E, bool)
	Key(E) string
}

func AddIfGreater[E interfaces.Lessor[E]](
	c AddGetKeyer[E],
	e E,
) (ok bool) {
	k := c.Key(e)
	var old E

	if old, ok = c.Get(k); !ok || old.Less(e) {
		c.Add(e)
	}

	return
}

func AddOrReplaceIfGreater[T interface {
	interfaces.Stringer
	interfaces.ValueLike
	interfaces.Lessor[T]
}](
	c interfaces.MutableSetLike[T],
	b T,
) (shouldAdd bool, err error) {
	a, ok := c.Get(c.Key(b))

	// 	if ok {
	// 		log.Debug().Print("less:", a.Less(b))
	// 	} else {
	// 		log.Debug().Print("ok:", ok)
	// 	}

	shouldAdd = !ok || a.Less(b)

	if shouldAdd {
		err = c.Add(b)
	}

	return
}

func AddStringPtr[E any, EPtr interfaces.SetterPtr[E]](
	c interfaces.AdderPtr[E, EPtr],
	p interfaces.Pool[E, EPtr],
	v string,
) (err error) {
	e := p.Get()

	if err = e.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.AddPtr(e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func AddString[E any, EPtr interfaces.SetterPtr[E]](
	c interfaces.Adder[E],
	v string,
) (err error) {
	var e E

	if err = EPtr(&e).Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.Add(e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
