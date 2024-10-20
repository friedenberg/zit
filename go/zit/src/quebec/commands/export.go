package commands

import (
	"bufio"
	"flag"
	"io"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/inventory_list"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Export struct {
	AgeIdentity     age.Identity
	CompressionType immutable_config.CompressionType
}

func init() {
	registerCommand(
		"export",
		func(f *flag.FlagSet) CommandWithResult {
			c := &Export{
				CompressionType: immutable_config.CompressionTypeEmpty,
			}

			f.Var(&c.AgeIdentity, "age-identity", "")
			c.CompressionType.AddToFlagSet(f)

			return c
		},
	)
}

func (c Export) Run(u *env.Env, args ...string) (result Result) {
	list := inventory_list.MakeInventoryList()
	var l sync.Mutex

	if result.Error = u.GetStore().QueryPrimitive(
		sku.MakePrimitiveQueryGroup(),
		func(sk *sku.Transacted) (err error) {
			l.Lock()
			defer l.Unlock()

			list.Add(sk.CloneTransacted())

			return
		},
	); result.Error != nil {
		result.Error = errors.Wrap(result.Error)
		return
	}

	var ag age.Age

	if result.Error = ag.AddIdentity(c.AgeIdentity); result.Error != nil {
		result.Error = errors.Wrapf(result.Error, "age-identity: %q", &c.AgeIdentity)
		return
	}

	var wc io.WriteCloser

	// setup inventory list reader
	{
		o := fs_home.WriteOptions{
			Age:             &ag,
			CompressionType: c.CompressionType,
			Writer:          u.Out(),
		}

		if wc, result.Error = fs_home.NewWriter(o); result.Error != nil {
			result.Error = errors.Wrap(result.Error)
			return
		}

		defer errors.DeferredCloser(&result.Error, wc)
	}

	po := (print_options.General{}).
		WithPrintShas(true).
		WithDescriptionInBox(true)

	boxFormat := u.StringFormatWriterSkuBox(
		po,
		u.FormatColorOptionsOut(),
		string_format_writer.CliFormatTruncation66CharEllipsis,
	)

	printer := string_format_writer.MakeDelim(
		"\n",
		u.Out(),
		string_format_writer.MakeFunc(
			func(w interfaces.WriterAndStringWriter, o *sku.Transacted) (n int64, err error) {
				return boxFormat.WriteStringFormat(w, o)
			},
		),
	)

	var sk *sku.Transacted
	var hasMore bool
	bw := bufio.NewWriter(wc)
	defer errors.DeferredFlusher(&result.Error, bw)

	for {
		sk, hasMore = list.Pop()

		if !hasMore {
			break
		}

		if result.Error = printer(sk); result.Error != nil {
			result.Error = errors.Wrap(result.Error)
			return
		}
	}

	return
}
