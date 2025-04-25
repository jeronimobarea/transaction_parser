package parser

import "github.com/jeronimobarea/transaction_parser/internal/pkg/evm"

type Transaction struct {
	Hash        string
	From        evm.Address
	To          evm.Address
	Value       string
	BlockNumber string
}
