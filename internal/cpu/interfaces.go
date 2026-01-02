package cpu

//go:generate mockgen -source=interfaces.go -destination=mocks/interfaces.go -package=mocks
type iBus interface {
	Read(addr uint16) (byte, error)
	Write(addr uint16, val byte) error
}
