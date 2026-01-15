// Package debugger provides CLI-debugger for CPU
package debugger

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kourtnet/GoBoy/internal/bus"
	"github.com/kourtnet/GoBoy/internal/cpu"
)

const (
	entriesOnScreen = 16
	entryFormat     = "%02X    %02X\t  %-10s%s"
	entryWidth      = 50
)

var (
	entriesBlankField  = strings.Repeat(strings.Repeat(" ", entryWidth)+"\n", entriesOnScreen)
	entriesFieldReturn = "\r\033[16A"
)

type entry struct {
	str  string
	next *entry
}

type Debugger struct {
	cpu *cpu.CPU
	bus *bus.Bus

	// Required to check which cpu parameters have
	// changed after the Step()
	regs cpu.Registers

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

	// Invalid opcode. Required for first step proper processing
	deb.regs.IR = 0xD3

	deb.initInstructions()
	deb.nextInstructionAddr = deb.cpu.Registers.PC()

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

func (deb *Debugger) printData() {
	fmt.Print(entriesFieldReturn)
	fmt.Print(entriesBlankField)
	fmt.Print(entriesFieldReturn)

	entry := deb.entriesTail
	for entry != nil {
		fmt.Println(entry.str)
		entry = entry.next
	}
	time.Sleep(time.Second * 1)
}

func (deb *Debugger) firstPrint() {
	fmt.Print(entriesBlankField)
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
		deb.newEntry()
		deb.printData()
		deb.regs = *deb.cpu.Registers
	}

	return cpuEnd, nil
}
