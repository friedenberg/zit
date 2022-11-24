package kennung

type setExpanded EtikettSet

func newSetExpanded(es ...Etikett) setExpanded {
	return setExpanded(MakeSet(es...))
}

func (_ setExpanded) IsExpanded() bool {
	return true
}
