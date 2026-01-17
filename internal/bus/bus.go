// Package bus contains Bus struct which implements memory read/write
// operations and different memory blocks mapping (ROM, WRAM, VRAM, e.t.c.)
package bus

import "fmt"

const (
	ramSize = 64 * 1024
)

type Bus struct {
	RAM [ramSize]byte
}

type OutOfRange struct {
	str string
}

func NewOutOfRange(addr uint16) OutOfRange {
	return OutOfRange{
		str: fmt.Sprintf(
			"address %04X is out or memory range",
			addr,
		),
	}
}

func (o OutOfRange) Error() string {
	return o.str
}

func (b *Bus) Read(addr uint16) (byte, error) {
	if int(addr) >= len(b.RAM) {
		return 0, NewOutOfRange(addr)
	}

	return b.RAM[addr], nil
}

func (b *Bus) Write(addr uint16, val byte) error {
	if int(addr) >= len(b.RAM) {
		return NewOutOfRange(addr)
	}

	b.RAM[addr] = val

	return nil
}

// TODO: check where loop call of Read is used to replace
// it with batch
func (b *Bus) ReadBatch(start, end uint16) ([]byte, error) {
	if start > end {
		return []byte{}, fmt.Errorf("start address is bigger than the end")
	}

	if int(start) >= len(b.RAM) {
		return []byte{}, NewOutOfRange(start)
	}

	if int(end) >= len(b.RAM) {
		end = uint16(len(b.RAM) - 1)
	}

	end++

	return b.RAM[start:end], nil
}
