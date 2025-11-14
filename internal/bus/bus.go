// Package bus contains Bus struct which implements memory read/write
// operations and different memory blocks mapping (ROM, WRAM, VRAM, e.t.c.)
package bus

const (
	ramSize = 64 * 1024
)

type Bus struct {
	// TODO: implement separate memory blocks
	RAM [ramSize]byte
}

func (b *Bus) Read(addr uint16) byte {
	// TODO: index out of range handling
	return b.RAM[addr]
}

func (b *Bus) Write(addr uint16, val byte) {
	// TODO: index out of range handling
	b.RAM[addr] = val
}
