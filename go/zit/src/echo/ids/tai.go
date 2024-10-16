package ids

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/charlie/delim_io"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/delta/thyme"
	chai "github.com/brandondube/tai"
)

type tai = chai.TAI

func init() {
	register(Tai{})
	collections_value.RegisterGobValue[Tai](nil)
}

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

func TaiFromTime(t1 thyme.Time) (t2 Tai) {
	t2 = TaiFromTimeWithIndex(t1, 0)
	return
}

func TaiFromTime1(t1 time.Time) (t2 Tai) {
	t2.wasSet = true
	t2.tai = chai.FromTime(t1)

	return
}

func TaiFromTimeWithIndex(t1 thyme.Time, n int) (t2 Tai) {
	t2.wasSet = true
	t2.tai = chai.FromTime(t1.GetTime())
	t2.Asec += int64(n * chai.Attosecond)

	return
}

func (t Tai) AsTime() (t1 thyme.Time) {
	t1 = thyme.Tyme(t.tai.AsTime().Local())
	return
}

func (a Tai) Before(b Tai) bool {
	return a.tai.Before(b.tai)
}

func (a Tai) After(b Tai) bool {
	return a.tai.After(b.tai)
}

func (t Tai) GetGenre() interfaces.Genre {
	return genres.InventoryList
}

func (t Tai) Parts() [3]string {
	a := strings.TrimRight(fmt.Sprintf("%018d", t.Asec), "0")

	if a == "" {
		a = "0"
	}

	return [3]string{strconv.FormatInt(t.Sec, 10), ".", a}
}

func (t Tai) String() (v string) {
	a := strings.TrimRight(fmt.Sprintf("%018d", t.Asec), "0")

	if a == "" {
		a = "0"
	}

	v = fmt.Sprintf("%s.%s", strconv.FormatInt(t.Sec, 10), a)

	// if v == "0.0" {
	// 	panic("empty tai")
	// }

	return
}

func (t Tai) StringDefaultFormat() string {
	f := string_format_writer.StringFormatDateTime + ".000000000"
	return t.Format(f)
}

func (t Tai) Format(v string) string {
	return t.AsTime().Format(v)
}

func (t *Tai) SetFromRFC3339(v string) (err error) {
	t.wasSet = true

	var t1 time.Time

	if t1, err = thyme.Parse(thyme.RFC3339, v); err != nil {
		err = errors.Wrap(err)
		return
	}

	*t = TaiFromTime1(t1)

	return
}

func (t *Tai) Set(v string) (err error) {
	t.wasSet = true

	dr := delim_io.Make('.', strings.NewReader(v))
	defer delim_io.PutReader(dr)

	idx := 0
	var val string

	for {
		val, err = dr.ReadOneString()

		switch idx {
		case 0:
			if err != nil {
				err = errors.Wrap(err)
				return
			}

			val = strings.TrimSpace(val)

			if len(val) == 0 {
				break
			}

			if t.Sec, err = strconv.ParseInt(val, 10, 64); err != nil {
				err = errors.Wrapf(err, "failed to parse Sec time: %s", v)
				return
			}

		case 1:
			if err != nil {
				if err == io.EOF {
					err = nil
				} else {
					err = errors.Wrap(err)
					return
				}
			}

			val = strings.TrimSpace(val)
			val = strings.TrimRight(val, "0")

			if val == "" {
				break
			}

			var pre int64

			if pre, err = strconv.ParseInt(val, 10, 64); err != nil {
				err = errors.Wrapf(err, "failed to parse Asec time: %s", val)
				return
			}

			t.Asec = pre * int64(math.Pow10(18-len(val)))

		default:
			if err == io.EOF {
				err = nil
			} else {
				err = errors.Errorf("expected no more elements but got %s", val)
			}

			return
		}

		idx++
	}
}

// func (t Tai) Sha() sha.Sha {
// 	hash := sha256.New()
// 	sr := strings.NewReader(t.String())

// 	if _, err := io.Copy(hash, sr); err != nil {
// 		errors.PanicIfError(err)
// 	}

// 	return sha.FromHash(hash)
// }

func (t Tai) IsZero() (ok bool) {
	ok = (t.Sec == 0 && t.Asec == 0) || !t.wasSet
	return
}

func (t Tai) IsEmpty() (ok bool) {
	ok = t.IsZero()
	return
}

func (t Tai) GetTai() Tai {
	return t
}

func (t *Tai) Reset() {
	t.Sec = 0
	t.Asec = 0
	t.wasSet = false
}

func (t *Tai) ResetWith(b Tai) {
	t.Sec = b.Sec
	t.Asec = b.Asec
	t.wasSet = b.wasSet
}

func (t Tai) WriteTo(w io.Writer) (n int64, err error) {
	b := make([]byte, binary.MaxVarintLen64*2)
	binary.PutVarint(b[:binary.MaxVarintLen64], t.Sec)
	binary.PutVarint(b[binary.MaxVarintLen64:], t.Asec)
	var n1 int
	n1, err = ohio.WriteAllOrDieTrying(w, b)
	n += int64(n1)
	return
}

func (t *Tai) ReadFrom(r io.Reader) (n int64, err error) {
	b := make([]byte, binary.MaxVarintLen64*2)

	var n1 int
	n1, err = ohio.ReadAllOrDieTrying(r, b)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	t.wasSet = true
	t.Sec, _ = binary.Varint(b[:binary.MaxVarintLen64])
	t.Asec, _ = binary.Varint(b[binary.MaxVarintLen64:])

	return
}

func (t Tai) MarshalText() (text []byte, err error) {
	ui.Err().Printf(t.String())
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

func (a Tai) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (t Tai) Equals(t1 Tai) bool {
	if !t.Eq(t1.tai) {
		return false
	}

	return true
}

func (t Tai) Less(t1 Tai) bool {
	return t.Before(t1)
}

type TaiRFC3339Value Tai

func (t *TaiRFC3339Value) Set(v string) (err error) {
	t.wasSet = true

	var t1 time.Time

	if t1, err = thyme.Parse(thyme.RFC3339, v); err != nil {
		err = errors.Wrap(err)
		return
	}

	*t = TaiRFC3339Value(TaiFromTime1(t1))

	return
}

func (t *TaiRFC3339Value) String() string {
	return t.Format(thyme.RFC3339)
}
