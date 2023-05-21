package commands

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/november/umwelt"
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
	ms kennung.MetaSet,
	cwdFiles *cwd.CwdFiles,
) (err error) {
	// e := zettel_external.MakeFileEncoderJustOpen(
	// 	u.StoreObjekten(),
	// 	u.Konfig(),
	// )

	fInline := metadatei.MakeTextFormatterMetadateiInlineAkte(
		u.StoreObjekten(),
		nil,
	)

	fMetadatei := metadatei.MakeTextFormatterMetadateiOnly(
		u.StoreObjekten(),
		nil,
	)

	if err = u.StoreWorkingDirectory().ReadFiles(
		cwdFiles,
		ms,
		iter.MakeChain(
			objekte.MakeFilterFromMetaSet(ms),
			func(co objekte.CheckedOutLike) (err error) {
				var zco *zettel.CheckedOut
				ok := false

				if zco, ok = co.(*zettel.CheckedOut); !ok {
					todo.Change("add support for other gattung")
					return
				}

				wg := iter.MakeErrorWaitGroup()

				var rLeft, wLeft *os.File

				if rLeft, wLeft, err = os.Pipe(); err != nil {
					err = errors.Wrap(err)
					return
				}

				defer errors.DeferredCloser(&err, rLeft)

				var rRight, wRight *os.File

				if rRight, wRight, err = os.Pipe(); err != nil {
					err = errors.Wrap(err)
					return
				}

				defer errors.DeferredCloser(&err, rRight)

				todo.Change("support checkout mode")
				wg.Do(
					func() (err error) {
						defer errors.DeferredCloser(&err, wLeft)

						formatFunc := fInline.FormatMetadatei

						if !u.Konfig().IsInlineTyp(zco.Internal.GetTyp()) {
							formatFunc = fMetadatei.FormatMetadatei
						}

						if _, err = formatFunc(wLeft, zco.Internal); err != nil {
							err = errors.Wrap(err)
							return
						}

						return
					},
				)

				wg.Do(
					func() (err error) {
						defer errors.DeferredCloser(&err, wRight)

						formatFunc := fInline.FormatMetadatei

						if !u.Konfig().IsInlineTyp(zco.External.GetTyp()) {
							formatFunc = fMetadatei.FormatMetadatei
						}

						if _, err = formatFunc(wRight, zco.External); err != nil {
							err = errors.Wrap(err)
							return
						}

						return
					},
				)

				todo.Change("disambiguate internal and external, and objekte / akte")
				cmd := exec.Command(
					"diff",
					"--color=always",
					"-u",
					"--label", fmt.Sprintf("%s@zettel", zco.Internal.Sku.Kennung),
					"--label", fmt.Sprintf("%s@zettel", zco.Internal.Sku.Kennung),
					"/dev/fd/3",
					"/dev/fd/4",
				)

				cmd.ExtraFiles = []*os.File{rLeft, rRight}
				cmd.Stdout = u.Out()
				cmd.Stderr = u.Err()

				wg.Do(
					func() (err error) {
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
