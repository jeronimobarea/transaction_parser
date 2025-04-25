package ethereumtest

import (
	"context"

	"github.com/jeronimobarea/transaction_parser/internal/chains/ethereum/client"
)

type FakeClient struct {
	GetCurrentBlockResp string
	GetCurrentBlockErr  error
	GetBlockResp        []client.TransactionResponse
	GetBlockErr         error
}

func (f *FakeClient) GetBlock(_ context.Context, _ string) ([]client.TransactionResponse, error) {
	return f.GetBlockResp, f.GetBlockErr
}

func (f *FakeClient) GetCurrentBlock(_ context.Context) (string, error) {
	return f.GetCurrentBlockResp, f.GetCurrentBlockErr
}
