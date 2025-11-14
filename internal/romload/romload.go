// Package romload provides simple helper function to load ROM data into the RAM
package romload

import (
	"io"
	"os"

	"github.com/kourtnet/GoBoy/internal/bus"
)

func load(filename string, bus bus.Bus) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	rom, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	for i := uint16(0); i < uint16(len(rom)); i++ {
		bus.Write(i, rom[i])
	}

	return nil
}
