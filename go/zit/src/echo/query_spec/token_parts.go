package query_spec

type TokenParts struct {
	Left, Right []byte
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