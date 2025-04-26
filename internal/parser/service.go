package parser

import (
	"context"
	"errors"
	"log"
	"sync"
)

const EthereumChainID = 1

var (
	ErrLoadingParser = errors.New("error loading chain parser")
	ErrCastingParser = errors.New("error casting chain parser")
)

type Service interface {
	Parser
	Register(chainID int, parser Parser)
}

type service struct {
	chainParsers sync.Map
	logger       *log.Logger
}

func NewService(logger *log.Logger) Service {
	return &service{
		logger: logger,
	}
}

func (svc *service) Register(chainID int, chainParser Parser) {
	svc.chainParsers.Store(chainID, chainParser)
}

func (svc *service) getParser(_ int) (Parser, error) {
	// By default we only load the Ethereum parser so there is no need to propagate the chainID
	// In case we want to add support for more chains we will just need to propagate the chainID across methods.
	// Usage of sync.Map in this case is just to showcase a possible method of modularity.
	// We could just use the ethereum parser as a variable instead of loading it from the map for this excercise.
	v, ok := svc.chainParsers.Load(EthereumChainID)
	if !ok {
		svc.logger.Printf("error loading parser: %d\n", EthereumChainID)
		return nil, ErrLoadingParser
	}

	parser, ok := v.(Parser)
	if !ok {
		svc.logger.Printf("error casting parser wanted (Parser) got %T\n", v)
		return nil, ErrCastingParser
	}

	return parser, nil
}

func (svc *service) GetCurrentBlock(ctx context.Context) (int64, error) {
	parser, err := svc.getParser(EthereumChainID)
	if err != nil {
		return -1, err
	}

	return parser.GetCurrentBlock(ctx)
}

func (svc *service) GetTransactions(ctx context.Context, address string) ([]Transaction, error) {
	parser, err := svc.getParser(EthereumChainID)
	if err != nil {
		return nil, err
	}

	return parser.GetTransactions(ctx, address)
}

func (svc *service) Subscribe(ctx context.Context, address string) error {
	parser, err := svc.getParser(EthereumChainID)
	if err != nil {
		return err
	}

	return parser.Subscribe(ctx, address)
}
