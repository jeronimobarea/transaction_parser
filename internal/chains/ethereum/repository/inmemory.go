package repository

import (
	"sync"

	"github.com/jeronimobarea/transaction_parser/internal/chains/ethereum"
	"github.com/jeronimobarea/transaction_parser/internal/parser"
	"github.com/jeronimobarea/transaction_parser/internal/pkg/evm"
)

type repository struct {
	mu              sync.RWMutex
	addresses       map[evm.Address]struct{}
	txs             map[evm.Address][]parser.Transaction
	lastParsedBlock string
}

func NewMemoryStorage() ethereum.Repository {
	return &repository{
		addresses:       make(map[evm.Address]struct{}),
		txs:             make(map[evm.Address][]parser.Transaction),
		lastParsedBlock: "0x0",
	}
}

func (r *repository) GetLastParsedBlock() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.lastParsedBlock
}

func (r *repository) AddAddress(address evm.Address) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.addresses[address]
	if !exists {
		r.addresses[address] = struct{}{}
		return nil
	}
	return ethereum.ErrAddressConflict
}

func (r *repository) HasAddress(address evm.Address) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.addresses[address]
	return exists
}

func (r *repository) SaveTransaction(address evm.Address, tx parser.Transaction) {
	r.mu.Lock()
	defer r.mu.Unlock()

	txs := r.txs[address]
	for _, savedTx := range txs {
		if savedTx.Hash == tx.Hash {
			return
		}
	}

	r.lastParsedBlock = tx.BlockNumber
	r.txs[address] = append(txs, tx)
}

func (r *repository) GetTransactions(address evm.Address) []parser.Transaction {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.txs[address]
}
