package evm

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrInvalidAddress = errors.New("error validating address")

	evmAddressRegex = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
)

type Address string

func (a Address) Validate() error {
	if !evmAddressRegex.MatchString(string(a)) {
		return fmt.Errorf("%w: %s", ErrInvalidAddress, a)
	}
	return nil
}

func (a Address) String() string {
	return string(a)
}
