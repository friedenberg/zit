package collections

type ValueSet2[E ValueSetElement] struct {
	SetLike[E]
}

func (v ValueSet2[E]) String() string {
	return String[E](v.SetLike)
}
