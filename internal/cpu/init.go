package cpu

func (cpu *CPU) initInstructions() {
	cpu.initNOP()

	cpu.initLD()
	cpu.InitPush()
	cpu.InitPop()

	cpu.initArithmetic()
	cpu.initLogic()
	cpu.initCall()
}

func (cpu *CPU) initNOP() {
	cpu.instructions[0x00] = []func(){cpu.nop}
}

func (cpu *CPU) initLD() {
	// LD r8, r8
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

	// LD r8, n8
	cpu.instructions[0x06] = []func(){cpu.readN8, cpu.ldR8R8('B', 'T')}
	cpu.instructions[0x16] = []func(){cpu.readN8, cpu.ldR8R8('D', 'T')}
	cpu.instructions[0x26] = []func(){cpu.readN8, cpu.ldR8R8('H', 'T')}
	cpu.instructions[0x0E] = []func(){cpu.readN8, cpu.ldR8R8('C', 'T')}
	cpu.instructions[0x1E] = []func(){cpu.readN8, cpu.ldR8R8('E', 'T')}
	cpu.instructions[0x2E] = []func(){cpu.readN8, cpu.ldR8R8('L', 'T')}
	cpu.instructions[0x3E] = []func(){cpu.readN8, cpu.ldR8R8('A', 'T')}

	// LD r8, [R16]
	cpu.instructions[0x46] = []func(){cpu.readR16Addr("HL"), cpu.ldR8R8('B', 'T')}
	cpu.instructions[0x56] = []func(){cpu.readR16Addr("HL"), cpu.ldR8R8('D', 'T')}
	cpu.instructions[0x66] = []func(){cpu.readR16Addr("HL"), cpu.ldR8R8('H', 'T')}
	cpu.instructions[0x4E] = []func(){cpu.readR16Addr("HL"), cpu.ldR8R8('C', 'T')}
	cpu.instructions[0x5E] = []func(){cpu.readR16Addr("HL"), cpu.ldR8R8('E', 'T')}
	cpu.instructions[0x6E] = []func(){cpu.readR16Addr("HL"), cpu.ldR8R8('L', 'T')}
	cpu.instructions[0x7E] = []func(){cpu.readR16Addr("HL"), cpu.ldR8R8('A', 'T')}
	cpu.instructions[0x0A] = []func(){cpu.readR16Addr("BC"), cpu.ldR8R8('A', 'T')}
	cpu.instructions[0x1A] = []func(){cpu.readR16Addr("DE"), cpu.ldR8R8('A', 'T')}

	// LD [R16], r8
	cpu.instructions[0x02] = []func(){cpu.ldAddrR8("BC", 'A'), cpu.nop}
	cpu.instructions[0x12] = []func(){cpu.ldAddrR8("DE", 'A'), cpu.nop}
	cpu.instructions[0x70] = []func(){cpu.ldAddrR8("HL", 'B'), cpu.nop}
	cpu.instructions[0x71] = []func(){cpu.ldAddrR8("HL", 'C'), cpu.nop}
	cpu.instructions[0x72] = []func(){cpu.ldAddrR8("HL", 'D'), cpu.nop}
	cpu.instructions[0x73] = []func(){cpu.ldAddrR8("HL", 'E'), cpu.nop}
	cpu.instructions[0x74] = []func(){cpu.ldAddrR8("HL", 'H'), cpu.nop}
	cpu.instructions[0x75] = []func(){cpu.ldAddrR8("HL", 'L'), cpu.nop}
	cpu.instructions[0x77] = []func(){cpu.ldAddrR8("HL", 'A'), cpu.nop}

	// LD [HL], n8
	cpu.instructions[0x36] = []func(){cpu.readN8, cpu.ldAddrR8("HL", 'T'), cpu.nop}

	// LD [N16], A
	cpu.instructions[0xEA] = []func(){cpu.readN16Lsb, cpu.readN16Msb, cpu.ldAddrR8("Temp16", 'A'), cpu.nop}

	// LD A, [N16]
	cpu.instructions[0xFA] = []func(){cpu.readN16Lsb, cpu.readN16Msb, cpu.readR16Addr("Temp16"), cpu.ldR8R8('A', 'T')}

	// LDH (those are not tested by ROM tests)
	cpu.instructions[0xE0] = []func(){cpu.readN8, cpu.ldhAddrR8('T', 'A'), cpu.nop}
	cpu.instructions[0xE2] = []func(){cpu.ldhAddrR8('C', 'A'), cpu.nop}
	cpu.instructions[0xF0] = []func(){cpu.readN8, cpu.readR8Addr('T'), cpu.ldR8R8('A', 'T')}
	cpu.instructions[0xF2] = []func(){cpu.readR8Addr('C'), cpu.ldR8R8('A', 'T')}

	// LD [HL+], A
	cpu.instructions[0x22] = []func(){cpu.ldHLPlusACycle1, cpu.nop}

	// LD A, [HL+]
	cpu.instructions[0x2A] = []func(){cpu.readHLAddrInc, cpu.ldR8R8('A', 'T')}

	// LD [HL-], A
	cpu.instructions[0x32] = []func(){cpu.ldR16AddrDecR8("HL", 'A'), cpu.nop}

	// LD A, [HL-]
	cpu.instructions[0x3A] = []func(){cpu.readHLAddrDec, cpu.ldR8R8('A', 'T')}

	// LD SP, HL
	cpu.instructions[0xF9] = []func(){cpu.ldR16R16("SP", "HL"), cpu.nop}

	// LD [n16], SP
	cpu.instructions[0x08] = []func(){cpu.readN16Lsb, cpu.readN16Msb, cpu.ldN16AddrSPCycle3, cpu.ldN16AddrSPCycle4, cpu.nop}

	// LD HL, SP + n8
	cpu.instructions[0xF8] = []func(){cpu.readN8, cpu.LDHLSPPlusN8Cycle2, cpu.LDHLSPPlusN8Cycle3}

	// LD R16, N16
	cpu.instructions[0x01] = []func(){cpu.readN16Lsb, cpu.readN16Msb, cpu.ldR16R16("BC", "Temp16")}
	cpu.instructions[0x11] = []func(){cpu.readN16Lsb, cpu.readN16Msb, cpu.ldR16R16("DE", "Temp16")}
	cpu.instructions[0x21] = []func(){cpu.readN16Lsb, cpu.readN16Msb, cpu.ldR16R16("HL", "Temp16")}
	cpu.instructions[0x31] = []func(){cpu.readN16Lsb, cpu.readN16Msb, cpu.ldR16R16("SP", "Temp16")}
}

func (cpu *CPU) InitPush() {
	cpu.instructions[0xC5] = []func(){cpu.decSp, cpu.ldR16AddrDecR8("SP", 'B'), cpu.ldAddrR8("SP", 'C'), cpu.nop}
	cpu.instructions[0xD5] = []func(){cpu.decSp, cpu.ldR16AddrDecR8("SP", 'D'), cpu.ldAddrR8("SP", 'E'), cpu.nop}
	cpu.instructions[0xE5] = []func(){cpu.decSp, cpu.ldR16AddrDecR8("SP", 'H'), cpu.ldAddrR8("SP", 'L'), cpu.nop}
	cpu.instructions[0xF5] = []func(){cpu.decSp, cpu.ldR16AddrDecR8("SP", 'A'), cpu.ldAddrR8("SP", 'F'), cpu.nop}
}

func (cpu *CPU) InitPop() {
	cpu.instructions[0xC1] = []func(){cpu.popLsb, cpu.popMsb, cpu.ldR16R16("BC", "Temp16")}
	cpu.instructions[0xD1] = []func(){cpu.popLsb, cpu.popMsb, cpu.ldR16R16("DE", "Temp16")}
	cpu.instructions[0xE1] = []func(){cpu.popLsb, cpu.popMsb, cpu.ldR16R16("HL", "Temp16")}
	cpu.instructions[0xF1] = []func(){cpu.popLsb, cpu.popMsb, cpu.ldR16R16("AF", "Temp16")}
}

func (cpu *CPU) initArithmetic() {
	// CCF
	cpu.instructions[0x3F] = []func(){cpu.ccf}
	// SCF
	cpu.instructions[0x37] = []func(){cpu.scf}
	// DAA
	cpu.instructions[0x27] = []func(){cpu.daa}
	// CPL
	cpu.instructions[0x2F] = []func(){cpu.cpl}

	// ADD [HL]
	cpu.instructions[0x86] = []func(){cpu.readR16Addr("HL"), cpu.addR8('T')}

	// ADD n8
	cpu.instructions[0xC6] = []func(){cpu.readN8, cpu.addR8('T')}

	// ADC [HL]
	cpu.instructions[0x8E] = []func(){cpu.readR16Addr("HL"), cpu.adcR8('T')}

	// ADC n8
	cpu.instructions[0xCE] = []func(){cpu.readN8, cpu.adcR8('T')}

	// ADD SP, n8
	cpu.instructions[0xE8] = []func(){cpu.readN8, cpu.addSPN8Cycle2, cpu.addSPN8Cycle3, cpu.ldR16R16("SP", "Temp16")}

	// ADD r8
	cpu.instructions[0x80] = []func(){cpu.addR8('B')}
	cpu.instructions[0x81] = []func(){cpu.addR8('C')}
	cpu.instructions[0x82] = []func(){cpu.addR8('D')}
	cpu.instructions[0x83] = []func(){cpu.addR8('E')}
	cpu.instructions[0x84] = []func(){cpu.addR8('H')}
	cpu.instructions[0x85] = []func(){cpu.addR8('L')}
	cpu.instructions[0x87] = []func(){cpu.addR8('A')}

	// ADD R16
	cpu.instructions[0x09] = []func(){cpu.addR16LSB('C'), cpu.addR16MSB('B')}
	cpu.instructions[0x19] = []func(){cpu.addR16LSB('E'), cpu.addR16MSB('D')}
	cpu.instructions[0x29] = []func(){cpu.addR16LSB('L'), cpu.addR16MSB('H')}
	cpu.instructions[0x39] = []func(){cpu.addR16LSB('P'), cpu.addR16MSB('S')}

	// ADD r8
	cpu.instructions[0x88] = []func(){cpu.adcR8('B')}
	cpu.instructions[0x89] = []func(){cpu.adcR8('C')}
	cpu.instructions[0x8A] = []func(){cpu.adcR8('D')}
	cpu.instructions[0x8B] = []func(){cpu.adcR8('E')}
	cpu.instructions[0x8C] = []func(){cpu.adcR8('H')}
	cpu.instructions[0x8D] = []func(){cpu.adcR8('L')}
	cpu.instructions[0x8F] = []func(){cpu.adcR8('A')}

	// SUB [HL]
	cpu.instructions[0x96] = []func(){cpu.readR16Addr("HL"), cpu.subR8('T')}

	// SUB n8
	cpu.instructions[0xD6] = []func(){cpu.readN8, cpu.subR8('T')}

	// SBC [HL]
	cpu.instructions[0x9E] = []func(){cpu.readR16Addr("HL"), cpu.sbcR8('T')}

	// SBC n8
	cpu.instructions[0xDE] = []func(){cpu.readN8, cpu.sbcR8('T')}

	// SUB r8
	cpu.instructions[0x90] = []func(){cpu.subR8('B')}
	cpu.instructions[0x91] = []func(){cpu.subR8('C')}
	cpu.instructions[0x92] = []func(){cpu.subR8('D')}
	cpu.instructions[0x93] = []func(){cpu.subR8('E')}
	cpu.instructions[0x94] = []func(){cpu.subR8('H')}
	cpu.instructions[0x95] = []func(){cpu.subR8('L')}
	cpu.instructions[0x97] = []func(){cpu.subR8('A')}

	// SBC r8
	cpu.instructions[0x98] = []func(){cpu.sbcR8('B')}
	cpu.instructions[0x99] = []func(){cpu.sbcR8('C')}
	cpu.instructions[0x9a] = []func(){cpu.sbcR8('D')}
	cpu.instructions[0x9b] = []func(){cpu.sbcR8('E')}
	cpu.instructions[0x9c] = []func(){cpu.sbcR8('H')}
	cpu.instructions[0x9d] = []func(){cpu.sbcR8('L')}
	cpu.instructions[0x9F] = []func(){cpu.sbcR8('A')}

	// CP [HL]
	cpu.instructions[0xBE] = []func(){cpu.readR16Addr("HL"), cpu.cpR8('T')}
	// CP n8
	cpu.instructions[0xFE] = []func(){cpu.readN8, cpu.cpR8('T')}

	// CP r8
	cpu.instructions[0xB8] = []func(){cpu.cpR8('B')}
	cpu.instructions[0xB9] = []func(){cpu.cpR8('C')}
	cpu.instructions[0xBa] = []func(){cpu.cpR8('D')}
	cpu.instructions[0xBb] = []func(){cpu.cpR8('E')}
	cpu.instructions[0xBc] = []func(){cpu.cpR8('H')}
	cpu.instructions[0xBd] = []func(){cpu.cpR8('L')}
	cpu.instructions[0xBF] = []func(){cpu.cpR8('A')}

	// INC [HL]
	cpu.instructions[0x34] = []func(){cpu.readR16Addr("HL"), cpu.incHLAddrCycle2, cpu.nop}

	// INC r8
	cpu.instructions[0x04] = []func(){cpu.incR8('B')}
	cpu.instructions[0x14] = []func(){cpu.incR8('D')}
	cpu.instructions[0x24] = []func(){cpu.incR8('H')}
	cpu.instructions[0x0C] = []func(){cpu.incR8('C')}
	cpu.instructions[0x1C] = []func(){cpu.incR8('E')}
	cpu.instructions[0x2C] = []func(){cpu.incR8('L')}
	cpu.instructions[0x3C] = []func(){cpu.incR8('A')}

	// INC R16
	cpu.instructions[0x03] = []func(){cpu.determine16RegInc("BC"), cpu.nop}
	cpu.instructions[0x13] = []func(){cpu.determine16RegInc("DE"), cpu.nop}
	cpu.instructions[0x23] = []func(){cpu.determine16RegInc("HL"), cpu.nop}
	cpu.instructions[0x33] = []func(){cpu.determine16RegInc("SP"), cpu.nop}

	// DEC [HL]
	cpu.instructions[0x35] = []func(){cpu.readR16Addr("HL"), cpu.decHLAddrCycle2, cpu.nop}

	// DEC r8
	cpu.instructions[0x05] = []func(){cpu.decR8('B')}
	cpu.instructions[0x15] = []func(){cpu.decR8('D')}
	cpu.instructions[0x25] = []func(){cpu.decR8('H')}
	cpu.instructions[0x0D] = []func(){cpu.decR8('C')}
	cpu.instructions[0x1D] = []func(){cpu.decR8('E')}
	cpu.instructions[0x2D] = []func(){cpu.decR8('L')}
	cpu.instructions[0x3D] = []func(){cpu.decR8('A')}

	// DEC R16
	cpu.instructions[0x0B] = []func(){cpu.determine16RegDec("BC"), cpu.nop}
	cpu.instructions[0x1B] = []func(){cpu.determine16RegDec("DE"), cpu.nop}
	cpu.instructions[0x2B] = []func(){cpu.determine16RegDec("HL"), cpu.nop}
	cpu.instructions[0x3B] = []func(){cpu.determine16RegDec("SP"), cpu.nop}
}

func (cpu *CPU) initLogic() {
	// AND [HL]
	cpu.instructions[0xA6] = []func(){cpu.readR16Addr("HL"), cpu.andR8('T')}
	// AND n8
	cpu.instructions[0xE6] = []func(){cpu.readN8, cpu.andR8('T')}

	// AND r8
	cpu.instructions[0xA0] = []func(){cpu.andR8('B')}
	cpu.instructions[0xA1] = []func(){cpu.andR8('C')}
	cpu.instructions[0xA2] = []func(){cpu.andR8('D')}
	cpu.instructions[0xA3] = []func(){cpu.andR8('E')}
	cpu.instructions[0xA4] = []func(){cpu.andR8('H')}
	cpu.instructions[0xA5] = []func(){cpu.andR8('L')}
	cpu.instructions[0xA7] = []func(){cpu.andR8('A')}

	// OR [HL]
	cpu.instructions[0xB6] = []func(){cpu.readR16Addr("HL"), cpu.orR8('T')}
	// OR n8
	cpu.instructions[0xF6] = []func(){cpu.readN8, cpu.orR8('T')}

	// OR r8
	cpu.instructions[0xB0] = []func(){cpu.orR8('B')}
	cpu.instructions[0xB1] = []func(){cpu.orR8('C')}
	cpu.instructions[0xB2] = []func(){cpu.orR8('D')}
	cpu.instructions[0xB3] = []func(){cpu.orR8('E')}
	cpu.instructions[0xB4] = []func(){cpu.orR8('H')}
	cpu.instructions[0xB5] = []func(){cpu.orR8('L')}
	cpu.instructions[0xB7] = []func(){cpu.orR8('A')}

	// OR [HL]
	cpu.instructions[0xAE] = []func(){cpu.readR16Addr("HL"), cpu.xorR8('T')}
	// OR n8
	cpu.instructions[0xEE] = []func(){cpu.readN8, cpu.xorR8('T')}

	// OR r8
	cpu.instructions[0xA8] = []func(){cpu.xorR8('B')}
	cpu.instructions[0xA9] = []func(){cpu.xorR8('C')}
	cpu.instructions[0xAA] = []func(){cpu.xorR8('D')}
	cpu.instructions[0xAB] = []func(){cpu.xorR8('E')}
	cpu.instructions[0xAC] = []func(){cpu.xorR8('H')}
	cpu.instructions[0xAD] = []func(){cpu.xorR8('L')}
	cpu.instructions[0xAF] = []func(){cpu.xorR8('A')}
}

func (cpu *CPU) initCall() {
	// jump n16
	cpu.instructions[0xC3] = []func(){cpu.readN16Lsb, cpu.readN16Msb, cpu.ldR16R16("PC", "Temp16"), cpu.nop}
}
