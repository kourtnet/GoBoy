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

	ReadNewInstruction bool

	Registers *Registers
	bus       iBus
}

func New(bus iBus) (CPU, error) {
	cpu := CPU{
		Registers: &Registers{
			pc: 0x100,
			S:  0xFF,
			P:  0xFE,
		},
		bus: bus,
	}

	// instruction set
	cpu.initInstructions()

	// Can happen after init functions if determine8/16Reg would
	// get invalid register name
	if cpu.internalErr != nil {
		return CPU{}, cpu.internalErr
	}

	// Required for initial cpu step of reading opcode, basically a NOP
	cpu.currInstructionLen = len(cpu.instructions[cpu.Registers.IR])

	return cpu, nil
}

func (cpu *CPU) fetchOpcode() error {
	var err error
	cpu.Registers.IR, err = cpu.bus.Read(cpu.Registers.PC())
	if err != nil {
		return err
	}

	cpu.opNum = 0

	cpu.currInstructionLen = len(cpu.instructions[cpu.Registers.IR])
	if cpu.currInstructionLen == 0 {
		return fmt.Errorf("unknown opcode at: 0x%x", cpu.Registers.PC())
	}

	cpu.Registers.IncPC()
	cpu.ReadNewInstruction = true

	return nil
}

func (cpu *CPU) execute() {
	cpu.ReadNewInstruction = false
	cpu.instructions[cpu.Registers.IR][cpu.opNum]()
	cpu.opNum++
}

func (cpu *CPU) Step() (bool, error) {
	cpu.step++
	// fmt.Printf("Step %d:\n", cpu.step)
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
	// fmt.Printf("Regs:\n%s\n", cpu.Registers)
	// fmt.Println("16-bit regs:")
	// fmt.Printf("AF: %#x BC: %#x DE: %#x HL: %#x\n", cpu.Registers.AF(), cpu.Registers.BC(), cpu.Registers.DE(), cpu.Registers.HL())

	return false, nil
}
