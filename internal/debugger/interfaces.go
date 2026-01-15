package debugger

type bus interface {
	Read(addr uint16) (byte, error)
}
