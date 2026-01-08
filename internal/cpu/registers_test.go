package cpu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordRegisterGetters(t *testing.T) {
	cases := []struct {
		a    byte
		b    byte
		want uint16
	}{
		{
			0x00,
			0x00,
			0x0000,
		},
		{
			0x01,
			0x00,
			0x0100,
		},
		{
			0x01,
			0x02,
			0x0102,
		},
		{
			0xFF,
			0xFF,
			0xFFFF,
		},
		{
			0xFE,
			0xFA,
			0xFEFA,
		},
	}

	regs := registers{}

	funcs := []struct {
		name string
		foo  func() uint16
		a    *byte
		b    *byte
	}{
		{
			"AF",
			regs.AF,
			&regs.A,
			&regs.F,
		},
		{
			"BC",
			regs.BC,
			&regs.B,
			&regs.C,
		},
		{
			"DE",
			regs.DE,
			&regs.D,
			&regs.E,
		},
		{
			"HL",
			regs.HL,
			&regs.H,
			&regs.L,
		},
	}

	for _, f := range funcs {
		for _, c := range cases {
			t.Run(fmt.Sprintf("%s with values %d and %d", f.name, c.a, c.b), func(t *testing.T) {
				*f.a = c.a
				*f.b = c.b

				res := f.foo()
				assert.Equal(t, c.want, res)
			})
		}
	}
}

func TestIncHL(t *testing.T) {
	cases := []struct {
		name string
		val  uint16
	}{
		{"HL = 0x0 (without overflow)", 0x0},
		{"HL = 0x1 (without overflow)", 0x1},
		{"HL = 0xFE (without overflow)", 0xFE},
		{"HL = 0xFE (without overflow)", 0xFFFE},
		{"HL = 0xFE (without overflow)", 0xFF00},
		{"HL = 0xFE (overflow of L)", 0xFF},
		{"HL = 0xFE (overflow of HL)", 0xFFFF},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			regs := registers{}
			regs.H = byte(c.val >> 8)
			regs.L = byte(c.val & 0x00FF)

			regs.IncHL()
			assert.Equal(t, c.val+1, regs.HL())
		})
	}
}

func TestDecHL(t *testing.T) {
	cases := []struct {
		name string
		val  uint16
	}{
		{"HL = 0x1 (without overflow)", 0x1},
		{"HL = 0xFE (without overflow)", 0xFF},
		{"HL = 0xFE (overflow of H)", 0xFF00},
		{"HL = 0x0 (overflow of HL)", 0x0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			regs := registers{}
			regs.H = byte(c.val >> 8)
			regs.L = byte(c.val & 0x00FF)

			regs.DecHL()
			assert.Equal(t, c.val-1, regs.HL())
		})
	}
}
