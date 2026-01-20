// Package debugger provides CLI-debugger for CPU
package debugger

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kourtnet/GoBoy/internal/bus"
	"github.com/kourtnet/GoBoy/internal/cpu"
)

type entry struct {
	str  string
	next *entry
}

// struct required for easier registers compare
type reg8Pair struct {
	debReg *byte
	cpuReg *byte
	name   string
}

func (d *reg8Pair) compare() bool {
	return *d.cpuReg != *d.debReg
}

type dataPair[T comparable] struct {
	debReg func() T
	cpuReg func() T
	name   string
}

func (d *dataPair[T]) compare() bool {
	return d.cpuReg() != d.debReg()
}

type Debugger struct {
	cpu *cpu.CPU
	bus *bus.Bus

	// Required to check which cpu parameters have
	// changed after the Step()
	regs       cpu.Registers
	reg8Pairs  [8]reg8Pair
	reg16Pairs [2]dataPair[uint16]
	flagPairs  [4]dataPair[bool]

	instructions        [instructionsNum]instruction
	entriesTail         *entry
	entriesHead         *entry
	entriesNum          int
	nextInstructionAddr uint16
	isNotFirstStep      bool
}

func New(cpu *cpu.CPU, b *bus.Bus) (*Debugger, error) {
	deb := Debugger{
		cpu:  cpu,
		bus:  b,
		regs: *cpu.Registers,
	}

	deb.initCPU()

	for deb.entriesNum < entriesOnScreen {
		err := deb.newEntry()
		if err != nil {
			var ourErr bus.OutOfRange
			if !errors.As(err, &ourErr) {
				return nil, err
			}

			break
		}
	}

	return &deb, nil
}

func (deb *Debugger) initCPU() {
	// Invalid opcode. Required for first step proper processing
	deb.regs.IR = 0xD3

	deb.initInstructions()
	deb.nextInstructionAddr = deb.cpu.Registers.PC()

	deb.reg8Pairs[0] = reg8Pair{&deb.regs.A, &deb.cpu.Registers.A, "a"}
	deb.reg8Pairs[1] = reg8Pair{&deb.regs.F, &deb.cpu.Registers.F, "f"}
	deb.reg8Pairs[2] = reg8Pair{&deb.regs.B, &deb.cpu.Registers.B, "b"}
	deb.reg8Pairs[3] = reg8Pair{&deb.regs.C, &deb.cpu.Registers.C, "c"}
	deb.reg8Pairs[4] = reg8Pair{&deb.regs.D, &deb.cpu.Registers.D, "d"}
	deb.reg8Pairs[5] = reg8Pair{&deb.regs.E, &deb.cpu.Registers.E, "e"}
	deb.reg8Pairs[6] = reg8Pair{&deb.regs.H, &deb.cpu.Registers.H, "h"}
	deb.reg8Pairs[7] = reg8Pair{&deb.regs.L, &deb.cpu.Registers.L, "l"}

	deb.reg16Pairs[0] = dataPair[uint16]{deb.regs.PC, deb.cpu.Registers.PC, "pc"}
	deb.reg16Pairs[1] = dataPair[uint16]{deb.regs.SP, deb.cpu.Registers.SP, "sp"}

	deb.flagPairs[0] = dataPair[bool]{deb.regs.FlagZ, deb.cpu.Registers.FlagZ, "z"}
	deb.flagPairs[1] = dataPair[bool]{deb.regs.FlagN, deb.cpu.Registers.FlagN, "n"}
	deb.flagPairs[2] = dataPair[bool]{deb.regs.FlagH, deb.cpu.Registers.FlagH, "h"}
	deb.flagPairs[3] = dataPair[bool]{deb.regs.FlagC, deb.cpu.Registers.FlagC, "c"}
}

func (deb *Debugger) newEntry() error {
	str, err := deb.newEntryStr()
	if err != nil {
		return err
	}

	ent := &entry{str: str}
	if deb.entriesTail == nil {
		// dummy for first call
		deb.entriesTail = &entry{next: ent}
		deb.entriesHead = ent
		deb.entriesNum = 2
		return nil
	}

	deb.entriesHead.next = ent
	deb.entriesHead = deb.entriesHead.next

	if deb.entriesNum >= entriesOnScreen {
		deb.entriesTail = deb.entriesTail.next
	} else {
		deb.entriesNum++
	}

	return nil
}

func (deb *Debugger) newEntryStr() (string, error) {
	opcode, err := deb.bus.Read(deb.nextInstructionAddr)
	if err != nil {
		return "", err
	}

	instructionAddr := deb.nextInstructionAddr
	deb.nextInstructionAddr++

	instruction := instruction{name: "UNKNOWN"}
	if int(opcode) <= len(deb.instructions) && deb.instructions[opcode].name != "" {
		instruction = deb.instructions[opcode]
	}

	instruction.ArgsFormat = strings.ReplaceAll(instruction.ArgsFormat, "n16", "n8n8")
	byteSeq := make([]byte, 1, 3)
	byteSeq[0] = opcode

	for strings.Contains(instruction.ArgsFormat, "n8") {
		arg, err := deb.bus.Read(deb.nextInstructionAddr)
		if err != nil {
			var oufErr bus.OutOfRange
			if !errors.As(err, &oufErr) {
				return "", err
			}

			instruction.ArgsFormat = strings.ReplaceAll(instruction.ArgsFormat, "n8", "??")
			instruction.ArgsFormat += " INCOMPLETE INSTRUCTION"
		}

		deb.nextInstructionAddr++

		idx := strings.LastIndex(instruction.ArgsFormat, "n8")
		if idx == -1 {
			break
		}

		instruction.ArgsFormat = instruction.ArgsFormat[:idx] + fmt.Sprintf("%02X", arg) + instruction.ArgsFormat[idx+len("n8"):]
		byteSeq = append(byteSeq, arg)
	}

	str := fmt.Sprintf(
		entryFormat,
		instructionAddr,
		byteSeq,
		instruction.name,
		instruction.ArgsFormat,
	)

	return str, nil
}

func (deb *Debugger) Step() (bool, error) {
	if !deb.isNotFirstStep {
		deb.firstPrint()
		deb.isNotFirstStep = true
	}

	cpuEnd, err := deb.cpu.Step()
	if err != nil {
		return cpuEnd, err
	}

	// Checking if new instruction has been read
	if deb.cpu.ReadNewInstruction {
		err := deb.newEntry()
		if err != nil {
			return false, err
		}

		err = deb.printData()
		if err != nil {
			return false, err
		}

		deb.regs = *deb.cpu.Registers
	}

	return cpuEnd, nil
}
