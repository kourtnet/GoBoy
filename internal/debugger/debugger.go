// Package debugger provides CLI-debugger for CPU
package debugger

import (
	"fmt"

	"github.com/kourtnet/GoBoy/internal/cpu"
)

const (
	entriesOnScreen = 16
	entryFormat     = "%02X    %02X\t  %-10s%s"
)

type entry struct {
	str  string
	next *entry
}

type Debugger struct {
	cpu *cpu.CPU
	bus bus

	instructions        [instructionsNum]instruction
	entriesTail         *entry
	entriesHead         *entry
	entriesNum          int
	nextInstructionAddr uint16
}

func New(cpu *cpu.CPU, bus bus) (*Debugger, error) {
	deb := Debugger{
		cpu: cpu,
		bus: bus,
	}

	deb.initInstructions()
	deb.nextInstructionAddr = deb.cpu.Registers.PC()

	for range entriesOnScreen {
		err := deb.newEntry()
		if err != nil {
			return nil, err
		}
	}

	return &deb, nil
}

func (deb *Debugger) newEntry() error {
	str, err := deb.newEntryStr()
	if err != nil {
		return err
	}

	ent := &entry{str: str}
	if deb.entriesTail == nil {
		deb.entriesTail = ent
		deb.entriesHead = ent
		return nil
	}

	deb.entriesHead.next = ent
	deb.entriesHead = deb.entriesHead.next

	if deb.entriesNum >= entriesOnScreen {
		deb.entriesTail = deb.entriesTail.next
	}

	deb.entriesNum++

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

	byteSeq := make([]byte, instruction.argsNum+1)
	byteSeq[0] = opcode

	for i := range instruction.argsNum {
		arg, err := deb.bus.Read(deb.nextInstructionAddr)
		if err != nil {
			return "", err
		}

		byteSeq[i+1] = arg
		deb.nextInstructionAddr++
	}

	var argsStr string
	switch instruction.argsNum {
	case 0:
		argsStr = instruction.ArgsFormat
	case 1:
		argsStr = fmt.Sprintf(instruction.ArgsFormat, byteSeq[1])
	case 2:
		argsStr = fmt.Sprintf(instruction.ArgsFormat, byteSeq[2], byteSeq[1])
	default:
		return "", fmt.Errorf("invalid number of arguments for %x opcode", byteSeq[0])
	}

	str := fmt.Sprintf(
		entryFormat,
		instructionAddr,
		byteSeq,
		instruction.name,
		argsStr)

	return str, nil
}

func (deb *Debugger) PrintData() {
	entry := deb.entriesTail
	for entry != nil {
		fmt.Println(entry.str)
		entry = entry.next
	}
}
