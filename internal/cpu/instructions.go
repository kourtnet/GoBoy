package cpu

import "errors"

func (cpu *CPU) readPCAddrAndInc() {
	var err error

	cpu.registers.temp, err = cpu.bus.Read(cpu.registers.PC())
	if err != nil {
		cpu.internalErr = err
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
		var err error
		cpu.registers.temp, err = cpu.bus.Read(r())
		if err != nil {
			cpu.internalErr = err
			return
		}
	}
}

func (cpu *CPU) readR8Addr(rName byte) func() {
	r := cpu.determine8Reg(rName)
	if cpu.internalErr != nil {
		return func() {}
	}

	return func() {
		var err error
		cpu.registers.temp, err = cpu.bus.Read(0xFF00 + uint16(*r))
		if err != nil {
			cpu.internalErr = err
			return
		}
	}
}

func (cpu *CPU) readAddrLsb() {
	cpu.readPCAddrAndInc()
	cpu.registers.SetTempAddrLsb(cpu.registers.temp)
}

func (cpu *CPU) readAddrMsb() {
	cpu.readPCAddrAndInc()
	cpu.registers.SetTempAddrMsb(cpu.registers.temp)
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
		return &cpu.registers.temp
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

func (cpu *CPU) nop() {
}

// TODO: ADD FLAGS AND OUT OF RANGE ARITHMETICS
func (cpu *CPU) incA() {
	cpu.registers.A++
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

// Requires fetchData as first procedure in instruction
func (cpu *CPU) ldR8Temp(rName byte) func() {
	r := cpu.determine8Reg(rName)

	return func() {
		*r = cpu.registers.temp
	}
}

func (cpu *CPU) ldR16AddrR8(r16Name string, r8Name byte) func() {
	r16 := cpu.determine16Reg(r16Name)
	r8 := cpu.determine8Reg(r8Name)

	return func() {
		cpu.bus.Write(r16(), *r8)
	}
}

func (cpu *CPU) ldR8AddrR8(r1Name byte, r2Name byte) func() {
	r1 := cpu.determine8Reg(r1Name)
	r2 := cpu.determine8Reg(r2Name)

	return func() {
		cpu.bus.Write(0xFF00+uint16(*r1), *r2)
	}
}

// there's no LD [R16], N8 instructions rather then with HL
func (cpu *CPU) ldHLAddrTemp() {
	cpu.bus.Write(cpu.registers.HL(), cpu.registers.temp)
}

func (cpu *CPU) read_HL_Addr_Dec() {
	cpu.readR16Addr("HL")
	cpu.registers.decHL()
}

func (cpu *CPU) read_HL_Addr_Inc() {
	cpu.readR16Addr("HL")
	cpu.registers.incHL()
}

func (cpu *CPU) ld_HL_Addr_A_Dec() {
	cpu.bus.Write(cpu.registers.HL(), cpu.registers.A)
	cpu.registers.incHL()
}

func (cpu *CPU) ld_HL_Addr_A_Inc() {
	cpu.bus.Write(cpu.registers.HL(), cpu.registers.A)
	cpu.registers.decHL()
}
