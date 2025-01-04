package box_scanner

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/box"
)

type Token struct {
	Contents []byte
	box.TokenType
}

func (token Token) String() string {
	return fmt.Sprintf("%s:%s", token.TokenType, token.Contents)
}

func (src Token) Clone() (dst Token) {
  dst = src
  dst.Contents = make([]byte, len(src.Contents))
  copy(dst.Contents, src.Contents)
  return
}
