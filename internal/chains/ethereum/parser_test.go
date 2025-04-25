package ethereum

import (
	"context"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/jeronimobarea/transaction_parser/internal/parser"
	"github.com/jeronimobarea/transaction_parser/internal/pkg/evm"
	"github.com/jeronimobarea/transaction_parser/internal/test"
	"github.com/jeronimobarea/transaction_parser/internal/test/ethereumtest"
	"github.com/jeronimobarea/transaction_parser/internal/test/evmtest"
)

func TestHexToDecimal(t *testing.T) {
	testCases := []struct {
		name     string
		hexValue string
		want     int64
		errMsg   string
	}{
		{
			name:     "happy path",
			hexValue: "0x4b7",
			want:     1207,
		},
		{
			name:     "happy path zero value",
			hexValue: "0x0",
			want:     0,
		},
		{
			name:     "missing 0x prefix",
			hexValue: "4b7",
			want:     -1,
			errMsg:   "unexpected result format",
		},
		{
			name:     "proper prefix wrong format",
			hexValue: "0xZZZ",
			want:     -1,
			errMsg:   "failed to parse hex",
		},
		{
			name:     "prefix only",
			hexValue: "0x",
			want:     -1,
			errMsg:   "failed to parse hex",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := hexToDecimal(tc.hexValue)
			if tc.errMsg != "" {
				if err == nil {
					t.Errorf("hexToDecimal(%q): expected error, got none", tc.hexValue)
				}
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("hexToDecimal(%q): error %q does not contain %q", tc.hexValue, err.Error(), tc.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("hexToDecimal(%q): unexpected error: %v", tc.hexValue, err)
				}
				if got != tc.want {
					t.Errorf("hexToDecimal(%q): want %d, got %d", tc.hexValue, tc.want, got)
				}
			}
		})
	}
}

func TestService_GetCurrentBlock(t *testing.T) {
	var (
		ctx = context.Background()

		logger = log.Default()
		repo   = &ethereumtest.FakeRepo{}
	)

	t.Run("happy path", func(t *testing.T) {
		var (
			fc = &ethereumtest.FakeClient{
				GetCurrentBlockResp: "0x1a",
			}

			p = NewEthereumParser(fc, repo, logger)
		)

		got, err := p.GetCurrentBlock(ctx)
		if err != nil {
			t.Fatalf("GetCurrentBlock: unexpected error: %v", err)
		}
		if want := int64(0x1a); got != want {
			t.Errorf("GetCurrentBlock: want %d, got %d", want, got)
		}
	})

	t.Run("block error", func(t *testing.T) {
		var (
			ctx = context.Background()

			fc = &ethereumtest.FakeClient{
				GetCurrentBlockErr: test.DummyErr,
			}
			p = NewEthereumParser(fc, repo, logger)
		)

		got, err := p.GetCurrentBlock(ctx)
		if err == nil {
			t.Fatal("GetCurrentBlock: expected error, got none")
		}
		if got != -1 {
			t.Errorf("GetCurrentBlock on error: want -1, got %d", got)
		}
	})
}

func TestParser_GetTransactions(t *testing.T) {
	logger := log.Default()

	t.Run("happy path", func(t *testing.T) {
		var (
			ctx = context.Background()

			want = []parser.Transaction{
				{Hash: "h1", From: evmtest.EVMZeroValueAddress, To: "0x0000000000000000000000000000000000000001", Value: "10", BlockNumber: "0x2"},
				{Hash: "h2", From: "0x0000000000000000000000000000000000000001", To: evmtest.EVMZeroValueAddress, Value: "20", BlockNumber: "0x2"},
			}

			fc   = &ethereumtest.FakeClient{}
			repo = &ethereumtest.FakeRepo{
				GetTransactionsResp: want,
				HasAddressResp:      true,
			}

			p = NewEthereumParser(fc, repo, logger)
		)

		got, err := p.GetTransactions(ctx, evmtest.EVMZeroValueAddress)
		if err != nil {
			t.Fatalf("GetTransactions: unexpected error: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetTransactions: want %+v, got %+v", want, got)
		}
	})

	t.Run("invalid address", func(t *testing.T) {
		var (
			ctx = context.Background()

			repo = &ethereumtest.FakeRepo{}

			p = NewEthereumParser(&ethereumtest.FakeClient{}, repo, logger)
		)

		_, err := p.GetTransactions(ctx, evm.Address("not-an-address"))
		if err == nil {
			t.Fatal("GetTransactions with invalid address: expected error, got none")
		}
	})
}
