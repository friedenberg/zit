package schnittstellen

type KonfigGetter interface {
	GetKonfig() Konfig
}

type AngeborenGetter interface {
	GetAngeboren() Angeboren
}

type Angeboren interface {
	GetStoreVersion() StoreVersion
}

type Konfig interface {
	Angeboren
	UsePredictableHinweisen() bool
	UsePrintTime() bool
	GetFilters() map[string]string
	IsDryRun() bool
}
