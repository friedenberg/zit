package zettel_formats

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/echo/zettel"
)

const (
	MetadateiBoundary = "---"
)

type Text struct {
	DoNotWriteEmptyBezeichnung bool
}

type textStateReadField int

const (
	textStateReadFieldEmpty = textStateReadField(iota)
	textStateReadFieldFirstBoundary
	textStateReadFieldSecondBoundary
	textStateReadFieldAkteBody
)

type textStateReadAkte int

const (
	// no akte file or ext, therefore it's inline
	// yes akte just ext and it's an inline type, therefore it's inline
	// yes akte just ext and it's not an inline type, therefore error
	// yes akte file and ext, therefore it's a file
	// yes akte file and ext and content inline, therefore error
	textStateReadAkteInline            = textStateReadAkte(iota)
	textStateReadAkteFileWithExtension = textStateReadAkte(iota)
	textStateReadAkteJustExtension     = textStateReadAkte(iota)
)

type textStateRead struct {
	context                 *_ZettelFormatContextRead
	field                   textStateReadField
	lastFieldWasBezeichnung bool
	didReadAkte             bool
	metadataiAkteSha        sha.Sha
	readAkteSha             sha.Sha
	akteWriter              _ObjekteWriter
}

func (s *textStateRead) Close() (err error) {
	if s.akteWriter != nil {
		if err = s.akteWriter.Close(); err != nil {
			err = _Error(err)
			return
		}

		s.readAkteSha = s.akteWriter.Sha()
	}

	return
}

type textStateWrite struct {
	zettel.Zettel
}

func (f Text) ReadFrom(c *_ZettelFormatContextRead) (n int64, err error) {
	r := bufio.NewReader(c.In)

	c.Zettel.Etiketten = etikett.MakeSet()

	state := &textStateRead{
		context: c,
	}

	defer func() {
		err1 := state.Close()

		if err == nil {
			err = err1
		}

		if !state.context.Zettel.Akte.IsNull() {
			return
		}

		//TODO log the following states
		if !state.metadataiAkteSha.IsNull() {
			state.context.Zettel.Akte = state.metadataiAkteSha
			return
		}

		if !state.readAkteSha.IsNull() {
			state.context.Zettel.Akte = state.readAkteSha
			return
		}
	}()

	if c.AkteWriterFactory == nil {
		err = _Errorf("akte writer factory is nil")
		return
	}

	if state.akteWriter, err = c.AkteWriter(); err != nil {
		err = _Error(err)
		return
	}

	if state.akteWriter == nil {
		err = _Errorf("akte writer is nil")
		return
	}

	for {
		var line string

		line, err = r.ReadString('\n')
		n += int64(len(line))

		if err != nil && err != io.EOF {
			err = _Error(err)
			return
		}

		if err == io.EOF {
			err = nil
			break
		}

		line = strings.TrimSuffix(line, "\n")

		switch state.field {
		case textStateReadFieldEmpty:
			if line != MetadateiBoundary {
				err = _Errorf("expected %q but got %q", MetadateiBoundary, line)
			}

			state.field += 1

		case textStateReadFieldFirstBoundary:
			if line == MetadateiBoundary {
				state.field += 1
			} else if err = f.readMetadateiLine(state, line); err != nil {
				err = _Error(err)
				return
			}

		case textStateReadFieldSecondBoundary:
			if line != "" {
				err = _Errorf("expected empty line after metadatei boundary, but got %q", line)
				return
			}

			state.field += 1

		case textStateReadFieldAkteBody:

			if c.AktePath != "" {
				c.RecoverableError = ErrHasInlineAkteAndFilePath{
					Zettel:             c.Zettel,
					_AkteWriterFactory: c,
					Sha:                state.readAkteSha,
					FilePath:           c.AktePath,
				}

				c.AktePath = ""
			}

			var n1 int
			n1, err = io.WriteString(state.akteWriter, fmt.Sprintln(line))

			if err != nil {
				err = _Error(err)
				break
			}

			if n1 != len(line)+1 {
				err = _Errorf("wanted to write %d but only wrote %d", len(line), n1)
				return
			}

		default:
			err = _Errorf("impossible state for field %d", state.field)
			return
		}
	}

	//TODO outsource this to a context method to allow for injection
	if c.AktePath != "" {
		var f *os.File

		if f, err = _Open(c.AktePath); err != nil {
			err = _Error(err)
			return
		}

		defer _Close(f)

		if _, err = io.Copy(state.akteWriter, f); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}

func (f Text) readMetadateiLine(state *textStateRead, line string) (err error) {
	var head, tail string

	switch len(line) {
	case 0:
	case 1:
		head = line[:1] + " "
	case 2:
		head = line[:2]
	default:
		head = line[:2]
		tail = line[2:]
	}

	switch head {
	case "- ":
		err = state.context.Zettel.Etiketten.AddString(tail)
		state.lastFieldWasBezeichnung = false

	case "! ":
		err = f.readAkteDesc(state, tail)
		state.lastFieldWasBezeichnung = false

	case "# ":
		err = state.context.Zettel.Bezeichnung.Set(tail)
		state.lastFieldWasBezeichnung = true

		// 		if state.lastFieldWasBezeichnung {
		// 			err = state.context.Zettel.Bezeichnung.Set(tail)
		// 			state.lastFieldWasBezeichnung = true
		// 			break
		// 		}

		// fallthrough

	default:
		if strings.TrimSpace(head) != "" || strings.TrimSpace(tail) != "" {
			err = _Errorf(
				"unsupported metadatei prefix for format (%q): %q",
				reflect.TypeOf(f).Name(),
				head,
			)
		}
	}

	if err != nil {
		err = _Error(err)
		return
	}

	return
}

func (f Text) readAkteDesc(state *textStateRead, desc string) (err error) {
	if desc == "" {
		return
	}

	tail := path.Ext(desc)
	head := desc[:len(desc)-len(tail)]

	// path
	if _FilesExists(desc) {
		logz.Print("valid path", desc)

		if err = state.context.Zettel.AkteExt.Set(tail); err != nil {
			err = _Error(err)
			return
		}

		state.context.AktePath = desc

		return
	}

	//TODO handl akte descs that are invalid files

	shaError := state.metadataiAkteSha.Set(head)

	logz.Print(head)
	logz.Print(tail)

	if tail == "" {
		//sha or ext
		if shaError != nil {
			if err = state.context.Zettel.AkteExt.Set(head); err != nil {
				err = _Error(err)
				return
			}
		}
	} else {
		//sha.ext or error
		if shaError != nil {
			err = _Error(err)
			return
		}

		if err = state.context.Zettel.AkteExt.Set(tail); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}
