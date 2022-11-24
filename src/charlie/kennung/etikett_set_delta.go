package kennung

// TODO rename
type Delta struct {
	Added, Removed EtikettSet
}

// TODO rename
func MakeSetDelta(s1, s2 EtikettSet) (d Delta) {
	added := MakeEtikettMutableSet()
	removed := s1.MutableCopy()

	for _, e := range s2.Elements() {
		if s1.Contains(e) {
			//zettel had etikett previously
		} else {
			//zettel did not have etikett previously
			added.Add(e)
		}

		removed.Del(e)
	}

	d.Added = added.Copy()
	d.Removed = removed.Copy()

	return
}
