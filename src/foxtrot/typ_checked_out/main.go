package typ_checked_out

import (
	"github.com/friedenberg/zit/src/foxtrot/cwd_files"
	"github.com/friedenberg/zit/src/golf/typ"
)

type Typ struct {
	cwd_files.CwdTyp
	typ.Named
}
