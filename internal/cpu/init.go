package cpu

import "github.com/kourtnet/GoBoy/internal/bus"

func New(bus *bus.Bus) CPU {
	cpu := CPU{
		registers: registers{
			PC: 0x100,
		},

		bus: bus,
	}

	cpu.instructions[0x0] = []func(){cpu.nop}
	cpu.instructions[0x3C] = []func(){cpu.incA}

	// Required for initial cpu step, basically a NOP
	cpu.instructionLen = len(cpu.instructions[cpu.registers.IR])

	return cpu
}
