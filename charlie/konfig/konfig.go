package konfig

import (
	"io/ioutil"
	"os"
)

type KonfigTag struct {
	AddToNewZettels bool `toml:"add-to-new-zettels"`
	Hide            bool `toml:"hide"`
}

type Konfig struct {
	Cli
	Toml
	Logger _Logger
}

func LoadKonfig(p string) (c Konfig, err error) {
	// c = DefaultKonfig()

	var f *os.File

	if f, err = _Open(p); err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}

		err = _Error(err)
		return
	}

	defer _Close(f)

	doc, err := ioutil.ReadAll(f)

	defer func() {
		if r := recover(); r != nil {
			c = Konfig{}
			err = _Errorf("toml unmarshalling panicked: %q", r)
		}
	}()

	var tc Toml
	err = _TomlUnmarshal([]byte(doc), &tc)

	c.Toml = tc

	if err != nil {
		err = _Errorf("failed to parse config: %s", err)
		return
	}

	return
}
