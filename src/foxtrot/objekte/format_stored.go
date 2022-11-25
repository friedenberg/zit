package objekte

//type FormatStored struct {
//	arf metadatei_io.AkteIOFactory
//}

////TODO add support for metadata like with Zettel
//func MakeFormat(arf metadatei_io.AkteIOFactory) *FormatStored {
//	return &FormatStored{
//		arf: arf,
//	}
//}

//func (f FormatStored) ReadFormat(r1 io.Reader, t StoredLike) (n int64, err error) {
//	r := bufio.NewReader(r1)

//	for {
//		var lineOriginal string
//		lineOriginal, err = r.ReadString('\n')

//		if err == io.EOF {
//			err = nil
//			break
//		} else if err != nil {
//			return
//		}

//		// line := strings.TrimSpace(lineOriginal)
//		line := lineOriginal

//		n += int64(len(lineOriginal))

//		loc := strings.Index(line, " ")

//		if line == "" {
//			//TODO this should be cleaned up
//		}

//		var g gattung.Gattung

//		switch {
//		case line == "":
//			err = errors.Errorf("found empty line: %q", lineOriginal)
//			return

//		case line != "" && loc == -1:
//			if err = g.Set(line[:loc]); err != nil {
//				err = errors.Errorf("%s: %s", err, line[:loc])
//				return
//			}

//			if g != t.Gattung() {
//				err = errors.Errorf(
//					"expected objekte to have gattung '%s' but got '%s'",
//					gattung.Typ,
//					g,
//				)

//				return
//			}

//			continue

//		case lineOriginal == "\n" && loc == -1:
//			continue

//		case loc == -1:
//			err = errors.Errorf("expected at least one space, but found none: %q", lineOriginal)
//			return
//		}

//		if err = g.Set(line[:loc]); err != nil {
//			err = errors.Errorf("%s: %s", err, line[:loc])
//			return
//		}

//		v := line[loc+1:]

//		switch g {
//		case gattung.Akte:
//			if err = t.setSha(v); err != nil {
//				err = errors.Wrap(err)
//				return
//			}

//		default:
//			err = errors.Errorf("unsupported gattung: %s", g)
//			return
//		}
//	}

//	return
//}

//func (f FormatStored) WriteFormat(w1 io.Writer, t StoredLike) (n int64, err error) {
//	w := line_format.NewWriter()

//	w.WriteFormat("%s", t.Gattung())
//	w.WriteFormat("%s %s", gattung.Akte, t.Sha())

//	if n, err = w.WriteTo(w1); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	return
//}
