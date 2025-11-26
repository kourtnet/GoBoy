package cpu

import "github.com/kourtnet/GoBoy/internal/bus"

func New(bus *bus.Bus) CPU {
	cpu := CPU{
		registers: &registers{
			pc: 0x100,
		},

		bus: bus,
	}

	cpu.instructions[0x0] = []func(){cpu.nop}
	cpu.instructions[0x3C] = []func(){cpu.incA}

	// instruction set
	// nop
	cpu.initNOP()

	// ld
	cpu.initLDR8R8()
	cpu.initLDR8N8()
	cpu.initLDR8R16Addr()
	cpu.initLDR16AddrR8()
	cpu.initLDHLAddrN()
	cpu.initLDR8N16Addr()
	cpu.initLDN16AddrR8()
	cpu.init_LD_A_C_Addr()
	cpu.init_LD_C_Addr_A()

	// Required for initial cpu step, basically a NOP
	cpu.instructionLen = len(cpu.instructions[cpu.registers.IR])

	return cpu
}

func (cpu *CPU) initNOP() {
	cpu.instructions[0x00] = []func(){cpu.nop}
}

func (cpu *CPU) initLDR8R8() {
	cpu.instructions[0x40] = []func(){cpu.ldR8R8('B', 'B')}
	cpu.instructions[0x41] = []func(){cpu.ldR8R8('B', 'C')}
	cpu.instructions[0x42] = []func(){cpu.ldR8R8('B', 'D')}
	cpu.instructions[0x43] = []func(){cpu.ldR8R8('B', 'E')}
	cpu.instructions[0x44] = []func(){cpu.ldR8R8('B', 'H')}
	cpu.instructions[0x45] = []func(){cpu.ldR8R8('B', 'L')}
	cpu.instructions[0x47] = []func(){cpu.ldR8R8('B', 'A')}
	cpu.instructions[0x48] = []func(){cpu.ldR8R8('C', 'B')}
	cpu.instructions[0x49] = []func(){cpu.ldR8R8('C', 'C')}
	cpu.instructions[0x4A] = []func(){cpu.ldR8R8('C', 'D')}
	cpu.instructions[0x4B] = []func(){cpu.ldR8R8('C', 'E')}
	cpu.instructions[0x4C] = []func(){cpu.ldR8R8('C', 'H')}
	cpu.instructions[0x4D] = []func(){cpu.ldR8R8('C', 'L')}
	cpu.instructions[0x4F] = []func(){cpu.ldR8R8('C', 'A')}
	cpu.instructions[0x50] = []func(){cpu.ldR8R8('D', 'B')}
	cpu.instructions[0x51] = []func(){cpu.ldR8R8('D', 'C')}
	cpu.instructions[0x52] = []func(){cpu.ldR8R8('D', 'D')}
	cpu.instructions[0x53] = []func(){cpu.ldR8R8('D', 'E')}
	cpu.instructions[0x54] = []func(){cpu.ldR8R8('D', 'H')}
	cpu.instructions[0x55] = []func(){cpu.ldR8R8('D', 'L')}
	cpu.instructions[0x57] = []func(){cpu.ldR8R8('D', 'A')}
	cpu.instructions[0x58] = []func(){cpu.ldR8R8('E', 'B')}
	cpu.instructions[0x59] = []func(){cpu.ldR8R8('E', 'C')}
	cpu.instructions[0x5A] = []func(){cpu.ldR8R8('E', 'D')}
	cpu.instructions[0x5B] = []func(){cpu.ldR8R8('E', 'E')}
	cpu.instructions[0x5C] = []func(){cpu.ldR8R8('E', 'H')}
	cpu.instructions[0x5D] = []func(){cpu.ldR8R8('E', 'L')}
	cpu.instructions[0x5F] = []func(){cpu.ldR8R8('E', 'A')}
	cpu.instructions[0x60] = []func(){cpu.ldR8R8('H', 'B')}
	cpu.instructions[0x61] = []func(){cpu.ldR8R8('H', 'C')}
	cpu.instructions[0x62] = []func(){cpu.ldR8R8('H', 'D')}
	cpu.instructions[0x63] = []func(){cpu.ldR8R8('H', 'E')}
	cpu.instructions[0x64] = []func(){cpu.ldR8R8('H', 'H')}
	cpu.instructions[0x65] = []func(){cpu.ldR8R8('H', 'L')}
	cpu.instructions[0x67] = []func(){cpu.ldR8R8('H', 'A')}
	cpu.instructions[0x68] = []func(){cpu.ldR8R8('L', 'B')}
	cpu.instructions[0x69] = []func(){cpu.ldR8R8('L', 'C')}
	cpu.instructions[0x6A] = []func(){cpu.ldR8R8('L', 'D')}
	cpu.instructions[0x6B] = []func(){cpu.ldR8R8('L', 'E')}
	cpu.instructions[0x6C] = []func(){cpu.ldR8R8('L', 'H')}
	cpu.instructions[0x6D] = []func(){cpu.ldR8R8('L', 'L')}
	cpu.instructions[0x6F] = []func(){cpu.ldR8R8('L', 'A')}
	cpu.instructions[0x78] = []func(){cpu.ldR8R8('A', 'B')}
	cpu.instructions[0x79] = []func(){cpu.ldR8R8('A', 'C')}
	cpu.instructions[0x7A] = []func(){cpu.ldR8R8('A', 'D')}
	cpu.instructions[0x7B] = []func(){cpu.ldR8R8('A', 'E')}
	cpu.instructions[0x7C] = []func(){cpu.ldR8R8('A', 'H')}
	cpu.instructions[0x7D] = []func(){cpu.ldR8R8('A', 'L')}
	cpu.instructions[0x7F] = []func(){cpu.ldR8R8('A', 'A')}
}

func (cpu *CPU) initLDR8N8() {
	cpu.instructions[0x06] = []func(){cpu.readPCMemAndInc, cpu.ldR8Temp('B')}
	cpu.instructions[0x16] = []func(){cpu.readPCMemAndInc, cpu.ldR8Temp('D')}
	cpu.instructions[0x26] = []func(){cpu.readPCMemAndInc, cpu.ldR8Temp('H')}
	cpu.instructions[0x0E] = []func(){cpu.readPCMemAndInc, cpu.ldR8Temp('C')}
	cpu.instructions[0x1E] = []func(){cpu.readPCMemAndInc, cpu.ldR8Temp('E')}
	cpu.instructions[0x2E] = []func(){cpu.readPCMemAndInc, cpu.ldR8Temp('L')}
	cpu.instructions[0x3E] = []func(){cpu.readPCMemAndInc, cpu.ldR8Temp('A')}
}

func (cpu *CPU) initLDR8R16Addr() {
	cpu.instructions[0x46] = []func(){cpu.readR16Addr("HL"), cpu.ldR8Temp('B')}
	cpu.instructions[0x56] = []func(){cpu.readR16Addr("HL"), cpu.ldR8Temp('D')}
	cpu.instructions[0x66] = []func(){cpu.readR16Addr("HL"), cpu.ldR8Temp('H')}
	cpu.instructions[0x4E] = []func(){cpu.readR16Addr("HL"), cpu.ldR8Temp('C')}
	cpu.instructions[0x5E] = []func(){cpu.readR16Addr("HL"), cpu.ldR8Temp('E')}
	cpu.instructions[0x6E] = []func(){cpu.readR16Addr("HL"), cpu.ldR8Temp('L')}
	cpu.instructions[0x7E] = []func(){cpu.readR16Addr("HL"), cpu.ldR8Temp('A')}

	cpu.instructions[0x0A] = []func(){cpu.readR16Addr("BC"), cpu.ldR8Temp('A')}
	cpu.instructions[0x1A] = []func(){cpu.readR16Addr("DE"), cpu.ldR8Temp('A')}
}

// Here memory write happens on the first M-cycle. I find it strange, but doc
// states that it's true. Gonna examine it once I will run test-ROMs
func (cpu *CPU) initLDR16AddrR8() {
	cpu.instructions[0x02] = []func(){cpu.ldR16AddrR8("BC", 'A'), cpu.nop}
	cpu.instructions[0x12] = []func(){cpu.ldR16AddrR8("DE", 'A'), cpu.nop}

	cpu.instructions[0x70] = []func(){cpu.ldR16AddrR8("HL", 'B'), cpu.nop}
	cpu.instructions[0x71] = []func(){cpu.ldR16AddrR8("HL", 'C'), cpu.nop}
	cpu.instructions[0x72] = []func(){cpu.ldR16AddrR8("HL", 'D'), cpu.nop}
	cpu.instructions[0x73] = []func(){cpu.ldR16AddrR8("HL", 'E'), cpu.nop}
	cpu.instructions[0x74] = []func(){cpu.ldR16AddrR8("HL", 'H'), cpu.nop}
	cpu.instructions[0x75] = []func(){cpu.ldR16AddrR8("HL", 'L'), cpu.nop}
	cpu.instructions[0x77] = []func(){cpu.ldR16AddrR8("HL", 'A'), cpu.nop}
}

// Here memory write happens on the first M-cycle. I find it strange, but doc
// states that it's true. Gonna examine it once I will run test-ROMs
func (cpu *CPU) initLDHLAddrN() {
	cpu.instructions[0x36] = []func(){cpu.readPCMemAndInc, cpu.ldHLAddrTemp, cpu.nop}
}

func (cpu *CPU) initLDR8N16Addr() {
	cpu.instructions[0xFA] = []func(){cpu.readAddrLsb, cpu.readAddrMsb, cpu.readR16Addr("TempAddr"), cpu.ldR8Temp('A')}
}

func (cpu *CPU) initLDN16AddrR8() {
	cpu.instructions[0xEA] = []func(){cpu.readAddrLsb, cpu.readAddrMsb, cpu.readR16Addr("TempAddr"), cpu.ldR16AddrR8("TempAddr", 'A')}
}

func (cpu *CPU) init_LD_A_C_Addr() {
	cpu.instructions[0xF2] = []func(){cpu.readR8Addr('C'), cpu.ldR8Temp('A')}
}

func (cpu *CPU) init_LD_C_Addr_A() {
	cpu.instructions[0xE2] = []func(){cpu.ldR8AddrR8('C', 'A'), cpu.nop}
}
