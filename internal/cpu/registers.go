package cpu

import "fmt"

const (
	flagZ uint8 = 0b10000000
	flagN uint8 = 0b01000000
	flagH uint8 = 0b00100000
	flagC uint8 = 0b00010000
)

type Registers struct {
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
	S  byte
	P  byte

	// temporary regs
	Temp8  byte
	temp16 uint16
}

func (r *Registers) String() string {
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
	res += fmt.Sprintf("SP: %#x", r.SP())

	return res
}

func (r *Registers) r8R8ToR16(r1, r2 byte) uint16 {
	return uint16(r1)<<8 + uint16(r2)
}

func (r *Registers) r16Msb(reg uint16) byte {
	return byte(reg >> 8)
}

func (r *Registers) r16Lsb(reg uint16) byte {
	return byte(reg & 0x00FF)
}

func (r *Registers) AF() uint16 {
	return r.r8R8ToR16(r.A, r.F)
}

func (r *Registers) SetAF(val uint16) {
	r.A = r.r16Msb(val)
	r.F = r.r16Lsb(val)
}

func (r *Registers) BC() uint16 {
	return r.r8R8ToR16(r.B, r.C)
}

func (r *Registers) IncBC() {
	BC := r.BC()
	BC++

	r.SetBC(BC)
}

func (r *Registers) DecBC() {
	BC := r.BC()
	BC--

	r.SetBC(BC)
}

func (r *Registers) SetBC(val uint16) {
	r.B = r.r16Msb(val)
	r.C = r.r16Lsb(val)
}

func (r *Registers) DE() uint16 {
	return r.r8R8ToR16(r.D, r.E)
}

func (r *Registers) IncDE() {
	DE := r.DE()
	DE++

	r.SetDE(DE)
}

func (r *Registers) DecDE() {
	DE := r.DE()
	DE--

	r.SetDE(DE)
}

func (r *Registers) SetDE(val uint16) {
	r.D = r.r16Msb(val)
	r.E = r.r16Lsb(val)
}

func (r *Registers) HL() uint16 {
	return r.r8R8ToR16(r.H, r.L)
}

func (r *Registers) SetHL(val uint16) {
	r.H = r.r16Msb(val)
	r.L = r.r16Lsb(val)
}

func (r *Registers) IncHL() {
	HL := r.HL()
	HL++

	r.SetHL(HL)
}

func (r *Registers) DecHL() {
	HL := r.HL()
	HL--

	r.SetHL(HL)
}

func (r *Registers) PC() uint16 {
	return r.pc
}

func (r *Registers) SP() uint16 {
	return r.r8R8ToR16(r.S, r.P)
}

func (r *Registers) SetSP(val uint16) {
	r.S = r.r16Msb(val)
	r.P = r.r16Lsb(val)
}

func (r *Registers) IncSP() {
	SP := r.SP()
	SP++

	r.SetSP(SP)
}

func (r *Registers) DecSP() {
	SP := r.SP()
	SP--

	r.SetSP(SP)
}

func (r *Registers) Temp16() uint16 {
	return r.temp16
}

func (r *Registers) IncPC() {
	r.pc++
}

func (r *Registers) setFlag(mask uint8, val bool) {
	if val {
		r.F |= mask
	} else {
		r.F &^= mask
	}
}

func (r *Registers) getFlag(mask uint8) bool {
	return (r.F & mask) != 0
}

func (r *Registers) SetFlagZ(val bool) {
	r.setFlag(flagZ, val)
}

func (r *Registers) FlagZ() bool {
	return r.getFlag(flagZ)
}

func (r *Registers) SetFlagN(val bool) {
	r.setFlag(flagN, val)
}

func (r *Registers) FlagN() bool {
	return r.getFlag(flagN)
}

func (r *Registers) SetFlagH(val bool) {
	r.setFlag(flagH, val)
}

func (r *Registers) FlagH() bool {
	return r.getFlag(flagH)
}

func (r *Registers) SetFlagC(val bool) {
	r.setFlag(flagC, val)
}

func (r *Registers) FlagC() bool {
	return r.getFlag(flagC)
}

func (r *Registers) SetTemp16Msb(v byte) {
	r.temp16 = r.r8R8ToR16(v, r.r16Lsb(r.temp16))
}

// Lsb set also clears prev tempAddr value and writes FF into Msb
// it is required for indirect addresses where Msb part is always FF
func (r *Registers) SetTemp16Lsb(v byte) {
	r.temp16 = r.r8R8ToR16(0xFF, v)
}
