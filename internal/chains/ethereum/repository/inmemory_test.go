package repository_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jeronimobarea/transaction_parser/internal/chains/ethereum"
	"github.com/jeronimobarea/transaction_parser/internal/chains/ethereum/repository"
	"github.com/jeronimobarea/transaction_parser/internal/parser"
	"github.com/jeronimobarea/transaction_parser/internal/test/evmtest"
)

func TestRepository_NewMemoryStorage(t *testing.T) {
	repo := repository.NewMemoryStorage()

	if repo.HasAddress(evmtest.EVMZeroValueAddress) {
		t.Errorf("HasAddress(%q): expected false, got true", evmtest.EVMZeroValueAddress)
	}

	txs := repo.GetTransactions(evmtest.EVMZeroValueAddress)
	if txs != nil && len(txs) != 0 {
		t.Errorf("GetTransactions(%q): expected empty slice, got %v", evmtest.EVMZeroValueAddress, txs)
	}
}

func TestRepository_AddAddress(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		repo := repository.NewMemoryStorage()

		err := repo.AddAddress(evmtest.EVMZeroValueAddress)
		if err != nil {
			t.Fatalf("AddAddress(%q): unexpected error: %v", evmtest.EVMZeroValueAddress, err)
		}

		if !repo.HasAddress(evmtest.EVMZeroValueAddress) {
			t.Errorf("HasAddress(%q): expected true after AddAddress, got false", evmtest.EVMZeroValueAddress)
		}
	})

	t.Run("trying to add an existing address", func(t *testing.T) {
		repo := repository.NewMemoryStorage()

		if err := repo.AddAddress(evmtest.EVMZeroValueAddress); err != nil {
			t.Fatalf("first AddAddress(%q): unexpected error: %v", evmtest.EVMZeroValueAddress, err)
		}

		err := repo.AddAddress(evmtest.EVMZeroValueAddress)
		if err == nil {
			t.Fatalf("second AddAddress(%q): expected error %v, got nil", evmtest.EVMZeroValueAddress, ethereum.ErrAddressConflict)
		}
		if !errors.Is(err, ethereum.ErrAddressConflict) {
			t.Errorf("second AddAddress(%q): expected error %v, got %v", evmtest.EVMZeroValueAddress, ethereum.ErrAddressConflict, err)
		}
	})
}

func TestRepository_SaveAndGetTransactions(t *testing.T) {
	var (
		repo = repository.NewMemoryStorage()

		tx1 = parser.Transaction{
			Hash:        "h1",
			From:        evmtest.EVMZeroValueAddress,
			To:          "0xabc",
			Value:       "10",
			BlockNumber: "0x1",
		}
		tx2 = parser.Transaction{
			Hash:        "h2",
			From:        "0xdef",
			To:          evmtest.EVMZeroValueAddress,
			Value:       "20",
			BlockNumber: "0x2",
		}
	)

	t.Run("insert transactions", func(t *testing.T) {
		repo.SaveTransaction(evmtest.EVMZeroValueAddress, tx1)
		repo.SaveTransaction(evmtest.EVMZeroValueAddress, tx2)

		t.Run("retrieve transactions", func(t *testing.T) {
			txs := repo.GetTransactions(evmtest.EVMZeroValueAddress)

			want := []parser.Transaction{tx1, tx2}
			if !reflect.DeepEqual(txs, want) {
				t.Errorf("GetTransactions(%q):\n got %#v\nwant %#v", evmtest.EVMZeroValueAddress, txs, want)
			}
		})
	})
}

func TestRepository_ConcurrentAccess(t *testing.T) {
	repo := repository.NewMemoryStorage()

	done := make(chan struct{})
	go func() {
		for i := 0; i < 1000; i++ {
			_ = repo.AddAddress(evmtest.EVMZeroValueAddress)
			repo.SaveTransaction(evmtest.EVMZeroValueAddress, parser.Transaction{Hash: "h"})
		}
		close(done)
	}()

	for i := 0; i < 1000; i++ {
		_ = repo.HasAddress(evmtest.EVMZeroValueAddress)
		_ = repo.GetTransactions(evmtest.EVMZeroValueAddress)
	}
	<-done

	if !repo.HasAddress(evmtest.EVMZeroValueAddress) {
		t.Errorf("HasAddress(%q): expected true after concurrent adds, got false", evmtest.EVMZeroValueAddress)
	}
}
