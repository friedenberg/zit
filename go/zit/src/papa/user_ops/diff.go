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
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

// TODO move to store_fs
type Diff struct {
	*env.Local

	object_metadata.TextFormatterFamily
}

func (op Diff) Run(
	remoteCheckedOut sku.SkuType,
	options object_metadata.TextFormatterOptions,
) (err error) {
	var localCheckedOut sku.SkuType

	{
		if localCheckedOut, err = op.GetStore().GetStoreFS().CheckoutOne(
			checkout_options.Options{
				CheckoutMode: checkout_mode.MetadataAndBlob,
				OptionsWithoutMode: checkout_options.OptionsWithoutMode{
					Path: checkout_options.PathTempLocal,
				},
			},
			remoteCheckedOut.GetSku(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, func() (err error) {
			if err = op.GetStore().GetStoreFS().DeleteCheckedOutInternal(
				localCheckedOut,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		})
	}

	wg := quiter.MakeErrorWaitGroupParallel()

	var mode checkout_mode.Mode

	local := localCheckedOut.GetSku()
	localContext := object_metadata.TextFormatterContext{
		PersistentFormatterContext: local,
		TextFormatterOptions:       options,
	}

	remote := remoteCheckedOut.GetSkuExternal()
	remoteCtx := object_metadata.TextFormatterContext{
		PersistentFormatterContext: remote,
		TextFormatterOptions:       options,
	}

	if mode, err = op.GetStore().GetStoreFS().GetCheckoutModeOrError(
		remote,
	); err != nil {
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
	internalInline := op.GetConfig().IsInlineType(local.GetType())
	externalInline := op.GetConfig().IsInlineType(remote.GetType())

	var fds *sku.FSItem

	if fds, err = op.GetStore().GetStoreFS().ReadFSItemFromExternal(remote); err != nil {
		err = errors.Wrap(err)
		return
	}

	var externalFD *fd.FD

	switch {
	case mode.IncludesMetadata():
		if internalInline && externalInline {
			wg.Do(op.makeDo(wLeft, op.InlineBlob, localContext))
			wg.Do(op.makeDo(wRight, op.InlineBlob, remoteCtx))
		} else {
			wg.Do(op.makeDo(wLeft, op.MetadataOnly, localContext))
			wg.Do(op.makeDo(wRight, op.MetadataOnly, remoteCtx))
		}

		externalFD = &fds.Object

	case internalInline && externalInline:
		wg.Do(op.makeDoBlob(wLeft, op.GetDirectoryLayout(), local.GetBlobSha()))
		wg.Do(op.makeDoFD(wRight, &fds.Blob))
		externalFD = &fds.Blob

	default:
		wg.Do(op.makeDo(wLeft, op.MetadataOnly, localContext))
		wg.Do(op.makeDo(wRight, op.MetadataOnly, remoteCtx))
		externalFD = &fds.Blob
	}

	internalLabel := fmt.Sprintf(
		"%s:%s",
		local.GetObjectId(),
		strings.ToLower(local.GetGenre().GetGenreString()),
	)

	externalLabel := op.GetDirectoryLayout().Rel(externalFD.GetPath())

	colorOptions := op.FormatColorOptionsOut()
	colorString := "always"

	if colorOptions.OffEntirely {
		colorString = "never"
	}

	todo.Change("disambiguate internal and external, and object / blob")
	cmd := exec.Command(
		"diff",
		fmt.Sprintf("--color=%s", colorString),
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

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
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

func (c Diff) makeDoBlob(
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
