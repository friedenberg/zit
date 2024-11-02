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

type MutableStoredConfig interface {
	ConfigGetFilters
}

type MutableConfig interface {
	MutableStoredConfig
	UsePrintTime() bool
	UsePredictableHinweisen() bool
	ConfigDryRun
}

type Config interface {
	MutableConfig
	ImmutableConfig
	GetTypeStringFromExtension(t string) string
}
