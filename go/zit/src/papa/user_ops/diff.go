package user_ops

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Diff struct {
	*env.Env
	Inline   object_metadata.TextFormatter
	Metadata object_metadata.TextFormatter
}

func (op Diff) Run(col sku.CheckedOutLike) (err error) {
	cofs, ok := col.(*store_fs.CheckedOut)

	if !ok {
		if col, err = op.GetStore().GetCwdFiles().CheckoutOne(
			checkout_options.Options{
				Path:         checkout_options.PathTempLocal,
				CheckoutMode: checkout_mode.ModeMetadataAndBlob,
			},
			col.GetSku(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		cofs = col.(*store_fs.CheckedOut)

		defer errors.Deferred(&err, func() (err error) {
			if err = op.GetStore().GetCwdFiles().DeleteCheckout(cofs); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		})
	}

	wg := iter.MakeErrorWaitGroupParallel()
	var mode checkout_mode.Mode

	il := &cofs.Internal
	el := &cofs.External

	if mode, err = el.GetFDs().GetCheckoutModeOrError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var rLeft, wLeft *os.File

	if rLeft, wLeft, err = os.Pipe(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var rRight, wRight *os.File

	if rRight, wRight, err = os.Pipe(); err != nil {
		err = errors.Wrap(err)
		return
	}

	// sameTyp := il.GetTyp().Equals(el.GetTyp())
	internalInline := op.GetConfig().IsInlineType(il.GetType())
	externalInline := op.GetConfig().IsInlineType(el.GetType())

	var externalFD *fd.FD

	switch {
	case mode.IncludesObject():
		if internalInline && externalInline {
			wg.Do(op.makeDo(wLeft, op.Inline, il))
			wg.Do(op.makeDo(wRight, op.Inline, el))
		} else {
			wg.Do(op.makeDo(wLeft, op.Metadata, il))
			wg.Do(op.makeDo(wRight, op.Metadata, el))
		}

		externalFD = el.GetObjekteFD()

	case internalInline && externalInline:
		wg.Do(op.makeDoAkte(wLeft, op.GetFSHome(), il.GetBlobSha()))
		wg.Do(op.makeDoFD(wRight, el.GetAkteFD()))
		externalFD = el.GetAkteFD()

	default:
		wg.Do(op.makeDo(wLeft, op.Metadata, il))
		wg.Do(op.makeDo(wRight, op.Metadata, el))
		externalFD = el.GetAkteFD()
	}

	internalLabel := fmt.Sprintf(
		"%s:%s",
		il.GetObjectId(),
		strings.ToLower(il.GetGenre().GetGenreString()),
	)

	externalLabel := op.GetFSHome().Rel(externalFD.GetPath())

	todo.Change("disambiguate internal and external, and objekte / akte")
	cmd := exec.Command(
		"diff",
		"--color=always",
		"-u",
		"--label", internalLabel,
		"--label", externalLabel,
		"/dev/fd/3",
		"/dev/fd/4",
	)

	cmd.ExtraFiles = []*os.File{rLeft, rRight}
	cmd.Stdout = op.Out()
	cmd.Stderr = op.Err()

	wg.Do(
		func() (err error) {
			defer errors.DeferredCloser(&err, rLeft)
			defer errors.DeferredCloser(&err, rRight)

			if err = cmd.Run(); err != nil {
				if cmd.ProcessState.ExitCode() == 1 {
					todo.Change("return non-zero exit code")
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}

			return
		},
	)

	return wg.GetError()
}

func (c Diff) makeDo(
	w io.WriteCloser,
	mf object_metadata.TextFormatter,
	m object_metadata.TextFormatterContext,
) interfaces.FuncError {
	return func() (err error) {
		defer errors.DeferredCloser(&err, w)

		if _, err = mf.FormatMetadata(w, m); err != nil {
			if errors.IsBrokenPipe(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		return
	}
}

func (c Diff) makeDoAkte(
	w io.WriteCloser,
	arf interfaces.BlobReaderFactory,
	sh interfaces.Sha,
) interfaces.FuncError {
	return func() (err error) {
		defer errors.DeferredCloser(&err, w)

		var ar sha.ReadCloser

		if ar, err = arf.BlobReader(sh); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, ar)

		if _, err = io.Copy(w, ar); err != nil {
			if errors.IsBrokenPipe(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		return
	}
}

func (c Diff) makeDoFD(
	w io.WriteCloser,
	fd *fd.FD,
) interfaces.FuncError {
	return func() (err error) {
		defer errors.DeferredCloser(&err, w)

		var f *os.File

		if f, err = files.OpenExclusiveReadOnly(fd.GetPath()); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, f)

		if _, err = io.Copy(w, f); err != nil {
			if errors.IsBrokenPipe(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		return
	}
}
