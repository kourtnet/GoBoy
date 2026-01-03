// Package emu provides Run function for running whole emulator application
package emu

import (
	"fmt"
	"os"

	"github.com/kourtnet/GoBoy/internal/bus"
	"github.com/kourtnet/GoBoy/internal/cpu"
	"github.com/kourtnet/GoBoy/internal/flags"
	"github.com/kourtnet/GoBoy/internal/romload"
)

func Run() {
	flags := flags.Flags{}
	flags.MustParse()

	bus := &bus.Bus{}

	err := romload.Load(flags.RomFile, bus)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	cpu, err := cpu.New(bus)
	if err != nil {
		// TODO: hide internal error in prod
		fmt.Println(err.Error())
		os.Exit(2)
	}

	var cpuEnd bool
	for ; err == nil && !cpuEnd; cpuEnd, err = cpu.Step() {
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}
