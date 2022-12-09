package ts

import (
	"crypto/sha256"
	"io"
	"strconv"
	"strings"
	tyme "time"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
)

const (
	Epoch = 1660007128
)

type Time struct {
	time
}

type time = tyme.Time

func Now() Time {
	return Time{
		time: tyme.Now(),
	}
}

func Tyme(t tyme.Time) Time {
	return Time{
		time: t,
	}
}

func (t *Time) MoveForwardIota() {
	t.time = t.time.Add(tyme.Second)
}

func (t Time) String() string {
	return strconv.FormatInt(t.Unix(), 10)
}

func (t *Time) Set(v string) (err error) {
	var n int64

	if n, err = strconv.ParseInt(v, 10, 64); err != nil {
		err = errors.Wrapf(err, "failed to parse time: %s", v)
		return
	}

	t.time = tyme.Unix(n, 0)

	return
}

func (t Time) Kopf() string {
	return strconv.FormatInt((t.Unix()-Epoch)/(60*60*24*30), 10)
}

func (t Time) Schwanz() string {
	return strconv.FormatInt(t.Unix()-Epoch, 10)
}

func (t Time) Sha() sha.Sha {
	hash := sha256.New()
	sr := strings.NewReader(t.String())

	if _, err := io.Copy(hash, sr); err != nil {
		errors.PanicIfError(err)
	}

	return sha.FromHash(hash)
}

func (t Time) MarshalText() (text []byte, err error) {
	errors.Err().Printf(t.String())
	text = []byte(t.String())

	return
}

func (t *Time) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		return
	}

	return
}

func (t Time) Equals(t1 Time) bool {
	return t.Unix() == t1.Unix()
}

func (t Time) Less(t1 Time) bool {
	return t.Unix() < t1.Unix()
}

func (t Time) IsEmpty() bool {
	return t.time.IsZero()
}
