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

type ConfigDryRun interface {
	IsDryRun() bool
	SetDryRun(bool)
}

type ConfigGetFilters interface {
	GetFilters() map[string]string
}

type Config interface {
	ImmutableConfig
	UsePredictableHinweisen() bool
	UsePrintTime() bool
	GetFilters() map[string]string
	GetTypeStringFromExtension(t string) string
	ConfigDryRun
	ConfigGetFilters
}
