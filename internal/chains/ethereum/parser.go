package ethereum

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/jeronimobarea/transaction_parser/internal/parser"
	"github.com/jeronimobarea/transaction_parser/internal/pkg/evm"
	"github.com/jeronimobarea/transaction_parser/internal/pkg/svcerrors"
)

var (
	ErrAddressNotSubscribed = fmt.Errorf("%w: error address not subscribed", svcerrors.ErrNotFound)
	ErrAddressConflict      = fmt.Errorf("%w: error adress already exists", svcerrors.ErrConflict)
)

type ethereumParser struct {
	repo   Repository
	logger *log.Logger
}

func NewEthereumParser(repo Repository, logger *log.Logger) parser.Parser {
	return &ethereumParser{
		repo:   repo,
		logger: logger,
	}
}

func (p *ethereumParser) GetCurrentBlock(_ context.Context) (int64, error) {
	hexValue := p.repo.GetLastParsedBlock()
	return hexToDecimal(hexValue)
}

func hexToDecimal(hexValue string) (int64, error) {
	if len(hexValue) >= 2 && hexValue[0:2] == "0x" {
		hexValue = hexValue[2:]
	} else {
		return -1, fmt.Errorf("unexpected result format: %s", hexValue)
	}

	value, err := strconv.ParseInt(hexValue, 16, 64)
	if err != nil {
		return -1, fmt.Errorf("failed to parse hex %q: %v", hexValue, err)
	}

	return value, nil
}

func (p *ethereumParser) GetTransactions(ctx context.Context, address string) ([]parser.Transaction, error) {
	addr := evm.Address(address)
	if err := addr.Validate(); err != nil {
		p.logger.Printf("error validating address: %v\n", err)
		return nil, err
	}

	if !p.repo.HasAddress(addr) {
		return nil, ErrAddressNotSubscribed
	}
	return p.repo.GetTransactions(addr), nil
}

func (p *ethereumParser) Subscribe(_ context.Context, address string) error {
	addr := evm.Address(address)
	if err := addr.Validate(); err != nil {
		p.logger.Printf("error validating address: %v\n", err)
		return err
	}

	if p.repo.HasAddress(addr) {
		return ErrAddressConflict
	}

	return p.repo.AddAddress(addr)
}
