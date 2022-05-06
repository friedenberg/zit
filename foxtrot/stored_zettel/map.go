package stored_zettel

import "github.com/friedenberg/zit/charlie/hinweis"

type SetNamed map[string]Named

func (s SetNamed) Hinweisen() (h []hinweis.Hinweis) {
	h = make([]hinweis.Hinweis, 0, len(s))

	for _, z := range s {
		h = append(h, z.Hinweis)
	}

	return
}

func (s SetNamed) HinweisStrings() (h []string) {
	h = make([]string, 0, len(s))

	for i, _ := range s {
		h = append(h, i)
	}

	return
}
