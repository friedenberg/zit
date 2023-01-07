package ts

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	chai "github.com/brandondube/tai"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/format"
)

type tai = chai.TAI

type Tai struct {
	tai
}

func NowTai() Tai {
	return Tai{
		tai: chai.Now(),
	}
}

func TaiFromTimeWithIndex(t1 Time, n int) (t2 Tai) {
	t2.tai = chai.FromTime(t1.time)
	t2.tai.Asec += int64(n * chai.Attosecond)

	return
}

func (t Tai) String() string {
	return fmt.Sprintf(
		"%s.%s",
		strconv.FormatInt(t.tai.Sec, 10),
		strconv.FormatInt(t.tai.Asec, 10),
	)
}

func (t *Tai) Set(v string) (err error) {
	r := bufio.NewReader(strings.NewReader(v))

	if _, err = format.ReadSep(
		'.',
		r,
		func(v string) (err error) {
			if t.tai.Sec, err = strconv.ParseInt(v, 10, 64); err != nil {
				err = errors.Wrapf(err, "failed to parse Sec time: %s", v)
				return
			}

			return
		},
		func(v string) (err error) {
			if t.tai.Asec, err = strconv.ParseInt(v, 10, 64); err != nil {
				err = errors.Wrapf(err, "failed to parse Asec time: %s", v)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// func (t Tai) Sha() sha.Sha {
// 	hash := sha256.New()
// 	sr := strings.NewReader(t.String())

// 	if _, err := io.Copy(hash, sr); err != nil {
// 		errors.PanicIfError(err)
// 	}

// 	return sha.FromHash(hash)
// }

func (t Tai) MarshalText() (text []byte, err error) {
	errors.Err().Printf(t.String())
	text = []byte(t.String())

	return
}

func (t *Tai) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		return
	}

	return
}

func (t Tai) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())

	return
}

func (t *Tai) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		return
	}

	return
}

func (t Tai) Equals(t1 *Tai) bool {
	if t1 == nil {
		return false
	}

	if t != *t1 {
		return false
	}

	return true
}

func (t Tai) Less(t1 Tai) bool {
	return t.tai.Before(t1.tai)
}
