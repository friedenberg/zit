package umwelt

type (
	Options       int
	OptionsGetter interface {
		GetUmweltInitializeOptions() Options
	}
)

const (
	OptionsEmpty                = Options(iota)
	OptionsAllowKonfigReadError = Options(1 << iota)
)

func (o Options) GetAllowKonfigReadError() bool {
	return o&OptionsAllowKonfigReadError != 0
}
