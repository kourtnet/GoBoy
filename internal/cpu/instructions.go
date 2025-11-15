package cpu

func (cpu *CPU) nop() {
}

// TODO: ADD FLAGS AND OUT OF RANGE ARITHMETICS
func (cpu *CPU) incA() {
	cpu.registers.A++
}

// TODO: ADD FLAGS AND OUT OF RANGE ARITHMETICS
func (cpu *CPU) incPC() {
	cpu.registers.PC++
}

// LD
func (cpu *CPU) ldR8R8(r1Name, r2Name byte) func() {
	if r1Name == r2Name {
		return func() {}
	}

	var r1, r2 *byte

	switch r1Name {
	case 'A':
		r1 = &cpu.registers.A
	case 'B':
		r1 = &cpu.registers.B
	case 'C':
		r1 = &cpu.registers.C
	case 'D':
		r1 = &cpu.registers.D
	case 'E':
		r1 = &cpu.registers.E
	case 'H':
		r1 = &cpu.registers.H
	case 'L':
		r1 = &cpu.registers.L
	default:
		panic("unknown register name for r1")
	}

	switch r2Name {
	case 'A':
		r2 = &cpu.registers.A
	case 'B':
		r2 = &cpu.registers.B
	case 'C':
		r2 = &cpu.registers.C
	case 'D':
		r2 = &cpu.registers.D
	case 'E':
		r2 = &cpu.registers.E
	case 'H':
		r2 = &cpu.registers.H
	case 'L':
		r2 = &cpu.registers.L
	default:
		panic("unknown register name for r2")
	}

	return func() {
		*r1 = *r2
	}
}
