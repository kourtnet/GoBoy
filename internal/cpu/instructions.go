package cpu

import (
	"fmt"
)

func (cpu *CPU) readAddr(addr uint16) {
	var err error

	cpu.Registers.Temp8, err = cpu.bus.Read(addr)
	if err != nil {
		cpu.internalErr = err
		return
	}
}

func (cpu *CPU) readN8() {
	cpu.readAddr(cpu.Registers.pc)
	if cpu.internalErr != nil {
		return
	}

	cpu.Registers.IncPC()
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

func (cpu *CPU) readN16Lsb() {
	cpu.readN8()
	cpu.Registers.SetTemp16Lsb(cpu.Registers.Temp8)
}

func (cpu *CPU) readN16Msb() {
	cpu.readN8()
	cpu.Registers.SetTemp16Msb(cpu.Registers.Temp8)
}

func (cpu *CPU) determine8Reg(rName byte) *byte {
	var errReg byte

	switch rName {
	case 'A':
		return &cpu.Registers.A
	case 'F':
		return &cpu.Registers.F
	case 'B':
		return &cpu.Registers.B
	case 'C':
		return &cpu.Registers.C
	case 'D':
		return &cpu.Registers.D
	case 'E':
		return &cpu.Registers.E
	case 'H':
		return &cpu.Registers.H
	case 'L':
		return &cpu.Registers.L
	case 'S':
		return &cpu.Registers.S
	case 'P':
		return &cpu.Registers.P
	// Special processing for easier memory reading with immediate addr
	case 'T':
		return &cpu.Registers.Temp8
	default:
		cpu.internalErr = fmt.Errorf("unknown register name '%s'", string(rName))
		return &errReg
	}
}

func (cpu *CPU) determine16Reg(rName string) func() uint16 {
	switch rName {
	case "AF":
		return cpu.Registers.AF
	case "BC":
		return cpu.Registers.BC
	case "DE":
		return cpu.Registers.DE
	case "HL":
		return cpu.Registers.HL
	case "PC":
		return cpu.Registers.PC
	case "SP":
		return cpu.Registers.SP
	// Special processing for easier memory reading with immediate addr
	case "Temp16":
		return cpu.Registers.Temp16
	default:
		cpu.internalErr = fmt.Errorf("unknown register name '%s'", rName)
		return func() uint16 { return 0 }
	}
}

func (cpu *CPU) determine16RegSetter(rName string) func(uint16) {
	switch rName {
	case "BC":
		return cpu.Registers.SetBC
	case "DE":
		return cpu.Registers.SetDE
	case "HL":
		return cpu.Registers.SetHL
	case "SP":
		return cpu.Registers.SetSP
	case "AF":
		return cpu.Registers.SetAF
	default:
		cpu.internalErr = fmt.Errorf("unknown register name '%s'", rName)
		return func(uint16) {}
	}
}

func (cpu *CPU) determine16RegInc(rName string) func() {
	switch rName {
	case "BC":
		return cpu.Registers.IncBC
	case "DE":
		return cpu.Registers.IncDE
	case "HL":
		return cpu.Registers.IncHL
	case "SP":
		return cpu.Registers.IncSP
	default:
		cpu.internalErr = fmt.Errorf("unknown register name '%s'", rName)
		return func() {}
	}
}

func (cpu *CPU) determine16RegDec(rName string) func() {
	switch rName {
	case "BC":
		return cpu.Registers.DecBC
	case "DE":
		return cpu.Registers.DecDE
	case "HL":
		return cpu.Registers.DecHL
	case "SP":
		return cpu.Registers.DecSP
	default:
		cpu.internalErr = fmt.Errorf("unknown register name '%s'", rName)
		return func() {}
	}
}

func (cpu *CPU) nop() {
}

func (cpu *CPU) ldR8R8(r1Name, r2Name byte) func() {
	if r1Name == r2Name {
		return func() {}
	}

	r1 := cpu.determine8Reg(r1Name)
	r2 := cpu.determine8Reg(r2Name)

	return func() {
		*r1 = *r2
	}
}

func (cpu *CPU) ldAddrR8(r16Name string, r8Name byte) func() {
	r16 := cpu.determine16Reg(r16Name)
	r8 := cpu.determine8Reg(r8Name)

	return func() {
		cpu.bus.Write(r16(), *r8)
	}
}

// LDH [n8], A and LDH [C], A 2nd cycle
func (cpu *CPU) ldhAddrR8(r1Name, r2Name byte) func() {
	r1 := cpu.determine8Reg(r1Name)
	r2 := cpu.determine8Reg(r2Name)

	return func() {
		cpu.bus.Write(0xFF00+uint16(*r1), *r2)
	}
}

func (cpu *CPU) readHLAddrDec() {
	cpu.readR16Addr("HL")()
	cpu.Registers.DecHL()
}

func (cpu *CPU) readHLAddrInc() {
	cpu.readR16Addr("HL")()
	cpu.Registers.IncHL()
}

func (cpu *CPU) ldR16R16(r1Name, r2Name string) func() {
	r1 := cpu.determine16RegSetter(r1Name)
	r2 := cpu.determine16Reg(r2Name)

	return func() {
		r1(r2())
	}
}

func (cpu *CPU) ldN16AddrSPCycle3() {
	cpu.bus.Write(cpu.Registers.Temp16(), cpu.Registers.P)
}

func (cpu *CPU) ldN16AddrSPCycle4() {
	cpu.bus.Write(cpu.Registers.Temp16()+1, cpu.Registers.S)
}

func (cpu *CPU) decSp() {
	cpu.Registers.DecSP()
}

// This func is used both by LD [HL-], A and PUSH instructions. That's why
// it can be used for different regs, unlike ldHLPlusACycle2 that is used
// only by LD [HL+], A
func (cpu *CPU) ldR16AddrDecR8(r1Name string, r2Name byte) func() {
	r1 := cpu.determine16Reg(r1Name)
	r1Dec := cpu.determine16RegDec(r1Name)
	r2 := cpu.determine8Reg(r2Name)

	return func() {
		cpu.bus.Write(r1(), *r2)
		r1Dec()
	}
}

func (cpu *CPU) ldHLPlusACycle1() {
	r1 := cpu.Registers.HL()
	r2 := cpu.Registers.A

	cpu.bus.Write(r1, r2)
	cpu.Registers.IncHL()
}

func (cpu *CPU) popLsb() {
	cpu.readAddr(cpu.Registers.SP())
	cpu.Registers.SetTemp16Lsb(cpu.Registers.Temp8)
	cpu.Registers.IncSP()
}

func (cpu *CPU) popMsb() {
	cpu.readAddr(cpu.Registers.SP())
	cpu.Registers.SetTemp16Msb(cpu.Registers.Temp8)
	cpu.Registers.IncSP()
}

func (cpu *CPU) determineFlagH(a, b uint8, isSub, carry bool) {
	a &= 0x0F
	b &= 0x0F

	var res byte
	if isSub {
		res = a - b
		if carry {
			res--
		}
	} else {
		res = a + b
		if carry {
			res++
		}
	}

	if res > 0x0F {
		cpu.Registers.SetFlagH(true)
		return
	}

	cpu.Registers.SetFlagH(false)
}

func (cpu *CPU) determineFlagC(a, b uint8, isSub, carry bool) {
	a16 := uint16(a)
	b16 := uint16(b)

	var res uint16
	if isSub {
		res = a16 - b16
		if carry {
			res--
		}
	} else {
		res = a16 + b16
		if carry {
			res++
		}
	}

	if res > 0x0FF {
		cpu.Registers.SetFlagC(true)
		return
	}

	cpu.Registers.SetFlagC(false)
}

// LDHLSPPlusN8Cycle2 I'm so sorry for this abomination of a name
func (cpu *CPU) LDHLSPPlusN8Cycle2() {
	SPL := cpu.Registers.P
	e := cpu.Registers.Temp8
	result := uint16(SPL) + uint16(e)

	cpu.Registers.SetFlagZ(false)
	cpu.Registers.SetFlagN(false)
	cpu.determineFlagH(SPL, e, false, false)
	cpu.determineFlagC(SPL, e, false, false)

	cpu.Registers.L = uint8(result)
}

func (cpu *CPU) signAdjust(val uint8) uint8 {
	if (val & 0b10000000) == 0 {
		return 0x00
	}

	return 0xFF
}

// LDHLSPPlusN8Cycle3 And this one
func (cpu *CPU) LDHLSPPlusN8Cycle3() {
	adj := cpu.signAdjust(cpu.Registers.Temp8)
	carry := uint8(0)
	if cpu.Registers.FlagC() {
		carry = 1
	}

	cpu.Registers.H = cpu.Registers.S + adj + carry
}

func (cpu *CPU) determineFlagZ(val byte) {
	cpu.Registers.SetFlagZ(val == 0)
}

func (cpu *CPU) determineFlags(op1, op2 byte, isSub, carry bool) {
	var res byte
	if isSub {
		res = op1 - op2
		if carry {
			res--
		}
	} else {
		res = op1 + op2
		if carry {
			res++
		}
	}

	cpu.determineFlagZ(res)
	cpu.Registers.SetFlagN(isSub)
	cpu.determineFlagH(op1, op2, isSub, carry)
	cpu.determineFlagC(op1, op2, isSub, carry)
}

func (cpu *CPU) addVal(val byte, carry bool) {
	cpu.determineFlags(cpu.Registers.A, val, false, carry)

	cpu.Registers.A += val
	if carry {
		cpu.Registers.A++
	}
}

func (cpu *CPU) addR8(rName byte) func() {
	r := cpu.determine8Reg(rName)

	return func() {
		cpu.addVal(*r, false)
	}
}

func (cpu *CPU) adcR8(rName byte) func() {
	r := cpu.determine8Reg(rName)

	return func() {
		cpu.addVal(*r, cpu.Registers.FlagC())
	}
}

func (cpu *CPU) subVal(val byte, carry bool) {
	cpu.determineFlags(cpu.Registers.A, val, true, carry)

	cpu.Registers.A -= val
	if carry {
		cpu.Registers.A--
	}
}

func (cpu *CPU) subR8(rName byte) func() {
	r := cpu.determine8Reg(rName)

	return func() {
		cpu.subVal(*r, false)
	}
}

func (cpu *CPU) sbcR8(rName byte) func() {
	r := cpu.determine8Reg(rName)

	return func() {
		cpu.subVal(*r, cpu.Registers.FlagC())
	}
}

func (cpu *CPU) cpR8(rName byte) func() {
	r := cpu.determine8Reg(rName)

	return func() {
		A := cpu.Registers.A
		cpu.subVal(*r, false)
		cpu.Registers.A = A
	}
}

func (cpu *CPU) incR8(rName byte) func() {
	r := cpu.determine8Reg(rName)

	return func() {
		cpu.determineFlagH(*r, 1, false, false)
		cpu.Registers.SetFlagN(false)

		*r++
		cpu.determineFlagZ(*r)
	}
}

// really goofy name, but it's just a second
// operation of INC [HL] instruction
func (cpu *CPU) incHLAddrCycle2() {
	cpu.incR8('T')()
	cpu.ldAddrR8("HL", 'T')()
}

func (cpu *CPU) decR8(rName byte) func() {
	r := cpu.determine8Reg(rName)

	return func() {
		cpu.determineFlagH(*r, 1, true, false)
		cpu.Registers.SetFlagN(true)

		*r--
		cpu.determineFlagZ(*r)
	}
}

func (cpu *CPU) decHLAddrCycle2() {
	cpu.decR8('T')()
	cpu.ldAddrR8("HL", 'T')()
}

func (cpu *CPU) andR8(rName byte) func() {
	r := cpu.determine8Reg(rName)

	return func() {
		cpu.Registers.A &= *r

		cpu.determineFlagZ(cpu.Registers.A)
		cpu.Registers.SetFlagN(false)
		cpu.Registers.SetFlagH(true)
		cpu.Registers.SetFlagC(false)
	}
}

func (cpu *CPU) orR8(rName byte) func() {
	r := cpu.determine8Reg(rName)

	return func() {
		cpu.Registers.A |= *r

		cpu.determineFlagZ(cpu.Registers.A)
		cpu.Registers.SetFlagN(false)
		cpu.Registers.SetFlagH(false)
		cpu.Registers.SetFlagC(false)
	}
}

func (cpu *CPU) xorR8(rName byte) func() {
	r := cpu.determine8Reg(rName)

	return func() {
		cpu.Registers.A ^= *r

		cpu.determineFlagZ(cpu.Registers.A)
		cpu.Registers.SetFlagN(false)
		cpu.Registers.SetFlagH(false)
		cpu.Registers.SetFlagC(false)
	}
}

func (cpu *CPU) ccf() {
	cpu.Registers.SetFlagN(false)
	cpu.Registers.SetFlagH(false)
	cpu.Registers.SetFlagC(!cpu.Registers.FlagC())
}

func (cpu *CPU) scf() {
	cpu.Registers.SetFlagN(false)
	cpu.Registers.SetFlagH(false)
	cpu.Registers.SetFlagC(true)
}

func (cpu *CPU) daa() {
	var offset uint8

	if cpu.Registers.FlagN() {
		if cpu.Registers.FlagH() {
			offset += 0x06
		}
		if cpu.Registers.FlagC() {
			offset += 0x60
		}

		cpu.Registers.A -= offset
	} else {
		A := cpu.Registers.A
		if cpu.Registers.FlagH() || (A&0x0F) > 0x9 {
			offset += 0x06
		}
		if cpu.Registers.FlagC() || A > 0x99 {
			offset += 0x60
			cpu.Registers.SetFlagC(true)
		}

		cpu.Registers.A += offset
	}

	cpu.Registers.SetFlagH(false)
	cpu.determineFlagZ(cpu.Registers.A)
}

func (cpu *CPU) cpl() {
	cpu.Registers.A = ^cpu.Registers.A

	cpu.Registers.SetFlagN(true)
	cpu.Registers.SetFlagH(true)
}

func (cpu *CPU) addR16LSB(rName byte) func() {
	r := cpu.determine8Reg(rName)

	return func() {
		cpu.Registers.SetFlagN(false)
		cpu.determineFlagH(cpu.Registers.L, *r, false, false)
		cpu.determineFlagC(cpu.Registers.L, *r, false, false)

		cpu.Registers.L += *r
	}
}

func (cpu *CPU) addR16MSB(rName byte) func() {
	r := cpu.determine8Reg(rName)

	return func() {
		carry := byte(0)
		if cpu.Registers.FlagC() {
			carry = 1
		}

		cpu.Registers.SetFlagN(false)
		cpu.determineFlagH(cpu.Registers.H, *r, false, cpu.Registers.FlagC())
		cpu.determineFlagC(cpu.Registers.H, *r, false, cpu.Registers.FlagC())

		cpu.Registers.H += *r + carry
	}
}

func (cpu *CPU) addSPN8Cycle2() {
	cpu.Registers.SetFlagZ(false)
	cpu.Registers.SetFlagN(false)
	cpu.determineFlagH(cpu.Registers.P, cpu.Registers.Temp8, false, false)
	cpu.determineFlagC(cpu.Registers.P, cpu.Registers.Temp8, false, false)

	cpu.Registers.SetTemp16Lsb(cpu.Registers.P + cpu.Registers.Temp8)
}

func (cpu *CPU) addSPN8Cycle3() {
	res := cpu.Registers.S + cpu.signAdjust(cpu.Registers.Temp8)
	if cpu.Registers.FlagC() {
		res++
	}

	cpu.Registers.SetTemp16Msb(res)
}
