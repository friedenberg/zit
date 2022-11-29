package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/script_value"
)

type WriterScript struct {
	script    script_value.ScriptValue
	scriptIn  io.WriteCloser
	scriptOut io.Reader
	enc       WriterJson
}

func MakeWriterScript(s script_value.ScriptValue) (w WriterScript, err error) {
	w = WriterScript{
		script: s,
	}

	if w.scriptIn, w.scriptOut, err = w.script.RunWithInput(); err != nil {
		err = errors.Wrap(err)
		return
	}

	w.enc = MakeWriterJson(w.scriptIn)

	return
}

func (w WriterScript) Reader() io.Reader {
	return w.scriptOut
}

func (w WriterScript) WriteZettelTransacted(z *Transacted) (err error) {
	errors.Log().Printf("writing zettel: %v", z)
	if err = w.enc.WriteZettelTransacted(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (w WriterScript) Close() (err error) {
	if err = w.scriptIn.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = w.script.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
