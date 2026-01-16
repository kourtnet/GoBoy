// Package debugger provides CLI-debugger for CPU
package debugger

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kourtnet/GoBoy/internal/bus"
	"github.com/kourtnet/GoBoy/internal/cpu"
)

const (
	entriesOnScreen     = 16
	memoryLinesOnScreen = 8
	entryFormat         = "%02X    %02X\t  %-10s%s"
	entryWidth          = 50
	registerWidth       = 11
	flagWidth           = 4
	stackWidth          = 11
	memoryWidth         = width - stackWidth - 3

	height = 1 + entriesOnScreen + 1 + memoryLinesOnScreen + 1
	width  = 1 + entryWidth + 1 + registerWidth + 1 + flagWidth + 1

	instructionsLabel = "Instructions"
	registersLabel    = "Registers"
	flagsLabel        = "Fl"
	memoryLabel       = "Memory"
	stackLabel        = "Stck"
)

var (
	hideCursor = "\033[?25l"

	blankField = strings.Repeat(strings.Repeat(" ", width)+"\n", height) + "\r\033[" + strconv.Itoa(height+1) + "A"

	topBorder = ("╔═" + colorizeStr(instructionsLabel, "red") +
		strings.Repeat("═", entryWidth-1-len(instructionsLabel)) +
		"╦═" + colorizeStr(registersLabel, "yellow") +
		strings.Repeat("═", registerWidth-1-len(registersLabel)) +
		"╦═" + colorizeStr(flagsLabel, "blue") +
		strings.Repeat("═", flagWidth-1-len(flagsLabel)) +
		"╗\n")

	topSideBorders = strings.Repeat("║"+"\033["+strconv.Itoa(entryWidth)+
		"C│\033["+strconv.Itoa(registerWidth)+"C│\033["+
		strconv.Itoa(flagWidth)+"C║\n", entriesOnScreen)

	middleBorder = ("╠─" + colorizeStr(memoryLabel, "green") +
		strings.Repeat("─", entryWidth-1-len(memoryLabel)) +
		"┴" + strings.Repeat("─", memoryWidth-1-entryWidth) +
		"┬─" + colorizeStr(stackLabel, "purple") +
		strings.Repeat("─", stackWidth-2-flagWidth-len(stackLabel)) +
		"┴" + strings.Repeat("─", flagWidth) + "╣\n")

	bottomSideBorders = strings.Repeat("║"+
		strings.Repeat(" ", memoryWidth)+
		"│"+strings.Repeat(" ", registerWidth)+
		"║\n", memoryLinesOnScreen)

	bottomBorder = ("╚" + strings.Repeat("═", memoryWidth) + "╩" +
		strings.Repeat("═", stackWidth) + "╝\n\r\033[" +
		strconv.Itoa(height+1)) + "A"

	entriesBlankField = strings.Repeat(entriesOffset+strings.Repeat(" ", entryWidth), entriesOnScreen)
	entriesReturn     = "\r\033[" + strconv.Itoa(entriesOnScreen) + "A"
	entriesOffset     = "\r\033[1C\033[1B"
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

func colorizeStr(str, color string) string {
	var colorStr string

	switch color {
	case "red":
		colorStr = "31"
	case "green":
		colorStr = "32"
	case "yellow":
		colorStr = "33"
	case "blue":
		colorStr = "34"
	case "purple":
		colorStr = "35"
	default:
		colorStr = "30"
	}

	return fmt.Sprintf("\033[%sm%s\033[0m", colorStr, str)
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
	fmt.Print(entriesBlankField + entriesReturn)

	entry := deb.entriesTail
	for entry != nil {
		fmt.Print(entriesOffset + entry.str)
		entry = entry.next
	}
	fmt.Print(entriesReturn)
	time.Sleep(time.Hour)
}

func (deb *Debugger) firstPrint() {
	fmt.Print(hideCursor)
	fmt.Println(blankField)

	fmt.Println(
		topBorder +
			topSideBorders +
			middleBorder +
			bottomSideBorders +
			bottomBorder,
	)
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
