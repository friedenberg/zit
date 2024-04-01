package organize_text

import (
	"code.linenisgreat.com/zit/src/bravo/organize_text_mode"
	"code.linenisgreat.com/zit/src/golf/compare_map"
	"code.linenisgreat.com/zit/src/lima/changes"
)

type (
	CompareMap        = compare_map.CompareMap
	SetKeyToMetadatei = compare_map.SetKeyToMetadatei
	Changes           = changes.Changes2
	Change            = changes.Change
)

const (
	ModeOutputOnly     = organize_text_mode.ModeOutputOnly
	ModeInteractive    = organize_text_mode.ModeInteractive
	ModeCommitDirectly = organize_text_mode.ModeCommitDirectly
)

// ChangesFrom            = changes2.ChangesFrom
var ChangesFrom = changes.ChangesFrom2
