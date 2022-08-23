package stored_zettel

import "github.com/friedenberg/zit/src/delta/hinweis"

type SetTransacted map[string]Transacted

func NewSetTransacted() *SetTransacted {
	s := MakeSetTransacted()
	return &s
}

func MakeSetTransacted() SetTransacted {
	return make(SetTransacted)
}

func (s *SetTransacted) Add(z Transacted) {
	(*s)[z.Hinweis.String()] = z
}

func (s SetTransacted) Get(h hinweis.Hinweis) (z Transacted, ok bool) {
	z, ok = s[h.String()]
	return
}

func (a SetTransacted) Merge(b SetTransacted) {
	for _, z := range b {
		a.Add(z)
	}
}

func (a SetTransacted) Contains(z Transacted) bool {
	_, ok := a[z.Hinweis.String()]
	return ok
}

func (s SetTransacted) Hinweisen() (h []hinweis.Hinweis) {
	h = make([]hinweis.Hinweis, 0, len(s))

	for _, z := range s {
		h = append(h, z.Hinweis)
	}

	return
}

func (s SetTransacted) HinweisStrings() (h []string) {
	h = make([]string, 0, len(s))

	for i, _ := range s {
		h = append(h, i)
	}

	return
}
