package cpu

import (
	"errors"
	"fmt"
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

	cpu.registers.IncPC()
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
	case 'F':
		return &cpu.registers.F
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

// TODO: test
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
	case "AF":
		return cpu.registers.SetAF
	default:
		cpu.internalErr = errors.New("unknown register name " + rName)
		return func(uint16) {}
	}
}

func (cpu *CPU) determine16RegInc(rName string) func() {
	switch rName {
	case "BC":
		return cpu.registers.IncBC
	case "DE":
		return cpu.registers.IncDE
	case "HL":
		return cpu.registers.IncHL
	case "SP":
		return cpu.registers.IncSP
	default:
		cpu.internalErr = errors.New("unknown register name " + rName)
		return func() {}
	}
}

func (cpu *CPU) determine16RegDec(rName string) func() {
	switch rName {
	case "BC":
		return cpu.registers.DecBC
	case "DE":
		return cpu.registers.DecDE
	case "HL":
		return cpu.registers.DecHL
	case "SP":
		return cpu.registers.DecSP
	default:
		cpu.internalErr = errors.New("unknown register name " + rName)
		return func() {}
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
	cpu.registers.DecHL()
}

func (cpu *CPU) readHLAddrInc() {
	cpu.readR16Addr("HL")()
	cpu.registers.IncHL()
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

func (cpu *CPU) decSp() {
	cpu.registers.DecSP()
}

func (cpu *CPU) ld_R16_Addr_R8_Dec(r1Name string, r2Name byte) func() {
	r1 := cpu.determine16Reg(r1Name)
	r1Dec := cpu.determine16RegDec(r1Name)
	r2 := cpu.determine8Reg(r2Name)

	return func() {
		cpu.bus.Write(r1(), *r2)
		r1Dec()
	}
}

func (cpu *CPU) ld_R16_Addr_R8_Inc(r1Name string, r2Name byte) func() {
	r1 := cpu.determine16Reg(r1Name)
	r1Inc := cpu.determine16RegInc(r1Name)
	r2 := cpu.determine8Reg(r2Name)

	return func() {
		cpu.bus.Write(r1(), *r2)
		r1Inc()
	}
}

func (cpu *CPU) popLsb() {
	cpu.readAddr(cpu.registers.sp)
	cpu.registers.SetTempAddrLsb(cpu.registers.Temp)
	cpu.registers.IncSP()
}

func (cpu *CPU) popMsb() {
	cpu.readAddr(cpu.registers.sp)
	cpu.registers.SetTempAddrMsb(cpu.registers.Temp)
	cpu.registers.IncSP()
}

func (cpu *CPU) hCarry(a, b uint8) {
	if ((a & 0x0F) + (b & 0x0F)) > 0x0F {
		cpu.registers.SetFlagH(true)
		return
	}

	cpu.registers.SetFlagH(false)
}

func (cpu *CPU) cCarry(a, b uint8) {
	if (uint16(a) + uint16(b)) > 0x0FF {
		cpu.registers.SetFlagC(true)
		return
	}

	cpu.registers.SetFlagC(false)
}

func (cpu *CPU) ld_L_SP_plus_N8() {
	SPL := cpu.registers.SPL()
	e := cpu.registers.Temp
	result := uint16(SPL) + uint16(e)
	fmt.Println(SPL, e, result, uint8(result))

	cpu.registers.SetFlagZ(false)
	cpu.registers.SetFlagN(false)
	cpu.hCarry(SPL, e)
	cpu.cCarry(SPL, e)

	cpu.registers.L = uint8(result)
}

func (cpu *CPU) signAdjust(val uint8) uint8 {
	if (val & 0b10000000) == 0 {
		return 0x00
	}

	return 0xFF
}

func (cpu *CPU) ld_H_SP_plus_N8() {
	adj := cpu.signAdjust(cpu.registers.Temp)
	carry := uint8(0)
	if cpu.registers.GetFlagC() {
		carry = 1
	}

	cpu.registers.H = cpu.registers.SPH() + adj + carry
}
