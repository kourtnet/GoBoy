// Package cpu contains CPU struct that implements all the Sharp CPU logic you need
package cpu

import (
	"fmt"
)

const instructionsNum = 256

type CPU struct {
	instructions       [instructionsNum][]func()
	currInstructionLen int
	opNum              int
	step               int

	// required for operations that work with external structs (like bus)
	// and can get an error. Because CPU operations are func(), we can't
	// process errors properly. So we have to store them into this var
	// and process later in the Step().
	internalErr error

	registers *registers
	bus       iBus
}

// WARNING: full test done
func (cpu *CPU) fetchOpcode() error {
	var err error
	cpu.registers.IR, err = cpu.bus.Read(cpu.registers.PC())
	if err != nil {
		return err
	}

	cpu.opNum = 0

	cpu.currInstructionLen = len(cpu.instructions[cpu.registers.IR])
	if cpu.currInstructionLen == 0 {
		return fmt.Errorf("unknown opcode at: 0x%x", cpu.registers.PC())
	}

	cpu.registers.IncPC()

	return nil
}

// TODO: integration test
func (cpu *CPU) execute() {
	cpu.instructions[cpu.registers.IR][cpu.opNum]()
	cpu.opNum++
}

// TODO: integration test
func (cpu *CPU) Step() (bool, error) {
	cpu.step++
	fmt.Printf("Step %d:\n", cpu.step)
	cpu.execute()
	if cpu.internalErr != nil {
		return true, cpu.internalErr
	}

	if cpu.opNum == cpu.currInstructionLen {
		err := cpu.fetchOpcode()
		if err != nil {
			return true, err
		}
	}

	// TODO: write a debugger instead of this
	fmt.Printf("Regs:\n%s\n", cpu.registers)
	fmt.Println("16-bit regs:")
	fmt.Printf("AF: %#x BC: %#x DE: %#x HL: %#x\n", cpu.registers.AF(), cpu.registers.BC(), cpu.registers.DE(), cpu.registers.HL())

	return false, nil
}
