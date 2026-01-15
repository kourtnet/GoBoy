package cpu

import (
	"errors"
	"testing"

	"github.com/kourtnet/GoBoy/internal/cpu/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_fetchOpcode(t *testing.T) {
	cases := []struct {
		name    string
		PC      uint16
		wantErr bool
	}{
		{
			"no error",
			0,
			false,
		},
		{
			"error from bus",
			1,
			true,
		},
		{
			"error from unknown opcode",
			2,
			true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	busMock := mocks.NewMockiBus(ctrl)

	busMock.EXPECT().
		Read(cases[0].PC).
		Return(uint8(1), nil)
	busMock.EXPECT().
		Read(cases[1].PC).
		Return(uint8(0), errors.New("error"))
	busMock.EXPECT().
		Read(cases[2].PC).
		Return(uint8(0), nil)

	cpu := CPU{}
	cpu.Registers = &Registers{}
	cpu.bus = busMock
	cpu.instructions[1] = []func(){func() {}}
	cpu.opNum = 1

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cpu.Registers.pc = c.PC
			err := cpu.fetchOpcode()
			if c.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)

				assert.Equal(t, cpu.opNum, 0)
				assert.Equal(t, cpu.Registers.IR, uint8(1))
			}
		})
	}
}
