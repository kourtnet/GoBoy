package cpu

import "fmt"

type registers struct {
	IR byte
	IE byte
	A  byte
	F  byte
	B  byte
	C  byte
	D  byte
	E  byte
	H  byte
	L  byte
	PC uint16
	SP uint16

	temp byte
}

func (r *registers) String() string {
	var res string

	res += fmt.Sprintf("IR: %#x ", r.IR)
	res += fmt.Sprintf("IE: %#x ", r.IE)
	res += fmt.Sprintf("A: %#x ", r.A)
	res += fmt.Sprintf("F: %#x ", r.F)
	res += fmt.Sprintf("B: %#x ", r.B)
	res += fmt.Sprintf("C: %#x ", r.C)
	res += fmt.Sprintf("D: %#x ", r.D)
	res += fmt.Sprintf("E: %#x ", r.E)
	res += fmt.Sprintf("H %#x ", r.H)
	res += fmt.Sprintf("L: %#x ", r.L)
	res += fmt.Sprintf("PC: %#x ", r.PC)
	res += fmt.Sprintf("SP: %#x", r.SP)

	return res
}

func (r *registers) AF() uint16 {
	return uint16(r.A)<<8 + uint16(r.F)
}

func (r *registers) BC() uint16 {
	return uint16(r.B)<<8 + uint16(r.C)
}

func (r *registers) DE() uint16 {
	return uint16(r.D)<<8 + uint16(r.E)
}

func (r *registers) HL() uint16 {
	return uint16(r.H)<<8 + uint16(r.L)
}
