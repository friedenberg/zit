package user_ops

import (
	"io"
	"os/exec"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/script_value"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type FilterZettelsWithScript struct {
	Set    collections.MutableSetLike[*zettel.Transacted]
	Filter script_value.ScriptValue
}

func (op FilterZettelsWithScript) Run() (err error) {
	if op.Filter.IsEmpty() {
		errors.Log().Print("no filter")
		return
	}

	cmd := exec.Command(op.Filter.String())

	var w io.WriteCloser

	if w, err = cmd.StdinPipe(); err != nil {
		errors.Wrap(err)
		return
	}

	var r io.Reader

	if r, err = cmd.StdoutPipe(); err != nil {
		errors.Wrap(err)
		return
	}

	enc := zettel.MakeWriterJson(w)

	chDone, chErr := op.runGetHinweisen(r)

	if err = cmd.Start(); err != nil {
		err = errors.Wrap(err)
		return
	}

	go func() {
		defer w.Close()
		op.Set.Each(enc.WriteZettelVerzeichnisse)
	}()

	select {
	case err = <-chErr:
		err = errors.Wrap(err)
		return

	case hinweisen := <-chDone:

		errors.Log().Printf("%#v", hinweisen)
		op.Set.Each(
			collections.MakeChain(
				func(z *zettel.Transacted) (err error) {
					ok := hinweisen.Contains(z.Sku.Kennung)

					if ok {
						err = collections.MakeErrStopIteration()
						return
					}

					return
				},
				op.Set.Del,
			),
		)
	}

	if err = cmd.Wait(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op FilterZettelsWithScript) runGetHinweisen(
	r io.Reader,
) (chDone <-chan kennung.HinweisSet, chErr <-chan error) {
	doneBoth := make(chan kennung.HinweisSet)
	chDone = doneBoth

	errBoth := make(chan error)
	chErr = errBoth

	go func() {
		irl := kennung.ReaderLine{}

		if _, err := irl.ReadFrom(r); err != nil {
			err = errors.Wrap(err)
			errBoth <- err
			return
		}

		doneBoth <- irl.Set.Hinweisen.Copy()
	}()

	return
}
