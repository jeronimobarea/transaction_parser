package parser_test

import (
	"context"
	"errors"
	"log"
	"reflect"
	"testing"

	"github.com/jeronimobarea/transaction_parser/internal/parser"
	"github.com/jeronimobarea/transaction_parser/internal/test"
	"github.com/jeronimobarea/transaction_parser/internal/test/evmtest"
	"github.com/jeronimobarea/transaction_parser/internal/test/parsertest"
)

func TestService_GetCurrentBlock(t *testing.T) {
	logger := log.Default()

	t.Run("happy path", func(t *testing.T) {
		svc := parser.NewService(logger)
		svc.Register(1, &parsertest.FakeParserSvc{GetCurrentBlockResp: 42})

		got, err := svc.GetCurrentBlock(context.Background())
		if err != nil {
			t.Fatalf("GetCurrentBlock success: unexpected error %v", err)
		}
		if got != 42 {
			t.Errorf("GetCurrentBlock success: expected 42, got %d", got)
		}
	})

	t.Run("no parser", func(t *testing.T) {
		svc := parser.NewService(logger)

		_, err := svc.GetCurrentBlock(context.Background())
		if !errors.Is(err, parser.ErrLoadingParser) {
			t.Errorf("GetCurrentBlock without parser: expected %v, got %v", parser.ErrLoadingParser, err)
		}
	})
}

func TestService_GetTransactions(t *testing.T) {
	logger := log.Default()

	t.Run("happy path", func(t *testing.T) {
		var (
			expected = []parser.Transaction{{Hash: "h1", From: evmtest.EVMZeroValueAddress, To: evmtest.EVMZeroValueAddress, Value: "10", BlockNumber: "0x1"}}

			svc = parser.NewService(logger)
		)

		svc.Register(1, &parsertest.FakeParserSvc{GetTransactionsResp: expected})

		got, err := svc.GetTransactions(context.Background(), evmtest.EVMZeroValueAddress.String())
		if err != nil {
			t.Fatalf("GetTransactions success: unexpected error %v", err)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("GetTransactions success: expected %v, got %v", expected, got)
		}
	})

	t.Run("no parser", func(t *testing.T) {
		svc := parser.NewService(logger)

		_, err := svc.GetTransactions(context.Background(), evmtest.EVMZeroValueAddress.String())
		if !errors.Is(err, parser.ErrLoadingParser) {
			t.Errorf("GetTransactions no parser: expected %v, got %v", parser.ErrLoadingParser, err)
		}
	})

	t.Run("not subscribed", func(t *testing.T) {
		svc := parser.NewService(logger)

		svc.Register(1, &parsertest.FakeParserSvc{
			GetTransactionsErr: test.DummyErr,
		})

		_, err := svc.GetTransactions(context.Background(), evmtest.EVMZeroValueAddress.String())
		if !errors.Is(err, test.DummyErr) {
			t.Errorf("GetTransactions not subscribed: expected %v, got %v", test.DummyErr, err)
		}
	})
}

func TestService_Subscribe(t *testing.T) {
	var (
		logger = log.Default()

		svc = parser.NewService(logger)
	)
	svc.Register(1, &parsertest.FakeParserSvc{})

	t.Run("happy path", func(t *testing.T) {
		err := svc.Subscribe(context.Background(), evmtest.EVMZeroValueAddress.String())
		if err != nil {
			t.Fatalf("Subscribe success: expected err=nil; got err=%v", err)
		}
	})

	t.Run("parser error", func(t *testing.T) {
		svc := parser.NewService(logger)
		svc.Register(1, &parsertest.FakeParserSvc{SubscribeErr: test.DummyErr})

		err := svc.Subscribe(context.Background(), evmtest.EVMZeroValueAddress.String())
		if err == nil {
			t.Errorf("Subscribe repo error: expected failure, got err=%v", err)
		}
		if !errors.Is(err, test.DummyErr) {
			t.Errorf("Subscribe repo error: expected %v, got err=%v", test.DummyErr, err)
		}
	})
}
