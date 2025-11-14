// Package flags provides Flags struct containing flag parser and all parameters required by emulator
package flags

import (
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	RomFile string
}

func (f *Flags) MustParse() {
	flag.StringVar(&f.RomFile, "rom", "", "full name of the ROM file")

	flag.Parse()

	if f.RomFile == "" {
		fmt.Println("ROM filepath is empty")
		os.Exit(2)
	}
}
