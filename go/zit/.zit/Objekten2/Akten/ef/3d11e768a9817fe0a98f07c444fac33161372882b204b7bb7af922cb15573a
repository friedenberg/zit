package query_spec

import "fmt"

type TokenParts struct {
	Left, Right []byte
}

func (tp TokenParts) String() string {
	return fmt.Sprintf("Left: %q, Right: %q", string(tp.Left), string(tp.Right))
}

func (tp *TokenParts) Reset() {
	tp.Left = nil
	tp.Right = nil
}

func (src TokenParts) Clone() (dst TokenParts) {
	dst = TokenParts{
		Left:  make([]byte, len(src.Left)),
		Right: make([]byte, len(src.Right)),
	}

	copy(dst.Left, src.Left)
	copy(dst.Right, src.Right)

	return
}
