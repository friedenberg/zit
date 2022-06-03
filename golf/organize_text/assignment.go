package organize_text

import (
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type assignment struct {
	etiketten etikett.Set
	named     stored_zettel.SetNamed
	unnamed   map[string]bool
}

// func (a assignment) WriteTo(out io.Writer) (n int64, err error) {
// 	w := line_format.NewWriter()

// 	for _, e := range a.etiketten.Sorted() {
// 		ezs := zs.etikettenToExisting[e]

// 		if e != "" {
// 			w.WriteLines(fmt.Sprintf("# %s", e))
// 			w.WriteEmpty()
// 		}

// 		for _, z := range ezs.sorted() {
// 			w.WriteStringers(z)
// 		}

// 		w.WriteEmpty()
// 	}

// 	n, err = w.WriteTo(out)

// 	return
// }

// func (a *assignment) ReadFrom(r1 io.Reader) (n int64, err error) {
// 	r := bufio.NewReader(r1)

// 	var currentEtikettString string

// 	lineNo := 0

// 	for {
// 		var s string
// 		s, err = r.ReadString('\n')

// 		if err == io.EOF {
// 			err = nil
// 			break
// 		}

// 		if err != nil {
// 			err = errors.Error(err)
// 			return
// 		}

// 		n += int64(len(s))

// 		s = strings.TrimSuffix(s, "\n")
// 		slen := len(s)

// 		if slen < 1 {
// 			continue
// 		}

// 		p := s[0]
// 		v := ""

// 		if slen > 1 {
// 			v = strings.TrimSpace(s[1:])
// 		}

// 		switch p {

// 		case '#':
// 			currentEtikett := etikett.NewSet()

// 			if v == "" {
// 				currentEtikettString = ""
// 			} else {
// 				if err = currentEtikett.Set(v); err != nil {
// 					err = ErrorRead{
// 						error:  err,
// 						line:   lineNo,
// 						column: 2,
// 					}

// 					return
// 				}

// 				currentEtikettString = currentEtikett.String()
// 			}

// 		case '-':
// 			var z zettel

// 			err = z.Set(s)

// 			if err == nil {
// 				zs.Add(currentEtikettString, z.hinweis, z.bezeichnung)
// 			} else {
// 				var nz newZettel
// 				var errNz error

// 				if errNz = nz.Set(s); errNz != nil {
// 					err = ErrorRead{
// 						error:  err,
// 						line:   lineNo,
// 						column: 2,
// 					}

// 					return
// 				} else {
// 					zs.addNew(currentEtikettString, nz)
// 				}
// 			}

// 		default:
// 			err = ErrorRead{
// 				error:  errors.Errorf("unsupported verb %q, %q", p, s),
// 				line:   lineNo,
// 				column: 0,
// 			}

// 			return
// 		}

// 		lineNo++
// 	}

// 	return
// }

// func (a assignments) Copy() (b assignments) {
// 	b = newEtikettToZettels()

// 	for k, v := range a.etikettenToExisting {
// 		for z, _ := range v {
// 			b.Add(k, z.hinweis, z.bezeichnung)
// 		}
// 	}

// 	return
// }
