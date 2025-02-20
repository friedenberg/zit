package commands

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/xdg"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_immutable_io"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("info-repo", &InfoRepo{})
}

type InfoRepo struct {
	command_components.EnvRepo
}

func (cmd InfoRepo) Run(req command.Request) {
	args := req.PopArgs()
	repo := cmd.MakeEnvRepo(req, false)

  // TODO should this be the private config flavor?
	configLoaded := repo.GetConfigPublic()
	c := configLoaded.ImmutableConfig

	if len(args) == 0 {
		args = []string{"store-version"}
	}

	for _, arg := range args {
		switch strings.ToLower(arg) {
		default:
			repo.CancelWithBadRequestf("unsupported info key: %q", arg)

		case "config-immutable":
			if _, err := (config_immutable_io.CoderPublic{}).EncodeTo(
				&configLoaded,
				repo.GetUIFile(),
			); err != nil {
				repo.CancelWithError(err)
			}

		case "store-version":
			repo.GetUI().Print(c.GetStoreVersion())

		case "type":
			repo.GetUI().Print(c.GetRepoType())

		case "id":
			repo.GetUI().Print(c.GetRepoId())

		case "compression-type":
			repo.GetUI().Print(c.GetBlobStoreConfigImmutable().GetBlobCompression())

		case "age-encryption":
			for _, i := range c.GetBlobStoreConfigImmutable().GetBlobEncryption().(*age.Age).Identities {
				repo.GetUI().Print(i)
			}

		case "xdg":
			ecksDeeGee := repo.GetXDG()

			dotenv := xdg.Dotenv{
				XDG: &ecksDeeGee,
			}

			if _, err := dotenv.WriteTo(repo.GetUIFile()); err != nil {
				repo.CancelWithError(err)
			}
		}
	}
}
