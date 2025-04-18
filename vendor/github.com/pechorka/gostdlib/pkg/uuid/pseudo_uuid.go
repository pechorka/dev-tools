package uuid

import (
	"sync"
	"time"
)

func NewV4Pseudo() (UUID, error) {
	pcg := pcgPool.Get().(*pcg)
	uuid, err := NewV4(pcg)
	pcgPool.Put(pcg)
	return uuid, err
}

func MustV4Pseudo() UUID {
	pcg := pcgPool.Get().(*pcg)
	uuid := MustV4(pcg)
	pcgPool.Put(pcg)
	return uuid
}

func NewV4PseudoString() (string, error) {
	pcg := pcgPool.Get().(*pcg)
	uuid, err := NewV4String(pcg)
	pcgPool.Put(pcg)
	return uuid, err
}

func MustV4PseudoString() string {
	pcg := pcgPool.Get().(*pcg)
	uuid := MustV4String(pcg)
	pcgPool.Put(pcg)
	return uuid
}

func NewV7Pseudo() (UUID, error) {
	pcg := pcgPool.Get().(*pcg)
	uuid, err := NewV7(pcg)
	pcgPool.Put(pcg)
	return uuid, err
}

func MustV7Pseudo() UUID {
	pcg := pcgPool.Get().(*pcg)
	uuid := MustV7(pcg)
	pcgPool.Put(pcg)
	return uuid
}

func NewV7PseudoString() (string, error) {
	pcg := pcgPool.Get().(*pcg)
	uuid, err := NewV7String(pcg)
	pcgPool.Put(pcg)
	return uuid, err
}

func MustV7PseudoString() string {
	pcg := pcgPool.Get().(*pcg)
	uuid := MustV7String(pcg)
	pcgPool.Put(pcg)
	return uuid
}

var pcgPool = sync.Pool{
	New: func() any {
		return newPCG()
	},
}

// https://www.pcg-random.org/download.html#minimal-c-implementation
type pcg struct {
	State uint64
	Inc   uint64
}

func newPCG() *pcg {
	initState := uint64(time.Now().UnixNano())
	initSeq := uint64(initState>>32) ^ uint64(initState<<32)
	return &pcg{
		State: initState,
		Inc:   initSeq,
	}
}

func (p *pcg) Read(b []byte) (int, error) {
	val := p.random()
	pos := 4

	for i := range b {
		if pos == 0 {
			val = p.random()
			pos = 4
		}
		b[i] = byte(val)
		val >>= 8
		pos--
	}

	return len(b), nil
}

func random() uint32 {
	pcg := pcgPool.Get().(*pcg)
	val := pcg.random()
	pcgPool.Put(pcg)
	return val
}

func (p *pcg) random() uint32 {
	oldstate := p.State
	p.State = oldstate*6364136223846793005 + p.Inc
	xorshifted := uint32(((oldstate >> 18) ^ oldstate) >> 27)
	rot := uint32(oldstate >> 59)
	return (xorshifted >> rot) | (xorshifted << ((-rot) & 31))
}
