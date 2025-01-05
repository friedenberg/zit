package box

type TokenMatcher interface {
	Match(Token) bool
}

type TokenMatcherOp byte

func (tokenMatcher TokenMatcherOp) Match(token Token) bool {
	if token.TokenType != TokenTypeOperator {
		return false
	}

	if token.Contents[0] != byte(tokenMatcher) {
		return false
	}

	return true
}
