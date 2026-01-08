package cpu

// TODO: integration test
func New(bus iBus) (CPU, error) {
	cpu := CPU{
		registers: &registers{
			pc: 0x100,
		},
		bus: bus,
	}

	// instruction set
	cpu.initInstructions()

	// Required for initial cpu step, basically a NOP
	cpu.currInstructionLen = len(cpu.instructions[cpu.registers.IR])

	// Can happen after init functions if determine8/16Reg would
	// get invalid register name
	if cpu.internalErr != nil {
		return CPU{}, cpu.internalErr
	}

	// Required for initial cpu step, basically a NOP
	cpu.currInstructionLen = len(cpu.instructions[cpu.registers.IR])

	return cpu, nil
}

func (cpu *CPU) initInstructions() {
	// nop
	cpu.init_NOP()

	// ld
	cpu.init_LD8()
	cpu.init_LD16()
}

func (cpu *CPU) init_NOP() {
	cpu.instructions[0x00] = []func(){cpu.nop}
}

func (cpu *CPU) init_LD8() {
	cpu.init_LD_R8_R8()
	cpu.init_LD_R8_N8()
	cpu.init_LD_R8_R16_Addr()
	cpu.init_LD_R16_Addr_R8()
	cpu.init_LD_N16_Addr()
	cpu.init_LDH()
	cpu.init_LD_HL_Inc_Dec()
}

func (cpu *CPU) init_LD_R8_R8() {
	cpu.instructions[0x40] = []func(){cpu.ld_R8_R8('B', 'B')}
	cpu.instructions[0x41] = []func(){cpu.ld_R8_R8('B', 'C')}
	cpu.instructions[0x42] = []func(){cpu.ld_R8_R8('B', 'D')}
	cpu.instructions[0x43] = []func(){cpu.ld_R8_R8('B', 'E')}
	cpu.instructions[0x44] = []func(){cpu.ld_R8_R8('B', 'H')}
	cpu.instructions[0x45] = []func(){cpu.ld_R8_R8('B', 'L')}
	cpu.instructions[0x47] = []func(){cpu.ld_R8_R8('B', 'A')}
	cpu.instructions[0x48] = []func(){cpu.ld_R8_R8('C', 'B')}
	cpu.instructions[0x49] = []func(){cpu.ld_R8_R8('C', 'C')}
	cpu.instructions[0x4A] = []func(){cpu.ld_R8_R8('C', 'D')}
	cpu.instructions[0x4B] = []func(){cpu.ld_R8_R8('C', 'E')}
	cpu.instructions[0x4C] = []func(){cpu.ld_R8_R8('C', 'H')}
	cpu.instructions[0x4D] = []func(){cpu.ld_R8_R8('C', 'L')}
	cpu.instructions[0x4F] = []func(){cpu.ld_R8_R8('C', 'A')}
	cpu.instructions[0x50] = []func(){cpu.ld_R8_R8('D', 'B')}
	cpu.instructions[0x51] = []func(){cpu.ld_R8_R8('D', 'C')}
	cpu.instructions[0x52] = []func(){cpu.ld_R8_R8('D', 'D')}
	cpu.instructions[0x53] = []func(){cpu.ld_R8_R8('D', 'E')}
	cpu.instructions[0x54] = []func(){cpu.ld_R8_R8('D', 'H')}
	cpu.instructions[0x55] = []func(){cpu.ld_R8_R8('D', 'L')}
	cpu.instructions[0x57] = []func(){cpu.ld_R8_R8('D', 'A')}
	cpu.instructions[0x58] = []func(){cpu.ld_R8_R8('E', 'B')}
	cpu.instructions[0x59] = []func(){cpu.ld_R8_R8('E', 'C')}
	cpu.instructions[0x5A] = []func(){cpu.ld_R8_R8('E', 'D')}
	cpu.instructions[0x5B] = []func(){cpu.ld_R8_R8('E', 'E')}
	cpu.instructions[0x5C] = []func(){cpu.ld_R8_R8('E', 'H')}
	cpu.instructions[0x5D] = []func(){cpu.ld_R8_R8('E', 'L')}
	cpu.instructions[0x5F] = []func(){cpu.ld_R8_R8('E', 'A')}
	cpu.instructions[0x60] = []func(){cpu.ld_R8_R8('H', 'B')}
	cpu.instructions[0x61] = []func(){cpu.ld_R8_R8('H', 'C')}
	cpu.instructions[0x62] = []func(){cpu.ld_R8_R8('H', 'D')}
	cpu.instructions[0x63] = []func(){cpu.ld_R8_R8('H', 'E')}
	cpu.instructions[0x64] = []func(){cpu.ld_R8_R8('H', 'H')}
	cpu.instructions[0x65] = []func(){cpu.ld_R8_R8('H', 'L')}
	cpu.instructions[0x67] = []func(){cpu.ld_R8_R8('H', 'A')}
	cpu.instructions[0x68] = []func(){cpu.ld_R8_R8('L', 'B')}
	cpu.instructions[0x69] = []func(){cpu.ld_R8_R8('L', 'C')}
	cpu.instructions[0x6A] = []func(){cpu.ld_R8_R8('L', 'D')}
	cpu.instructions[0x6B] = []func(){cpu.ld_R8_R8('L', 'E')}
	cpu.instructions[0x6C] = []func(){cpu.ld_R8_R8('L', 'H')}
	cpu.instructions[0x6D] = []func(){cpu.ld_R8_R8('L', 'L')}
	cpu.instructions[0x6F] = []func(){cpu.ld_R8_R8('L', 'A')}
	cpu.instructions[0x78] = []func(){cpu.ld_R8_R8('A', 'B')}
	cpu.instructions[0x79] = []func(){cpu.ld_R8_R8('A', 'C')}
	cpu.instructions[0x7A] = []func(){cpu.ld_R8_R8('A', 'D')}
	cpu.instructions[0x7B] = []func(){cpu.ld_R8_R8('A', 'E')}
	cpu.instructions[0x7C] = []func(){cpu.ld_R8_R8('A', 'H')}
	cpu.instructions[0x7D] = []func(){cpu.ld_R8_R8('A', 'L')}
	cpu.instructions[0x7F] = []func(){cpu.ld_R8_R8('A', 'A')}
}

func (cpu *CPU) init_LD_R8_N8() {
	cpu.instructions[0x06] = []func(){cpu.readPCAddrAndInc, cpu.ld_R8_Temp('B')}
	cpu.instructions[0x16] = []func(){cpu.readPCAddrAndInc, cpu.ld_R8_Temp('D')}
	cpu.instructions[0x26] = []func(){cpu.readPCAddrAndInc, cpu.ld_R8_Temp('H')}
	cpu.instructions[0x0E] = []func(){cpu.readPCAddrAndInc, cpu.ld_R8_Temp('C')}
	cpu.instructions[0x1E] = []func(){cpu.readPCAddrAndInc, cpu.ld_R8_Temp('E')}
	cpu.instructions[0x2E] = []func(){cpu.readPCAddrAndInc, cpu.ld_R8_Temp('L')}
	cpu.instructions[0x3E] = []func(){cpu.readPCAddrAndInc, cpu.ld_R8_Temp('A')}
}

func (cpu *CPU) init_LD_R8_R16_Addr() {
	cpu.instructions[0x46] = []func(){cpu.readR16Addr("HL"), cpu.ld_R8_Temp('B')}
	cpu.instructions[0x56] = []func(){cpu.readR16Addr("HL"), cpu.ld_R8_Temp('D')}
	cpu.instructions[0x66] = []func(){cpu.readR16Addr("HL"), cpu.ld_R8_Temp('H')}
	cpu.instructions[0x4E] = []func(){cpu.readR16Addr("HL"), cpu.ld_R8_Temp('C')}
	cpu.instructions[0x5E] = []func(){cpu.readR16Addr("HL"), cpu.ld_R8_Temp('E')}
	cpu.instructions[0x6E] = []func(){cpu.readR16Addr("HL"), cpu.ld_R8_Temp('L')}
	cpu.instructions[0x7E] = []func(){cpu.readR16Addr("HL"), cpu.ld_R8_Temp('A')}

	cpu.instructions[0x0A] = []func(){cpu.readR16Addr("BC"), cpu.ld_R8_Temp('A')}
	cpu.instructions[0x1A] = []func(){cpu.readR16Addr("DE"), cpu.ld_R8_Temp('A')}
}

// Here memory write happens on the first M-cycle. I find it strange, but doc
// states that it's true. Gonna examine it once I will run test-ROMs
func (cpu *CPU) init_LD_R16_Addr_R8() {
	cpu.instructions[0x02] = []func(){cpu.ld_R16_Addr_R8("BC", 'A'), cpu.nop}
	cpu.instructions[0x12] = []func(){cpu.ld_R16_Addr_R8("DE", 'A'), cpu.nop}

	cpu.instructions[0x70] = []func(){cpu.ld_R16_Addr_R8("HL", 'B'), cpu.nop}
	cpu.instructions[0x71] = []func(){cpu.ld_R16_Addr_R8("HL", 'C'), cpu.nop}
	cpu.instructions[0x72] = []func(){cpu.ld_R16_Addr_R8("HL", 'D'), cpu.nop}
	cpu.instructions[0x73] = []func(){cpu.ld_R16_Addr_R8("HL", 'E'), cpu.nop}
	cpu.instructions[0x74] = []func(){cpu.ld_R16_Addr_R8("HL", 'H'), cpu.nop}
	cpu.instructions[0x75] = []func(){cpu.ld_R16_Addr_R8("HL", 'L'), cpu.nop}
	cpu.instructions[0x77] = []func(){cpu.ld_R16_Addr_R8("HL", 'A'), cpu.nop}
}

func (cpu *CPU) init_LD_N16_Addr() {
	// Here memory write happens on the first M-cycle. I find it strange, but doc
	// states that it's true. Gonna examine it once I will run test-ROMs
	cpu.instructions[0x36] = []func(){cpu.readPCAddrAndInc, cpu.ldHLAddrTemp, cpu.nop}
	cpu.instructions[0xEA] = []func(){cpu.readAddrLsb, cpu.readAddrMsb, cpu.ld_R16_Addr_R8("TempAddr", 'A'), cpu.nop}
	cpu.instructions[0xFA] = []func(){cpu.readAddrLsb, cpu.readAddrMsb, cpu.readR16Addr("TempAddr"), cpu.ld_R8_Temp('A')}
}

func (cpu *CPU) init_LDH() {
	// WARNING: not tested by ROMs
	cpu.instructions[0xE0] = []func(){cpu.readPCAddrAndInc, cpu.ld_R8_Addr_R8('T', 'A'), cpu.nop}
	cpu.instructions[0xE2] = []func(){cpu.ld_R8_Addr_R8('C', 'A'), cpu.nop}
	cpu.instructions[0xF0] = []func(){cpu.readPCAddrAndInc, cpu.readR8Addr('T'), cpu.ld_R8_Temp('A')}
	cpu.instructions[0xF2] = []func(){cpu.readR8Addr('C'), cpu.ld_R8_Temp('A')}
}

// TODO: get rid of temp func
func (cpu *CPU) init_LD_HL_Inc_Dec() {
	cpu.instructions[0x22] = []func(){cpu.ld_HL_Addr_A_Inc, cpu.nop}
	cpu.instructions[0x2A] = []func(){cpu.readHLAddrInc, cpu.ld_R8_Temp('A')}
	cpu.instructions[0x32] = []func(){cpu.ld_HL_Addr_A_Dec, cpu.nop}
	cpu.instructions[0x3A] = []func(){cpu.readHLAddrDec, cpu.ld_R8_Temp('A')}
}

func (cpu *CPU) init_LD16() {
	// LD SP, HL
	cpu.instructions[0xF9] = []func(){cpu.ld_R16_R16("SP", "HL"), cpu.nop}
	// LD [n16], SP
	cpu.instructions[0x08] = []func(){cpu.readAddrLsb, cpu.readAddrMsb, cpu.ld_TempAddr_SPL, cpu.ld_TempAddr_SPH, cpu.nop}
	cpu.init_LD_R16_N16()
}

func (cpu *CPU) init_LD_R16_N16() {
	cpu.instructions[0x01] = []func(){cpu.readAddrLsb, cpu.readAddrMsb, cpu.ld_R16_R16("BC", "TempAddr")}
	cpu.instructions[0x11] = []func(){cpu.readAddrLsb, cpu.readAddrMsb, cpu.ld_R16_R16("DE", "TempAddr")}
	cpu.instructions[0x21] = []func(){cpu.readAddrLsb, cpu.readAddrMsb, cpu.ld_R16_R16("HL", "TempAddr")}
	cpu.instructions[0x31] = []func(){cpu.readAddrLsb, cpu.readAddrMsb, cpu.ld_R16_R16("SP", "TempAddr")}
}
