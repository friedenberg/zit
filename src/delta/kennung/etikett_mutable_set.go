package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type EtikettMutableSet = schnittstellen.MutableSet[Etikett]

func MakeEtikettMutableSet(hs ...Etikett) EtikettMutableSet {
	return EtikettMutableSet(
		collections.MakeMutableSet[Etikett](
			(Etikett).String,
			hs...,
		),
	)
}

func AddNormalized(es EtikettMutableSet, e Etikett) {
	e.Expanded(ExpanderRight).Each(es.Add)
	es.Add(e)

	c := es.ImmutableClone()
	es.Reset()
	WithRemovedCommonPrefixes(c).Each(es.Add)
}

func RemovePrefixes(es EtikettMutableSet, needle Etikett) {
	for _, haystack := range es.Elements() {
		// TODO make more efficient
		if strings.HasPrefix(haystack.String(), needle.String()) {
			es.Del(haystack)
		}
	}
}

func Withdraw(s1 EtikettMutableSet, e Etikett) (s2 EtikettSet) {
	s3 := MakeEtikettMutableSet()

	for _, e1 := range s1.Elements() {
		if e1.Contains(e) {
			s3.Add(e1)
		}
	}

	s3.Each(s1.Del)
	s2 = s3.ImmutableClone()

	return
}
