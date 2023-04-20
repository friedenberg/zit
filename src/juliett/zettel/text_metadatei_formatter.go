package zettel

// type Metadatei struct {
// 	Objekte
// 	AktePath string
// }

// type TextMetadateiFormatter struct {
// 	DoNotWriteEmptyBezeichnung bool
// 	IncludeAkteSha             bool
// }

// func (f *TextMetadateiFormatter) Format(w1 io.Writer, m *Metadatei) (n int64, err error) {
// 	errors.TodoP3("turn *Objekte into an interface to allow zettel_external to use this")

// 	w := format.NewLineWriter()

// 	if m.Metadatei.Bezeichnung.String() != "" || !f.DoNotWriteEmptyBezeichnung {
// 		w.WriteLines(
// 			fmt.Sprintf("# %s", m.Metadatei.Bezeichnung),
// 		)
// 	}

// 	for _, e := range collections.SortedValues(m.Metadatei.Etiketten) {
// 		errors.TodoP3("determine how to handle this")

// 		if e.IsEmpty() {
// 			continue
// 		}

// 		w.WriteFormat("- %s", e)
// 	}

// 	switch {
// 	case m.AktePath != "":
// 		w.WriteLines(
// 			fmt.Sprintf("! %s", m.AktePath),
// 		)

// 	case f.IncludeAkteSha:
// 		w.WriteLines(
// 			fmt.Sprintf("! %s.%s", m.Metadatei.AkteSha, m.GetTyp()),
// 		)

// 	default:
// 		w.WriteLines(
// 			fmt.Sprintf("! %s", m.GetTyp()),
// 		)
// 	}

// 	if n, err = w.WriteTo(w1); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }
