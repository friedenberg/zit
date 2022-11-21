package zettel

//type DecoderText struct {
//	in                         io.Reader
//	DoNotWriteEmptyBezeichnung bool
//	TypError                   error
//	awf                        AkteWriterFactory
//}

//func (f DecoderText) Decode(z *Zettel) (n int64, err error) {
//	r := bufio.NewReader(f.in)

//	var aktePath string

//	z.Etiketten = etikett.MakeSet()

//	state := &textStateRead{
//		etiketten: etikett.MakeMutableSet(),
//	}

//	defer func() {
//		err1 := state.Close()

//		if err == nil {
//			err = err1
//		}

//		if !state.context.Zettel.Akte.IsNull() {
//			return
//		}

//		//TODO log the following states
//		if !state.metadataiAkteSha.IsNull() {
//			state.context.Zettel.Akte = state.metadataiAkteSha
//			return
//		}

//		if !state.readAkteSha.IsNull() {
//			state.context.Zettel.Akte = state.readAkteSha
//			return
//		}
//	}()

//	if f.awf == nil {
//		err = errors.Errorf("akte writer factory is nil")
//		return
//	}

//	if state.akteWriter, err = f.awf.AkteWriter(); err != nil {
//		err = errors.Wrap(err)
//		return
//	}

//	if state.akteWriter == nil {
//		err = errors.Errorf("akte writer is nil")
//		return
//	}

//	for {
//		var line string

//		line, err = r.ReadString('\n')
//		n += int64(len(line))

//		if err != nil && err != io.EOF {
//			err = errors.Wrap(err)
//			return
//		}

//		if err == io.EOF {
//			err = nil
//			break
//		}

//		line = strings.TrimSuffix(line, "\n")

//		switch state.field {
//		case textStateReadFieldEmpty:
//			if line != MetadateiBoundary {
//				err = errors.Errorf("expected %q but got %q", MetadateiBoundary, line)
//			}

//			state.field += 1

//		case textStateReadFieldFirstBoundary:
//			if line == MetadateiBoundary {
//				state.field += 1
//			} else if err = f.readMetadateiLine(state, line); err != nil {
//				err = errors.Wrap(err)
//				return
//			}

//		case textStateReadFieldSecondBoundary:
//			if line != "" {
//				err = errors.Errorf("expected empty line after metadatei boundary, but got %q", line)
//				return
//			}

//			state.field += 1

//		case textStateReadFieldAkteBody:

//			if aktePath != "" {
//				c.RecoverableErrors.Add(
//					ErrHasInlineAkteAndFilePath{
//						Zettel:            z,
//						AkteWriterFactory: f.awf,
//						Sha:               state.readAkteSha,
//						FilePath:          aktePath,
//					},
//				)

//				c.AktePath = ""
//			}

//			var n1 int
//			n1, err = io.WriteString(state.akteWriter, fmt.Sprintln(line))

//			if err != nil {
//				err = errors.Wrap(err)
//				break
//			}

//			if n1 != len(line)+1 {
//				err = errors.Errorf("wanted to write %d but only wrote %d", len(line), n1)
//				return
//			}

//		default:
//			err = errors.Errorf("impossible state for field %d", state.field)
//			return
//		}
//	}

//	//TODO outsource this to a context method to allow for injection
//	if c.AktePath != "" {
//		var f *os.File

//		if f, err = files.Open(c.AktePath); err != nil {
//			err = errors.Wrap(err)
//			return
//		}

//		defer files.Close(f)

//		if _, err = io.Copy(state.akteWriter, f); err != nil {
//			err = errors.Wrap(err)
//			return
//		}
//	}

//	return
//}

//func (f DecoderText) readMetadateiLine(state *textStateRead, line string) (err error) {
//	var head, tail string

//	switch len(line) {
//	case 0:
//	case 1:
//		head = line[:1] + " "
//	case 2:
//		head = line[:2]
//	default:
//		head = line[:2]
//		tail = line[2:]
//	}

//	switch head {
//	case "- ":
//		if err = state.etiketten.AddString(tail); err != nil {
//			err = errors.Wrap(err)
//			return
//		}

//		state.lastFieldWasBezeichnung = false

//	case "! ":
//		if err = f.readTyp(state, tail); err != nil {
//			err = errors.Wrap(err)
//			return
//		}

//		state.lastFieldWasBezeichnung = false

//	case "# ":
//		if err = state.context.Zettel.Bezeichnung.Set(tail); err != nil {
//			err = errors.Wrap(err)
//			return
//		}

//		state.lastFieldWasBezeichnung = true

//		// 		if state.lastFieldWasBezeichnung {
//		// 			err = state.context.Zettel.Bezeichnung.Set(tail)
//		// 			state.lastFieldWasBezeichnung = true
//		// 			break
//		// 		}

//		// fallthrough

//	default:
//		if strings.TrimSpace(head) != "" || strings.TrimSpace(tail) != "" {
//			err = errors.Errorf(
//				"unsupported metadatei prefix for format (%q): %q",
//				reflect.TypeOf(f).Name(),
//				head,
//			)
//			return
//		}
//	}

//	return
//}

//func (f DecoderText) readTyp(state *textStateRead, desc string) (err error) {
//	if desc == "" {
//		return
//	}

//	tail := path.Ext(desc)
//	head := desc[:len(desc)-len(tail)]

//	// path
//	if files.Exists(desc) {
//		errors.Print("valid path", desc)

//		if err = state.context.Zettel.Typ.Set(tail); err != nil {
//			err = errors.Wrap(err)
//			return
//		}

//		state.context.AktePath = desc

//		return
//	}

//	//TODO handl akte descs that are invalid files

//	shaError := state.metadataiAkteSha.Set(head)

//	if tail == "" {
//		//sha or ext
//		if shaError != nil {
//			if err = state.context.Zettel.Typ.Set(head); err != nil {
//				err = errors.Wrap(err)
//				return
//			}
//		}
//	} else {
//		//sha.ext or error
//		if shaError != nil {
//			state.context.RecoverableErrors.Add(
//				errors.Wrap(
//					ErrHasInvalidAkteShaOrFilePath{
//						Value: head,
//					},
//				),
//			)

//			return
//		}

//		if err = state.context.Zettel.Typ.Set(tail); err != nil {
//			err = errors.Wrap(err)
//			return
//		}
//	}

//	return
//}
