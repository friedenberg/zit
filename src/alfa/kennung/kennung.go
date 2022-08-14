package kennung

import (
	"math"
	"strconv"
)

type Int = uint64
type Float = float64

type Kennung struct {
	Left, Right Int
}

func (p Kennung) Equals(p1 Kennung) bool {
	return p.Left == p1.Left && p.Right == p1.Right
}

func Extrema(n Int) Kennung {
	n1 := Float(n)
	return Kennung{
		Left:  Int(((n1 - 1) * n1 / 2) + 1),
		Right: Int(n1 * ((n1 + 1) / 2)),
	}
}

func (p *Kennung) SetCoordinates(left, right string) (err error) {
	l, err := strconv.ParseUint(left, 10, 64)

	if err != nil {
		return
	}

	r, err := strconv.ParseUint(right, 10, 64)

	if err != nil {
		return
	}

	p.Left = l
	p.Right = r

	return
}

func (p *Kennung) Set(id string) (err error) {
	i, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		return
	}

	p.SetInt(i)

	return
}

func (p *Kennung) SetInt(id Int) {
	n := math.Round(math.Sqrt(Float(id) * 2))
	ext := Extrema(Int(n))

	p.Left = id - ext.Left
	p.Right = ext.Right - id
}

func (p Kennung) Id() Int {
	n := p.Left + p.Right + 1
	ext := Extrema(n)
	return ext.Left + p.Left
}
