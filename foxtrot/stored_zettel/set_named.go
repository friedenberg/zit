package stored_zettel

import "github.com/friedenberg/zit/charlie/hinweis"

type SetNamed map[string]Named

func MakeSetNamed() SetNamed {
	return make(SetNamed)
}

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

// func (s SetNamed) Slice() (slice []string) {
// 	slice = make([]string, len(zs.etikettenToExisting))
// 	i := 0

// 	for e, _ := range zs.etikettenToExisting {
// 		sorted[i] = e
// 		i++
// 	}

// 	sort.Slice(sorted, func(i, j int) bool {
// 		return sorted[i] < sorted[j]
// 	})

// 	return
// }
