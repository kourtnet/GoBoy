// Package emu provides Run function for running whole emulator application
package emu

import (
	"github.com/kourtnet/GoBoy/internal/bus"
	"github.com/kourtnet/GoBoy/internal/cpu"
	"github.com/kourtnet/GoBoy/internal/die"
	"github.com/kourtnet/GoBoy/internal/flags"
	"github.com/kourtnet/GoBoy/internal/romload"
)

func Run() {
	flags := flags.Flags{}
	flags.MustParse()

	bus := &bus.Bus{}

	err := romload.Load(flags.RomFile, bus)
	if err != nil {
		die.Die(err.Error())
	}

	cpu, err := cpu.New(bus)
	if err != nil {
		// TODO: hide internal error in prod
		die.Die(err.Error())
	}

	var cpuEnd bool
	// TODO: move error check into loop
	for ; err == nil && !cpuEnd; cpuEnd, err = cpu.Step() {
	}

	if err != nil {
		die.Die(err.Error())
	}
}
