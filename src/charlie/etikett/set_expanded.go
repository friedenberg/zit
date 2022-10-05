package etikett

type setExpanded Set

func newSetExpanded() setExpanded {
	return setExpanded(MakeSet())
}

func (_ setExpanded) IsExpanded() bool {
	return true
}
