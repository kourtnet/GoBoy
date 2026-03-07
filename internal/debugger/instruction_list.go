package debugger

import (
	"fmt"

	"github.com/rivo/tview"
)

var entryFormat = "%-6s %-12s %-10s %s"

// custom tview.List that can work with instruction
// addresses instead of line indexes
type instructionList struct {
	*tview.List

	addrs    map[uint16]int
	prevInd  int
	prevText string

	isPrevItem bool
}

func newInstructionList() *instructionList {
	res := &instructionList{
		List:  tview.NewList(),
		addrs: map[uint16]int{},
	}

	return res
}

func (il *instructionList) addItem(pc uint16, byteSeq []byte, inst instruction) {
	il.addrs[pc] = il.GetItemCount()

	text := fmt.Sprintf(
		entryFormat,
		fmt.Sprintf("%04X", pc),
		fmt.Sprintf("%02X", byteSeq),
		inst.name,
		inst.ArgsFormat,
	)

	il.AddItem(text, "", 0, nil)
}

func (il *instructionList) getText(pc uint16) string {
	res, _ := il.GetItemText(il.addrs[pc])
	return res
}

func (il *instructionList) setCurrItem(pc uint16) {
	if il.isPrevItem {
		il.SetItemText(il.prevInd, il.prevText, "")
	}

	il.isPrevItem = true
	il.prevInd = il.addrs[pc]
	il.prevText = il.getText(pc)

	res := il.prevText[:5] + "►" + il.prevText[6:]
	il.SetItemText(il.prevInd, res, "")
	il.SetCurrentItem(il.prevInd)
}
