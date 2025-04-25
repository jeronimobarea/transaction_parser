package ethereumtest

import (
	"context"

	"github.com/jeronimobarea/transaction_parser/internal/chains/ethereum/client"
)

type FakeClient struct {
	GetBlockResp []client.TransactionResponse
	GetBlockErr  error
}

func (f *FakeClient) GetBlock(_ context.Context, _ string) ([]client.TransactionResponse, error) {
	return f.GetBlockResp, f.GetBlockErr
}
