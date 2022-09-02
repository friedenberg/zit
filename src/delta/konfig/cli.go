package konfig

import (
	"flag"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Cli struct {
	BasePath             string
	Debug                bool
	Verbose              bool
	DryRun               bool
	AllowMissingHinweis  bool
	CheckoutCacheEnabled bool
	IncludeHidden        bool
	PredictableHinweisen bool
}

func (c *Cli) AddToFlags(f *flag.FlagSet) {
	f.StringVar(&c.BasePath, "dir-zit", "", "")
	f.BoolVar(&c.Debug, "debug", false, "")
	f.BoolVar(&c.Verbose, "verbose", false, "")
	f.BoolVar(&c.DryRun, "dry-run", false, "")
	f.BoolVar(&c.CheckoutCacheEnabled, "checkout-cache-enabled", false, "")
	f.BoolVar(&c.AllowMissingHinweis, "allow-missing-hinweis", false, "")
	f.BoolVar(&c.IncludeHidden, "include-hidden", false, "include zettels that have hidden etiketten")
	f.BoolVar(&c.PredictableHinweisen, "predictable-hinweisen", false, "don't randomly select new hinweisen")
}

func (c Cli) DirZit() (p string, err error) {
	if c.BasePath == "" {
		if p, err = os.Getwd(); err != nil {
			err = errors.Error(err)
			return
		}
	} else {
		p = c.BasePath
	}

	return
}

func (c Cli) KonfigPath() (p string, err error) {
	// var usr *user.User

	// if usr, err = user.Current(); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	// p = path.Join(
	// 	usr.HomeDir,
	// 	".config",
	// 	"zettelkasten",
	// 	"config.toml",
	// )

	if p, err = c.DirZit(); err != nil {
		err = errors.Error(err)
		return
	}

	p = path.Join(p, ".zit", "Konfig")

	return
}

func (c Cli) Konfig() (k Konfig, err error) {
	if c.Verbose {
		errors.SetVerbose()
	} else {
		// logz.SetOutput(ioutil.Discard)
	}

	var p string

	if p, err = c.KonfigPath(); err != nil {
		err = errors.Error(err)
		return
	}

	if k, err = LoadKonfig(p); err != nil {
		err = errors.Error(err)
		return
	}

	k.Cli = c

	return
}

func DefaultCli() (c Cli) {
	return
}
