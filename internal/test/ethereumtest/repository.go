package ethereumtest

import (
	"github.com/jeronimobarea/transaction_parser/internal/parser"
	"github.com/jeronimobarea/transaction_parser/internal/pkg/evm"
)

type FakeRepo struct {
	GetTransactionsResp []parser.Transaction
	HasAddressResp      bool
	AddAddressErr       error
}

func (r FakeRepo) GetTransactions(_ evm.Address) []parser.Transaction {
	return r.GetTransactionsResp
}

func (r FakeRepo) SaveTransaction(_ evm.Address, _ parser.Transaction) {}

func (r FakeRepo) HasAddress(_ evm.Address) bool {
	return r.HasAddressResp
}

func (r FakeRepo) AddAddress(_ evm.Address) error {
	return r.AddAddressErr
}
