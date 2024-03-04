package organize_text

import (
	_ "code.linenisgreat.com/zit/src/kilo/organize_text"
	old "code.linenisgreat.com/zit/src/kilo/organize_text"
	organize_text "code.linenisgreat.com/zit/src/kilo/organize_text2"
	"code.linenisgreat.com/zit/src/lima/changes2"
)

type (
	CompareMap        = changes2.CompareMap
	SetKeyToMetadatei = changes2.SetKeyToMetadatei
	Text              = organize_text.Text
	Options           = organize_text.Options
	Flags             = organize_text.Flags
	Mode              = old.Mode
	ErrorRead         = organize_text.ErrorRead
	Changes           = changes2.Changes
)

const (
	ModeOutputOnly     = old.ModeOutputOnly
	ModeInteractive    = old.ModeInteractive
	ModeCommitDirectly = old.ModeCommitDirectly
)

var (
	MakeFlags              = organize_text.MakeFlags
	MakeFlagsWithMetadatei = organize_text.MakeFlagsWithMetadatei
	New                    = organize_text.New
)
