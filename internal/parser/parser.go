package parser

import (
	"context"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock(ctx context.Context) (int64, error)

	// add address to observer
	Subscribe(ctx context.Context, address string) error

	// list of inbound or outbound transactions for an address
	GetTransactions(ctx context.Context, address string) ([]Transaction, error)
}
