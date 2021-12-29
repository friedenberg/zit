package mime

// import (
// 	"mime"
// 	"path"
// 	"strings"
// )

// var (
// 	Markdown = Mime{value: "md"}
// )

// type Mime struct {
// 	value string
// }

// func MakeMime(v string) (m Mime, err error) {
// 	if err = m.Set(v); err != nil {
// 		err = _Error(err)
// 		return
// 	}

// 	return
// }

// func (m *Mime) Set(v string) (err error) {
// 	if v == "" {
// 		return
// 	}

// 	var mt string

// 	s1 := v

// 	if !strings.HasPrefix(s1, ".") {
// 		s1 = "." + s1
// 	}

// 	mt = mime.TypeByExtension(path.Ext(s1))

// 	if mt == "" {
// 		err = _Errorf("unknown mime type: %s", v)
// 		return
// 	}

// 	m.value = mt

// 	return
// }

// func (m Mime) IsEmpty() bool {
// 	return m.value == ""
// }

// func (m Mime) IsTextType() bool {
// 	return strings.Split(m.String(), "/")[0] == "text"
// }

// func (m Mime) String() string {
// 	if m.IsEmpty() {
// 		return Markdown.String()
// 	}

// 	return m.value
// }

// func (m Mime) PreferredExtension() (e string) {
// 	var es []string
// 	var err error

// 	if es, err = mime.ExtensionsByType(m.String()); err != nil || es == nil || len(es) < 1 {
// 		panic(err)
// 	}

// 	if err != nil {
// 		panic(_Error(err))
// 	} else if es == nil || len(es) < 1 {
// 		panic(_Errorf("expected at least one mime type, but go none"))
// 	}

// 	return es[0]
// }
