package cpu

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

const memSizeInBytes = 16

type suiteBus struct {
	memory []byte
}

func newSuiteBus(ROMPath string, spareMemory []byte) (*suiteBus, error) {
	fileStat, err := os.Stat(ROMPath)
	if err != nil {
		return nil, err
	}

	size := fileStat.Size() + memSizeInBytes

	sb := &suiteBus{
		memory: make([]byte, size),
	}

	file, err := os.Open(ROMPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = file.Read(sb.memory[:memSizeInBytes])
	if err != nil {
		return nil, err
	}

	copy(sb.memory[memSizeInBytes:], spareMemory)

	return sb, nil
}

func (b *suiteBus) Read(addr uint16) (byte, error) {
	if int(addr) >= len(b.memory) {
		return 0, fmt.Errorf("address: %#x if out of memory range", addr)
	}

	return b.memory[addr], nil
}

func (b *suiteBus) Write(addr uint16, val byte) error {
	if int(addr) >= len(b.memory) {
		return fmt.Errorf("address: %#x if out of memory range", addr)
	}

	b.memory[addr] = val

	return nil
}

type snapshot struct {
	regs   registers
	memory []byte
}

func newSnapshot(A, F, B, C, D, E, H, L byte, PC, SP uint16, memory []byte) snapshot {
	return snapshot{
		regs: registers{
			A:  A,
			F:  F,
			B:  B,
			C:  C,
			D:  D,
			E:  E,
			H:  H,
			L:  L,
			pc: PC,
			sp: SP,
		},
		memory: memory,
	}
}

func (s *snapshot) Equal(cpu CPU) (bool, error) {
	s.regs.Temp = cpu.registers.Temp
	s.regs.tempAddr = cpu.registers.tempAddr
	s.regs.IR = cpu.registers.IR

	if *cpu.registers != s.regs {
		return false, nil
	}

	//	for i := range s.memory {
	//		val, err := cpu.bus.Read(uint16(i))
	//		if err != nil {
	//			return false, err
	//		}
	//
	//		if s.memory[i] != val {
	//			return false, nil
	//		}
	//	}

	return true, nil
}

func hexValue(b byte) byte {
	switch {
	case '0' <= b && b <= '9':
		return b - '0'
	case 'A' <= b && b <= 'F':
		return 10 + b - 'A'
	case 'a' <= b && b <= 'f':
		return 10 + b - 'a'
	}

	// invalid
	return 16
}

func hexBytesToByte(high, low byte) byte {
	h := hexValue(high)
	l := hexValue(low)

	return byte(h<<4 | l)
}

func hexBytesToWord(high1, high2, low1, low2 byte) uint16 {
	h1 := hexValue(high1)
	h2 := hexValue(high2)
	l1 := hexValue(low1)
	l2 := hexValue(low2)

	high := byte(h1<<4 | h2)
	low := byte(l1<<4 | l2)
	return uint16(high)<<8 | uint16(low)
}

func hexSliceToBytes(hex []byte) []byte {
	res := make([]byte, len(hex)/2)

	for i := 0; i < len(hex); i += 2 {
		res[0] = hexBytesToByte(hex[i], hex[i+1])
	}

	return res
}

func newSnapshots(snapPath string) ([]snapshot, error) {
	file, err := os.Open(snapPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	snaps := []snapshot{}

	regexpPattern := `^[0-9A-Fa-f]{2}( [0-9A-Fa-f]{2}){7} [0-9A-Fa-f]{4} [0-9A-Fa-f]{4} [0-9A-Fa-f]{32}$`
	re := regexp.MustCompile(regexpPattern)

	lineNum := 1
	for scanner.Scan() {
		row := scanner.Bytes()

		if row[0] == '\n' || row[0] == '#' {
			continue
		}

		if !re.Match(row) {
			return nil, fmt.Errorf("line: %d wrong snapshot syntax", lineNum)
		}

		snap := newSnapshot(
			hexBytesToByte(row[0], row[1]),
			hexBytesToByte(row[3], row[4]),
			hexBytesToByte(row[6], row[7]),
			hexBytesToByte(row[9], row[10]),
			hexBytesToByte(row[12], row[13]),
			hexBytesToByte(row[15], row[16]),
			hexBytesToByte(row[18], row[19]),
			hexBytesToByte(row[21], row[22]),
			hexBytesToWord(row[24], row[25], row[26], row[27]),
			hexBytesToWord(row[29], row[30], row[31], row[32]),
			hexSliceToBytes(row[:memSizeInBytes*2]),
		)

		snaps = append(snaps, snap)

		lineNum++
	}

	return snaps, nil
}

type testSuite struct {
	cpu     CPU
	bus     *suiteBus
	snaps   []snapshot
	stepNum int
	err     error
}

func NewTestSuite(ROMPath, SnapPath string) (*testSuite, error) {
	snaps, err := newSnapshots(SnapPath)
	if err != nil {
		return nil, err
	}

	bus, err := newSuiteBus(ROMPath, snaps[0].memory)
	if err != nil {
		return nil, err
	}

	cpu, err := New(bus)
	if err != nil {
		return nil, err
	}

	cpu.registers.A = snaps[0].regs.A
	cpu.registers.F = snaps[0].regs.F
	cpu.registers.B = snaps[0].regs.B
	cpu.registers.C = snaps[0].regs.C
	cpu.registers.D = snaps[0].regs.D
	cpu.registers.E = snaps[0].regs.E
	cpu.registers.H = snaps[0].regs.H
	cpu.registers.L = snaps[0].regs.L
	cpu.registers.pc = snaps[0].regs.pc
	cpu.registers.sp = snaps[0].regs.sp

	suite := &testSuite{
		cpu:   cpu,
		bus:   bus,
		snaps: snaps,
	}

	return suite, nil
}

func (s *testSuite) step() bool {
	s.err = nil
	s.stepNum++

	if s.stepNum >= len(s.snaps) {
		return false
	}

	_, err := s.cpu.Step()
	if err != nil {
		s.err = err
		return true
	}

	eq, err := s.snaps[s.stepNum].Equal(s.cpu)
	if err != nil {
		s.err = err
		return true
	}

	if !eq {
		errFormat := "data mismatch at step %d\n\tcpu registers:   %v\n\tsuite registers: %v"
		s.err = fmt.Errorf(errFormat,
			s.stepNum, *s.cpu.registers, s.snaps[s.stepNum].regs)
		return true
	}

	return true
}
