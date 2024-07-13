package interfaces

type ConfigGetter interface {
	GetKonfig() Config
}

type ImmutableConfigGetter interface {
	GetImmutableConfig() ImmutableConfig
}

type ImmutableConfig interface {
	GetStoreVersion() StoreVersion
}

type Config interface {
	ImmutableConfig
	UsePredictableHinweisen() bool
	UsePrintTime() bool
	GetFilters() map[string]string
	IsDryRun() bool
	GetTypStringFromExtension(t string) string
}
