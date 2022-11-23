package typ_checked_out

import (
	"github.com/friedenberg/zit/src/echo/typ"
	"github.com/friedenberg/zit/src/foxtrot/cwd_files"
)

type Typ struct {
	cwd_files.CwdTyp
	typ.Named
}
