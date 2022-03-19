package organize_text

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

type assignments struct {
	etikettenToExisting map[string]zettelSet
	etikettenToNew      map[string]newZettelSet
}

func newEtikettToZettels() assignments {
	return assignments{
		etikettenToExisting: make(map[string]zettelSet),
		etikettenToNew:      make(map[string]newZettelSet),
	}
}

func (zs assignments) AddStored(e string, z _NamedZettel) {
	d := z.Zettel.Description()
	// d = fmt.Sprintf("%s %s", z.Sha.String()[:7], d)
	zs.Add(e, z.Hinweis.String(), d)
}

func (zs assignments) Add(e string, h, b string) {
	zs.add(
		e,
		zettel{
			hinweis:     h,
			bezeichnung: b,
		},
	)
}

func (zs assignments) add(e string, z zettel) {
	if _, ok := zs.etikettenToExisting[e]; !ok {
		zs.etikettenToExisting[e] = makeZettelSet()
	}

	zs.etikettenToExisting[e].Add(z)
}

func (zs assignments) addNew(e string, z newZettel) {
	if _, ok := zs.etikettenToNew[e]; !ok {
		zs.etikettenToNew[e] = makeNewZettelSet()
	}

	zs.etikettenToNew[e].Add(z)
}

func (zs assignments) sorted() (sorted []string) {
	sorted = make([]string, len(zs.etikettenToExisting))
	i := 0

	for e, _ := range zs.etikettenToExisting {
		sorted[i] = e
		i++
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	return
}

func (zs assignments) WriteTo(out io.Writer) (n int64, err error) {
	w := _LineFormatNewWriter()

	for _, e := range zs.sorted() {
		ezs := zs.etikettenToExisting[e]

		if e != "" {
			w.WriteLines(fmt.Sprintf("# %s", e))
			w.WriteEmpty()
		}

		for _, z := range ezs.sorted() {
			w.WriteStringers(z)
		}

		w.WriteEmpty()
	}

	n, err = w.WriteTo(out)

	return
}

func (zs *assignments) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	var currentEtikettString string

	lineNo := 0

	for {
		var s string
		s, err = r.ReadString('\n')

		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			err = _Error(err)
			return
		}

		n += int64(len(s))

		s = strings.TrimSuffix(s, "\n")
		slen := len(s)

		if slen < 1 {
			continue
		}

		p := s[0]
		v := ""

		if slen > 1 {
			v = strings.TrimSpace(s[1:])
		}

		switch p {

		case '#':
			currentEtikett := _EtikettNewSet()

			if v == "" {
				currentEtikettString = ""
			} else {
				if err = currentEtikett.Set(v); err != nil {
					err = ErrorRead{
						error:  err,
						line:   lineNo,
						column: 2,
					}

					return
				}

				currentEtikettString = currentEtikett.String()
			}

		case '-':
			var z zettel

			err = z.Set(s)

			if err == nil {
				zs.Add(currentEtikettString, z.hinweis, z.bezeichnung)
			} else {
				var nz newZettel
				var errNz error

				if errNz = nz.Set(s); errNz != nil {
					err = ErrorRead{
						error:  err,
						line:   lineNo,
						column: 2,
					}

					return
				} else {
					zs.addNew(currentEtikettString, nz)
				}
			}

		default:
			err = ErrorRead{
				error:  _Errorf("unsupported verb %q, %q", p, s),
				line:   lineNo,
				column: 0,
			}

			return
		}

		lineNo++
	}

	return
}

func (a assignments) Copy() (b assignments) {
	b = newEtikettToZettels()

	for k, v := range a.etikettenToExisting {
		for z, _ := range v {
			b.Add(k, z.hinweis, z.bezeichnung)
		}
	}

	return
}
