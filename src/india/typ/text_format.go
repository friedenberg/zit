package typ

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/etikett_rule"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/schnittstellen"
)

type TextFormat struct {
	arf              schnittstellen.AkteIOFactory
	ignoreTomlErrors bool
}

func MakeFormatText(arf schnittstellen.AkteIOFactory) *TextFormat {
	return &TextFormat{
		arf: arf,
	}
}

// TODO-P4 remove
func MakeFormatTextIgnoreTomlErrors(arf schnittstellen.AkteIOFactory) *TextFormat {
	return &TextFormat{
		arf:              arf,
		ignoreTomlErrors: true,
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

		defer func() {
			if r := recover(); r != nil {
				if f.ignoreTomlErrors {
					err = nil
				} else {
					err = toml.MakeError(errors.Errorf("panicked during toml decoding: %s", r))
					pr.CloseWithError(errors.Wrap(err))
				}
			}
		}()

		if err = td.Decode(&t.Akte); err != nil {
			switch {
			case !errors.IsEOF(err) && !f.ignoreTomlErrors:
				err = errors.Wrap(toml.MakeError(err))
				pr.CloseWithError(err)

			case !errors.IsEOF(err) && f.ignoreTomlErrors:
				err = nil
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
		err = errors.Wrap(err)
		return
	}

	t.Sha = sha.Make(aw.Sha())

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
