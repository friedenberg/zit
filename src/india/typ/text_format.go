package typ

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/etikett_rule"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
)

// TODO-P1 rename to TextFormat
type FormatText struct {
	arf gattung.AkteIOFactory
}

func MakeFormatText(arf gattung.AkteIOFactory) *FormatText {
	return &FormatText{
		arf: arf,
	}
}

func (f FormatText) Parse(r io.Reader, t *Objekte) (n int64, err error) {
	return f.ReadFormat(r, t)
}

func (f FormatText) ReadFormat(r io.Reader, t *Objekte) (n int64, err error) {
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
			t.Akte.Actions = make(map[string]Action)
		}

		if t.Akte.EtikettenRules == nil {
			t.Akte.EtikettenRules = make(map[string]etikett_rule.Rule)
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

	t.Sha = aw.Sha()

	return
}

func (f FormatText) WriteFormat(w io.Writer, t *Objekte) (n int64, err error) {
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
