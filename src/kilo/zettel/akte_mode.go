package zettel

type AkteMode int

const (
	AkteModeShaOnly = AkteMode(iota)
	AkteModeInlineText
	AkteModeExternalFile
)
