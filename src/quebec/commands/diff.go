package commands

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/lima/cwd"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Diff struct{}

func init() {
	registerCommandWithCwdQuery(
		"diff",
		func(f *flag.FlagSet) CommandWithCwdQuery {
			c := &Diff{}

			return c
		},
	)
}

func (c Diff) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(gattung.TrueGattung()...)
}

func (c Diff) RunWithCwdQuery(
	u *umwelt.Umwelt,
	ms matcher.Query,
	cwdFiles *cwd.CwdFiles,
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
		cwdFiles,
		matcher.MakeFuncReaderTransactedLikePtr(ms, u.StoreObjekten().Query),
		iter.MakeChain(
			matcher.MakeFilterFromQuery(ms),
			func(co *sku.CheckedOut) (err error) {
				wg := iter.MakeErrorWaitGroup()

				il := co.Internal
				el := co.External

				var mode checkout_mode.Mode

				if mode, err = el.GetFDs().GetCheckoutMode(); err != nil {
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

				var externalFD kennung.FD

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
					il.GetKennungLike(),
					strings.ToLower(il.GetGattung().GetGattungString()),
				)

				externalLabel := u.Standort().Rel(externalFD.Path)

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
	fd kennung.FD,
) schnittstellen.FuncError {
	return func() (err error) {
		defer errors.DeferredCloser(&err, w)

		var f *os.File

		if f, err = files.OpenExclusiveReadOnly(fd.Path); err != nil {
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
