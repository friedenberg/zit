package kennung

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	chai "github.com/brandondube/tai"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/ohio"
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
	if t.wasSet && !t.tai.Eq(tai{}) {
		t1 = Time{time: t.tai.AsTime().Local()}
		errors.Log().Printf("non empty tai")
	} else {
		errors.Log().Printf("empty tai")
	}

	return
}

func (t Tai) GetGattung() schnittstellen.GattungLike {
	return gattung.Bestandsaufnahme
}

func (t Tai) KennungSansGattungClone() KennungSansGattung {
	return t
}

func (t Tai) KennungSansGattungPtrClone() KennungSansGattungPtr {
	return &t
}

func (t Tai) KennungClone() Kennung {
	return t
}

func (t Tai) KennungPtrClone() KennungPtr {
	return &t
}

func (t Tai) Parts() [3]string {
	a := strings.TrimRight(fmt.Sprintf("%018d", t.tai.Asec), "0")

	if a == "" {
		a = "0"
	}

	return [3]string{strconv.FormatInt(t.tai.Sec, 10), a}
}

func (t Tai) String() string {
	a := strings.TrimRight(fmt.Sprintf("%018d", t.tai.Asec), "0")

	if a == "" {
		a = "0"
	}

	return fmt.Sprintf("%s.%s", strconv.FormatInt(t.tai.Sec, 10), a)
}

func (t Tai) Format(v string) string {
	return t.AsTime().Format(v)
}

func (t *Tai) Set(v string) (err error) {
	t.wasSet = true

	dr := ohio.MakeDelimReader('.', strings.NewReader(v))
	defer ohio.PutDelimReader(dr)

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

			if val == "" {
				break
			}

			if t.tai.Sec, err = strconv.ParseInt(val, 10, 64); err != nil {
				err = errors.Wrapf(err, "failed to parse Sec time: %s", v)
				return
			}

		case 1:
			if err != nil {
				if errors.IsEOF(err) {
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

			t.tai.Asec = pre * int64(math.Pow10(18-len(val)))

		default:
			if errors.IsEOF(err) {
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
	ok = (t.tai.Sec == 0 && t.tai.Asec == 0) || !t.wasSet
	return
}

func (t Tai) IsEmpty() (ok bool) {
	ok = t.IsZero()
	return
}

func (t *Tai) Reset() {
	t.tai.Sec = 0
	t.tai.Asec = 0
	t.wasSet = false
}

func (t *Tai) ResetWith(b Tai) {
	t.tai.Sec = b.tai.Sec
	t.tai.Asec = b.tai.Asec
	t.wasSet = b.wasSet
}

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

func (a Tai) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (t Tai) Equals(t1 Tai) bool {
	if !t.tai.Eq(t1.tai) {
		return false
	}

	return true
}

func (t Tai) Less(t1 Tai) bool {
	return t.tai.Before(t1.tai)
}