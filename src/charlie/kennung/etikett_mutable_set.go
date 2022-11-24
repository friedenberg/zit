package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/bravo/collections"
)

type EtikettMutableSet = collections.MutableValueSet[Etikett, *Etikett]

func MakeMutableSet(hs ...Etikett) EtikettMutableSet {
	return EtikettMutableSet(collections.MakeMutableValueSet[Etikett, *Etikett](hs...))
}

func AddNormalized(es EtikettMutableSet, e Etikett) {
	e.Expanded(ExpanderRight).Each(es.Add)
	es.Add(e)

	es.Reset(WithRemovedCommonPrefixes(es.Copy()))
}

func RemovePrefixes(es EtikettMutableSet, needle Etikett) {
	for _, haystack := range es.Elements() {
		//TODO make more efficient
		if strings.HasPrefix(haystack.String(), needle.String()) {
			es.Del(haystack)
		}
	}
}

func Withdraw(s1 EtikettMutableSet, e Etikett) (s2 EtikettSet) {
	s3 := MakeMutableSet()

	for _, e1 := range s1.Elements() {
		if e1.Contains(e) {
			s3.Add(e1)
		}
	}

	s3.Each(s1.Del)
	s2 = s3.Copy()

	return
}
