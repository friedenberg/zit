package organize_text

import "strings"

type EtikettZettel struct {
	Etikett, Hinweis string
}

type CompareMap map[EtikettZettel]bool

func (in organizeText) ToCompareMap() (out CompareMap) {
	out = make(CompareMap)

	for e, zs := range in.ZettelsExisting() {
		for z, _ := range zs {
			// individual etiketten
			for _, e1 := range strings.Split(e, ", ") {
				// root etiketten have an empty string representation
				if e1 != "" {
					out[EtikettZettel{Etikett: e1, Hinweis: z.Hinweis}] = true
				}
			}

			// root etiketten
			for _, e2 := range in.Etiketten() {
				out[EtikettZettel{Etikett: e2.String(), Hinweis: z.Hinweis}] = true
			}
		}
	}

	return
}
