package ids

type ObjectIdGetter interface {
	GetObjectId() *ObjectId
}

type ObjectIdKeyer[
	T ObjectIdGetter,
] struct{}

func (sk ObjectIdKeyer[T]) GetKey(e T) string {
	return e.GetObjectId().String()
}
