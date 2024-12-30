package env

type OptionsGetter interface {
	GetEnvOptions() Options
}

type Options struct {
	UIFileIsStderr bool
}
