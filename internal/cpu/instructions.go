package cpu

// have to panic here, cause func is used in cpu.instructions array
// TODO: figure out proper error handling
func (cpu *CPU) readPCMemAndInc() {
	var err error
	cpu.registers.temp, err = cpu.bus.Read(cpu.registers.PC())
	if err != nil {
		panic(err)
	}

	cpu.registers.incPC()
}

func (cpu *CPU) readR16Addr(rName string) func() {
	r := cpu.determine16Reg(rName)

	return func() {
		var err error
		cpu.registers.temp, err = cpu.bus.Read(r())
		if err != nil {
			panic(err)
		}
	}
}

func (cpu *CPU) determine8Reg(rName byte) *byte {
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
	default:
		panic("unknown register name " + string(rName))
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
	default:
		panic("unknown register name " + rName)
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
