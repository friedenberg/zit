package konfig

import (
	"bufio"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/files"
)

type Konfig struct {
	Cli
	tomlKonfig
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

	br := bufio.NewReader(f)

	if err = c.tryParseToml(br); err != nil {
		err = errors.Wrap(err)
		return
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				c = Konfig{}
				err = errors.Errorf("toml unmarshalling panicked: %q", r)
			}
		}()

		td := toml.NewDecoder(br)

		if err = td.Decode(&c.tomlKonfig); err != nil {
			err = errors.Errorf("failed to parse config: %s", err)
			return
		}
	}()

	if c.Compiled, err = makeCompiled(c.tomlKonfig); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Konfig) tryParseToml(br *bufio.Reader) (err error) {
	defer func() {
		if r := recover(); r != nil {
			c = &Konfig{}
			err = errors.Errorf("toml unmarshalling panicked: %q", r)
		}
	}()

	td := toml.NewDecoder(br)

	if err = td.Decode(&c.tomlKonfig); err != nil {
		err = errors.Errorf("failed to parse config: %s", err)
		return
	}

	return
}
