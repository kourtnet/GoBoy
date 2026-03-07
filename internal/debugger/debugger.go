// Package debugger provides CLI-debugger for CPU
package debugger

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/kourtnet/GoBoy/internal/bus"
	"github.com/kourtnet/GoBoy/internal/cpu"
	"github.com/rivo/tview"
)

type regType interface {
	~byte | ~uint16 | ~bool
}

var (
	instructionsTitle = "Instructions"
	registersTitle    = "Regs"
	flagsTitle        = "Flg"
	cyclesTitle       = "Cycls"
	memoryTitle       = "Memory"
	stackTitle        = "Stck"
)

const (
	romAddrEnd    uint16 = 0x7FFF
	regTableLen          = 11
	flagTableLen         = 7
	cyclesNumLen         = 3
	stackTableLen        = 14
)

type tui struct {
	app *tview.Application

	instList       *instructionList
	regTable       *tview.Table
	flagTable      *tview.Table
	cyclesNum      *tview.TextView
	flagCyclesFlex *tview.Flex
	topFlex        *tview.Flex

	memTable     *tview.Table
	memContent   *memoryContent
	stackTable   *tview.Table
	stackContent *stackContent
	bottomFlex   *tview.Flex

	mainFlex *tview.Flex
}

func newTUI(memory []byte, sp uint16) *tui {
	tui := &tui{}
	tui.app = tview.NewApplication()

	tui.instList = newInstructionList()
	tui.instList.
		SetSelectedFocusOnly(true).
		SetHighlightFullLine(true).
		ShowSecondaryText(false).
		SetBorder(true).
		SetTitle(instructionsTitle).
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorGreen)

	tui.regTable = tview.NewTable()
	tui.regTable.
		SetSelectable(false, false).
		SetTitle(registersTitle).
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(tcell.ColorYellow).
		SetBorder(true)

	tui.regTable.SetCellSimple(0, 0, " a")
	tui.regTable.SetCellSimple(1, 0, " f")
	tui.regTable.SetCellSimple(2, 0, " b")
	tui.regTable.SetCellSimple(3, 0, " c")
	tui.regTable.SetCellSimple(4, 0, " d")
	tui.regTable.SetCellSimple(5, 0, " e")
	tui.regTable.SetCellSimple(6, 0, " h")
	tui.regTable.SetCellSimple(7, 0, " l")
	tui.regTable.SetCellSimple(8, 0, " pc")
	tui.regTable.SetCellSimple(9, 0, " sp")

	tui.flagTable = tview.NewTable()
	tui.flagTable.
		SetSelectable(false, false).
		SetSeparator('=').
		SetTitle(flagsTitle).
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(tcell.ColorRed).
		SetBorder(true)

	tui.flagTable.SetCellSimple(0, 0, " z")
	tui.flagTable.SetCellSimple(1, 0, " n")
	tui.flagTable.SetCellSimple(2, 0, " h")
	tui.flagTable.SetCellSimple(3, 0, " c")

	tui.cyclesNum = tview.NewTextView()
	tui.cyclesNum.
		SetTextAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitle(cyclesTitle).
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(tcell.ColorBlue)

	tui.flagCyclesFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tui.flagTable, 0, 1, false).
		AddItem(tui.cyclesNum, cyclesNumLen, 1, false)

	tui.topFlex = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tui.instList, 0, 1, true).
		AddItem(tui.regTable, regTableLen, 1, false).
		AddItem(tui.flagCyclesFlex, flagTableLen, 0, false)

	tui.memContent = newMemoryTable(memory)
	tui.memTable = tview.NewTable()
	tui.memTable.
		SetSelectable(false, false).
		SetTitle(memoryTitle).
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorBlue).
		SetBorder(true)
	tui.memTable.SetContent(tui.memContent)

	tui.stackContent = newStackContent(memory, sp)
	tui.stackTable = tview.NewTable()
	tui.stackTable.
		SetBorder(true).
		SetTitle(stackTitle).
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(tcell.ColorPurple)
	tui.stackTable.SetContent(tui.stackContent)
	tui.stackTable.SetOffset(int(tui.stackContent.sp), 0)

	tui.bottomFlex = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tui.memTable, 0, 1, true).
		AddItem(
			tui.stackTable,
			stackTableLen,
			1,
			true,
		)

	tui.mainFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tui.topFlex, 0, 3, true).
		AddItem(tui.bottomFlex, 0, 2, false)

	tui.app.SetRoot(tui.mainFlex, true)

	return tui
}

type Debugger struct {
	cpu *cpu.CPU
	bus *bus.Bus

	// Required to check which cpu parameters have
	// changed after the Step()
	regs         cpu.Registers
	instructions [instructionsNum]instruction
	tui          *tui
	doneCh       chan struct{}
	err          error
}

func New(cpu *cpu.CPU, b *bus.Bus) (*Debugger, error) {
	deb := &Debugger{
		cpu:    cpu,
		bus:    b,
		regs:   *cpu.Registers,
		tui:    newTUI(b.Memory(), cpu.Registers.SP()),
		doneCh: make(chan struct{}),
	}

	deb.initInstructions()
	deb.loadROM()
	deb.initControls()

	return deb, nil
}

func (deb *Debugger) loadROM() error {
	data := deb.bus.Memory()

	for i := 0; i < len(data); {
		pc := uint16(i)
		byteSeq := make([]byte, 1, 3)

		byteSeq[0] = data[i]

		inst := deb.instructions[byteSeq[0]]
		if inst.name == "" {
			inst.name = "UNKNOWN"
		}

		inst.ArgsFormat = strings.ReplaceAll(inst.ArgsFormat, "n16", "n8n8")

		for range strings.Count(inst.ArgsFormat, "n8") {
			i++
			if i >= len(data) {
				inst.ArgsFormat = strings.ReplaceAll(inst.ArgsFormat, "n8", "??")
				inst.ArgsFormat += " (INCOMPLETE)"
				break
			}

			byteSeq = append(byteSeq, data[i])

			arg := fmt.Sprintf("%02X", data[i])
			substrInd := strings.LastIndex(inst.ArgsFormat, "n8")

			inst.ArgsFormat = (inst.ArgsFormat[:substrInd] +
				arg + inst.ArgsFormat[substrInd+len("n8"):])
		}

		deb.tui.instList.addItem(pc, byteSeq, inst)
		i++
	}

	return nil
}

func compareRegs[T regType](table *tview.Table, row int, r1, r2 T) {
	var cellText string
	switch any(r2).(type) {
	case byte:
		cellText = fmt.Sprintf("%02X", r2)
	case uint16:
		cellText = fmt.Sprintf("%04X", r2)
	case bool:
		b2 := any(r2).(bool)
		if b2 {
			cellText = "1 "
		} else {
			cellText = "0 "
		}
	}

	cell := table.GetCell(row, 1).SetText(cellText)
	table.SetCell(row, 1, cell)

	if r1 != r2 {
		table.GetCell(row, 0).SetTextColor(tcell.ColorRed)
		table.GetCell(row, 1).SetTextColor(tcell.ColorRed)
	} else {
		table.GetCell(row, 0).SetTextColor(tcell.ColorWhite)
		table.GetCell(row, 1).SetTextColor(tcell.ColorWhite)
	}
}

func (deb *Debugger) setRegs() {
	table := deb.tui.regTable
	regs := deb.cpu.Registers
	debRegs := deb.regs

	compareRegs(table, 0, debRegs.A, regs.A)
	compareRegs(table, 1, debRegs.F, regs.F)
	compareRegs(table, 2, debRegs.B, regs.B)
	compareRegs(table, 3, debRegs.C, regs.C)
	compareRegs(table, 4, debRegs.D, regs.D)
	compareRegs(table, 5, debRegs.E, regs.E)
	compareRegs(table, 6, debRegs.H, regs.H)
	compareRegs(table, 7, debRegs.L, regs.L)
	compareRegs(table, 8, debRegs.PC(), regs.PC())
	compareRegs(table, 9, debRegs.SP(), regs.SP())

	table = deb.tui.flagTable

	compareRegs(table, 0, debRegs.FlagZ(), regs.FlagZ())
	compareRegs(table, 1, debRegs.FlagN(), regs.FlagN())
	compareRegs(table, 2, debRegs.FlagH(), regs.FlagH())
	compareRegs(table, 3, debRegs.FlagC(), regs.FlagC())
}

func (deb *Debugger) stepInst() {
	for {
		select {
		case <-deb.doneCh:
			return
		default:
		}

		cpuEnd, err := deb.cpu.Step()
		if err != nil {
			deb.err = err
			close(deb.doneCh)
			return
		}

		if cpuEnd {
			close(deb.doneCh)
			return
		}

		if deb.cpu.ReadNewInstruction {
			deb.tui.instList.setCurrItem(deb.cpu.Registers.PC() - 1)
			deb.tui.stackContent.sp = deb.cpu.Registers.SP()
			deb.setRegs()
			deb.regs = *deb.cpu.Registers
			deb.tui.cyclesNum.SetText(fmt.Sprintf("%d", deb.cpu.MCycle))

			break
		}
	}
}

func (deb *Debugger) initControls() {
	deb.tui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'b':
			deb.stepInst()
			return nil
		case 's':
			deb.tui.stackTable.SetOffset(
				int(deb.cpu.Registers.SP()),
				0,
			)
			return nil
		case '1':
			deb.tui.app.SetFocus(deb.tui.instList)
			return nil
		case '2':
			deb.tui.app.SetFocus(deb.tui.regTable)
			return nil
		case '3':
			deb.tui.app.SetFocus(deb.tui.flagTable)
			return nil
		case '4':
			deb.tui.app.SetFocus(deb.tui.memTable)
			return nil
		case '5':
			deb.tui.app.SetFocus(deb.tui.stackTable)
			return nil
		}

		return event
	})

	deb.tui.app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		width, _ := screen.Size()
		bytesInLine := (width - 10 - stackTableLen) / 3
		if bytesInLine != deb.tui.memContent.bytesInLine {
			deb.tui.memContent.bytesInLine = bytesInLine
		}
		return false
	})
}

func (deb *Debugger) Run() error {
	cpuEnd, err := deb.cpu.Step()
	if err != nil {
		return err
	}

	if cpuEnd {
		return nil
	}

	deb.tui.instList.setCurrItem(deb.cpu.Registers.PC() - 1)
	deb.setRegs()
	deb.regs = *deb.cpu.Registers
	deb.tui.cyclesNum.SetText(fmt.Sprintf("%d", deb.cpu.MCycle))

	go func() {
		deb.err = deb.tui.app.Run()

		select {
		case <-deb.doneCh:
		default:
			close(deb.doneCh)
		}
	}()

	<-deb.doneCh
	deb.tui.app.Stop()

	if deb.err != nil {
		return deb.err
	}

	return nil
}
