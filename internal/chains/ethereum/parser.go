package ethereum

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/jeronimobarea/transaction_parser/internal/chains/ethereum/client"
	"github.com/jeronimobarea/transaction_parser/internal/parser"
	"github.com/jeronimobarea/transaction_parser/internal/pkg/evm"
	"github.com/jeronimobarea/transaction_parser/internal/pkg/svcerrors"
)

var (
	ErrAddressNotSubscribed = fmt.Errorf("%w: error address not subscribed", svcerrors.ErrNotFound)
	ErrAddressConflict      = fmt.Errorf("%w: error adress already exists", svcerrors.ErrConflict)
)

type ethereumParser struct {
	client client.Client
	repo   Repository
	logger *log.Logger
}

func NewEthereumParser(client client.Client, repo Repository, logger *log.Logger) parser.Parser {
	return &ethereumParser{
		client: client,
		repo:   repo,
		logger: logger,
	}
}

func (p *ethereumParser) GetCurrentBlock(ctx context.Context) (int64, error) {
	hexValue, err := p.client.GetCurrentBlock(ctx)
	if err != nil {
		p.logger.Printf("error making call go get current block: %v\n", err)
		return -1, err
	}

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

func (p *ethereumParser) GetTransactions(ctx context.Context, address evm.Address) ([]parser.Transaction, error) {
	if err := address.Validate(); err != nil {
		p.logger.Printf("error validating address: %v\n", err)
		return nil, err
	}

	if !p.repo.HasAddress(address) {
		return nil, ErrAddressNotSubscribed
	}
	return p.repo.GetTransactions(address), nil
}

func (p *ethereumParser) Subscribe(_ context.Context, address evm.Address) error {
	if err := address.Validate(); err != nil {
		p.logger.Printf("error validating address: %v\n", err)
		return err
	}

	if p.repo.HasAddress(address) {
		return ErrAddressConflict
	}

	return p.repo.AddAddress(address)
}
