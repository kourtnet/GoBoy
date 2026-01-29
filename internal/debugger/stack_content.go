package debugger

import (
	"fmt"

	"github.com/rivo/tview"
)

type stackContent struct {
	*tview.TableContentReadOnly

	memory []byte
	sp     uint16
}

func newStackContent(memory []byte, sp uint16) *stackContent {
	return &stackContent{
		memory: memory,
		sp:     sp,
	}
}

func (st *stackContent) GetCell(row, col int) *tview.TableCell {
	if col == 0 {
		symb := ':'
		if st.sp == uint16(row) {
			symb = '►'
		}
		return tview.NewTableCell(
			fmt.Sprintf("0x%04X%c", row, symb),
		)
	}

	addr := uint16(row)

	if int(addr) >= len(st.memory) {
		return tview.NewTableCell("")
	}

	b1 := fmt.Sprintf("%02X", st.memory[addr])

	if int(addr+1) < len(st.memory) {
		return tview.NewTableCell(
			fmt.Sprintf("%s%02X", b1, st.memory[addr+1]),
		)
	}

	return tview.NewTableCell(b1)
}

func (st *stackContent) GetColumnCount() int {
	return 2
}

func (st *stackContent) GetRowCount() int {
	return len(st.memory)
}
