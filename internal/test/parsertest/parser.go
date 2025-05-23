package parsertest

import (
	"context"

	"github.com/jeronimobarea/transaction_parser/internal/parser"
)

type FakeParserSvc struct {
	GetCurrentBlockResp int64
	GetCurrentBlockErr  error
	GetTransactionsResp []parser.Transaction
	GetTransactionsErr  error
	SubscribeErr        error
}

func (f *FakeParserSvc) GetCurrentBlock(_ context.Context) (int64, error) {
	return f.GetCurrentBlockResp, f.GetCurrentBlockErr
}

func (f *FakeParserSvc) GetTransactions(_ context.Context, _ string) ([]parser.Transaction, error) {
	return f.GetTransactionsResp, f.GetTransactionsErr
}

func (f *FakeParserSvc) Subscribe(_ context.Context, _ string) error {
	return f.SubscribeErr
}

func (f *FakeParserSvc) Register(_ int, _ parser.Parser) {
	panic("unimplemented")
}
