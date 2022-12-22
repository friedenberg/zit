package typ

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/etikett_rule"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
)

type TextFormat struct {
	arf              gattung.AkteIOFactory
	IgnoreTomlErrors bool
}

func MakeFormatText(arf gattung.AkteIOFactory) *TextFormat {
	return &TextFormat{
		arf: arf,
	}
}

func MakeFormatTextIgnoreTomlErrors(arf gattung.AkteIOFactory) *TextFormat {
	return &TextFormat{
		arf:              arf,
		IgnoreTomlErrors: true,
	}
}

func (f TextFormat) Parse(r io.Reader, t *Objekte) (n int64, err error) {
	return f.ReadFormat(r, t)
}

func (f TextFormat) ReadFormat(r io.Reader, t *Objekte) (n int64, err error) {
	var aw sha.WriteCloser

	if aw, err = f.arf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, aw.Close)

	pr, pw := io.Pipe()
	td := toml.NewDecoder(pr)

	chDone := make(chan error)

	go func(pr *io.PipeReader) {
		var err error
		defer func() {
			chDone <- err
			close(chDone)
		}()

		if err = td.Decode(&t.Akte); err != nil {
			if !errors.IsEOF(err) {
				pr.CloseWithError(err)
			}
		}

		if t.Akte.Actions == nil {
			t.Akte.Actions = make(map[string]script_config.ScriptConfig)
		}

		if t.Akte.Formatters == nil {
			t.Akte.Formatters = make(map[string]script_config.ScriptConfigWithUTI)
		}

		if t.Akte.EtikettenRules == nil {
			t.Akte.EtikettenRules = make(map[string]etikett_rule.Rule)
		}

		if t.Akte.FormatterUTIGroups == nil {
			t.Akte.FormatterUTIGroups = make(map[string]FormatterUTIGroup)
		}
	}(pr)

	mw := io.MultiWriter(aw, pw)

	if n, err = io.Copy(mw, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = pw.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = <-chDone; err != nil {
		if f.IgnoreTomlErrors {
			errors.Err().Print(err)
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	t.Sha = aw.Sha()

	return
}

func (f TextFormat) Format(w io.Writer, t *Objekte) (n int64, err error) {
	return f.WriteFormat(w, t)
}

func (f TextFormat) WriteFormat(w io.Writer, t *Objekte) (n int64, err error) {
	var ar sha.ReadCloser

	if ar, err = f.arf.AkteReader(t.Sha); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.Deferred(&err, ar.Close)

	if n, err = io.Copy(w, ar); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
