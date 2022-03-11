package user_ops

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/script_value"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/friedenberg/zit/hotel/zettels"
)

type (
	_Store                 = zettels.Zettels
	_Umwelt                = *umwelt.Umwelt
	_ZettelCheckedOut      = stored_zettel.CheckedOut
	_ZettelsCheckinOptions = zettels.CheckinOptions
	_NamedZettel           = stored_zettel.Named
	_ZettelFormatsText     = zettel_formats.Text
	_Zettel                = zettel.Zettel
	_ScriptValue           = script_value.ScriptValue
)

var (
	_Error     = errors.Error
	_Errorf    = errors.Errorf
	_OpenFiles = open_file_guard.OpenFiles
	_Outf      = stdprinter.Outf
)
