package ts

import (
	"bufio"
	"fmt"
	"math"
	"strconv"
	"strings"

	chai "github.com/brandondube/tai"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/format"
)

type tai = chai.TAI

type Tai struct {
	wasSet bool
	tai
}

func NowTai() Tai {
	return Tai{
		wasSet: true,
		tai:    chai.Now(),
	}
}

func TaiFromTime(t1 Time) (t2 Tai) {
	t2 = TaiFromTimeWithIndex(t1, 0)
	return
}

func TaiFromTimeWithIndex(t1 Time, n int) (t2 Tai) {
	t2.wasSet = true
	t2.tai = chai.FromTime(t1.time)
	t2.tai.Asec += int64(n * chai.Attosecond)

	return
}

func (t Tai) AsTime() (t1 Time) {
	// if t.wasSet && !t.tai.Eq(tai{}) {
	t1 = Time{time: t.tai.AsTime()}
	// }

	errors.Log().Printf("TODO: %#v -> %#v -> %#v", t, t.tai.AsTime(), t1)

	return
}

func (t Tai) String() string {
	a := strings.TrimRight(fmt.Sprintf("%018d", t.tai.Asec), "0")

	if a == "" {
		a = "0"
	}

	return fmt.Sprintf("%s.%s", strconv.FormatInt(t.tai.Sec, 10), a)
}

func (t *Tai) Set(v string) (err error) {
	t.wasSet = true
	r := bufio.NewReader(strings.NewReader(v))

	if _, err = format.ReadSep(
		'.',
		r,
		func(v string) (err error) {
			v = strings.TrimSpace(v)

			if v == "" {
				return
			}

			if t.tai.Sec, err = strconv.ParseInt(v, 10, 64); err != nil {
				err = errors.Wrapf(err, "failed to parse Sec time: %s", v)
				return
			}

			return
		},
		func(v string) (err error) {
			v = strings.TrimSpace(v)
			v = strings.TrimRight(v, "0")

			if v == "" {
				return
			}

			var pre int64

			if pre, err = strconv.ParseInt(v, 10, 64); err != nil {
				err = errors.Wrapf(err, "failed to parse Asec time: %s", v)
				return
			}

			t.tai.Asec = pre * int64(math.Pow10(18-len(v)))

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

	if !t.tai.Eq(t1.tai) {
		return false
	}

	return true
}

func (t Tai) Less(t1 Tai) bool {
	return t.tai.Before(t1.tai)
}