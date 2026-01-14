package cpu

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kourtnet/GoBoy/internal/cpu/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_readAddr(t *testing.T) {
	cases := []struct {
		name    string
		addr    uint16
		wantErr bool
		retVal  uint8
	}{
		{
			name:    "Read completes succesfully",
			addr:    0,
			wantErr: false,
			retVal:  0xAE,
		},
		{
			name:    "Read returns error",
			addr:    1,
			wantErr: true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	busMock := mocks.NewMockiBus(ctrl)

	busMock.EXPECT().
		Read(cases[0].addr).
		Return(cases[0].retVal, nil)
	busMock.EXPECT().
		Read(cases[1].addr).
		Return(cases[1].retVal, errors.New("error"))

	cpu := CPU{}
	cpu.bus = busMock
	cpu.registers = &registers{}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cpu.registers.Temp8 = 0

			cpu.readAddr(c.addr)

			if c.wantErr {
				assert.Error(t, cpu.internalErr)
			} else {
				assert.NoError(t, cpu.internalErr)
				assert.Equal(t, cpu.registers.Temp8, c.retVal)
			}
		})
	}
}

func Test_determine8Reg(t *testing.T) {
	cpu := CPU{}
	cpu.registers = &registers{}

	cases := []struct {
		name    byte
		reg     *byte
		wantErr bool
	}{
		{
			name: 'A',
			reg:  &cpu.registers.A,
		},
		{
			name: 'B',
			reg:  &cpu.registers.B,
		},
		{
			name: 'D',
			reg:  &cpu.registers.D,
		},
		{
			name: 'H',
			reg:  &cpu.registers.H,
		},
		{
			name: 'C',
			reg:  &cpu.registers.C,
		},
		{
			name: 'E',
			reg:  &cpu.registers.E,
		},
		{
			name: 'L',
			reg:  &cpu.registers.L,
		},
		{
			name: 'T',
			reg:  &cpu.registers.Temp8,
		},
		{
			name:    'U',
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(string(c.name)+" register", func(t *testing.T) {
			res := cpu.determine8Reg(c.name)

			if c.wantErr {
				assert.Error(t, cpu.internalErr)
			} else {
				assert.NoError(t, cpu.internalErr)
				assert.True(t, c.reg == res)
			}
		})
	}
}

func Test_determine16Reg(t *testing.T) {
	cpu := CPU{}
	cpu.registers = &registers{}

	cases := []struct {
		name    string
		reg     func() uint16
		wantErr bool
	}{
		{
			name: "AF",
			reg:  cpu.registers.AF,
		},
		{
			name: "BC",
			reg:  cpu.registers.BC,
		},
		{
			name: "DE",
			reg:  cpu.registers.DE,
		},
		{
			name: "HL",
			reg:  cpu.registers.HL,
		},
		{
			name: "PC",
			reg:  cpu.registers.PC,
		},
		{
			name: "SP",
			reg:  cpu.registers.SP,
		},
		{
			name: "Temp16",
			reg:  cpu.registers.Temp16,
		},
		{
			name:    "Invalid Address",
			wantErr: true,
		},
	}

	cpu.registers.A = 0x11
	cpu.registers.F = 0x22
	cpu.registers.B = 0x33
	cpu.registers.C = 0x44
	cpu.registers.D = 0x55
	cpu.registers.E = 0x66
	cpu.registers.H = 0x77
	cpu.registers.L = 0x88
	cpu.registers.pc = 0x9999
	cpu.registers.S = 0xAA
	cpu.registers.P = 0xAA

	for _, c := range cases {
		t.Run(string(c.name)+" register", func(t *testing.T) {
			reg := cpu.determine16Reg(c.name)

			if c.wantErr {
				assert.Error(t, cpu.internalErr)
			} else {
				assert.NoError(t, cpu.internalErr)
				want := c.reg()
				get := reg()
				assert.Equal(t, want, get)
			}
		})
	}
}

func Test_determine16RegSetter(t *testing.T) {
	cpu := CPU{}
	cpu.registers = &registers{}

	cases := []struct {
		name    string
		reg     func() uint16
		wantErr bool
	}{
		{
			name: "AF",
			reg:  cpu.registers.AF,
		},
		{
			name: "BC",
			reg:  cpu.registers.BC,
		},
		{
			name: "DE",
			reg:  cpu.registers.DE,
		},
		{
			name: "HL",
			reg:  cpu.registers.HL,
		},
		{
			name: "SP",
			reg:  cpu.registers.SP,
		},
		{
			name:    "Invalid Address",
			wantErr: true,
		},
	}

	add := uint16(0x1234)
	setVal := add

	for _, c := range cases {
		t.Run(string(c.name)+" register setter", func(t *testing.T) {
			regSetter := cpu.determine16RegSetter(c.name)

			if c.wantErr {
				assert.Error(t, cpu.internalErr)
			} else {
				assert.NoError(t, cpu.internalErr)

				reg := cpu.determine16Reg(c.name)
				require.NoError(t, cpu.internalErr)

				regSetter(setVal)
				assert.Equal(t, setVal, reg())

				setVal += add
			}
		})
	}
}

func Test_determine16RegINC_DEC(t *testing.T) {
	cpu := CPU{}
	cpu.registers = &registers{}

	cases := []struct {
		name    string
		reg     func() uint16
		wantErr bool
	}{
		{
			name: "BC",
			reg:  cpu.registers.BC,
		},
		{
			name: "DE",
			reg:  cpu.registers.DE,
		},
		{
			name: "HL",
			reg:  cpu.registers.HL,
		},
		{
			name: "SP",
			reg:  cpu.registers.SP,
		},
		{
			name:    "Invalid Address",
			wantErr: true,
		},
	}

	add := uint16(0x1234)
	setVal := add

	for _, c := range cases {
		t.Run(string(c.name)+" register inc", func(t *testing.T) {
			regInc := cpu.determine16RegInc(c.name)

			if c.wantErr {
				assert.Error(t, cpu.internalErr)
			} else {
				assert.NoError(t, cpu.internalErr)

				reg := cpu.determine16Reg(c.name)
				regSetter := cpu.determine16RegSetter(c.name)
				require.NoError(t, cpu.internalErr)

				regSetter(setVal)

				regInc()

				assert.Equal(t, setVal+1, reg())

				setVal += add
			}
		})
	}

	cpu.internalErr = nil

	for _, c := range cases {
		t.Run(string(c.name)+" register dec", func(t *testing.T) {
			regDec := cpu.determine16RegDec(c.name)

			if c.wantErr {
				assert.Error(t, cpu.internalErr)
			} else {
				assert.NoError(t, cpu.internalErr)

				reg := cpu.determine16Reg(c.name)
				regSetter := cpu.determine16RegSetter(c.name)
				require.NoError(t, cpu.internalErr)

				regSetter(setVal)

				regDec()

				assert.Equal(t, setVal-1, reg())

				setVal += add
			}
		})
	}
}

func TestWithROMFiles(t *testing.T) {
	snapDir := "../../testing/snapshots/."
	files, err := os.ReadDir(snapDir)
	if err != nil {
		t.Fatal(err)
	}

	filenames := make([]string, 0, len(files))

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".snap" {
			name := strings.TrimSuffix(file.Name(), ".snap")
			filenames = append(filenames, name)
		}
	}

	ROMPathFormat := "../../testing/roms/%s.rom"
	snapPathFormat := "../../testing/snapshots/%s.snap"

	for _, name := range filenames {
		ROMPath := fmt.Sprintf(ROMPathFormat, name)
		snapPath := fmt.Sprintf(snapPathFormat, name)

		suite, err := NewTestSuite(ROMPath, snapPath)
		if err != nil {
			t.Fatal(err)
		}

		t.Run(fmt.Sprintf("test %s instructions", name), func(t *testing.T) {
			for suite.step() {
				if suite.err != nil {
					t.Error(suite.err)
				}
			}
		})
	}
}
