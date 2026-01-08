package cpu

import (
	"errors"
)

func (cpu *CPU) readAddr(addr uint16) {
	var err error

	cpu.registers.Temp, err = cpu.bus.Read(addr)
	if err != nil {
		cpu.internalErr = err
		return
	}
}

func (cpu *CPU) readPCAddrAndInc() {
	cpu.readAddr(cpu.registers.pc)
	if cpu.internalErr != nil {
		return
	}

	cpu.registers.incPC()
}

func (cpu *CPU) readR16Addr(rName string) func() {
	r := cpu.determine16Reg(rName)
	if cpu.internalErr != nil {
		return func() {}
	}

	return func() {
		cpu.readAddr(r())
	}
}

func (cpu *CPU) readR8Addr(rName byte) func() {
	r := 0xFF00 + uint16(*cpu.determine8Reg(rName))
	if cpu.internalErr != nil {
		return func() {}
	}

	return func() {
		cpu.readAddr(r)
	}
}

func (cpu *CPU) readAddrLsb() {
	cpu.readPCAddrAndInc()
	cpu.registers.SetTempAddrLsb(cpu.registers.Temp)
}

func (cpu *CPU) readAddrMsb() {
	cpu.readPCAddrAndInc()
	cpu.registers.SetTempAddrMsb(cpu.registers.Temp)
}

func (cpu *CPU) determine8Reg(rName byte) *byte {
	var errReg byte

	switch rName {
	case 'A':
		return &cpu.registers.A
	case 'B':
		return &cpu.registers.B
	case 'C':
		return &cpu.registers.C
	case 'D':
		return &cpu.registers.D
	case 'E':
		return &cpu.registers.E
	case 'H':
		return &cpu.registers.H
	case 'L':
		return &cpu.registers.L
	// Special processing for easier memory reading with immediate addr
	case 'T':
		return &cpu.registers.Temp
	default:
		cpu.internalErr = errors.New("unknown register name " + string(rName))
		return &errReg
	}
}

func (cpu *CPU) determine16Reg(rName string) func() uint16 {
	switch rName {
	case "AF":
		return cpu.registers.AF
	case "BC":
		return cpu.registers.BC
	case "DE":
		return cpu.registers.DE
	case "HL":
		return cpu.registers.HL
	case "PC":
		return cpu.registers.PC
	case "SP":
		return cpu.registers.SP
	// Special processing for easier memory reading with immediate addr
	case "TempAddr":
		return cpu.registers.TempAddr
	default:
		cpu.internalErr = errors.New("unknown register name " + rName)
		return func() uint16 { return 0 }
	}
}

func (cpu *CPU) determine16RegSetter(rName string) func(uint16) {
	switch rName {
	case "BC":
		return cpu.registers.SetBC
	case "DE":
		return cpu.registers.SetDE
	case "HL":
		return cpu.registers.SetHL
	case "SP":
		return cpu.registers.SetSP
	default:
		cpu.internalErr = errors.New("unknown register name " + rName)
		return func(uint16) {}
	}
}

func (cpu *CPU) nop() {
}

// TODO: ADD FLAGS AND OUT OF RANGE ARITHMETICS
// WARNING: IS UNUSED NOW
func (cpu *CPU) incA() {
	cpu.registers.A++
}

// this part of package contains funcs with undescore symbols in
// tneir names. It's intentional and is used only for functions
// that contain CPU instruction in their names.
// For example:
// LD R8, R8 -> ld_R8_R8
// LD [R8], R8 -> ld_R8_Addr_R8
func (cpu *CPU) ld_R8_R8(r1Name, r2Name byte) func() {
	if r1Name == r2Name {
		return func() {}
	}

	r1 := cpu.determine8Reg(r1Name)
	r2 := cpu.determine8Reg(r2Name)

	return func() {
		*r1 = *r2
	}
}

func (cpu *CPU) ld_R8_Temp(rName byte) func() {
	r := cpu.determine8Reg(rName)
	return func() {
		*r = cpu.registers.Temp
	}
}

func (cpu *CPU) ld_R16_Addr_R8(r16Name string, r8Name byte) func() {
	r16 := cpu.determine16Reg(r16Name)
	r8 := cpu.determine8Reg(r8Name)

	return func() {
		cpu.bus.Write(r16(), *r8)
	}
}

func (cpu *CPU) ld_R8_Addr_R8(r1Name, r2Name byte) func() {
	r1 := cpu.determine8Reg(r1Name)
	r2 := cpu.determine8Reg(r2Name)

	return func() {
		cpu.bus.Write(0xFF00+uint16(*r1), *r2)
	}
}

// there's no LD [R16], N8 instructions rather then with HL
func (cpu *CPU) ldHLAddrTemp() {
	cpu.bus.Write(cpu.registers.HL(), cpu.registers.Temp)
}

func (cpu *CPU) readHLAddrDec() {
	cpu.readR16Addr("HL")()
	cpu.registers.decHL()
}

func (cpu *CPU) readHLAddrInc() {
	cpu.readR16Addr("HL")()
	cpu.registers.incHL()
}

func (cpu *CPU) ld_HL_Addr_A_Dec() {
	cpu.bus.Write(cpu.registers.HL(), cpu.registers.A)
	cpu.registers.decHL()
}

func (cpu *CPU) ld_HL_Addr_A_Inc() {
	cpu.bus.Write(cpu.registers.HL(), cpu.registers.A)
	cpu.registers.incHL()
}

func (cpu *CPU) ld_R16_R16(r1Name, r2Name string) func() {
	r1 := cpu.determine16RegSetter(r1Name)
	r2 := cpu.determine16Reg(r2Name)

	return func() {
		r1(r2())
	}
}

func (cpu *CPU) ld_TempAddr_SPL() {
	cpu.bus.Write(cpu.registers.TempAddr(), cpu.registers.SPL())
}

// +1 in address means this func is used only with ld_TempAddr_SPL for
// writing SP value into memory
func (cpu *CPU) ld_TempAddr_SPH() {
	cpu.bus.Write(cpu.registers.TempAddr()+1, cpu.registers.SPH())
}
