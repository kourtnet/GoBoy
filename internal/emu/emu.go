// Package emu provides Run function for running whole emulator application
package emu

import (
	"time"

	"github.com/kourtnet/GoBoy/internal/bus"
	"github.com/kourtnet/GoBoy/internal/cpu"
	"github.com/kourtnet/GoBoy/internal/debugger"
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
		die.Die(err.Error())
	}

	// TODO: make cpu konstructor return pointer
	deb, err := debugger.New(&cpu, bus)
	if err != nil {
		die.Die(err.Error())
	}

	var cpuEnd bool

	// WARNING: remove later
	for err == nil && !cpuEnd {
		deb.PrintData()
		time.Sleep(time.Hour)
		cpuEnd, err = cpu.Step()
	}

	if err != nil {
		die.Die(err.Error())
	}
}
