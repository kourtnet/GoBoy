// Package bus contains Bus struct which implements memory read/write
// operations and different memory blocks mapping (ROM, WRAM, VRAM, e.t.c.)
package bus

import "fmt"

const (
	ramSize = 64 * 1024
)

type Bus struct {
	// TODO: implement separate memory blocks
	RAM [ramSize]byte
}

func (b *Bus) Read(addr uint16) (byte, error) {
	if int(addr) >= len(b.RAM) {
		return 0, fmt.Errorf("address: %#x if out of memory range", addr)
	}

	return b.RAM[addr], nil
}

func (b *Bus) Write(addr uint16, val byte) error {
	if int(addr) >= len(b.RAM) {
		return fmt.Errorf("address: %#x if out of memory range", addr)
	}

	b.RAM[addr] = val

	return nil
}
