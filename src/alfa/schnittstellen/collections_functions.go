package schnittstellen

func ImmutableCloneCollection[E any, B interface {
	Collection[E]
	ImmutableCloner[B]
}](o B) Collection[E] {
	return o.ImmutableClone()
}
