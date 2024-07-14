package env

type (
	Options       int
	OptionsGetter interface {
		GetEnvironmentInitializeOptions() Options
	}
)

const (
	OptionsEmpty                = Options(iota)
	OptionsAllowConfigReadError = Options(1 << iota)
)

func (o Options) GetAllowConfigReadError() bool {
	return o&OptionsAllowConfigReadError != 0
}
