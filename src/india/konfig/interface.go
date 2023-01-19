package konfig

type Getter interface {
	GetKonfig() Compiled
}

type PtrGetter interface {
	GetKonfigPtr() *Compiled
}
