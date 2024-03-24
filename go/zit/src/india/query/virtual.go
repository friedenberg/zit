package query

import "sync"

type VirtualStore interface {
	Init() error
	Matcher
}

type Virtual struct {
	init sync.Once
	VirtualStore
	Kennung
}

func (v *Virtual) ContainsMatchable() bool {
	return false
}
