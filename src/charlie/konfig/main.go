package konfig

import (
	"io/ioutil"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	toml_package "github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/files"
)

type Konfig struct {
	Cli
	toml
	Compiled
}

func Make(p string, kc Cli) (c Konfig, err error) {
	c.Compiled = MakeDefaultCompiled()
	c.Cli = kc
	// c = DefaultKonfig()

	var f *os.File

	if f, err = files.Open(p); err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}

		err = errors.Wrap(err)
		return
	}

	defer files.Close(f)

	doc, err := ioutil.ReadAll(f)

	defer func() {
		if r := recover(); r != nil {
			c = Konfig{}
			err = errors.Errorf("toml unmarshalling panicked: %q", r)
		}
	}()

	if err = toml_package.Unmarshal([]byte(doc), &c.toml); err != nil {
		err = errors.Errorf("failed to parse config: %s", err)
		return
	}

	if c.Compiled, err = makeCompiled(c.toml); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
