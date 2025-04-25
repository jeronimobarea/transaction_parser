package pollers_test

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/jeronimobarea/transaction_parser/internal/chains/ethereum/pollers"
	"github.com/jeronimobarea/transaction_parser/internal/test"
	"github.com/jeronimobarea/transaction_parser/internal/test/ethereumtest"
)

func TestPoller_ErrorGettingBlock(t *testing.T) {
	fc := &ethereumtest.FakeClient{
		GetBlockErr: test.DummyErr,
	}
	fr := ethereumtest.FakeRepo{HasAddressResp: false}
	p := pollers.NewPoller(fc, fr, log.Default())

	err := p.Poll(context.Background())
	if err == nil || !errors.Is(err, test.DummyErr) {
		t.Fatalf("Poll() error = %v; want rpc failure", err)
	}
}
