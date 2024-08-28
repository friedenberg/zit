package ids

import (
	"regexp"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
)

const RepoIdRegexString = `^(/)?[-a-z0-9_]+$`

var RepoIdRegex *regexp.Regexp

func init() {
	RepoIdRegex = regexp.MustCompile(RepoIdRegexString)
	register(RepoId{})
}

func MustRepoId(v string) (e *RepoId) {
	e = &RepoId{}

	if err := e.Set(v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func MakeRepoId(v string) (e *RepoId, err error) {
	e = &RepoId{}

	if err = e.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type RepoId struct {
	value string
}

func (k RepoId) IsEmpty() bool {
	return k.value == ""
}

func (k RepoId) GetRepoId() interfaces.RepoId {
	return k
}

func (k RepoId) EqualsRepoId(kg interfaces.RepoIdGetter) bool {
	return kg.GetRepoId().GetRepoIdString() == k.GetRepoIdString()
}

func (k RepoId) GetRepoIdString() string {
	return k.String()
}

func (e *RepoId) Reset() {
	e.value = ""
}

func (e *RepoId) ResetWith(e1 RepoId) {
	e.value = e1.value
}

func (a RepoId) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a RepoId) Equals(b RepoId) bool {
	return a.value == b.value
}

func (o RepoId) GetGenre() interfaces.Genre {
	return genres.Repo
}

func (k RepoId) String() string {
	return k.value
}

func (k RepoId) Parts() [3]string {
	return [3]string{"", "/", k.value}
}

func (k RepoId) GetQueryPrefix() string {
	return "/"
}

func (e *RepoId) Set(v string) (err error) {
	v = strings.TrimPrefix(v, "/")
	v = strings.ToLower(strings.TrimSpace(v))

	if err = ErrOnConfig(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if v == "" {
		return
	}

	if !RepoIdRegex.Match([]byte(v)) {
		err = errors.Errorf("not a valid Kasten: '%s'", v)
		return
	}

	e.value = v

	return
}

func (t RepoId) MarshalText() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *RepoId) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t RepoId) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *RepoId) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
