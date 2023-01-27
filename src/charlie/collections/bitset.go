package collections

import (
	"encoding/binary"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
)

const (
	MaxBitsetIdx = 100_000
)

type Bitset interface {
	Equals(Bitset) bool
	Len() int
	Cap() int
	Get(int) bool

	Add(int)
	Del(int)
	set(int, bool)
}

const (
	//For compatibility with 32 bit systems
	//storageInt = uint32
	intSize     = 32
	bytesPerInt = intSize / 4
)

type bitset struct {
	slice []uint32
	lock  *sync.Mutex
}

func MakeBitset(n int) Bitset {
	return makeBitset(n)
}

func makeBitset(n int) (bs *bitset) {
	bs = &bitset{
		slice: make([]uint32, (n+intSize-1)/intSize),
		lock:  &sync.Mutex{},
	}

	return
}

func (b bitset) Cap() int {
	b.lock.Lock()
	defer b.lock.Unlock()

	return int(intSize * cap(b.slice))
}

func (b bitset) len() int {
	return int(intSize * len(b.slice))
}

func (b bitset) Len() int {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.len()
}

func (b bitset) Get(idx int) bool {
	b.lock.Lock()
	defer b.lock.Unlock()

	pos := idx / intSize
	j := uint32(idx % intSize)
	return (b.slice[pos] & (uint32(1) << j)) != 0
}

func (a bitset) Equals(b Bitset) bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	errors.TodoP4("should lock b beyond len call")
	if b.Len() != a.len() {
		return false
	}

	if bb, ok := b.(*bitset); ok {
		bb.lock.Lock()
		defer bb.lock.Unlock()

		for i, av := range a.slice {
			if bb.slice[i] != av {
				return false
			}
		}
	} else {
		errors.TodoP4("improve performance of this")
		for i := 0; i < a.Len(); i++ {
			if a.Get(i) != b.Get(i) {
				return false
			}
		}
	}

	return true
}

func (b *bitset) Add(idx int) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.set(idx, true)
}

func (b *bitset) Del(idx int) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.set(idx, false)
}

func (b *bitset) growIfNecessary(idx int) {
	if b.len() > idx {
		return
	}

	newSize := int(float32(idx) * float32(1.5))

	if newSize > MaxBitsetIdx {
		panic("would grow too large")
	}

	a := makeBitset(int(float32(idx) * float32(1.5)))
	copy(a.slice, b.slice)
	b.slice = a.slice

	return
}

func (b *bitset) set(idx int, value bool) {
	b.growIfNecessary(idx)

	pos := idx / intSize
	j := uint32(idx % intSize)

	if value {
		b.slice[pos] |= (uint32(1) << j)
	} else {
		b.slice[pos] &= ^(uint32(1) << j)
	}
}

func (b bitset) MarshalBinary() (bs []byte, err error) {
	bs = make([]byte, len(b.slice)*bytesPerInt)

	for i, v := range b.slice {
		binary.BigEndian.PutUint32(bs[bytesPerInt*i:], v)
	}

	return
}

func (b *bitset) UnmarshalBinary(bs []byte) (err error) {
	b.slice = make([]uint32, len(bs)/bytesPerInt)
	b.lock = &sync.Mutex{}

	for i := range b.slice {
		b.slice[i] = binary.BigEndian.Uint32(bs[bytesPerInt*i:])
	}

	return
}
