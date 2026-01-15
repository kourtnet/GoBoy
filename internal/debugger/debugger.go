// Package debugger provides CLI-debugger for CPU
package debugger

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kourtnet/GoBoy/internal/bus"
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
	bus *bus.Bus

	instructions        [instructionsNum]instruction
	entriesTail         *entry
	entriesHead         *entry
	entriesNum          int
	nextInstructionAddr uint16
}

func New(cpu *cpu.CPU, b *bus.Bus) (*Debugger, error) {
	deb := Debugger{
		cpu: cpu,
		bus: b,
	}

	deb.initInstructions()
	deb.nextInstructionAddr = deb.cpu.Registers.PC()

	for range entriesOnScreen {
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

func (deb *Debugger) PrintData() {
	entry := deb.entriesTail
	for entry != nil {
		fmt.Println(entry.str)
		entry = entry.next
	}
}
