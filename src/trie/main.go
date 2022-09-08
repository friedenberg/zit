package trie

import (
	"fmt"
)

type Trie struct {
	root node
}

type node struct {
	Count    int
	Children map[byte]node
}

func Make(vs ...fmt.Stringer) (t *Trie) {
	t = &Trie{
		root: node{Children: make(map[byte]node)},
	}

	for _, v := range vs {
		t.Add(v)
	}

	return
}

func (t *Trie) Contains(v fmt.Stringer) bool {
	return t.root.Contains(v.String())
}

func (t *Trie) ShortestUnique(v fmt.Stringer) string {
	return t.root.ShortestUnique(v.String(), 0)
}

func (t *Trie) Add(v fmt.Stringer) {
	if t.Contains(v) {
		return
	}

	t.root.Add(v.String())
}

func (n *node) Add(v string) {
	if len(v) == 0 {
		return
	}

	n.Count += 1

	c := v[0]

	var child node
	ok := false
	child, ok = n.Children[c]

	if !ok {
		child = node{Children: make(map[byte]node)}
	}

	child.Add(v[1:])

	n.Children[c] = child
}

func (n node) Contains(v string) bool {
	if len(v) == 0 {
		return true
	}

	c := v[0]

	child, ok := n.Children[c]

	if ok {
		return child.Contains(v[1:])
	} else {
		return false
	}
}

func (n node) ShortestUnique(v string, loc int) string {
	if len(v)-1 < loc {
		return v
	}

	if n.Count == 1 && n.Contains(v[loc:]) {
		return v[0:loc]
	}

	c := v[loc]

	child, ok := n.Children[c]

	if ok {
		return child.ShortestUnique(v, loc+1)
	} else {
		if len(v)-1 < loc {
			return v
		} else {
			return v[0 : loc+1]
		}
	}
}
