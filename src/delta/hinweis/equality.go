package hinweis

func (a inner) Equals(b *Hinweis) bool {
	if a.Left != b.Left {
		return false
	}

	if a.Right != b.Right {
		return false
	}

	return true
}
