package string_format_writer

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/erworben_cli_print_options"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
)

func StringPrefixFromOptions(
	options erworben_cli_print_options.PrintOptions,
) string {
	if options.ZittishNewlines {
		return "\n  " + StringIndent
	} else {
		return " "
	}
}

func WriteStringPrefixFormat(
	w schnittstellen.WriterAndStringWriter,
	prefix, body string,
) (n int64, err error) {
	var n1 int

	n1, err = w.WriteString(prefix)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = w.WriteString(body)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
