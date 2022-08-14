package organize_text

// type metadateiReader struct {
// 	hasSetBaseEtikett bool
// 	within            bool
// }

// func (ot *organizeText) ReadFrom(r1 io.Reader) (n int64, err error) {
// 	r := bufio.NewReader(r1)

// 	ot.etiketten = etikett.NewSet()

// 	within := false
// 	line := 0

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

// 		if !within && s == zettel_formats.MetadateiBoundary {
// 			within = true
// 		} else if within && s != zettel_formats.MetadateiBoundary {
// 			slen := len(s)

// 			if slen < 1 {
// 				continue
// 			}

// 			p := s[0]
// 			v := ""

// 			if slen > 1 {
// 				v = strings.TrimSpace(s[1:])
// 			}

// 			switch p {

// 			case '*':

// 				if v == "" {
// 					continue
// 				}

// 				if err = ot.etiketten.AddString(v); err != nil {
// 					err = errors.Error(err)
// 					return
// 				}

// 			default:
// 				err = errors.Errorf("unsupported verb '%q', '%q'", p, s)
// 				return
// 			}

// 			line += 1

// 		} else if within && s == zettel_formats.MetadateiBoundary {
// 			within = false

// 		} else {
// 			var n1 int64

// 			if n1, err = ot.zettels.ReadFrom(r); err != nil {
// 				err = errors.Error(err)
// 				return
// 			}

// 			n += n1
// 		}
// 	}

// 	return
// }
