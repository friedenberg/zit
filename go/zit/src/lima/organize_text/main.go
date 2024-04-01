package organize_text

import (
	"code.linenisgreat.com/zit/src/bravo/organize_text_mode"
	"code.linenisgreat.com/zit/src/golf/compare_map"
	organize_text "code.linenisgreat.com/zit/src/kilo/organize_text2"
	"code.linenisgreat.com/zit/src/lima/changes2"
)

type (
	CompareMap        = compare_map.CompareMap
	SetKeyToMetadatei = compare_map.SetKeyToMetadatei
	Text              = organize_text.Text
	Options           = organize_text.Options
	Flags             = organize_text.Flags
	Mode              = organize_text_mode.Mode
	ErrorRead         = organize_text.ErrorRead
	Changes           = changes2.Changes2
	Change            = changes2.Change
)

const (
	ModeOutputOnly     = organize_text_mode.ModeOutputOnly
	ModeInteractive    = organize_text_mode.ModeInteractive
	ModeCommitDirectly = organize_text_mode.ModeCommitDirectly
)

var (
	MakeFlags              = organize_text.MakeFlags
	MakeFlagsWithMetadatei = organize_text.MakeFlagsWithMetadatei
	New                    = organize_text.New
	// ChangesFrom            = changes2.ChangesFrom
	ChangesFrom = changes2.ChangesFrom2
)
