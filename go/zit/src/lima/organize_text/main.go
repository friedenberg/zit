package organize_text

import (
	_ "code.linenisgreat.com/zit/src/kilo/organize_text"
	old "code.linenisgreat.com/zit/src/kilo/organize_text"
	organize_text "code.linenisgreat.com/zit/src/kilo/organize_text2"
)

type (
	CompareMap        = organize_text.CompareMap
	Text              = organize_text.Text
	SetKeyToMetadatei = organize_text.SetKeyToMetadatei
	Options           = organize_text.Options
	Flags             = organize_text.Flags
	Mode              = old.Mode
	ErrorRead         = organize_text.ErrorRead
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
