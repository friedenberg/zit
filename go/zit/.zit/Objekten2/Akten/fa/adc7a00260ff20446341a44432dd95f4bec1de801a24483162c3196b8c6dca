package query_spec

import (
	"io"
	"unicode"
	"unicode/utf8"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/token_types"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
)

type TokenScanner struct {
	io.RuneScanner
	tokenTypeProbably token_types.TokenType
	tokenType         token_types.TokenType
	token             catgut.String
	parts             TokenParts
	err               error
	unscan            []rune
	n                 int64
	lastRune          rune
}

func (ts *TokenScanner) Reset(r io.RuneScanner) {
	ts.RuneScanner = r
	ts.token.Reset()
	ts.tokenType = token_types.TypeIncomplete
	ts.tokenTypeProbably = token_types.TypeIncomplete
	ts.parts.Reset()
	ts.err = nil
	ts.unscan = nil
	ts.n = 0
}

func (ts *TokenScanner) ReadRune() (r rune, n int, err error) {
	if len(ts.unscan) > 0 {
		r = ts.unscan[0]
		n = utf8.RuneLen(r)
		ts.unscan = ts.unscan[1:]
		return
	}

	ts.lastRune, n, err = ts.RuneScanner.ReadRune()
	ts.n += int64(n)

	return ts.lastRune, n, err
}

// TODO add support for unscan
func (ts *TokenScanner) UnreadRune() (err error) {
	err = ts.RuneScanner.UnreadRune()

	if err == nil {
		ts.n -= int64(utf8.RuneLen(ts.lastRune))
	}

	return
}

func (ts *TokenScanner) Unscan() {
	ts.unscan = []rune(string(ts.token.Bytes()))
}

func (ts *TokenScanner) ScanOnly(tokenType token_types.TokenType) (ok bool) {
	ok = ts.Scan()

	if !ok {
		return
	}

	if ts.tokenType != tokenType {
		ok = false
		ts.unscan = []rune(string(ts.token.Bytes()))
		return
	}

	return
}

func (ts *TokenScanner) CanScan() (ok bool) {
	if len(ts.unscan) > 0 {
		return true
	}

	return ts.err == nil
}

func (ts *TokenScanner) ScanIdentifierLikeSkipSpaces() (ok bool) {
	if len(ts.unscan) > 0 {
		ok = true
		ts.unscan = nil
		return
	}

	if ts.err == io.EOF {
		return
	}

	afterFirst := false
	ok = true

	ts.token.Reset()
	ts.tokenType = token_types.TypeIncomplete
	ts.tokenTypeProbably = token_types.TypeIncomplete
	ts.parts.Reset()

	for {
		var r rune

		r, _, ts.err = ts.ReadRune()

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
			ts.tokenType = token_types.TypeLiteral

			if !ts.consumeLiteralOrFieldValue(
				r,
				token_types.TypeLiteral,
				&ts.parts.Left,
			) {
				ok = false
				return
			}

			return

		case !afterFirst && isOperator:
			if isSpace {
				if !ts.ConsumeSpacesOrErrorOnFalse() {
					ok = false
					return
				}

				continue
			} else {
				ts.token.WriteRune(r)
				ts.tokenType = token_types.TypeOperator
				return
			}

		case !isOperator:
			ts.tokenTypeProbably = token_types.TypeIdentifier
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

			if ts.err = ts.UnreadRune(); ts.err != nil {
				ts.err = errors.Wrapf(ts.err, "%c", r)
				ok = false
			}

			ts.tokenType = token_types.TypeIdentifier

			return
		}
	}
}

func (ts *TokenScanner) ScanSkipSpace() (ok bool) {
	if !ts.ConsumeSpacesOrErrorOnFalse() {
		return
	}

	ok = ts.Scan()

	return
}

func (ts *TokenScanner) Scan() (ok bool) {
	return ts.scan(true)
}

func (ts *TokenScanner) ScanDotAllowedInIdentifiers() (ok bool) {
	return ts.scan(false)
}

func (ts *TokenScanner) scan(dotOperatorAsSplit bool) (ok bool) {
	if len(ts.unscan) > 0 {
		ok = true
		ts.unscan = nil
		return
	}

	if ts.err == io.EOF {
		return
	}

	afterFirst := false
	ok = true

	ts.token.Reset()
	ts.tokenType = token_types.TypeIncomplete
	ts.tokenTypeProbably = token_types.TypeIncomplete
	ts.parts.Reset()

	for {
		var r rune

		r, _, ts.err = ts.ReadRune()

		if ts.err != nil {
			if ts.err == io.EOF {
				ok = ts.token.Len() > 0
				ts.tokenType = ts.tokenTypeProbably
				ts.parts.Left = ts.token.Bytes()
			}

			return
		}

		isOperator := isOperator(r, !dotOperatorAsSplit)
		isSpace := unicode.IsSpace(r)

		switch {
		case r == '"' || r == '\'':
			ts.tokenType = token_types.TypeLiteral

			if !ts.consumeLiteralOrFieldValue(
				r,
				token_types.TypeLiteral,
				&ts.parts.Left,
			) {
				ok = false
				return
			}

			return

		case !afterFirst && isOperator:
			ts.token.WriteRune(r)

			if isSpace {
				if !ts.ConsumeSpacesOrErrorOnFalse() {
					ok = false
					return
				}
			}

			ts.tokenType = token_types.TypeOperator

			return

		case !isOperator:
			ts.tokenTypeProbably = token_types.TypeIdentifier
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

			if ts.err = ts.UnreadRune(); ts.err != nil {
				ts.err = errors.Wrapf(ts.err, "%c", r)
				ok = false
			}

			ts.tokenType = token_types.TypeIdentifier

			return
		}
	}
}

// Consumes any spaces currently available in the underlying RuneReader. If this
// returns false, it means that a read error has occurred, not that no spaces
// were consumed.
func (ts *TokenScanner) ConsumeSpacesOrErrorOnFalse() (ok bool) {
	for _, r := range ts.unscan {
		if ts.err != nil {
			ok = false
			return
		}

		if unicode.IsSpace(r) {
			continue
		}

		if ts.err = ts.UnreadRune(); ts.err != nil {
			ok = false
			ts.err = errors.Wrapf(ts.err, "%c", r)
		}

		ok = true
	}

	ts.unscan = nil

	ok = true

	for {
		var r rune

		r, _, ts.err = ts.ReadRune()

		if ts.err != nil {
			ok = false
			return
		}

		if unicode.IsSpace(r) {
			continue
		}

		if ts.err = ts.UnreadRune(); ts.err != nil {
			ok = false
			ts.err = errors.Wrapf(ts.err, "%c", r)
		}

		return
	}
}

// TODO add support for ellipis
func (ts *TokenScanner) consumeLiteralOrFieldValue(
	start rune,
	tt token_types.TokenType,
	partLocation *[]byte,
) (ok bool) {
	ok = true

	ts.token.WriteRune(start)
	lastWasBackslash := false

	idx := ts.token.Len()

	for {
		var r rune

		r, _, ts.err = ts.ReadRune()

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
	ok := ts.consumeIdentifierLike(token_types.TypeField, &ts.parts.Right)
	return ok
}

// TODO add support for ellipsis
func (ts *TokenScanner) consumeIdentifierLike(
	tt token_types.TokenType,
	partLocation *[]byte,
) (ok bool) {
	ok = true

	idx := ts.token.Len()

	for {
		var r rune

		r, _, ts.err = ts.ReadRune()

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
			if !ts.consumeLiteralOrFieldValue(r, token_types.TypeLiteral, partLocation) {
				ok = false
				return
			}

			ts.tokenType = tt

			return

		case !isOperator:
			ts.token.WriteRune(r)
			continue

		default: // wasSplitRune && afterFirst
			*partLocation = ts.token.Bytes()[idx:ts.token.Len()]

			if ts.err = ts.UnreadRune(); ts.err != nil {
				ts.err = errors.Wrapf(ts.err, "%c", r)
				ok = false
			}

			ts.tokenType = tt

			return
		}
	}
}

func (ts *TokenScanner) GetToken() *catgut.String {
	return &ts.token
}

func (ts *TokenScanner) GetTokenAndType() (*catgut.String, token_types.TokenType) {
	return &ts.token, ts.tokenType
}

func (ts *TokenScanner) GetTokenAndTypeAndParts() (*catgut.String, token_types.TokenType, TokenParts) {
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
