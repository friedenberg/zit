package ts

import (
	"crypto/sha256"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/bravo/stdprinter"
)

const (
	Epoch = 1660007128
)

type Time struct {
	time.Time
}

func Now() Time {
	return Time{
		Time: time.Now(),
	}
}

func (t Time) String() string {
	return strconv.FormatInt(t.Unix(), 10)
}

func (t *Time) Set(v string) (err error) {
	var n int64

	if n, err = strconv.ParseInt(v, 10, 64); err != nil {
		err = errors.Wrapped(err, "failed to parse time: %s", v)
		return
	}

	t.Time = time.Unix(n, 0)

	return
}

func (t Time) Head() string {
	return strconv.FormatInt((t.Unix()-Epoch)/(60*60*24*30), 10)
}

func (t Time) Tail() string {
	return strconv.FormatInt(t.Unix()-Epoch, 10)
}

func (t Time) Sha() sha.Sha {
	hash := sha256.New()
	sr := strings.NewReader(t.String())

	if _, err := io.Copy(hash, sr); err != nil {
		stdprinter.PanicIfError(err)
	}

	return sha.FromHash(hash)
}

func (t Time) MarshalText() (text []byte, err error) {
	text = []byte(t.String())

	return
}

func (t *Time) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		return
	}

	return
}
