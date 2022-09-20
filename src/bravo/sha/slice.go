package sha

type Slice []Sha

func MakeSlice(c int) Slice {
	return make([]Sha, 0, c)
}

func (s *Slice) Append(sh ...Sha) {
	*s = append(*s, sh...)
}
