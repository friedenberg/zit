package interfaces

type ImmutableConfigGetter interface {
	GetImmutableConfig() ImmutableConfig
}

type ImmutableConfig interface {
	GetStoreVersion() StoreVersion
}

type MutableConfigDryRun interface {
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
	UsePredictableZettelIds() bool
	MutableConfigDryRun
}

type Config interface {
	MutableConfig
	ImmutableConfig
	GetTypeStringFromExtension(t string) string
}
