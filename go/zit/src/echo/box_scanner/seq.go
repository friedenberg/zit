package box_scanner

import (
	"strings"
)

//go:generate stringer -type=SeqType
type SeqType int

const (
	SeqTypeUnknown    = SeqType(iota)
	SeqTypeIdentifier // one/uno, tag, !type, /browser/bookmark-1, @abcd
	SeqTypeField      // url=blah blob=hello contents="wow" contents="wow with\" quote"
)

type Seq []Token

func (a Seq) EqualsSeq(b Seq) bool {
	return false
}

func (seq *Seq) Add(tokenType TokenType, contents []byte) {
	*seq = append(*seq, Token{TokenType: tokenType, Contents: contents})
}

func (seq *Seq) AddToken(token Token) {
	*seq = append(*seq, token)
}

func (seq Seq) String() string {
	var sb strings.Builder

	for _, t := range seq {
		sb.Write(t.Contents)
	}

	return sb.String()
}

func (src Seq) Clone() (dst Seq) {
	dst = make(Seq, len(src))

	for i := range src {
		dst[i] = src[i].Clone()
	}

	return
}

func (seq *Seq) Reset() {
	*seq = Seq([]Token{})
}
