package commands

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/delta/xdg"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("repo-info", &RepoInfo{})
}

type RepoInfo struct {
	command_components.RepoLayout
}

func (cmd RepoInfo) Run(dep command.Request) {
	args := dep.Args()
	repo := cmd.MakeRepoLayout(dep, false)
	c := repo.GetConfig()

	if len(args) == 0 {
		args = []string{"store-version"}
	}

	for _, arg := range args {
		switch strings.ToLower(arg) {
		default:
			repo.CancelWithBadRequestf("unsupported info key: %q", arg)

		case "store-version":
			repo.GetUI().Print(c.GetStoreVersion())

		case "compression-type":
			repo.GetUI().Print(c.GetBlobStoreImmutableConfig().GetCompressionType())

		case "age-encryption":
			for _, i := range c.GetBlobStoreImmutableConfig().GetAgeEncryption().Identities {
				repo.GetUI().Print(i)
			}

		case "xdg":
			ecksDeeGee := repo.GetDirLayout().GetXDG()

			dotenv := xdg.Dotenv{
				XDG: &ecksDeeGee,
			}

			if _, err := dotenv.WriteTo(repo.GetOutFile()); err != nil {
				repo.CancelWithError(err)
			}
		}
	}
}
