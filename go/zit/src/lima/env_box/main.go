package env_box

import (
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
)

type Env interface {
	env_repo.Env
}

type env struct {
	env_repo.Env
}

// func (u env) MakePrinterBoxArchive(
// 	out interfaces.WriterAndStringWriter,
// 	includeTai bool,
// ) interfaces.FuncIter[*sku.Transacted] {
// 	boxFormat := box_format.MakeBoxTransactedArchive(
// 		u.GetEnv(),
// 		u.GetConfig().PrintOptions.WithPrintTai(includeTai),
// 	)

// 	return string_format_writer.MakeDelim(
// 		"\n",
// 		out,
// 		string_format_writer.MakeFunc(
// 			func(w interfaces.WriterAndStringWriter, o *sku.Transacted) (n int64, err error) {
// 				return boxFormat.WriteStringFormat(w, o)
// 			},
// 		),
// 	)
// }
