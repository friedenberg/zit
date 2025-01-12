package commands

import (
	"flag"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/delta/xdg"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type RepoInfo struct{}

func init() {
	registerCommand(
		"repo-info",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &RepoInfo{}

			return c
		},
	)
}

// TODO disambiguate this from repo / env
func (cmd RepoInfo) RunWithRepo(repo *repo_local.Repo, args ...string) {
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
