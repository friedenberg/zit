package repo_local

type (
	Options       int
	OptionsGetter interface {
		GetLocalRepoOptions() Options
	}
)

const (
	OptionsEmpty                = Options(iota)
	OptionsAllowConfigReadError = Options(1 << iota)
	OptionsUIFileIsStderr
)

func (o Options) GetAllowConfigReadError() bool {
	return o&OptionsAllowConfigReadError != 0
}

func (o Options) GetUIFileIsStderr() bool {
	return o&OptionsUIFileIsStderr != 0
}
