package etikett

type Delta struct {
	Added, Removed Set
}

func MakeSetDelta(s1, s2 Set) (d Delta) {
	d.Added = MakeSet()
	d.Removed = *s1.Copy()

	for _, e := range s2 {
		if s1.Contains(e) {
			//zettel had etikett previously
		} else {
			//zettel did not have etikett previously
			d.Added.Add(e)
		}

		d.Removed.Remove(e)
	}

	return
}
