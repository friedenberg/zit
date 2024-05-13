package etiketten_path

import "strings"

type StringForward Path

func (p *StringForward) String() string {
	var sb strings.Builder

	afterFirst := false
	l := (*Path)(p).Len()
	for i := l - 1; i >= 0; i-- {
		if afterFirst {
			sb.WriteString(" -> ")
		}

		afterFirst = true

		s := (*p)[i]
		sb.Write(s.Bytes())
	}

	return sb.String()
}

type StringBackward Path

func (p *StringBackward) String() string {
	var sb strings.Builder

	afterFirst := false
	l := (*Path)(p).Len()

	for i := 0; i < l; i++ {
		if afterFirst {
			sb.WriteString(" -> ")
		}

		afterFirst = true

		s := (*p)[i]
		sb.Write(s.Bytes())
	}

	return sb.String()
}
