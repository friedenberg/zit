package commands

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit-go/src/bravo/files"
	"code.linenisgreat.com/zit-go/src/bravo/iter"
	"code.linenisgreat.com/zit-go/src/bravo/todo"
	"code.linenisgreat.com/zit-go/src/charlie/gattung"
	"code.linenisgreat.com/zit-go/src/charlie/sha"
	"code.linenisgreat.com/zit-go/src/delta/gattungen"
	"code.linenisgreat.com/zit-go/src/echo/fd"
	"code.linenisgreat.com/zit-go/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
	"code.linenisgreat.com/zit-go/src/india/matcher"
	"code.linenisgreat.com/zit-go/src/oscar/umwelt"
)

type Diff struct{}

func init() {
	registerCommandWithQuery(
		"diff",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Diff{}

			return c
		},
	)
}

func (c Diff) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(gattung.TrueGattung()...)
}

func (c Diff) RunWithQuery(
	u *umwelt.Umwelt,
	ms matcher.Query,
) (err error) {
	fInline := metadatei.MakeTextFormatterMetadateiInlineAkte(
		u.Standort(),
		nil,
	)

	fMetadatei := metadatei.MakeTextFormatterMetadateiOnly(
		u.Standort(),
		nil,
	)

	if err = u.StoreObjekten().ReadFiles(
		matcher.MakeFuncReaderTransactedLikePtr(ms, u.StoreObjekten().QueryWithoutCwd),
		iter.MakeChain(
			matcher.MakeFilterFromQuery(ms),
			func(co *sku.CheckedOut) (err error) {
				wg := iter.MakeErrorWaitGroupParallel()

				il := &co.Internal
				el := &co.External

				var mode checkout_mode.Mode

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
				internalInline := u.Konfig().IsInlineTyp(il.GetTyp())
				externalInline := u.Konfig().IsInlineTyp(el.GetTyp())

				var externalFD *fd.FD

				switch {
				case mode.IncludesObjekte():
					if internalInline && externalInline {
						wg.Do(c.makeDo(wLeft, fInline, il))
						wg.Do(c.makeDo(wRight, fInline, el))
					} else {
						wg.Do(c.makeDo(wLeft, fMetadatei, il))
						wg.Do(c.makeDo(wRight, fMetadatei, el))
					}

					externalFD = el.GetObjekteFD()

				case internalInline && externalInline:
					wg.Do(c.makeDoAkte(wLeft, u.Standort(), il.GetAkteSha()))
					wg.Do(c.makeDoFD(wRight, el.GetAkteFD()))
					externalFD = el.GetAkteFD()

				default:
					wg.Do(c.makeDo(wLeft, fMetadatei, il))
					wg.Do(c.makeDo(wRight, fMetadatei, el))
					externalFD = el.GetAkteFD()
				}

				internalLabel := fmt.Sprintf(
					"%s:%s",
					il.GetKennung(),
					strings.ToLower(il.GetGattung().GetGattungString()),
				)

				externalLabel := u.Standort().Rel(externalFD.GetPath())

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
				cmd.Stdout = u.Out()
				cmd.Stderr = u.Err()

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
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Diff) makeDo(
	w io.WriteCloser,
	mf metadatei.TextFormatter,
	m metadatei.TextFormatterContext,
) schnittstellen.FuncError {
	return func() (err error) {
		defer errors.DeferredCloser(&err, w)

		if _, err = mf.FormatMetadatei(w, m); err != nil {
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
	arf schnittstellen.AkteReaderFactory,
	sh schnittstellen.ShaLike,
) schnittstellen.FuncError {
	return func() (err error) {
		defer errors.DeferredCloser(&err, w)

		var ar sha.ReadCloser

		if ar, err = arf.AkteReader(sh); err != nil {
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
) schnittstellen.FuncError {
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
