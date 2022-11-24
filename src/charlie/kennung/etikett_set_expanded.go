package kennung

type etikettSetExpanded EtikettSet

func newEtikettSetExpanded(es ...Etikett) etikettSetExpanded {
	return etikettSetExpanded(MakeEtikettSet(es...))
}

func (_ etikettSetExpanded) IsExpanded() bool {
	return true
}
