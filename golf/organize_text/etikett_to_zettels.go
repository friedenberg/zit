package organize_text

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

type etikettToZettels map[string]zettelSet

func newEtikettToZettels() etikettToZettels {
	return make(etikettToZettels)
}

func (zs etikettToZettels) AddStored(e string, z _NamedZettel) {
	d := z.Zettel.Description()
	// d = fmt.Sprintf("%s %s", z.Sha.String()[:7], d)
	zs.Add(e, z.Hinweis.String(), d)
}

func (zs etikettToZettels) Add(e string, h, b string) {
	zs.add(
		e,
		zettel{
			hinweis:     h,
			bezeichnung: b,
		},
	)
}

func (zs etikettToZettels) add(e string, z zettel) {
	if _, ok := zs[e]; !ok {
		zs[e] = newZettelSet()
	}

	zs[e].Add(z)
}

func (zs etikettToZettels) sorted() (sorted []string) {
	sorted = make([]string, len(zs))
	i := 0

	for e, _ := range zs {
		sorted[i] = e
		i++
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	return
}

func (zs etikettToZettels) WriteTo(out io.Writer) (n int64, err error) {
	w := _LineFormatNewWriter()

	for _, e := range zs.sorted() {
		ezs := zs[e]

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

func (zs *etikettToZettels) ReadFrom(r1 io.Reader) (n int64, err error) {
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

			if err = z.Set(s); err != nil {
				err = ErrorRead{
					error:  err,
					line:   lineNo,
					column: 2,
				}

				return
			}

			zs.Add(currentEtikettString, z.hinweis, z.bezeichnung)

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

func (a etikettToZettels) Copy() (b etikettToZettels) {
	b = newEtikettToZettels()

	for k, v := range a {
		for z, _ := range v {
			b.Add(k, z.hinweis, z.bezeichnung)
		}
	}

	return
}
