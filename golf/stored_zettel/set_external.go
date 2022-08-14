package stored_zettel

import (
	"github.com/friedenberg/zit/delta/hinweis"
)

type SetExternal map[string]External

func MakeSetExternal() SetExternal {
	return make(SetExternal)
}

func (s SetExternal) Hinweisen() (h []hinweis.Hinweis) {
	h = make([]hinweis.Hinweis, 0, len(s))

	for _, z := range s {
		h = append(h, z.Hinweis)
	}

	return
}

func (s SetExternal) HinweisStrings() (h []string) {
	h = make([]string, 0, len(s))

	for i, _ := range s {
		h = append(h, i)
	}

	return
}

func (s SetExternal) Paths() (p []string) {
	p = make([]string, 0, len(s))

	for _, z := range s {
		p = append(p, z.Path)
	}

	return
}

// func (s SetExternal) Slice() (slice []string) {
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
