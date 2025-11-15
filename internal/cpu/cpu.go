// Package cpu contains CPU struct that implements all the Sharp CPU logic you need
package cpu

import (
	"fmt"

	"github.com/kourtnet/GoBoy/internal/bus"
)

type registers struct {
	IR byte
	IE byte
	A  byte
	F  byte
	B  byte
	C  byte
	D  byte
	E  byte
	H  byte
	L  byte
	PC uint16
	SP uint16
}

const instructionsNum = 256

type CPU struct {
	instructions   [instructionsNum][]func()
	instructionLen int
	opNum          int

	registers *registers
	bus       *bus.Bus
}

func (cpu *CPU) fetch() error {
	cpu.registers.IR = cpu.bus.Read(cpu.registers.PC)

	cpu.opNum = 0
	cpu.instructionLen = len(cpu.instructions[cpu.registers.IR])

	if cpu.instructionLen == 0 {
		return fmt.Errorf("unknown opcode at: 0x%x", cpu.registers.PC)
	}

	cpu.incPC()

	return nil
}

func (cpu *CPU) execute() {
	cpu.instructions[cpu.registers.IR][cpu.opNum]()
	cpu.opNum++
}

func (cpu *CPU) Step() (bool, error) {
	cpu.execute()

	if cpu.opNum == cpu.instructionLen {
		err := cpu.fetch()
		if err != nil {
			return false, err
		}
	}

	return false, nil
}
