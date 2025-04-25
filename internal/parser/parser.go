package parser

import (
	"context"

	"github.com/jeronimobarea/transaction_parser/internal/pkg/evm"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock(ctx context.Context) (int64, error)

	// add address to observer
	Subscribe(ctx context.Context, address evm.Address) error

	// list of inbound or outbound transactions for an address
	GetTransactions(ctx context.Context, address evm.Address) ([]Transaction, error)
}
