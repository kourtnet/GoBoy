// Package bus contains Bus struct which implements memory read/write
// operations and different memory blocks mapping (ROM, Wmemory, VRAM, e.t.c.)
package bus

import "fmt"

const (
	memorySize = 64 * 1024
)

type Bus struct {
	memory []byte
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

func New() *Bus {
	return &Bus{
		memory: make([]byte, memorySize),
	}
}

func (b *Bus) Memory() []byte {
	return b.memory
}

func (b *Bus) Read(addr uint16) (byte, error) {
	if int(addr) >= len(b.memory) {
		return 0, NewOutOfRange(addr)
	}

	return b.memory[addr], nil
}

func (b *Bus) Write(addr uint16, val byte) error {
	if int(addr) >= len(b.memory) {
		return NewOutOfRange(addr)
	}

	b.memory[addr] = val

	return nil
}

// TODO: check where loop call of Read is used to replace
// it with batch
func (b *Bus) ReadBatch(start, end uint16) ([]byte, error) {
	if start > end {
		return b.readBatchSplitted(start, end)
	}

	if int(start) >= len(b.memory) {
		return []byte{}, NewOutOfRange(start)
	}

	if int(end) >= len(b.memory) {
		end = uint16(len(b.memory) - 1)
	}

	return b.memory[start : uint32(end)+1], nil
}

func (b *Bus) readBatchSplitted(start, end uint16) ([]byte, error) {
	start1 := start
	end1 := uint16(len(b.memory) - 1)

	start2 := uint16(0)
	end2 := end

	mem1, err := b.ReadBatch(start1, end1)
	if err != nil {
		return []byte{}, nil
	}

	mem2, err := b.ReadBatch(start2, end2)
	if err != nil {
		return []byte{}, nil
	}

	res := make([]byte, len(mem1)+len(mem2))
	copy(res[:len(mem1)], mem1)
	copy(res[len(mem1):len(mem2)], mem2)

	return res, nil
}
