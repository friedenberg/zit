package read_write_repo_local

type (
	Options       int
	OptionsGetter interface {
		GetLocalRepoOptions() Options
	}
)

const (
	OptionsEmpty                = Options(iota)
	OptionsAllowConfigReadError = Options(1 << iota)
)

func (o Options) GetAllowConfigReadError() bool {
	return o&OptionsAllowConfigReadError != 0
}
