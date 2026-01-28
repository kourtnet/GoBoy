package debugger

import (
	"fmt"

	"github.com/rivo/tview"
)

type memoryContent struct {
	*tview.TableContentReadOnly

	memory      []byte
	bytesInLine int
}

func newMemoryTable(memory []byte) *memoryContent {
	return &memoryContent{
		memory: memory,
	}
}

func (mt *memoryContent) GetCell(row, col int) *tview.TableCell {
	if col == 0 {
		return tview.NewTableCell(
			fmt.Sprintf("0x%04X:", row*mt.bytesInLine),
		)
	}

	addr := uint16(row*mt.bytesInLine + col - 1)

	// TODO: mb make a cycle
	if int(addr) >= len(mt.memory) {
		return tview.NewTableCell("")
	}

	return tview.NewTableCell(
		fmt.Sprintf("%02X", addr),
	)
}
