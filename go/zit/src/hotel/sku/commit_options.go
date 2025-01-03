package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type CommitOptions struct {
	StoreOptions
	ids.RepoId
	ids.Clock
	Proto              *Transacted
	DontAddMissingTags bool
	DontAddMissingType bool
}

type StoreOptions struct {
	AddToInventoryList bool
	ApplyProto         bool
	MergeCheckedOut    bool
	RunHooks           bool
	UpdateTai          bool
	Validate           bool
}

func GetStoreOptionsRealizeWithProto() StoreOptions {
	return StoreOptions{
		ApplyProto: true,
		RunHooks:   true,
		UpdateTai:  true,
	}
}

func GetStoreOptionsRealizeSansProto() StoreOptions {
	return StoreOptions{
		RunHooks:  true,
		UpdateTai: true,
	}
}

func GetStoreOptionsReindex() StoreOptions {
	return StoreOptions{
		AddToInventoryList: true,
	}
}

func GetStoreOptionsImport() StoreOptions {
	return StoreOptions{
		AddToInventoryList: true,
		RunHooks:           true,
		UpdateTai:          true,
		Validate:           true,
	}
}

func GetStoreOptionsUpdate() StoreOptions {
	return StoreOptions{
		AddToInventoryList: true,
		RunHooks:           true,
		UpdateTai:          true,
		Validate:           true,
	}
}

func GetStoreOptionsCreate() StoreOptions {
	return StoreOptions{
		AddToInventoryList: true,
		RunHooks:           true,
		ApplyProto:         true,
		UpdateTai:          true,
		Validate:           true,
	}
}

// type Mode byte

// const (
// 	ModeEmpty              = Mode(iota)
// 	ModeAddToInventoryList = Mode(1 << iota) // proper commit
// 	ModeUpdateTai                            // update the tai
// 	ModeLatest                               // only features updates that have no retroactive effects
// 	ModeMergeCheckedOut
// 	ModeApplyProto
// 	ModeHooks

// 	ModeRealizeWithProto = ModeUpdateTai | ModeApplyProto | ModeHooks
// 	ModeRealizeSansProto = ModeUpdateTai | ModeHooks

// 	ModeReindex = ModeLatest
// 	ModeImport  = ModeReindex | ModeAddToInventoryList
// 	ModeCommit  = ModeImport | ModeHooks
// 	ModeCreate  = ModeCommit | ModeApplyProto
// )
