package debugger

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	entriesOnScreen     = 16
	regsNum             = 10
	flagsNum            = 4
	memoryLinesOnScreen = 8
	memoryBytesInLine   = 16
	memoryBytesOnScreen = memoryBytesInLine * memoryLinesOnScreen
	entryFormat         = "%02X    %02X\t  %-10s%s"
	entryWidth          = 52
	regWidth            = 11
	flagWidth           = 5
	stackWidth          = 11
	memoryWidth         = width - stackWidth - 3

	height = 1 + entriesOnScreen + 1 + memoryLinesOnScreen + 1
	width  = 1 + entryWidth + 1 + regWidth + 1 + flagWidth + 1

	instructionsLabel = "Instructions"
	registersLabel    = "Registers"
	flagsLabel        = "Flg"
	memoryLabel       = "Memory"
	stackLabel        = "Stck"
)

// a little helper for printing flags
var (
	startMemoryAddr = uint16(0xFFF0)
	boolNum         = map[bool]int{
		true:  1,
		false: 0,
	}

	hideCursor = "\033[?25l"

	blankField = strings.Repeat(strings.Repeat(" ", width)+"\n", height) + "\r\033[" + strconv.Itoa(height+1) + "A"

	topBorder = ("╔═" + colorizeStr(instructionsLabel, "red") +
		strings.Repeat("═", entryWidth-1-len(instructionsLabel)) +
		"╦═" + colorizeStr(registersLabel, "yellow") +
		strings.Repeat("═", regWidth-1-len(registersLabel)) +
		"╦═" + colorizeStr(flagsLabel, "blue") +
		strings.Repeat("═", flagWidth-1-len(flagsLabel)) +
		"╗\n")

	topSideBorders = strings.Repeat("║"+"\033["+strconv.Itoa(entryWidth)+
		"C│\033["+strconv.Itoa(regWidth)+"C│\033["+
		strconv.Itoa(flagWidth)+"C║\n", entriesOnScreen)

	middleBorder = ("╠─" + colorizeStr(memoryLabel, "green") +
		strings.Repeat("─", entryWidth-1-len(memoryLabel)) +
		"┴" + strings.Repeat("─", memoryWidth-1-entryWidth) +
		"┬─" + colorizeStr(stackLabel, "purple") +
		strings.Repeat("─", stackWidth-2-flagWidth-len(stackLabel)) +
		"┴" + strings.Repeat("─", flagWidth) + "╣\n")

	bottomSideBorders = strings.Repeat("║"+
		strings.Repeat(" ", memoryWidth)+
		"│"+strings.Repeat(" ", regWidth)+
		"║\n", memoryLinesOnScreen)

	bottomBorder = ("╚" + strings.Repeat("═", memoryWidth) + "╩" +
		strings.Repeat("═", stackWidth) + "╝\n\r\033[" +
		strconv.Itoa(height+1)) + "A"

	entriesBlankField = strings.Repeat(entriesOffset+strings.Repeat(" ", entryWidth), entriesOnScreen)
	entriesReturn     = "\r\033[" + strconv.Itoa(entriesOnScreen) + "A"
	entriesOffset     = "\r\033[1C\033[1B"
	regsOffset        = "\r\033[" + strconv.Itoa(3+entryWidth) + "C\033[1B"
	regsReturn        = "\r\033[" + strconv.Itoa(regsNum) + "A"
	flagsOffset       = "\r\033[" + strconv.Itoa(4+entryWidth+regWidth) + "C\033[1B"
	flagsReturn       = "\r\033[" + strconv.Itoa(flagsNum) + "A"
	memoryInitOffset  = "\r\033[" + strconv.Itoa(1+entriesOnScreen) + "B"
	memoryOffset      = entriesOffset
	memoryReturn      = "\r\033[" + strconv.Itoa(1+entriesOnScreen+memoryLinesOnScreen) + "A"
	stackInitOffset   = "\r\033[" + strconv.Itoa(2+entriesOnScreen+memoryLinesOnScreen) + "B"
	stackOffset       = "\r\033[1A\033[" + strconv.Itoa(2+memoryWidth) + "C"
	stackReturn       = "\r\033[" + strconv.Itoa(2+entriesOnScreen) + "A"
)

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

func (deb *Debugger) printData() error {
	deb.printEntries()
	deb.printRegs("red")
	deb.printFlags("red")

	err := deb.printMemory()
	if err != nil {
		return err
	}

	err = deb.printStack()
	if err != nil {
		return err
	}

	time.Sleep(time.Second)
	return err
}

func (deb *Debugger) printEntries() {
	fmt.Print(entriesBlankField + entriesReturn)

	entry := deb.entriesTail
	for entry != nil {
		fmt.Print(entriesOffset + entry.str)
		entry = entry.next
	}
	fmt.Print(entriesReturn)
}

func (deb *Debugger) printRegs(changedColor string) {
	for _, d := range deb.reg8Pairs {
		var val string
		if d.compare() {
			val = colorizeStr(fmt.Sprintf("%02X", *d.cpuReg), changedColor)
		} else {
			val = fmt.Sprintf("%02X", *d.cpuReg)
		}

		fmt.Printf("%s %s = %s", regsOffset, d.name, val)
	}

	for _, d := range deb.reg16Pairs {
		var val string
		if d.compare() {
			val = colorizeStr(fmt.Sprintf("%04X", d.cpuReg()), changedColor)
		} else {
			val = fmt.Sprintf("%04X", d.cpuReg())
		}

		fmt.Printf("%s%s = %s", regsOffset, d.name, val)
	}

	fmt.Print(regsReturn)
}

func (deb *Debugger) printFlags(changedColor string) {
	for _, d := range deb.flagPairs {
		var val string
		if d.compare() {
			val = colorizeStr(fmt.Sprintf("%d", boolNum[d.cpuReg()]), changedColor)
		} else {
			val = fmt.Sprintf("%d", boolNum[d.cpuReg()])
		}

		fmt.Printf("%s%s=%s", flagsOffset, d.name, val)
	}

	fmt.Print(flagsReturn)
}

func (deb *Debugger) printMemory() error {
	mem, err := deb.bus.ReadBatch(startMemoryAddr, startMemoryAddr+memoryBytesOnScreen-1)
	if err != nil {
		// TODO: process it properly
		return err
	}

	fmt.Print(memoryInitOffset)
	for i, v := range mem {
		if i%memoryBytesInLine == 0 {
			fmt.Printf("%s%04X:  ", memoryOffset, startMemoryAddr+uint16(i))
		}
		fmt.Printf("%02X ", v)
	}

	fmt.Print(memoryReturn)
	return nil
}

func (deb *Debugger) printStack() error {
	mem, err := deb.bus.ReadBatch(deb.regs.SP(), deb.regs.SP()+memoryLinesOnScreen*2-1)
	if err != nil {
		// TODO: process it properly
		return err
	}

	fmt.Print(stackInitOffset)
	for i := 0; i < len(mem); i += 2 {
		fmt.Printf(
			"%s%04X: %02X%02X",
			stackOffset,
			deb.regs.SP()+uint16(i),
			mem[i],
			mem[i+1],
		)
	}

	fmt.Print(stackReturn)
	return nil
}
