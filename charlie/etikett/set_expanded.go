package etikett

type setExpanded Set

func newSetExpanded() setExpanded {
	return make(setExpanded)
}

func (_ setExpanded) IsExpanded() bool {
	return true
}
