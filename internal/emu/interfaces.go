package emu

type stepper interface {
	Step() (bool, error)
}
