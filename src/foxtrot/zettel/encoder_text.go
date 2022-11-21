package zettel

// type EncoderText struct {
// 	out io.Writer
//   konfig konfig.Compiled

// 	DoNotWriteEmptyBezeichnung bool
// 	TypError                   error
// }

// func MakeEncoderText(
//   k konfig.Konfig,
// 	out io.Writer,
// ) *EncoderText {
// 	return &EncoderText{
// 		out: out,
//     konfig: k.Compiled,
// 	}
// }

// func (f EncoderText) Encode(z ZettelCoded) (n int64, err error) {
// 	switch {
// 	case c.IncludeAkte && c.ExternalAktePath == "":
// 		return f.writeToInlineAkte(z)

// 	case c.IncludeAkte:
// 		return f.writeToExternalAkte(z)

// 	default:
// 		return f.writeToOmitAkte(z)
// 	}
// }

// func (f EncoderText) writeToOmitAkte(z *ZettelCoded) (n int64, err error) {
// 	w := line_format.NewWriter()

// 	w.WriteLines(
// 		MetadateiBoundary,
// 	)

// 	if c.Zettel.Bezeichnung.String() != "" || !f.DoNotWriteEmptyBezeichnung {
// 		w.WriteLines(
// 			fmt.Sprintf("# %s", c.Zettel.Bezeichnung),
// 		)
// 	}

// 	for _, e := range c.Zettel.Etiketten.Sorted() {
// 		w.WriteFormat("- %s", e)
// 	}

// 	switch {
// 	case c.Zettel.Akte.IsNull() && c.Zettel.Typ.String() == "":
// 		break

// 	case c.Zettel.Akte.IsNull():
// 		w.WriteLines(
// 			fmt.Sprintf("! %s", c.Zettel.Typ),
// 		)

// 	case c.Zettel.Typ.String() == "":
// 		w.WriteLines(
// 			fmt.Sprintf("! %s", c.Zettel.Akte),
// 		)

// 	default:
// 		w.WriteLines(
// 			fmt.Sprintf("! %s.%s", c.Zettel.Akte, c.Zettel.Typ),
// 		)
// 	}

// 	w.WriteLines(
// 		MetadateiBoundary,
// 	)

// 	n, err = w.WriteTo(c.Out)

// 	return
// }

// func (f EncoderText) writeToInlineAkte(z ZettelCoded) (n int64, err error) {
// 	w := line_format.NewWriter()

// 	w.WriteLines(
// 		MetadateiBoundary,
// 		fmt.Sprintf("# %s", c.Zettel.Bezeichnung),
// 	)

// 	for _, e := range c.Zettel.Etiketten.Sorted() {
// 		w.WriteFormat("- %s", e)
// 	}

// 	w.WriteLines(
// 		fmt.Sprintf("! %s", c.Zettel.Typ),
// 	)

// 	w.WriteLines(
// 		MetadateiBoundary,
// 	)

// 	w.WriteEmpty()

// 	n, err = w.WriteTo(c.Out)

// 	if err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	var ar io.ReadCloser

// 	if c.AkteReaderFactory == nil {
// 		err = errors.Errorf("akte reader factory is nil")
// 		return
// 	}

// 	ar, err = c.AkteReader(c.Zettel.Akte)

// 	if err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if ar == nil {
// 		err = errors.Errorf("akte reader is nil")
// 		return
// 	}

// 	defer errors.Deferred(&err, ar.Close)

// 	in := ar

// 	var cmd *exec.Cmd

// 	if c.FormatScript != nil {
// 		if cmd, err = c.FormatScript.Cmd(); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}

// 	if cmd != nil {
// 		cmd.Stdin = ar
// 		cmd.Stderr = os.Stderr

// 		if in, err = cmd.StdoutPipe(); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}

// 		if err = cmd.Start(); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}

// 	var n1 int64

// 	n1, err = io.Copy(c.Out, in)
// 	n += n1

// 	if err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if cmd != nil {
// 		if err = cmd.Wait(); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}

// 	return
// }

// func (f EncoderText) writeToExternalAkte(z ZettelCoded) (n int64, err error) {
// 	w := line_format.NewWriter()

// 	w.WriteLines(
// 		MetadateiBoundary,
// 		fmt.Sprintf("# %s", c.Zettel.Bezeichnung),
// 	)

// 	for _, e := range c.Zettel.Etiketten.Sorted() {
// 		w.WriteFormat("- %s", e)
// 	}

// 	if strings.Index(c.ExternalAktePath, "\n") != -1 {
// 		panic(errors.Errorf("ExternalAktePath contains newline: %q", c.ExternalAktePath))
// 	}

// 	w.WriteLines(
// 		fmt.Sprintf("! %s", c.ExternalAktePath),
// 	)

// 	w.WriteLines(
// 		MetadateiBoundary,
// 	)

// 	n, err = w.WriteTo(c.Out)

// 	if err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	var ar io.ReadCloser

// 	if c.AkteReaderFactory == nil {
// 		err = errors.Errorf("akte reader factory is nil")
// 		return
// 	}

// 	if ar, err = c.AkteReader(c.Zettel.Akte); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if ar == nil {
// 		err = errors.Errorf("akte reader is nil")
// 		return
// 	}

// 	defer errors.Deferred(&err, ar.Close)

// 	var file *os.File

// 	if file, err = files.Create(c.ExternalAktePath); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	defer files.Close(file)

// 	var n1 int64

// 	n1, err = io.Copy(file, ar)
// 	n += n1

// 	if err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }
