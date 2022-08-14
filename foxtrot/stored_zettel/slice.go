package stored_zettel

import "github.com/friedenberg/zit/delta/hinweis"

type SliceNamed []Named

func (s SliceNamed) Hinweisen() (h []hinweis.Hinweis) {
	h = make([]hinweis.Hinweis, 0, len(s))

	for _, z := range s {
		h = append(h, z.Hinweis)
	}

	return
}

func (s SliceNamed) HinweisStrings() (h []string) {
	h = make([]string, 0, len(s))

	for _, z := range s {
		h = append(h, z.Hinweis.String())
	}

	return
}

// func (s SliceNamed) Sorted() (slice SliceNamed) {
// 	sorted = make(SliceNamed, len(s))

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
