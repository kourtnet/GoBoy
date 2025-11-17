// Package cpu contains CPU struct that implements all the Sharp CPU logic you need
package cpu

import (
	"fmt"

	"github.com/kourtnet/GoBoy/internal/bus"
)

const instructionsNum = 256

type CPU struct {
	instructions   [instructionsNum][]func()
	instructionLen int
	opNum          int
	step           int

	registers *registers
	bus       *bus.Bus
}

// have to panic here, cause func is used in cpu.instructions array
// TODO: figure out proper error handling
func (cpu *CPU) fetchData() {
	var err error
	cpu.registers.temp, err = cpu.bus.Read(cpu.registers.PC)
	if err != nil {
		panic(err)
	}

	cpu.incPC()
}

func (cpu *CPU) fetchOpcode() error {
	var err error
	cpu.registers.IR, err = cpu.bus.Read(cpu.registers.PC)
	if err != nil {
		return err
	}

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
	cpu.step++
	fmt.Printf("Step %d:\n", cpu.step)
	cpu.execute()

	if cpu.opNum == cpu.instructionLen {
		err := cpu.fetchOpcode()
		if err != nil {
			return false, err
		}
	}

	fmt.Printf("Regs:\n%s\n", cpu.registers)
	fmt.Println("16-bit regs:")
	fmt.Printf("AF: %#x BC: %#x DE: %#x HL: %#x\n", cpu.registers.AF(), cpu.registers.BC(), cpu.registers.DE(), cpu.registers.HL())

	return false, nil
}
