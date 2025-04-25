package pollers

import (
	"context"
	"log"

	"github.com/jeronimobarea/transaction_parser/internal/chains/ethereum"
	"github.com/jeronimobarea/transaction_parser/internal/chains/ethereum/client"
	"github.com/jeronimobarea/transaction_parser/internal/parser"
	"github.com/jeronimobarea/transaction_parser/internal/pkg/evm"
)

type Poller interface {
	Poll(ctx context.Context) error
}

type poller struct {
	ethClient client.Client
	repo      ethereum.Repository
	logger    *log.Logger
}

func NewPoller(ethClient client.Client, repo ethereum.Repository, logger *log.Logger) Poller {
	return &poller{
		ethClient: ethClient,
		repo:      repo,
		logger:    logger,
	}
}

func (p *poller) Poll(ctx context.Context) error {
	txs, err := p.ethClient.GetBlock(ctx, client.LatestBlock)
	if err != nil {
		p.logger.Printf("error retrieving the latest block: %v\n", err)
		return err
	}
	p.logger.Printf("[DEBUG] Latest block info: %+v\n", txs)

	for _, tx := range txs {
		var (
			from = evm.Address(tx.From)
			to   = evm.Address(tx.To)
		)

		if p.repo.HasAddress(from) {
			p.repo.SaveTransaction(from, parser.Transaction{
				Hash:        tx.Hash,
				From:        from,
				To:          to,
				Value:       tx.Value,
				BlockNumber: tx.BlockNumber,
			})
			p.logger.Printf("[INFO] new outbound transaction saved: %+v\n", tx)
			continue
		}

		if p.repo.HasAddress(to) {
			p.repo.SaveTransaction(to, parser.Transaction{
				Hash:        tx.Hash,
				From:        from,
				To:          to,
				Value:       tx.Value,
				BlockNumber: tx.BlockNumber,
			})
			p.logger.Printf("[INFO] new inbound transaction saved: %+v\n", tx)
		}
	}
	return nil
}
