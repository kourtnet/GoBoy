// Package emu provides Run function for running whole emulator application
package emu

import (
	"fmt"

	"github.com/kourtnet/GoBoy/internal/bus"
	"github.com/kourtnet/GoBoy/internal/flags"
	"github.com/kourtnet/GoBoy/internal/romload"
)

func Run() {
	flags := flags.Flags{}
	flags.MustParse()

	bus := bus.Bus{}

	err := romload.Load(flags.RomFile, bus)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
