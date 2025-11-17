package cpu

func (cpu *CPU) determineReg(rName byte) *byte {
	var r *byte

	switch rName {
	case 'A':
		r = &cpu.registers.A
	case 'B':
		r = &cpu.registers.B
	case 'C':
		r = &cpu.registers.C
	case 'D':
		r = &cpu.registers.D
	case 'E':
		r = &cpu.registers.E
	case 'H':
		r = &cpu.registers.H
	case 'L':
		r = &cpu.registers.L
	default:
		panic("unknown register name for r2")
	}

	return r
}

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

func (cpu *CPU) ldR8R8(r1Name, r2Name byte) func() {
	if r1Name == r2Name {
		return func() {}
	}

	r1 := cpu.determineReg(r1Name)
	r2 := cpu.determineReg(r2Name)

	return func() {
		*r1 = *r2
	}
}

// Requires fetchData as first procedure in instruction
func (cpu *CPU) ldR8N8(rName byte) func() {
	r := cpu.determineReg(rName)

	return func() {
		*r = cpu.registers.temp
	}
}
