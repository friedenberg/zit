package query_spec

import (
	"io"
	"unicode"
	"unicode/utf8"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
)

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

type TokenScanner struct {
	rs                io.RuneScanner
	tokenTypeProbably TokenType
	tokenType         TokenType
	token             catgut.String
	parts             TokenParts
	err               error
	unscan            bool
	n                 int64
}

func (ts *TokenScanner) Reset(r io.RuneScanner) {
	ts.rs = r
	ts.token.Reset()
	ts.tokenType = TokenTypeIncomplete
	ts.tokenTypeProbably = TokenTypeIncomplete
	ts.parts.Reset()
	ts.err = nil
	ts.unscan = false
	ts.n = 0
}

func (ts *TokenScanner) Unscan() {
	ts.unscan = true
}

func (ts *TokenScanner) ScanOnly(tokenType TokenType) (ok bool) {
	ok = ts.Scan()

	if !ok {
		return
	}

	if ts.tokenType != tokenType {
		ok = false
		ts.unscan = true
		return
	}

	return
}

func (ts *TokenScanner) CanScan() (ok bool) {
	if ts.unscan {
		return true
	}

	return ts.err == nil
}

func (ts *TokenScanner) ScanIdentifierLikeSkipSpaces() (ok bool) {
	if ts.unscan {
		ok = true
		ts.unscan = false
		return
	}

	if ts.err == io.EOF {
		return
	}

	afterFirst := false
	ok = true

	ts.token.Reset()
	ts.tokenType = TokenTypeIncomplete
	ts.tokenTypeProbably = TokenTypeIncomplete
	ts.parts.Reset()

	for {
		var r rune
		var n int

		r, n, ts.err = ts.rs.ReadRune()
		ts.n += int64(n)

		if ts.err != nil {
			if ts.err == io.EOF {
				ok = ts.token.Len() > 0
				ts.tokenType = ts.tokenTypeProbably
				ts.parts.Left = ts.token.Bytes()
			}

			return
		}

		isOperator := unicode.IsSpace(r) || r == '[' || r == ']'
		isSpace := unicode.IsSpace(r)

		switch {
		case r == '"' || r == '\'':
			ts.tokenType = TokenTypeLiteral

			if !ts.consumeLiteralOrFieldValue(r, TokenTypeLiteral, &ts.parts.Left) {
				ok = false
				return
			}

			return

		case !afterFirst && isOperator:
			if isSpace {
				if !ts.consumeSpaces() {
					ok = false
					return
				}

				continue
			} else {
				ts.token.WriteRune(r)
				ts.tokenType = TokenTypeOperator
				return
			}

		case !isOperator:
			ts.tokenTypeProbably = TokenTypeIdentifier
			ts.token.WriteRune(r)
			afterFirst = true
			continue

		default: // wasSplitRune && afterFirst
			ts.parts.Left = ts.token.Bytes()

			if r == '=' {
				if !ts.consumeField(r) {
					ok = false
					return
				}

				return
			}

			if ts.err = ts.rs.UnreadRune(); ts.err != nil {
				ts.err = errors.Wrapf(ts.err, "%c", r)
				ok = false
			}

			ts.n -= int64(utf8.RuneLen(r))
			ts.tokenType = TokenTypeIdentifier

			return
		}
	}
}

func (ts *TokenScanner) Scan() (ok bool) {
	if ts.unscan {
		ok = true
		ts.unscan = false
		return
	}

	if ts.err == io.EOF {
		return
	}

	afterFirst := false
	ok = true

	ts.token.Reset()
	ts.tokenType = TokenTypeIncomplete
	ts.tokenTypeProbably = TokenTypeIncomplete
	ts.parts.Reset()

	for {
		var r rune
		var n int

		r, n, ts.err = ts.rs.ReadRune()
		ts.n += int64(n)

		if ts.err != nil {
			if ts.err == io.EOF {
				ok = ts.token.Len() > 0
				ts.tokenType = ts.tokenTypeProbably
				ts.parts.Left = ts.token.Bytes()
			}

			return
		}

		isOperator := IsOperator(r)
		isSpace := unicode.IsSpace(r)

		switch {
		case r == '"' || r == '\'':
			ts.tokenType = TokenTypeLiteral

			if !ts.consumeLiteralOrFieldValue(r, TokenTypeLiteral, &ts.parts.Left) {
				ok = false
				return
			}

			return

		case !afterFirst && isOperator:
			ts.token.WriteRune(r)

			if isSpace {
				if !ts.consumeSpaces() {
					ok = false
					return
				}
			}

			ts.tokenType = TokenTypeOperator

			return

		case !isOperator:
			ts.tokenTypeProbably = TokenTypeIdentifier
			ts.token.WriteRune(r)
			afterFirst = true
			continue

		default: // wasSplitRune && afterFirst
			ts.parts.Left = ts.token.Bytes()

			if r == '=' {
				if !ts.consumeField(r) {
					ok = false
					return
				}

				return
			}

			if ts.err = ts.rs.UnreadRune(); ts.err != nil {
				ts.err = errors.Wrapf(ts.err, "%c", r)
				ok = false
			}

			ts.n -= int64(utf8.RuneLen(r))
			ts.tokenType = TokenTypeIdentifier

			return
		}
	}
}

func (ts *TokenScanner) consumeSpaces() (ok bool) {
	ok = true

	for {
		var r rune
		var n int

		r, n, ts.err = ts.rs.ReadRune()
		ts.n += int64(n)

		if ts.err != nil {
			ok = false
			return
		}

		if unicode.IsSpace(r) {
			continue
		}

		if ts.err = ts.rs.UnreadRune(); ts.err != nil {
			ok = false
			ts.err = errors.Wrapf(ts.err, "%c", r)
		}

		ts.n -= int64(utf8.RuneLen(r))

		return
	}
}

func (ts *TokenScanner) consumeLiteralOrFieldValue(
	start rune,
	tt TokenType,
	partLocation *[]byte,
) (ok bool) {
	ok = true

	ts.token.WriteRune(start)
	lastWasBackslash := false

	idx := ts.token.Len()

	for {
		var r rune
		var n int

		r, n, ts.err = ts.rs.ReadRune()
		ts.n += int64(n)

		if ts.err != nil {
			ok = false
			return
		}

		ts.token.WriteRune(r)

		if r != start || lastWasBackslash {
			lastWasBackslash = r == '\\'
			continue
		}

		ts.tokenType = tt
		*partLocation = ts.token.Bytes()[idx : ts.token.Len()-1]

		return
	}
}

func (ts *TokenScanner) consumeField(start rune) bool {
	ts.token.WriteRune(start)
	ok := ts.consumeIdentifierLike(TokenTypeField, &ts.parts.Right)
	return ok
}

func (ts *TokenScanner) consumeIdentifierLike(
	tt TokenType,
	partLocation *[]byte,
) (ok bool) {
	ok = true

	for {
		var r rune
		var n int

		r, n, ts.err = ts.rs.ReadRune()
		ts.n += int64(n)

		if ts.err != nil {
			if ts.err == io.EOF {
				ok = ts.token.Len() > 0
				ts.tokenType = tt
			}

			return
		}

		isOperator := IsOperator(r)

		switch {
		case r == '"' || r == '\'':
			if !ts.consumeLiteralOrFieldValue(r, TokenTypeLiteral, partLocation) {
				ok = false
				return
			}

			ts.tokenType = tt

			return

		case !isOperator:
			ts.token.WriteRune(r)
			continue

		default: // wasSplitRune && afterFirst
			if ts.err = ts.rs.UnreadRune(); ts.err != nil {
				ts.err = errors.Wrapf(ts.err, "%c", r)
				ok = false
			}

			ts.n -= int64(utf8.RuneLen(r))
			ts.tokenType = tt

			return
		}
	}
}

func (ts *TokenScanner) GetToken() *catgut.String {
	return &ts.token
}

func (ts *TokenScanner) GetTokenAndType() (*catgut.String, TokenType) {
	return &ts.token, ts.tokenType
}

func (ts *TokenScanner) GetTokenAndTypeAndParts() (*catgut.String, TokenType, TokenParts) {
	return &ts.token, ts.tokenType, ts.parts
}

func (ts *TokenScanner) N() int64 {
	return ts.n
}

func (ts *TokenScanner) Error() error {
	if ts.err == io.EOF {
		return nil
	}

	return ts.err
}
