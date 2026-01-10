package core

import "errors"

var (
	ErrOutOfRange = errors.New("register address out of range")
)

func checkBoolRange(mem []bool, addr int, count int) error {
	if addr < 0 || count < 1 || addr+count > len(mem) {
		return ErrOutOfRange
	}
	return nil
}

func checkUint16Range(mem []uint16, addr int, count int) error {
	if addr < 0 || count < 1 || addr+count > len(mem) {
		return ErrOutOfRange
	}
	return nil
}
