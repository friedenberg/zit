package collections

type ErrEmptyKey[T any] struct {
	Element T
}

func (e ErrEmptyKey[T]) Error() string {
	return "empty key"
}

func (e ErrEmptyKey[T]) Is(target error) (ok bool) {
	_, ok = target.(ErrEmptyKey[T])
	return
}

type ErrDoNotRepool struct{}

func (e ErrDoNotRepool) Error() string {
	return "should not repool this element"
}

func (e ErrDoNotRepool) Is(target error) (ok bool) {
	_, ok = target.(ErrDoNotRepool)
	return
}
