package metadatei

// type metadateiAndAkteTextParser struct {
// 	MetadateiParser TextParser
// }

// func MakeMetadateiAndAkteTextParser(
// 	akteFactory schnittstellen.AkteIOFactory,
// 	akteFormatter script_config.RemoteScript,
// ) TextParser {
// 	return metadateiAndAkteTextParser{
// 		AkteFactory:     akteFactory,
// 		AkteFormatter:   akteFormatter,
// 		MetadateiParser: MakeTextParser(akteFactory),
// 	}
// }

// func (f metadateiAndAkteTextParser) Parse(
// 	r io.Reader,
// 	c ParserContext,
// ) (n int64, err error) {
// 	var akteWriter sha.WriteCloser

// 	if f.AkteFactory == nil {
// 		err = errors.Errorf("akte factory is nil")
// 		return
// 	}

// 	if akteWriter, err = f.AkteFactory.AkteWriter(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if akteWriter == nil {
// 		err = errors.Errorf("akte writer is nil")
// 		return
// 	}

// 	mr := metadatei_io.Reader{
// 		Metadatei: format.MakeReaderFromInterface[ParserContext](
// 			f.MetadateiParser.Parse,
// 			c,
// 		),
// 		Akte: akteWriter,
// 	}

// 	if n, err = mr.ReadFrom(r); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if err = akteWriter.Close(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	inlineAkteSha := akteWriter.Sha()
// 	m := c.GetMetadateiPtr()

// 	switch {
// 	case m.AkteSha.IsNull() && !inlineAkteSha.IsNull():
// 		m.AkteSha = sha.Make(inlineAkteSha)

// 	case !m.AkteSha.IsNull() && inlineAkteSha.IsNull():
// 		// noop

// 	case !m.AkteSha.IsNull() && !inlineAkteSha.IsNull():
// 		err = ErrHasInlineAkteAndFilePath{
// 			Metadatei: *m,
// 		}

// 		return
// 	}

// 	return
// }
