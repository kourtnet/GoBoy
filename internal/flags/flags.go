// Package flags provides Flags struct containing flag parser and all parameters required by emulator
package flags

import (
	"flag"

	"github.com/kourtnet/GoBoy/internal/die"
)

type Flags struct {
	RomFile string
}

func (f *Flags) MustParse() {
	flag.StringVar(&f.RomFile, "rom", "", "path to the ROM file")

	flag.Parse()

	if f.RomFile == "" {
		die.Die("ROM filepath is empty")
	}
}
