package cpu

import "fmt"

type registers struct {
	IR byte
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
	Temp byte
	// TODO: rename. Now tempAddr is used not only as an adress,
	// but as a temporal register also. So the name is a bit
	// confusing. Suggestion: rename temp to temp8 and tempAddr
	// to temp16
	tempAddr uint16
}

func (r *registers) String() string {
	var res string

	res += fmt.Sprintf("IR: %#x ", r.IR)
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

func (r *registers) r8R8ToR16(r1, r2 byte) uint16 {
	return uint16(r1)<<8 + uint16(r2)
}

func (r *registers) r16Msb(reg uint16) byte {
	return byte(reg >> 8)
}

func (r *registers) r16Lsb(reg uint16) byte {
	return byte(reg & 0x00FF)
}

func (r *registers) AF() uint16 {
	return r.r8R8ToR16(r.A, r.F)
}

func (r *registers) BC() uint16 {
	return r.r8R8ToR16(r.B, r.C)
}

func (r *registers) IncBC() {
	BC := r.BC()
	BC++

	r.SetBC(BC)
}

func (r *registers) DecBC() {
	BC := r.BC()
	BC--

	r.SetBC(BC)
}

func (r *registers) SetBC(val uint16) {
	r.B = r.r16Msb(val)
	r.C = r.r16Lsb(val)
}

func (r *registers) DE() uint16 {
	return r.r8R8ToR16(r.D, r.E)
}

func (r *registers) IncDE() {
	DE := r.DE()
	DE++

	r.SetBC(DE)
}

func (r *registers) DecDE() {
	DE := r.DE()
	DE--

	r.SetBC(DE)
}

func (r *registers) SetDE(val uint16) {
	r.D = r.r16Msb(val)
	r.E = r.r16Lsb(val)
}

func (r *registers) HL() uint16 {
	return r.r8R8ToR16(r.H, r.L)
}

func (r *registers) SetHL(val uint16) {
	r.H = r.r16Msb(val)
	r.L = r.r16Lsb(val)
}

func (r *registers) IncHL() {
	HL := r.HL()
	HL++

	r.SetHL(HL)
}

func (r *registers) DecHL() {
	HL := r.HL()
	HL--

	r.SetHL(HL)
}

func (r *registers) PC() uint16 {
	return r.pc
}

func (r *registers) SP() uint16 {
	return r.sp
}

func (r *registers) SetSP(val uint16) {
	r.sp = val
}

func (r *registers) DecSP() {
	r.sp--
}

func (r *registers) IncSP() {
	r.sp++
}

func (r *registers) SPL() byte {
	return r.r16Lsb(r.sp)
}

func (r *registers) SPH() byte {
	return r.r16Msb(r.sp)
}

func (r *registers) TempAddr() uint16 {
	return r.tempAddr
}

func (r *registers) IncTempAddr() {
	r.tempAddr++
}

func (r *registers) IncPC() {
	r.pc++
}

func (r *registers) SetTempAddrMsb(v byte) {
	r.tempAddr = r.r8R8ToR16(v, r.r16Lsb(r.tempAddr))
}

// Lsb set also clears prev tempAddr value and writes FF into Msb
// it is required for indirect addresses where Msb part is always FF
func (r *registers) SetTempAddrLsb(v byte) {
	r.tempAddr = r.r8R8ToR16(0xFF, v)
}
