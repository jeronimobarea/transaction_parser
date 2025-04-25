package ethereum

import (
	"github.com/jeronimobarea/transaction_parser/internal/parser"
	"github.com/jeronimobarea/transaction_parser/internal/pkg/evm"
)

type Repository interface {
	AddAddress(address evm.Address) error
	HasAddress(address evm.Address) bool
	SaveTransaction(address evm.Address, tx parser.Transaction)
	GetTransactions(address evm.Address) []parser.Transaction
}
