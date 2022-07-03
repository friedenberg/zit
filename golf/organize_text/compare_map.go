package organize_text

import "strings"

type etikettZettel struct {
	etikett, Hinweis string
}

type compareMap map[etikettZettel]bool

func (in organizeText) ToCompareMap() (out compareMap) {
	out = make(compareMap)

	for e, zs := range in.ZettelsExisting() {
		for z, _ := range zs {
			// individual etiketten
			for _, e1 := range strings.Split(e, ", ") {
				// root etiketten have an empty string representation
				if e1 != "" {
					out[etikettZettel{etikett: e1, Hinweis: z.Hinweis}] = true
				}
			}

			// root etiketten
			for _, e2 := range in.Etiketten() {
				out[etikettZettel{etikett: e2.String(), Hinweis: z.Hinweis}] = true
			}
		}
	}

	return
}
