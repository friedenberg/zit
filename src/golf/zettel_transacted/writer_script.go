package zettel_transacted

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/script_value"
)

type WriterScript struct {
	script    script_value.ScriptValue
	scriptOut io.Reader
	enc       WriterJson
}

func MakeWriterScript(s script_value.ScriptValue) (w WriterScript, err error) {
	w = WriterScript{
		script: s,
	}

	r, w1 := io.Pipe()

	w.enc = MakeWriterJson(w1)

	if w.scriptOut, err = w.script.RunWithInput(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (w WriterScript) WriteZettelTransacted(z *Zettel) (err error) {
	if err = w.enc.WriteZettelTransacted(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
