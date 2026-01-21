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

var instructionsTitle = "Instructions"

const romAddrEnd uint16 = 0x7FFF

type tui struct {
	app      *tview.Application
	instList *instructionList
	topFlex  *tview.Flex
}

// TODO: mb change to tui instead of deb?
func newTUI() *tui {
	tui := &tui{}
	tui.app = tview.NewApplication()

	tui.instList = newInstructionList()

	tui.instList.
		SetHighlightFullLine(true).
		ShowSecondaryText(false).
		SetBorder(true).
		SetTitle(instructionsTitle).
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorGreen)

	tui.topFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tui.instList, 0, 1, true)

	tui.app.SetRoot(tui.topFlex, true)

	return tui
}

type Debugger struct {
	cpu *cpu.CPU
	bus *bus.Bus

	// Required to check which cpu parameters have
	// changed after the Step()
	regs cpu.Registers
	//	reg8Pairs  [8]reg8Pair
	//	reg16Pairs [2]dataPair[uint16]
	//	flagPairs  [4]dataPair[bool]
	//
	instructions [instructionsNum]instruction
	//	entriesNum          int
	//	nextInstructionAddr uint16
	//	isNotFirstStep      bool

	tui *tui
	// initPC            uint16
	// lastInstAddr      uint16
	// oldRow            string
	// instructionsAddrs map[uint16]int
	doneCh chan struct{}
	err    error
}

func New(cpu *cpu.CPU, b *bus.Bus) (*Debugger, error) {
	deb := &Debugger{
		cpu:    cpu,
		bus:    b,
		tui:    newTUI(),
		doneCh: make(chan struct{}),
	}

	deb.initInstructions()
	deb.loadInstructions()
	deb.initControls()

	return deb, nil
}

func (deb *Debugger) loadInstructions() error {
	data, err := deb.bus.ReadBatch(
		deb.cpu.Registers.PC(),
		romAddrEnd,
	)
	if err != nil {
		return err
	}

	for i := 0; i < len(data); {
		pc := uint16(i) + deb.cpu.Registers.PC()
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

		deb.tui.instList.addItemByPC(pc, byteSeq, inst)
		i++
	}

	return nil
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
			deb.tui.instList.setCurrItemByPC(deb.cpu.Registers.PC() - 1)
			break
		}
	}
}

func (deb *Debugger) initControls() {
	deb.tui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'b' {
			deb.stepInst()
			return nil
		}

		return event
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

	deb.tui.instList.setCurrItemByPC(deb.cpu.Registers.PC() - 1)

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
