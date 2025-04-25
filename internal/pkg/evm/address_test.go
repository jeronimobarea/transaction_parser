package evm_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/jeronimobarea/transaction_parser/internal/pkg/evm"
)

func TestAddressValidate(t *testing.T) {
	testCases := []struct {
		name    string
		address evm.Address
		fails   bool
	}{
		{name: "happy path", address: "0xabcdefABCDEF0123456789abcdefABCDEF012345", fails: false},
		{name: "happy path zero value address", address: "0x0000000000000000000000000000000000000000", fails: false},
		{name: "empty", address: "", fails: true},
		{name: "too short (just prefix)", address: "0x", fails: true},
		{name: "too short", address: "0x123", fails: true},
		{name: "missing prefix", address: "1234567890abcdef1234567890abcdef12345678", fails: true},
		{name: "invalid hex characters", address: "0xGHIJKL0000000000000000000000000000000000", fails: true},
		{name: "missing characters", address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcde", fails: true},
		{name: "too many charachers", address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef", fails: true},
		{name: "wrong prefix format", address: "0Xabcdefabcdefabcdefabcdefabcdefabcdefabcdef", fails: true},
	}

	for _, tc := range testCases {
		err := tc.address.Validate()
		if !tc.fails {
			if err != nil {
				t.Errorf("Validate(%q): expected no error, got %v", tc.address, err)
			}
		} else {
			if err == nil {
				t.Errorf("Validate(%q): expected error, got nil", tc.address)
				continue
			}
			if !errors.Is(err, evm.ErrInvalidAddress) {
				t.Errorf("Validate(%q): expected error to wrap %v, got %v", tc.address, evm.ErrInvalidAddress, err)
			}
			if !strings.Contains(err.Error(), string(tc.address)) {
				t.Errorf("Validate(%q): error message %q does not contain the address", tc.address, err.Error())
			}
		}
	}
}
