package konfig

import (
	"io/ioutil"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
)

type KonfigTag struct {
	AddToNewZettels bool `toml:"add-to-new-zettels"`
	Hide            bool `toml:"hide"`
}

type KonfigTyp struct {
	FormatScript ScriptConfig `toml:"format-script"`
	InlineAkte   bool         `toml:"inline-akte" default:"true"`
	ExecCommand  ScriptConfig `toml:"exec-command"`
}

type Konfig struct {
	Cli
	Toml
	Logger errors.Logger
}

func LoadKonfig(p string) (c Konfig, err error) {
	// c = DefaultKonfig()

	var f *os.File

	if f, err = open_file_guard.Open(p); err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}

		err = errors.Wrap(err)
		return
	}

	defer open_file_guard.Close(f)

	doc, err := ioutil.ReadAll(f)

	defer func() {
		if r := recover(); r != nil {
			c = Konfig{}
			err = errors.Errorf("toml unmarshalling panicked: %q", r)
		}
	}()

	var tc Toml
	err = toml.Unmarshal([]byte(doc), &tc)

	c.Toml = tc

	if err != nil {
		err = errors.Errorf("failed to parse config: %s", err)
		return
	}

	return
}
