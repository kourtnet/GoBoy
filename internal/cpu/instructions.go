package cpu

func (cpu *CPU) nop() {
}

func (cpu *CPU) incA() {
	cpu.registers.A++
}

func (cpu *CPU) incPC() {
	cpu.registers.PC++
}
