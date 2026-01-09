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
	cpu.initNOP()

	// ld
	cpu.initLD8()
	cpu.initLD16()

	// arithmetic
	cpu.init_ADD_ADC()
	cpu.init_SUB_SBC()
	cpu.initCP()
}

func (cpu *CPU) initNOP() {
	cpu.instructions[0x00] = []func(){cpu.nop}
}

func (cpu *CPU) initLD8() {
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
	cpu.instructions[0x22] = []func(){cpu.ld_R16_Addr_R8_Inc("HL", 'A'), cpu.nop}
	cpu.instructions[0x2A] = []func(){cpu.readHLAddrInc, cpu.ld_R8_Temp('A')}
	cpu.instructions[0x32] = []func(){cpu.ld_R16_Addr_R8_Dec("HL", 'A'), cpu.nop}
	cpu.instructions[0x3A] = []func(){cpu.readHLAddrDec, cpu.ld_R8_Temp('A')}
}

func (cpu *CPU) initLD16() {
	// LD SP, HL
	cpu.instructions[0xF9] = []func(){cpu.ld_R16_R16("SP", "HL"), cpu.nop}
	// LD [n16], SP
	cpu.instructions[0x08] = []func(){cpu.readAddrLsb, cpu.readAddrMsb, cpu.ld_TempAddr_SPL, cpu.ld_TempAddr_SPH, cpu.nop}
	// LD HL, SP + n8
	cpu.instructions[0xF8] = []func(){cpu.readPCAddrAndInc, cpu.ld_L_SP_plus_N8, cpu.ld_H_SP_plus_N8}

	cpu.init_LD_R16_N16()
	cpu.InitPush()
	cpu.InitPop()
}

// TODO: rename readAddrLsb, readAddrMsb and readPCAddrAndInc. Smth like readN8, readN16Msb, readN16Lsb
func (cpu *CPU) init_LD_R16_N16() {
	cpu.instructions[0x01] = []func(){cpu.readAddrLsb, cpu.readAddrMsb, cpu.ld_R16_R16("BC", "TempAddr")}
	cpu.instructions[0x11] = []func(){cpu.readAddrLsb, cpu.readAddrMsb, cpu.ld_R16_R16("DE", "TempAddr")}
	cpu.instructions[0x21] = []func(){cpu.readAddrLsb, cpu.readAddrMsb, cpu.ld_R16_R16("HL", "TempAddr")}
	cpu.instructions[0x31] = []func(){cpu.readAddrLsb, cpu.readAddrMsb, cpu.ld_R16_R16("SP", "TempAddr")}
}

func (cpu *CPU) InitPush() {
	cpu.instructions[0xC5] = []func(){cpu.decSp, cpu.ld_R16_Addr_R8_Dec("SP", 'B'), cpu.ld_R16_Addr_R8("SP", 'C'), cpu.nop}
	cpu.instructions[0xD5] = []func(){cpu.decSp, cpu.ld_R16_Addr_R8_Dec("SP", 'D'), cpu.ld_R16_Addr_R8("SP", 'E'), cpu.nop}
	cpu.instructions[0xE5] = []func(){cpu.decSp, cpu.ld_R16_Addr_R8_Dec("SP", 'H'), cpu.ld_R16_Addr_R8("SP", 'L'), cpu.nop}
	cpu.instructions[0xF5] = []func(){cpu.decSp, cpu.ld_R16_Addr_R8_Dec("SP", 'A'), cpu.ld_R16_Addr_R8("SP", 'F'), cpu.nop}
}

func (cpu *CPU) InitPop() {
	cpu.instructions[0xC1] = []func(){cpu.popLsb, cpu.popMsb, cpu.ld_R16_R16("BC", "TempAddr")}
	cpu.instructions[0xD1] = []func(){cpu.popLsb, cpu.popMsb, cpu.ld_R16_R16("DE", "TempAddr")}
	cpu.instructions[0xE1] = []func(){cpu.popLsb, cpu.popMsb, cpu.ld_R16_R16("HL", "TempAddr")}
	cpu.instructions[0xF1] = []func(){cpu.popLsb, cpu.popMsb, cpu.ld_R16_R16("AF", "TempAddr")}
}

func (cpu *CPU) init_ADD_ADC() {
	// ADD [HL]
	cpu.instructions[0x86] = []func(){cpu.readR16Addr("HL"), cpu.addR8('T')}
	// ADD n8
	cpu.instructions[0xC6] = []func(){cpu.readPCAddrAndInc, cpu.addR8('T')}
	// ADC [HL]
	cpu.instructions[0x8E] = []func(){cpu.readR16Addr("HL"), cpu.adcR8('T')}
	// ADC n8
	cpu.instructions[0xCE] = []func(){cpu.readPCAddrAndInc, cpu.adcR8('T')}

	cpu.init_ADD_R8()
	cpu.init_ADC_R8()
}

func (cpu *CPU) init_ADD_R8() {
	cpu.instructions[0x80] = []func(){cpu.addR8('B')}
	cpu.instructions[0x81] = []func(){cpu.addR8('C')}
	cpu.instructions[0x82] = []func(){cpu.addR8('D')}
	cpu.instructions[0x83] = []func(){cpu.addR8('E')}
	cpu.instructions[0x84] = []func(){cpu.addR8('H')}
	cpu.instructions[0x85] = []func(){cpu.addR8('L')}
	cpu.instructions[0x87] = []func(){cpu.addR8('A')}
}

func (cpu *CPU) init_ADC_R8() {
	cpu.instructions[0x88] = []func(){cpu.adcR8('B')}
	cpu.instructions[0x89] = []func(){cpu.adcR8('C')}
	cpu.instructions[0x8A] = []func(){cpu.adcR8('D')}
	cpu.instructions[0x8B] = []func(){cpu.adcR8('E')}
	cpu.instructions[0x8C] = []func(){cpu.adcR8('H')}
	cpu.instructions[0x8D] = []func(){cpu.adcR8('L')}
	cpu.instructions[0x8F] = []func(){cpu.adcR8('A')}
}

func (cpu *CPU) init_SUB_SBC() {
	// SUB [HL]
	cpu.instructions[0x96] = []func(){cpu.readR16Addr("HL"), cpu.subR8('T')}
	// SUB n8
	cpu.instructions[0xD6] = []func(){cpu.readPCAddrAndInc, cpu.subR8('T')}
	// SBC [HL]
	cpu.instructions[0x9E] = []func(){cpu.readR16Addr("HL"), cpu.sbcR8('T')}
	// SBC n8
	cpu.instructions[0xDE] = []func(){cpu.readPCAddrAndInc, cpu.sbcR8('T')}

	cpu.init_SUB_R8()
	cpu.init_SBC_R8()
}

func (cpu *CPU) init_SUB_R8() {
	cpu.instructions[0x90] = []func(){cpu.subR8('B')}
	cpu.instructions[0x91] = []func(){cpu.subR8('C')}
	cpu.instructions[0x92] = []func(){cpu.subR8('D')}
	cpu.instructions[0x93] = []func(){cpu.subR8('E')}
	cpu.instructions[0x94] = []func(){cpu.subR8('H')}
	cpu.instructions[0x95] = []func(){cpu.subR8('L')}
	cpu.instructions[0x97] = []func(){cpu.subR8('A')}
}

func (cpu *CPU) init_SBC_R8() {
	cpu.instructions[0x98] = []func(){cpu.sbcR8('B')}
	cpu.instructions[0x99] = []func(){cpu.sbcR8('C')}
	cpu.instructions[0x9a] = []func(){cpu.sbcR8('D')}
	cpu.instructions[0x9b] = []func(){cpu.sbcR8('E')}
	cpu.instructions[0x9c] = []func(){cpu.sbcR8('H')}
	cpu.instructions[0x9d] = []func(){cpu.sbcR8('L')}
	cpu.instructions[0x9F] = []func(){cpu.sbcR8('A')}
}

func (cpu *CPU) initCP() {
	// CP [HL]
	cpu.instructions[0xBE] = []func(){cpu.readR16Addr("HL"), cpu.cpR8('T')}
	// CP n8
	cpu.instructions[0xFE] = []func(){cpu.readPCAddrAndInc, cpu.cpR8('T')}

	cpu.init_CP_R8()
}

func (cpu *CPU) init_CP_R8() {
	cpu.instructions[0xB8] = []func(){cpu.cpR8('B')}
	cpu.instructions[0xB9] = []func(){cpu.cpR8('C')}
	cpu.instructions[0xBa] = []func(){cpu.cpR8('D')}
	cpu.instructions[0xBb] = []func(){cpu.cpR8('E')}
	cpu.instructions[0xBc] = []func(){cpu.cpR8('H')}
	cpu.instructions[0xBd] = []func(){cpu.cpR8('L')}
	cpu.instructions[0xBF] = []func(){cpu.cpR8('A')}
}
