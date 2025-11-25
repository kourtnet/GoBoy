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
	pc uint16
	sp uint16

	// temporary regs
	temp     byte
	tempAddr uint16
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
	res += fmt.Sprintf("PC: %#x ", r.pc)
	res += fmt.Sprintf("SP: %#x", r.sp)

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

func (r *registers) PC() uint16 {
	return r.pc
}

func (r *registers) SP() uint16 {
	return r.sp
}

func (r *registers) TempAddr() uint16 {
	return r.tempAddr
}

func (r *registers) incPC() {
	r.pc++
}

func (r *registers) SetTempAddrMsb(v byte) {
	r.tempAddr += uint16(v) << 8
}

// Lsb set also clears tempAddr because there's no use in
// keeping a Msb part and writing new Lsb
func (r *registers) SetTempAddrLsb(v byte) {
	r.tempAddr = 0
	r.tempAddr = uint16(v)
}
